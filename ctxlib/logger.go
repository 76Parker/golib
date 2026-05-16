package ctxlib

import (
	"context"

	"github.com/76Parker/golib/loglib"
)

type ctxKeyLogger struct{}

func SetLoggerInContext(ctx context.Context, log loglib.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger{}, log)
}
func GetLoggerFromContext(ctx context.Context) loglib.Logger {
	logger, ok := ctx.Value(ctxKeyLogger{}).(loglib.Logger)
	if !ok {
		return nil
	}
	return logger
}
