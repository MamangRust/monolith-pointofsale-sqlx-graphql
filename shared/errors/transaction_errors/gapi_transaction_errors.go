package transaction_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidID         = response.NewGrpcError("error", "invalid ID", int(codes.InvalidArgument))
	ErrGrpcInvalidMonth      = response.NewGrpcError("error", "invalid month", int(codes.InvalidArgument))
	ErrGrpcInvalidYear       = response.NewGrpcError("error", "invalid year", int(codes.InvalidArgument))
	ErrGrpcInvalidMerchantId = response.NewGrpcError("error", "invalid merchant ID", int(codes.InvalidArgument))

	ErrGrpcValidateCreateTransaction = response.NewGrpcError("error", "validation failed: invalid create transaction request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateTransaction = response.NewGrpcError("error", "validation failed: invalid update transaction request", int(codes.InvalidArgument))
)
