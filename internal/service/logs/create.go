package logs

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/pkg/errors"
)

func (s *service) Create(ctx context.Context, logInfo *model.LogInfo) (*model.Log, error) {
	log := &model.Log{
		Data:          logInfo.Data,
		Service:       logInfo.Service,
		ContainerName: logInfo.ContainerName,
		ContainerId:   logInfo.ContainerId,
		NodeId:        logInfo.NodeId,
		Token:         logInfo.Token,
		Hash:          s.loghasher.Hash(logInfo.Data),
		CreatedAt:     s.chrono.Now(),
	}

	subscriptions, err := s.subscriptionsRepository.FindByToken(ctx, log.Token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscriptions")
	}
	if len(subscriptions) == 0 {
		return nil, errors.New("subscription not found")
	}

	now := s.chrono.Now()
	for _, subscription := range subscriptions {
		settings, err := s.chatSettingsRepository.FindOrDefaults(ctx, subscription.ChatId)
		if err != nil {
			s.logger.Error("failed to get chat settings", slogger.Err(err))
			continue
		}

		// Skip notification if chat is muted
		if ok, muteTimeRemaining := settings.IsMuted(now); ok {
			s.logger.Debug(
				"Notification wasn’t sent because this chat is muted",
				slog.Duration("mute_time_remaining", muteTimeRemaining),
			)
			continue
		}

		lastSentTimestamp, _ := s.lastSentRepository.LastSentOrDefault(ctx, &model.LastSentKey{
			ChatId: subscription.ChatId,
			Token:  log.Token,
			Hash:   log.Hash,
		})
		// Skip notification if it was sent recently
		if ok, timeSinceLastSent := settings.IsCollapsed(now, lastSentTimestamp); ok {
			s.logger.Debug(
				"Notification wasn’t sent because it falls within the collapse period",
				slog.Duration("time_since_last_sent", timeSinceLastSent),
			)
			continue
		}

		subscribers, err := s.labelsRepository.FindByChatIdAndLabel(ctx, subscription.ChatId, log.Service)
		if err != nil {
			s.logger.Error(
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

		isSilenced, silenceTimeRemaining := settings.IsSilenced(now)
		if isSilenced {
			s.logger.Debug("Notifications are silenced", slog.Duration("silence_time_remaining", silenceTimeRemaining))
		}

		message := s.tgErrorReportMessage.Create(log, subscription, subscribers)

		err = s.tgBot.SendMessage(subscription.ChatId, message, &gotgbot.SendMessageOpts{
			// Send silent notification
			DisableNotification: isSilenced,
			// ParseMode: "MarkdownV2",
			ParseMode: "Markdown",
		})
		if err != nil {
			s.logger.Error(
				"failed to send message to chat",
				slog.Int64("chat_id", subscription.ChatId),
				slogger.Err(err),
			)
		} else {
			_, err = s.lastSentRepository.UpdateOrCreate(ctx, &model.LastSent{
				ChatId:    subscription.ChatId,
				Token:     log.Token,
				Hash:      log.Hash,
				UpdatedAt: now,
			})
			if err != nil {
				s.logger.Error("failed to update last sent", slogger.Err(err))
			}
		}
	}

	return s.logsRepository.Create(ctx, log)
}
