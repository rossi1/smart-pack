package decorator

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/rossi1/smart-pack/pkg/logs"
)

type commandLoggingDecorator[C any] struct {
	base CommandHandler[C]
}

// Handle wraps the base command handler.
func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	defer func() {
		var stacktrace []byte
		if recoverr := recover(); recoverr != nil {
			err = fmt.Errorf("%v", recoverr)
			stacktrace = debug.Stack() //
		}
		logs.LogCommandExecution(ctx, generateActionName(cmd), cmd, err, stacktrace)
		if stacktrace != nil {
			panic(err.Error())
		}
	}()
	err = d.base.Handle(ctx, cmd)
	return
}

type queryLoggingDecorator[C any, R any] struct {
	base QueryHandler[C, R]
}

func (d queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	defer func() {
		var stacktrace []byte
		if recoverr := recover(); recoverr != nil {
			err = fmt.Errorf("%v", recoverr)
			stacktrace = debug.Stack()
		}
		logs.LogCommandExecution(ctx, generateActionName(cmd), cmd, err, stacktrace)
		if stacktrace != nil {
			panic(err.Error())
		}
	}()
	result, err = d.base.Handle(ctx, cmd)
	return
}
