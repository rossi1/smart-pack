package decorator

import (
	"context"
	"fmt"
	"strings"
)

func ApplyCommandDecorators[H any](handler CommandHandler[H]) CommandHandler[H] {
	return commandLoggingDecorator[H]{
		base: handler,
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
