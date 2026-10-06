package handler

import (
	"context"

	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
	pbuserrole "github.com/MamangRust/monolith-pointofsale-pb/user_role"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-role/service"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/role_errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userRoleHandleGrpc struct {
	pbuserrole.UnimplementedUserRoleServiceServer
	roleQuery   service.RoleQueryService
	roleCommand service.RoleCommandService
	logger      logger.LoggerInterface
}

func NewUserRoleHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbuserrole.UserRoleServiceServer {
	return &userRoleHandleGrpc{
		roleQuery:   service.RoleQuery,
		roleCommand: service.RoleCommand,
		logger:      logger,
	}
}

func (s *userRoleHandleGrpc) FindByUserId(ctx context.Context, req *pbuserrole.FindByIdUserRoleRequest) (*pbrole.ApiResponsesRole, error) {
	s.logger.Info("FindByUserId called", zap.Int32("userId", req.GetUserId()))

	userID := int(req.GetUserId())
	if userID <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	roles, err := s.roleQuery.FindByUserId(ctx, userID)
	if err != nil {
		s.logger.Error("FindByUserId failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindByUserId success")

	return &pbrole.ApiResponsesRole{
		Status:  "success",
		Message: "Successfully fetched role by user ID",
		Data:    mapResponsesRoleFromDB(roles),
	}, nil
}

func (s *userRoleHandleGrpc) AssignRoleToUser(ctx context.Context, reqPb *pbuserrole.AssignRoleToUserRequest) (*pbuserrole.ApiResponseUserRole, error) {
	s.logger.Info("AssignRoleToUser called",
		zap.Int32("userId", reqPb.GetUserId()),
		zap.Int32("roleId", reqPb.GetRoleId()),
	)

	if reqPb.GetUserId() <= 0 || reqPb.GetRoleId() <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	userRole, err := s.roleCommand.AssignRoleToUser(ctx, &requests.CreateUserRoleRequest{
		UserId: int(reqPb.GetUserId()),
		RoleId: int(reqPb.GetRoleId()),
	})
	if err != nil {
		s.logger.Error("AssignRoleToUser failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("AssignRoleToUser success")

	return &pbuserrole.ApiResponseUserRole{
		Status:  "success",
		Message: "Successfully assigned role to user",
		Data:    mapResponseUserRole(userRole),
	}, nil
}

func (s *userRoleHandleGrpc) RemoveRoleFromUser(ctx context.Context, reqPb *pbuserrole.RemoveRoleFromUserRequest) (*emptypb.Empty, error) {
	s.logger.Info("RemoveRoleFromUser called",
		zap.Int32("userId", reqPb.GetUserId()),
		zap.Int32("roleId", reqPb.GetRoleId()),
	)

	if reqPb.GetUserId() <= 0 || reqPb.GetRoleId() <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	if err := s.roleCommand.RemoveRoleFromUser(ctx, &requests.RemoveUserRoleRequest{
		UserId: int(reqPb.GetUserId()),
		RoleId: int(reqPb.GetRoleId()),
	}); err != nil {
		s.logger.Error("RemoveRoleFromUser failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RemoveRoleFromUser success")

	return &emptypb.Empty{}, nil
}
