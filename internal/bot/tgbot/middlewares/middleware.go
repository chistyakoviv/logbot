package middlewares

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

type TgMiddlewareHandler func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) error
type TgMiddleware func(next TgMiddlewareHandler) TgMiddlewareHandler

type TgMiddlewareChainInterface interface {
	Handler(context.Context, TgMiddlewareHandler) handlers.Response
	Pipe(TgMiddleware) TgMiddlewareChainInterface
}

type middleware struct {
	mw   TgMiddleware
	prev *middleware
}

func NewMiddleware() TgMiddlewareChainInterface {
	return (*middleware)(nil)
}

func (m *middleware) Pipe(mw TgMiddleware) TgMiddlewareChainInterface {
	return &middleware{
		prev: m,
		mw:   mw,
	}
}

func (m *middleware) Handler(ctx context.Context, handler TgMiddlewareHandler) handlers.Response {
	// Traverse the middleware chain once when creating the handler
	// chain := make([]TgMiddleware, 0)
	for cur := m; cur != nil; cur = cur.prev {
		handler = cur.mw(handler)
	}
	// slices.Reverse(chain)
	// for _, mw := range chain {
	// 	handler = mw(handler)
	// }
	return func(b *gotgbot.Bot, ectx *ext.Context) error {
		err := handler(ctx, b, ectx)
		if errors.Is(err, ErrMiddlewareCanceled) {
			return nil
		}
		return err
	}
}
