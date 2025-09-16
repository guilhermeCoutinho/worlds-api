package utils

import (
	"context"

	"github.com/sirupsen/logrus"
)

type LoggerCtxKeyType string

const LoggerCtxKey = LoggerCtxKeyType("logger_ctx_key")

func LoggerFromCtx(ctx context.Context) logrus.FieldLogger {
	return ctx.Value(LoggerCtxKey).(logrus.FieldLogger)
}
