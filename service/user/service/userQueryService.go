package service

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/user_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	mencache "github.com/MamangRust/monolith-pointofsale-user/cache"
	"github.com/MamangRust/monolith-pointofsale-user/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type userQueryDeps struct {
	Cache         mencache.UserQueryCache
	UserQuery     repository.UserQueryRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type userQueryService struct {
	mencache      mencache.UserQueryCache
	userQuery     repository.UserQueryRepository
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

func NewUserQueryService(params *userQueryDeps) UserQueryService {
	return &userQueryService{
		mencache:      params.Cache,
		userQuery:     params.UserQuery,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *userQueryService) FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersRow, *int, error) {
	const method = "FindAll"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUsersCache(ctx, req); found {
		logSuccess("Successfully fetched users from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	users, totalRecords, err := s.userQuery.FindAllUsers(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*db.GetUsersRow](
			s.logger,
			user_errors.ErrFailedFindAll.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedUsersCache(ctx, req, users, totalRecords)
	logSuccess("Successfully fetched users", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return users, totalRecords, nil
}

func (s *userQueryService) FindByID(ctx context.Context, id int) (*db.User, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user.id", id),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedUserCache(ctx, id); found {
		logSuccess("Successfully fetched user from cache", zap.Int("user.id", id))
		return data, nil
	}

	user, err := s.userQuery.FindById(ctx, id)
	if err != nil {
		status = "error"
		// Propagate the typed repository error unchanged: ErrUserNotFound
		// (404) for no-rows, ErrInternalServerError (500) for other DB
		// failures — do not mask them into a generic 404.
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			err,
			method,
			span,
			zap.Int("user.id", id),
			zap.Error(err),
		)
	}

	s.mencache.SetCachedUserCache(ctx, user)
	logSuccess("Successfully fetched user", zap.Int("user.id", id))
	return user, nil
}

func (s *userQueryService) FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersActiveRow, *int, error) {
	const method = "FindByActive"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUserActiveCache(ctx, req); found {
		logSuccess("Successfully fetched active users from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	users, totalRecords, err := s.userQuery.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*db.GetUsersActiveRow](
			s.logger,
			user_errors.ErrFailedFindActive.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedUserActiveCache(ctx, req, users, totalRecords)
	logSuccess("Successfully fetched active users", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return users, totalRecords, nil
}

func (s *userQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUserTrashedRow, *int, error) {
	const method = "FindByTrashed"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUserTrashedCache(ctx, req); found {
		logSuccess("Successfully fetched trashed users from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	users, totalRecords, err := s.userQuery.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*db.GetUserTrashedRow](
			s.logger,
			user_errors.ErrFailedFindTrashed.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedUserTrashedCache(ctx, req, users, totalRecords)
	logSuccess("Successfully fetched trashed users", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return users, totalRecords, nil
}

func (s *userQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

// FindByEmailWithPassword returns the account with its password hash so the auth
// service can verify credentials. Results are never cached because they carry a
// credential, and a missing account is reported as ErrUserNotFound.
func (s *userQueryService) FindByEmailWithPassword(ctx context.Context, email string) (*db.User, error) {
	const method = "FindByEmailWithPassword"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("email", email),
	)
	defer func() {
		end(status)
	}()

	user, err := s.userQuery.FindByEmail(ctx, email)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			err,
			method,
			span,
			zap.String("email", email),
			zap.Error(err),
		)
	}
	if user == nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrUserNotFound,
			method,
			span,
			zap.String("email", email),
		)
	}

	logSuccess("Successfully fetched user by email", zap.Int32("user.id", user.UserID))
	return user, nil
}

// FindByEmailAndVerify is the login-time lookup restricted to verified
// accounts; unverified users must fail credential verification.
func (s *userQueryService) FindByEmailAndVerify(ctx context.Context, email string) (*db.User, error) {
	const method = "FindByEmailAndVerify"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("email", email),
	)
	defer func() {
		end(status)
	}()

	user, err := s.userQuery.FindByEmailAndVerify(ctx, email)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			err,
			method,
			span,
			zap.String("email", email),
			zap.Error(err),
		)
	}
	if user == nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrUserNotFound,
			method,
			span,
			zap.String("email", email),
		)
	}

	logSuccess("Successfully fetched verified user by email", zap.Int32("user.id", user.UserID))
	return user, nil
}

// FindByVerificationCode exchanges the emailed code for the account it belongs
// to, backing the auth service's VerifyCode flow.
func (s *userQueryService) FindByVerificationCode(ctx context.Context, verification_code string) (*db.User, error) {
	const method = "FindByVerificationCode"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("verification_code", verification_code),
	)
	defer func() {
		end(status)
	}()

	user, err := s.userQuery.FindByVerificationCode(ctx, verification_code)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}
	if user == nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrUserNotFound,
			method,
			span,
		)
	}

	logSuccess("Successfully fetched user by verification code", zap.Int32("user.id", user.UserID))
	return user, nil
}
