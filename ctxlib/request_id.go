package ctxlib

import "context"

type ctxKeyRequestID struct{}

func RequestID(ctx context.Context) string {
	return ctx.Value(ctxKeyRequestID{}).(string)
}

func SetRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID{}, reqID)
}
