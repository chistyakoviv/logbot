package last_sent

import (
	"context"
	"time"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/repository/last_sent"
)

type ServiceInterface interface {
	Update(ctx context.Context, lastSentKey *model.LastSentKey) (*model.LastSent, error)
	Get(ctx context.Context, lastSentKey *model.LastSentKey) time.Time
}

type service struct {
	lastSentRepository last_sent.RepositoryInterface
	txManager          db.TxManager
}

func NewService(lastSentRepository last_sent.RepositoryInterface, txManager db.TxManager) ServiceInterface {
	return &service{
		lastSentRepository: lastSentRepository,
		txManager:          txManager,
	}
}
