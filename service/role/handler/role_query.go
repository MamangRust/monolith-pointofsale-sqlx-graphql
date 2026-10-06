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
)

type roleQueryHandleGrpc struct {
	pbrole.UnimplementedRoleQueryServiceServer
	roleQuery service.RoleQueryService
	logger    logger.LoggerInterface
}

func NewRoleQueryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbrole.RoleQueryServiceServer {
	return &roleQueryHandleGrpc{
		roleQuery: service.RoleQuery,
		logger:    logger,
	}
}

func (s *roleQueryHandleGrpc) FindAllRole(ctx context.Context, req *pbrole.FindAllRoleRequest) (*pbrole.ApiResponsePaginationRole, error) {
	s.logger.Info("FindAllRole called", zap.Int32("page", req.GetPage()))

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

	roles, totalRecords, err := s.roleQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAllRoles failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAllRoles success")

	return &pbrole.ApiResponsePaginationRole{
		Status:     "success",
		Message:    "Successfully fetched role records",
		Data:       mapResponsesRole(roles),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *roleQueryHandleGrpc) FindByIdRole(ctx context.Context, req *pbrole.FindByIdRoleRequest) (*pbrole.ApiResponseRole, error) {
	s.logger.Info("FindByIdRole called", zap.Int32("roleId", req.GetRoleId()))

	roleID := int(req.GetRoleId())
	if roleID <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleQuery.FindById(ctx, roleID)
	if err != nil {
		s.logger.Error("FindByIdRole failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindByIdRole success")

	return &pbrole.ApiResponseRole{
		Status:  "success",
		Message: "Successfully fetched role",
		Data:    mapResponseRole(role),
	}, nil
}

func (s *roleQueryHandleGrpc) FindByNameRole(ctx context.Context, req *pbrole.FindByNameRoleRequest) (*pbrole.ApiResponseRole, error) {
	s.logger.Info("FindByNameRole called", zap.String("name", req.GetName()))

	name := req.GetName()
	if name == "" {
		return nil, role_errors.ErrGrpcRoleInvalidName
	}

	role, err := s.roleQuery.FindByName(ctx, name)
	if err != nil {
		s.logger.Error("FindByNameRole failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindByNameRole success")

	return &pbrole.ApiResponseRole{
		Status:  "success",
		Message: "Successfully fetched role",
		Data:    mapResponseRole(role),
	}, nil
}

func (s *roleQueryHandleGrpc) FindByUserId(ctx context.Context, req *pbrole.FindByIdUserRoleRequest) (*pbrole.ApiResponsesRole, error) {
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

func (s *roleQueryHandleGrpc) FindByActive(ctx context.Context, req *pbrole.FindAllRoleRequest) (*pbrole.ApiResponsePaginationRoleDeleteAt, error) {
	s.logger.Info("FindByActive roles called", zap.Int32("page", req.GetPage()))

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

	roles, totalRecords, err := s.roleQuery.FindByActiveRole(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive roles failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive roles success")

	return &pbrole.ApiResponsePaginationRoleDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active roles",
		Data:       mapResponsesRoleFromActive(roles),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}
