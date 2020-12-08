package logger

import (
	"context"
	"strings"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryServerInterceptor attaches a logger to the context.
func UnaryServerInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(parent context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		service, method := parseServiceAndMethod(info.FullMethod)
		log := log.With(
			zap.String("grpc_service", service),
			zap.String("grpc_method", method))

		ctx := ToContext(parent, log)

		return handler(ctx, req)
	}
}

// StreamServerInterceptor attaches a logger to the context.
func StreamServerInterceptor(log *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service, method := parseServiceAndMethod(info.FullMethod)
		log := log.With(
			zap.String("grpc_service", service),
			zap.String("grpc_method", method))

		return handler(srv, &serverStream{
			log:      log,
			delegate: ss,
		})
	}
}

func parseServiceAndMethod(fullMethod string) (string, string) {
	parts := strings.Split(fullMethod, "/")
	return parts[1], parts[2]
}

type serverStream struct {
	log      *zap.Logger
	delegate grpc.ServerStream
}

func (s *serverStream) SetHeader(md metadata.MD) error {
	return s.delegate.SetHeader(md)
}

func (s *serverStream) SendHeader(md metadata.MD) error {
	return s.delegate.SendHeader(md)
}

func (s *serverStream) SetTrailer(md metadata.MD) {
	s.delegate.SetTrailer(md)
}

func (s *serverStream) Context() context.Context {
	ctx := s.delegate.Context()
	ctx = ToContext(ctx, s.log)
	return ctx
}

func (s *serverStream) SendMsg(m interface{}) error {
	return s.delegate.SendMsg(m)
}

func (s *serverStream) RecvMsg(m interface{}) error {
	return s.delegate.RecvMsg(m)
}

var _ grpc.ServerStream = &serverStream{}
