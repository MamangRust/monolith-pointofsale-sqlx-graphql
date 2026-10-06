package service

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	mencache "github.com/MamangRust/monolith-pointofsale-role/cache"
	"github.com/MamangRust/monolith-pointofsale-role/repository"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/role_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type roleQueryService struct {
	mencache      mencache.RoleQueryCache
	roleQuery     repository.RoleQueryRepository
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

func NewRoleQueryService(
	mencache mencache.RoleQueryCache,
	roleQuery repository.RoleQueryRepository,
	logger logger.LoggerInterface,
	obs observability.TraceLoggerObservability,
) *roleQueryService {
	return &roleQueryService{
		mencache:      mencache,
		roleQuery:     roleQuery,
		logger:        logger,
		observability: obs,
	}
}

func (s *roleQueryService) FindAll(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetRolesRow, *int, error) {
	const method = "FindAll"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedRoles(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	res, totalRecords, err := s.roleQuery.FindAllRoles(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*db.GetRolesRow](
			s.logger,
			role_errors.ErrFailedFindAll.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedRoles(ctx, req, res, totalRecords)
	logSuccess("Successfully fetched roles", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
	return res, totalRecords, nil
}

func (s *roleQueryService) FindById(ctx context.Context, id int) (*db.Role, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("role.id", id))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedRoleById(ctx, id); found {
		logSuccess("Data found in cache", zap.Int("role.id", id))
		return data, nil
	}

	res, err := s.roleQuery.FindById(ctx, id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Role](
			s.logger,
			role_errors.ErrRoleNotFoundRes.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedRoleById(ctx, res)
	logSuccess("Successfully fetched role", zap.Int("role.id", id), zap.Bool("success", true))
	return res, nil
}

func (s *roleQueryService) FindByName(ctx context.Context, name string) (*db.Role, error) {
	const method = "FindByName"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("role.name", name))
	defer func() {
		end(status)
	}()

	res, err := s.roleQuery.FindByName(ctx, name)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Role](
			s.logger,
			role_errors.ErrRoleNotFoundRes.WithInternal(err),
			method,
			span,
			zap.String("role.name", name),
			zap.Error(err),
		)
	}

	logSuccess("Successfully fetched role by name", zap.String("role.name", name), zap.Bool("success", true))
	return res, nil
}

func (s *roleQueryService) FindByUserId(ctx context.Context, id int) ([]*db.Role, error) {
	const method = "FindByUserId"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user.id", id))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedRoleByUserId(ctx, id); found {
		logSuccess("Data found in cache", zap.Int("user.id", id))
		return data, nil
	}

	res, err := s.roleQuery.FindByUserId(ctx, id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.Role](
			s.logger,
			role_errors.ErrRoleNotFoundRes.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedRoleByUserId(ctx, id, res)
	logSuccess("Successfully fetched role by user ID", zap.Int("user.id", id), zap.Bool("success", true))
	return res, nil
}

func (s *roleQueryService) FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetActiveRolesRow, *int, error) {
	const method = "FindByActiveRole"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedRoleActive(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	res, totalRecords, err := s.roleQuery.FindByActiveRole(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*db.GetActiveRolesRow](
			s.logger,
			role_errors.ErrFailedFindActive.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedRoleActive(ctx, req, res, totalRecords)
	logSuccess("Successfully fetched active role", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
	return res, totalRecords, nil
}

func (s *roleQueryService) FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetTrashedRolesRow, *int, error) {
	const method = "FindByTrashedRole"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedRoleTrashed(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	res, totalRecords, err := s.roleQuery.FindByTrashedRole(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*db.GetTrashedRolesRow](
			s.logger,
			role_errors.ErrFailedFindTrashed.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedRoleTrashed(ctx, req, res, totalRecords)
	logSuccess("Successfully fetched trashed role", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
	return res, totalRecords, nil
}

func (s *roleQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
