package category_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcFailedInvalidId         = response.NewGrpcError("error", "Invalid ID", int(codes.InvalidArgument))
	ErrGrpcFailedInvalidMerchantId = response.NewGrpcError("error", "Invalid merchant ID", int(codes.InvalidArgument))

	ErrGrpcFailedInvalidYear  = response.NewGrpcError("error", "Invalid year", int(codes.InvalidArgument))
	ErrGrpcFailedInvalidMonth = response.NewGrpcError("error", "Invalid month", int(codes.InvalidArgument))

	ErrGrpcValidateCreateCategory = response.NewGrpcError("error", "validation failed: invalid create category request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateCategory = response.NewGrpcError("error", "validation failed: invalid update category request", int(codes.InvalidArgument))
)
