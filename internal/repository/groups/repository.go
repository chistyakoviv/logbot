package groups

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/models"
)

type GroupsRepository interface {
	Create(chatID int64, token string) (*models.Group, error)
	FindByToken(token string) (*models.Group, error)
	DeleteByToken(token string) error
}

type Repository struct {
	db db.DB
	sq sq.StatementBuilderType
}

func NewRepository(db db.DB, sq sq.StatementBuilderType) *Repository {
	return &Repository{
		db: db,
		sq: sq,
	}
}
