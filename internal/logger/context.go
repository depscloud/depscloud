package logger

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"

	"go.uber.org/zap"
)

// FromContext pulls the logger off the context or returns the default.
func Extract(ctx context.Context) *zap.Logger {
	return ctxzap.Extract(ctx)
}

// PutContext puts the logger on the provided context.
func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return ctxzap.ToContext(ctx, logger)
}
