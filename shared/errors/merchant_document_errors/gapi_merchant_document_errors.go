package merchantdocument_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcMerchantInvalidID = response.NewGrpcError("merchant_document", "Invalid merchant id", int(codes.InvalidArgument))

	ErrGrpcFailedCreateMerchantDocument = response.NewGrpcError("merchant_document", "Failed to create merchant document", int(codes.Internal))
	ErrGrpcFailedUpdateMerchantDocument = response.NewGrpcError("merchant_document", "Failed to update merchant document", int(codes.Internal))

	ErrGrpcValidateCreateMerchantDocument = response.NewGrpcError("merchant_document", "Invalid input for create merchant document", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantDocument = response.NewGrpcError("merchant_document", "Invalid input for update merchant document", int(codes.InvalidArgument))
)
