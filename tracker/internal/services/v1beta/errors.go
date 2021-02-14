package v1beta

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrCancelled      = status.Errorf(codes.Canceled, "stream cancelled")
	ErrInvalidRequest = status.Errorf(codes.InvalidArgument, "invalid request")
	ErrQueryFailure   = status.Errorf(codes.Internal, "failed to query graph")
	ErrUpdateFailure  = status.Errorf(codes.Internal, "failed to update graph")
	ErrPruneFailure   = status.Errorf(codes.Internal, "failed to prune graph")
	ErrBFS            = status.Errorf(codes.InvalidArgument, "cannot call breadth-first search with another input request")
	ErrDFS            = status.Errorf(codes.InvalidArgument, "cannot call depth-first search with another input request")
)
