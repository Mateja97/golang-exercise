package codes

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidArgument = status.Error(codes.InvalidArgument, " request does not satisfy the requirements")
	ErrNotFound        = status.Error(codes.NotFound, "data not found")
	ErrInternal        = status.Error(codes.Internal, "internal error")
)
