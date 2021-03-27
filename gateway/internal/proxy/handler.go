package proxy

import (
	"io"

	"github.com/depscloud/depscloud/internal/logger"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var proxyDesc = &grpc.StreamDesc{
	ClientStreams: true,
	ServerStreams: true,
}

func forwardClientToServer(log *zap.Logger, in grpc.ServerStream, out grpc.ClientStream) chan error {
	errChan := make(chan error)
	go func() {
		f := &frame{}
		for i := 0; ; i++ {
			if err := in.RecvMsg(f); err != nil {
				errChan <- err
				return
			}

			log.Debug("c2s_msg_recv", zap.String("payload", f.String()))

			if err := out.SendMsg(f); err != nil {
				errChan <- err
				return
			}
		}
	}()
	return errChan
}

func forwardServerToClient(log *zap.Logger, in grpc.ClientStream, out grpc.ServerStream) chan error {
	errChan := make(chan error)
	go func() {
		// header blocks until it's received

		header, err := in.Header()
		if err != nil {
			errChan <- err
			return
		}

		err = out.SetHeader(header)
		if err != nil {
			errChan <- err
			return
		}

		f := &frame{}
		for {
			if err := in.RecvMsg(f); err != nil {
				errChan <- err
				return
			}

			log.Debug("s2c_msg_recv", zap.String("payload", f.String()))

			if err := out.SendMsg(f); err != nil {
				errChan <- err
				return
			}
		}
	}()
	return errChan
}

func passthru(log *zap.Logger, ss grpc.ServerStream, cs grpc.ClientStream) error {
	clientChan := forwardClientToServer(log, ss, cs)
	serverChan := forwardServerToClient(log, cs, ss)

	for i := 0; i < 2; i++ {
		select {
		case clientErr := <-clientChan:
			if clientErr != io.EOF {
				return clientErr
			}

			_ = cs.CloseSend()
		case serverErr := <-serverChan:
			ss.SetTrailer(cs.Trailer())

			// special case eof for streaming operations
			if serverErr == io.EOF {
				return nil
			} else if serverErr != nil {
				return serverErr
			}
		}
	}
	return nil
}

// UnknownServiceHandler returns a grpc.StreamHandler that uses the provided Router
// to direct messages to different backends.
func UnknownServiceHandler(router *Router) grpc.StreamHandler {
	return func(srv interface{}, ss grpc.ServerStream) error {
		ctx := ss.Context()
		log := logger.Extract(ctx)

		fullMethodName, ok := grpc.Method(ctx)
		if !ok {
			return status.Error(codes.Internal, "fullMethodName not present")
		}

		cc, err := router.Route(fullMethodName)
		if err != nil {
			return err
		}

		cs, err := grpc.NewClientStream(ctx, proxyDesc, cc, fullMethodName)
		if err != nil {
			return err
		}

		return passthru(log, ss, cs)
	}
}
