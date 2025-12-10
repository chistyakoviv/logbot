package log

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/chistyakoviv/logbot/internal/bot"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/http/handlers"
	"github.com/chistyakoviv/logbot/internal/lib/http/response"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/loghasher"
	"github.com/chistyakoviv/logbot/internal/markdown"
	"github.com/chistyakoviv/logbot/internal/model"
	srvChatSettings "github.com/chistyakoviv/logbot/internal/service/chat_settings"
	srvLabels "github.com/chistyakoviv/logbot/internal/service/labels"
	srvLastSent "github.com/chistyakoviv/logbot/internal/service/last_sent"
	srvLogs "github.com/chistyakoviv/logbot/internal/service/logs"
	srvSubscriptions "github.com/chistyakoviv/logbot/internal/service/subscriptions"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const stackTraceKey = "stack"

func New(
	ctx context.Context,
	logger *slog.Logger,
	validation handlers.Validator,
	tgBot bot.Bot,
	loghasher loghasher.HasherInterface,
	markdowner markdown.MarkdownerInterface,
	logs srvLogs.ServiceInterface,
	subscriptions srvSubscriptions.ServiceInterface,
	chatSettings srvChatSettings.ServiceInterface,
	labels srvLabels.ServiceInterface,
	lastSent srvLastSent.ServiceInterface,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LogRequest

		// TODO: implement middlewares to handle request parsing
		// in case there will be more than one handler
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Encounter such error if request body is empty
			// Handle it separately
			logger.Error("request body is empty")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty request"))

			return
		}
		if err != nil {
			logger.Error("failed to decode request body", slogger.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		logger.Debug("request body is decoded", slog.Any("request", req))

		if err := validation.Struct(req); err != nil {
			validationErr := err.(validator.ValidationErrors)

			logger.Error("invalid request", slogger.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validationErr))

			return
		}

		subscribedChats, err := subscriptions.FindChatsByToken(ctx, req.Token)
		if err != nil {
			logger.Error("failed to check subscription", slogger.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to check subscription"))

			return
		}
		if len(subscribedChats) == 0 {
			logger.Error("subscription not found")

			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("subscription not found"))

			return
		}

		logInfo := &model.LogInfo{
			Data:          req.Log,
			Service:       req.Labels.Service,
			ContainerName: req.Labels.ContainerName,
			ContainerId:   req.Labels.ContainerId,
			Node:          req.Labels.Node,
			NodeId:        req.Labels.NodeId,
			Token:         req.Token,
			Hash:          loghasher.Hash(req.Log),
		}
		log, err := logs.Create(ctx, logInfo)
		if err != nil {
			logger.Error("failed to create log", slogger.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create log"))

			return
		}

		now := time.Now().UTC()
		var settings *model.ChatSettings
		for _, chatId := range subscribedChats {
			settings, err = chatSettings.Find(ctx, chatId)
			if err != nil {
				if !errors.Is(err, db.ErrNotFound) {
					logger.Error("failed to find chat settings", slogger.Err(err))
					continue
				}
				// Use default chat settings
				settings = &model.ChatSettings{}
			}

			if !settings.MuteUntil.IsZero() && now.Before(settings.MuteUntil) {
				// Chat is silenced
				logger.Debug(
					"Notification wasn’t sent because this chat is muted",
					slog.Duration("mute_time_remaining",
						settings.MuteUntil.Sub(now)),
				)
				continue
			}

			lastSentTimestamp := lastSent.Get(ctx, &model.LastSentKey{
				ChatId: chatId,
				Token:  log.Token,
				Hash:   log.Hash,
			})
			if err != nil && !errors.Is(err, db.ErrNotFound) {
				logger.Error("failed to find former log with token and hash", slogger.Err(err))
				continue
			}
			timeSinceLastSent := now.Sub(lastSentTimestamp)
			if err == nil && settings.CollapsePeriod > 0 && timeSinceLastSent < settings.CollapsePeriod {
				// Notification falls within collapse period
				logger.Debug(
					"Notification wasn’t sent because it falls within the collapse period",
					slog.Duration("time_since_last_sent", timeSinceLastSent),
				)
				continue
			}

			subscribers, err := labels.FindByLabel(ctx, log.Service)
			if err != nil {
				logger.Error(
					"failed to find subscribers",
					slog.Attr{Key: "label", Value: slog.StringValue(log.Service)},
					slogger.Err(err),
				)

				continue
			}

			if len(subscribers) == 0 {
				// No subscribers for this label in this chat
				continue
			}

			var message bytes.Buffer
			for i, subscriber := range subscribers {
				message.WriteString("@")
				message.WriteString(subscriber.Username)
				if i < len(subscribers)-1 {
					message.WriteString(", ")
				}
			}

			message.WriteString("\n\n*Info*\n")

			message.WriteString("service: `")
			message.WriteString(log.Service)
			message.WriteString("`\n")

			if len(log.ContainerId) > 0 {
				message.WriteString("container id: `")
				message.WriteString(log.ContainerId)
				message.WriteString("`\n")
			}

			if len(log.NodeId) > 0 {
				message.WriteString("node id: `")
				message.WriteString(log.NodeId)
				message.WriteString("`\n")
			}
			message.WriteString("\n*Data*\n")

			var decodedData map[string]string
			err = json.Unmarshal([]byte(log.Data), &decodedData)

			if err == nil {
				keys := make([]string, 0, len(decodedData))
				for key := range decodedData {
					if key != stackTraceKey {
						keys = append(keys, key)
					}
				}
				// Always print entries in the same order
				sort.Strings(keys)
				for _, key := range keys {
					message.WriteString(markdowner.Escape(key))
					message.WriteString(": ")
					message.WriteString(markdowner.Escape(decodedData[key]))
					message.WriteString("\n")
				}

				// Print stack trace in code block
				if code, ok := decodedData[stackTraceKey]; ok {
					message.WriteString("```\n")
					message.WriteString(code)
					message.WriteString("\n```")
				}
			} else {
				// Log data is in an unexpected format, print it in code block
				message.WriteString("```\n")
				message.WriteString(log.Data)
				message.WriteString("\n```")
			}

			err = tgBot.SendMessage(chatId, message.String(), &gotgbot.SendMessageOpts{
				// ParseMode: "MarkdownV2",
				ParseMode: "Markdown",
			})
			if err != nil {
				logger.Error(
					"failed to send message to chat",
					slog.Int64("chat_id", chatId),
					slogger.Err(err),
				)
			} else {
				_, err = lastSent.Update(ctx, &model.LastSentKey{
					ChatId: chatId,
					Token:  log.Token,
					Hash:   log.Hash,
				})
				if err != nil {
					logger.Error("failed to update last sent", slogger.Err(err))
				}
			}
		}

		render.JSON(w, r, LogResponse{
			Response: response.OK(),
			Id:       log.Id,
		})
	}
}
