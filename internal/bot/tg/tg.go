package tg

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/bot"
	"github.com/chistyakoviv/logbot/internal/config"
)

type TgBot struct {
	cfg *config.Config
}

func New(cfg *config.Config) bot.Bot {
	return &TgBot{
		cfg: cfg,
	}
}

func (b *TgBot) Run(ctx context.Context) error {
	return nil
}
