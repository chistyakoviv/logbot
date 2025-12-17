package logs

import (
	"context"
	"log/slog"
	"time"

	"github.com/chistyakoviv/logbot/internal/bot"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/messages"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/lib/chrono"
	"github.com/chistyakoviv/logbot/internal/lib/loghasher"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/repository/chat_settings"
	"github.com/chistyakoviv/logbot/internal/repository/labels"
	"github.com/chistyakoviv/logbot/internal/repository/last_sent"
	"github.com/chistyakoviv/logbot/internal/repository/logs"
	"github.com/chistyakoviv/logbot/internal/repository/subscriptions"
	"github.com/google/uuid"
)

type ServiceInterface interface {
	Create(ctx context.Context, log *model.LogInfo) (*model.Log, error)
	FindAllByToken(ctx context.Context, token uuid.UUID) ([]*model.Log, error)
	Delete(ctx context.Context, id int) error
	DeleteByToken(ctx context.Context, token uuid.UUID) error
	DeleteByHash(ctx context.Context, hash string) error
	DeleteOlderThan(ctx context.Context, timestamp time.Time) error
}

type service struct {
	logsRepository          logs.RepositoryInterface
	subscriptionsRepository subscriptions.RepositoryInterface
	chatSettingsRepository  chat_settings.RepositoryInterface
	labelsRepository        labels.RepositoryInterface
	lastSentRepository      last_sent.RepositoryInterface
	loghasher               loghasher.HasherInterface
	chrono                  chrono.Chrono
	tgBot                   bot.Bot
	tgErrorReportMessage    messages.ErrorReportMessageInterface
	logger                  *slog.Logger
	txManager               db.TxManager
}

func NewService(
	logsRepository logs.RepositoryInterface,
	subscriptionsRepository subscriptions.RepositoryInterface,
	chatSettingsRepository chat_settings.RepositoryInterface,
	labelsRepository labels.RepositoryInterface,
	lastSentRepository last_sent.RepositoryInterface,
	loghasher loghasher.HasherInterface,
	chrono chrono.Chrono,
	tgBot bot.Bot,
	tgErrorReportMessage messages.ErrorReportMessageInterface,
	logger *slog.Logger,
	txManager db.TxManager,
) ServiceInterface {
	return &service{
		logsRepository:          logsRepository,
		subscriptionsRepository: subscriptionsRepository,
		chatSettingsRepository:  chatSettingsRepository,
		labelsRepository:        labelsRepository,
		lastSentRepository:      lastSentRepository,
		loghasher:               loghasher,
		chrono:                  chrono,
		tgBot:                   tgBot,
		tgErrorReportMessage:    tgErrorReportMessage,
		logger:                  logger,
		txManager:               txManager,
	}
}
