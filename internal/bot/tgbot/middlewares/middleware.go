package middlewares

import (
	"context"
	"errors"
	"slices"

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
	prev    *middleware
	handler TgMiddlewareHandler
}

func NewMiddleware() TgMiddlewareInterface {
	return (*middleware)(nil)
}

func (m *middleware) Pipe(next TgMiddlewareHandler) TgMiddlewareInterface {
	return &middleware{
		prev:    m,
		handler: next,
	}
}

func (m *middleware) Handler(ctx context.Context) handlers.Response {
	// Traverse the middleware chain once when creating the handler
	handlers := make([]TgMiddlewareHandler, 0)
	for cur := m; cur != nil; cur = cur.prev {
		handlers = append(handlers, cur.handler)
	}
	slices.Reverse(handlers)
	return func(b *gotgbot.Bot, ectx *ext.Context) error {
		var err error
		for _, handler := range handlers {
			if ctx, err = handler(ctx, b, ectx); err != nil {
				// Middleware can break the chain (e.g. if auth failed),
				// so return immediately in this case
				if errors.Is(err, ErrMiddlewareCanceled) {
					return nil
				}
				return err
			}
		}
		return nil
	}
}
