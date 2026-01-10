package command

import (
	"context"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
)

type CommandChainInterface interface {
	Add(middlewares ...middlewares.TgMiddlewareHandler) CommandChainInterface
	Build() middlewares.TgMiddlewareHandler
}

type commandChain struct {
	middlewares []middlewares.TgMiddlewareHandler
}

func NewCommandChain() CommandChainInterface {
	return &commandChain{}
}

func (c *commandChain) Add(middlewares ...middlewares.TgMiddlewareHandler) CommandChainInterface {
	c.middlewares = append(c.middlewares, middlewares...)
	return c
}

func (c *commandChain) Build() middlewares.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) error {
		for _, mw := range c.middlewares {
			if err := mw(ctx, b, ectx); err != nil {
				return err
			}
		}
		return nil
	}
}
