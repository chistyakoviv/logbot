package log

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/chistyakoviv/logbot/internal/http/handlers"
	"github.com/chistyakoviv/logbot/internal/lib/http/response"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/model"
	srvLogs "github.com/chistyakoviv/logbot/internal/service/logs"
	"github.com/go-chi/render"
)

func New(
	ctx context.Context,
	logger *slog.Logger,
	validation handlers.Validator,
	logs srvLogs.ServiceInterface,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			_ = r.Body.Close()
		}()

		var raw json.RawMessage

		// TODO: implement middlewares to handle request parsing
		// in case there will be more than one handler
		err := render.DecodeJSON(r.Body, &raw)
		if errors.Is(err, io.EOF) {
			// Encounter such error if request body is empty
			// Handle it separately
			logger.Error("request body is empty")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty request"))

			return
		}

		var reqBatch []LogRequest
		switch raw[0] {
		// Request from Fluend
		case '{':
			var req LogRequest
			if err := json.Unmarshal(raw, &req); err != nil {
				logger.Error("failed to decode request body", slogger.Err(err))

				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.Error("failed to decode request"))
				return
			}
			reqBatch = append(reqBatch, req)
		// Request from Fluent Bit
		case '[':
			if err := json.Unmarshal(raw, &reqBatch); err != nil {
				logger.Error("failed to decode request body", slogger.Err(err))

				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.Error("failed to decode request"))
				return
			}
		default:
			logger.Error("failed to decode request body", slogger.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		logger.Debug("request body is decoded", slog.Any("request", reqBatch))

		for _, req := range reqBatch {
			if err := validation.Struct(req); err != nil {
				// validationErr := err.(validator.ValidationErrors)
				logger.Error("invalid request entry", slogger.Err(err))
				continue
			}

			logInfo := &model.LogInfo{
				Data:          req.Log,
				Service:       req.Labels.Service,
				ContainerName: req.Labels.ContainerName,
				ContainerId:   req.Labels.ContainerId,
				Node:          req.Labels.Node,
				NodeId:        req.Labels.NodeId,
				Token:         req.Token,
			}
			_, err := logs.Create(ctx, logInfo)
			if err != nil {
				logger.Error("failed to create log", slogger.Err(err))
			}
		}

		render.JSON(w, r, LogResponse{
			Response: response.OK(),
		})
	}
}
