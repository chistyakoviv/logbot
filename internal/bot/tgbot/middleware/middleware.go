package middleware

import (
	"context"
	"errors"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

var (
	ErrMiddlewareCanceled = errors.New("middleware canceled")
)

type TgMiddlewareHandler func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) (context.Context, error)

type TgMiddlewareInterface interface {
	Handler(ctx context.Context) handlers.Response
	Pipe(next TgMiddlewareHandler) TgMiddlewareInterface
}

type middleware struct {
	handlers []TgMiddlewareHandler
}

func NewMiddleware() TgMiddlewareInterface {
	return &middleware{}
}

func (m middleware) Pipe(next TgMiddlewareHandler) TgMiddlewareInterface {
	m.handlers = append(m.handlers, next)
	return &m
}

func (m middleware) Handler(ctx context.Context) handlers.Response {
	return func(b *gotgbot.Bot, ectx *ext.Context) error {
		var err error
		for _, handler := range m.handlers {
			if ctx, err = handler(ctx, b, ectx); err != nil {
				if errors.Is(err, ErrMiddlewareCanceled) {
					return nil
				}
				return err
			}
		}
		return nil
	}
}
