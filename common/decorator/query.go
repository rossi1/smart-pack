package decorator

import (
	"context"
)

func ApplyQueryDecorators[H any, R any](handler QueryHandler[H, R]) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		base: handler,
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
