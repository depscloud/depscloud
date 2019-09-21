package main

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrModuleNotFound occurs when a module cannot be found in the graph
	ErrModuleNotFound = status.Error(codes.NotFound, "failed to locate module")
	// ErrPartialDeletion occurs when a partial deletion occurs during Put
	ErrPartialDeletion = status.Error(codes.Internal, "failed to delete removed edges")
	// ErrPartialInsertion occurs when a partial insertion occurs during Put
	ErrPartialInsertion = status.Error(codes.Internal, "failed to insert new edges")
	// ErrUnimplemented occurs when a method has not yet been implemented
	ErrUnimplemented = status.Error(codes.Unimplemented, "unimplemented")
	// ErrUnsupported occurs when calling a rw method on a read only service
	ErrUnsupported = status.Error(codes.NotFound, "read only")
)

func main() {}
