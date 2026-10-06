package product_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidID = response.NewGrpcError("error", "invalid ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateProduct = response.NewGrpcError("error", "validation failed: invalid create product request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateProduct = response.NewGrpcError("error", "validation failed: invalid update product request", int(codes.InvalidArgument))
)
