package handler

import (
	"context"
	"math"

	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-role/service"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/role_errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roleCommandHandleGrpc struct {
	pbrole.UnimplementedRoleCommandServiceServer
	roleQuery   service.RoleQueryService
	roleCommand service.RoleCommandService
	logger      logger.LoggerInterface
}

func NewRoleCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbrole.RoleCommandServiceServer {
	return &roleCommandHandleGrpc{
		roleQuery:   service.RoleQuery,
		roleCommand: service.RoleCommand,
		logger:      logger,
	}
}

func (s *roleCommandHandleGrpc) FindByTrashed(ctx context.Context, req *pbrole.FindAllRoleRequest) (*pbrole.ApiResponsePaginationRoleDeleteAt, error) {
	s.logger.Info("FindByTrashed roles called")

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindByTrashedRole(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed roles failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed roles success")

	return &pbrole.ApiResponsePaginationRoleDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed roles",
		Data:       mapResponsesRoleFromTrashed(roles),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *roleCommandHandleGrpc) CreateRole(ctx context.Context, reqPb *pbrole.CreateRoleRequest) (*pbrole.ApiResponseRole, error) {
	s.logger.Info("CreateRole called", zap.String("name", reqPb.Name))

	req := &requests.CreateRoleRequest{
		Name: reqPb.Name,
	}

	if err := req.Validate(); err != nil {
		return nil, role_errors.ErrGrpcValidateCreateRole
	}

	role, err := s.roleCommand.CreateRole(ctx, req)
	if err != nil {
		s.logger.Error("CreateRole failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("CreateRole success")

	return &pbrole.ApiResponseRole{
		Status:  "success",
		Message: "Successfully created role",
		Data:    mapResponseRole(role),
	}, nil
}

func (s *roleCommandHandleGrpc) UpdateRole(ctx context.Context, reqPb *pbrole.UpdateRoleRequest) (*pbrole.ApiResponseRole, error) {
	s.logger.Info("UpdateRole called", zap.Int32("id", reqPb.GetId()))

	roleID := int(reqPb.GetId())
	if roleID <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	name := reqPb.GetName()
	req := &requests.UpdateRoleRequest{
		ID:   &roleID,
		Name: name,
	}

	if err := req.Validate(); err != nil {
		return nil, role_errors.ErrGrpcValidateUpdateRole
	}

	role, err := s.roleCommand.UpdateRole(ctx, req)
	if err != nil {
		s.logger.Error("UpdateRole failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateRole success")

	return &pbrole.ApiResponseRole{
		Status:  "success",
		Message: "Successfully updated role",
		Data:    mapResponseRole(role),
	}, nil
}

func (s *roleCommandHandleGrpc) TrashedRole(ctx context.Context, req *pbrole.FindByIdRoleRequest) (*pbrole.ApiResponseRole, error) {
	s.logger.Info("TrashedRole called", zap.Int32("id", req.GetRoleId()))

	roleID := int(req.GetRoleId())
	if roleID <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommand.TrashedRole(ctx, roleID)
	if err != nil {
		s.logger.Error("TrashedRole failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedRole success")

	return &pbrole.ApiResponseRole{
		Status:  "success",
		Message: "Successfully trashed role",
		Data:    mapResponseRole(role),
	}, nil
}

func (s *roleCommandHandleGrpc) RestoreRole(ctx context.Context, req *pbrole.FindByIdRoleRequest) (*pbrole.ApiResponseRole, error) {
	s.logger.Info("RestoreRole called", zap.Int32("id", req.GetRoleId()))

	roleID := int(req.GetRoleId())
	if roleID <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommand.RestoreRole(ctx, roleID)
	if err != nil {
		s.logger.Error("RestoreRole failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreRole success")

	return &pbrole.ApiResponseRole{
		Status:  "success",
		Message: "Successfully restored role",
		Data:    mapResponseRole(role),
	}, nil
}

func (s *roleCommandHandleGrpc) DeleteRolePermanent(ctx context.Context, req *pbrole.FindByIdRoleRequest) (*pbrole.ApiResponseRoleDelete, error) {
	s.logger.Info("DeleteRolePermanent called", zap.Int32("id", req.GetRoleId()))

	id := int(req.GetRoleId())
	if id <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	_, err := s.roleCommand.DeleteRolePermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteRolePermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteRolePermanent success")

	return &pbrole.ApiResponseRoleDelete{
		Status:  "success",
		Message: "Successfully deleted role permanently",
	}, nil
}

func (s *roleCommandHandleGrpc) RestoreAllRole(ctx context.Context, req *emptypb.Empty) (*pbrole.ApiResponseRoleAll, error) {
	s.logger.Info("RestoreAllRole called")

	_, err := s.roleCommand.RestoreAllRole(ctx)
	if err != nil {
		s.logger.Error("RestoreAllRole failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllRole success")

	return &pbrole.ApiResponseRoleAll{
		Status:  "success",
		Message: "Successfully restored all roles",
	}, nil
}

func (s *roleCommandHandleGrpc) DeleteAllRolePermanent(ctx context.Context, req *emptypb.Empty) (*pbrole.ApiResponseRoleAll, error) {
	s.logger.Info("DeleteAllRolePermanent called")

	_, err := s.roleCommand.DeleteAllRolePermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllRolePermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllRolePermanent success")

	return &pbrole.ApiResponseRoleAll{
		Status:  "success",
		Message: "Successfully deleted all roles",
	}, nil
}
