package logger

import "context"

type CtxKeyRequestID struct{}

func GetRequestID(ctx context.Context) string {
	if val, ok := ctx.Value(CtxKeyRequestID{}).(string); ok {
		return val
	}
	return ""
}
