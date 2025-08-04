package logs

import (
	"context"

	"github.com/sirupsen/logrus"
)

func LogCommandExecution(ctx context.Context, commandName string, cmd interface{}, err error, stacktrace []byte) {
	if err != nil {
		logFields := logrus.Fields{
			"cmd": cmd,
		}
		if stacktrace != nil {
			logFields["stack"] = string(stacktrace)
		}
		logFields["error"] = err.Error()

		logrus.WithFields(logFields).Error(ctx, commandName+" command failed")
		return
	}

	logrus.WithFields(logrus.Fields{
		"cmd": cmd,
	}).Info(ctx, commandName+" command succeeded")
}
