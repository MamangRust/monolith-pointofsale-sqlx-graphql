package user_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcUserNotFound  = response.NewGrpcError("error", "User not found", int(codes.NotFound))
	ErrGrpcUserInvalidId = response.NewGrpcError("error", "Invalid User ID", int(codes.NotFound))

	ErrGrpcFailedFindAll     = response.NewGrpcError("error", "Failed to fetch Users", int(codes.Internal))
	ErrGrpcFailedFindActive  = response.NewGrpcError("error", "Failed to fetch active Users", int(codes.Internal))
	ErrGrpcFailedFindTrashed = response.NewGrpcError("error", "Failed to fetch trashed Users", int(codes.Internal))

	ErrGrpcFailedCreateUser   = response.NewGrpcError("error", "Failed to create User", int(codes.Internal))
	ErrGrpcFailedUpdateUser   = response.NewGrpcError("error", "Failed to update User", int(codes.Internal))
	ErrGrpcValidateCreateUser = response.NewGrpcError("error", "validation failed: invalid create User request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateUser = response.NewGrpcError("error", "validation failed: invalid update User request", int(codes.InvalidArgument))

	ErrGrpcFailedTrashedUser     = response.NewGrpcError("error", "Failed to move User to trash", int(codes.Internal))
	ErrGrpcFailedRestoreUser     = response.NewGrpcError("error", "Failed to restore User", int(codes.Internal))
	ErrGrpcFailedDeletePermanent = response.NewGrpcError("error", "Failed to delete User permanently", int(codes.Internal))

	ErrGrpcFailedRestoreAll = response.NewGrpcError("error", "Failed to restore all Users", int(codes.Internal))
	ErrGrpcFailedDeleteAll  = response.NewGrpcError("error", "Failed to delete all Users permanently", int(codes.Internal))
)
