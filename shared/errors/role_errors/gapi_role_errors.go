package role_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcRoleNotFound  = response.NewGrpcError("error", "Role not found", int(codes.NotFound))
	ErrGrpcRoleInvalidId = response.NewGrpcError("error", "Invalid Role ID", int(codes.NotFound))

	ErrGrpcValidateCreateRole = response.NewGrpcError("error", "validation failed: invalid create Role request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateRole = response.NewGrpcError("error", "validation failed: invalid update Role request", int(codes.InvalidArgument))

	ErrGrpcRoleInvalidName = response.NewGrpcError("error", "Invalid Role name", int(codes.InvalidArgument))
)
