package codes

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidArgument = status.Error(codes.InvalidArgument, " request does not satisfy the requirements")
	NotFound        = status.Error(codes.NotFound, "data not found")
	Internal        = status.Error(codes.NotFound, "internal error")
)
