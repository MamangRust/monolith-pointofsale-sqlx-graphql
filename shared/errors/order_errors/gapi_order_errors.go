package order_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidYear             = response.NewGrpcError("error", "Invalid year", int(codes.InvalidArgument))
	ErrGrpcInvalidMonth            = response.NewGrpcError("error", "Invalid month", int(codes.InvalidArgument))
	ErrGrpcFailedInvalidMerchantId = response.NewGrpcError("error", "Invalid merchant ID", int(codes.InvalidArgument))
	ErrGrpcFailedInvalidId         = response.NewGrpcError("error", "Invalid ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateOrder = response.NewGrpcError("error", "validation failed: invalid create order request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateOrder = response.NewGrpcError("error", "validation failed: invalid update order request", int(codes.InvalidArgument))
)
