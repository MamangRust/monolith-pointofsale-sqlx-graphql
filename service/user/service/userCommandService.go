package service

import (
	"context"

	roleadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/role"
	userroleadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/user_role"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/hash"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/role_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/user_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	mencache "github.com/MamangRust/monolith-pointofsale-user/cache"
	"github.com/MamangRust/monolith-pointofsale-user/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type userCommandDeps struct {
	Cache           mencache.UserCommandCache
	UserQuery       repository.UserQueryRepository
	UserCommand     repository.UserCommandRepository
	RoleAdapter     roleadapter.QueryRepository
	UserRoleAdapter userroleadapter.CommandRepository
	Logger          logger.LoggerInterface
	Hashing         hash.HashPassword
	Observability   observability.TraceLoggerObservability
}

type userCommandService struct {
	mencache              mencache.UserCommandCache
	userQueryRepository   repository.UserQueryRepository
	userCommandRepository repository.UserCommandRepository
	roleAdapter           roleadapter.QueryRepository
	userRoleAdapter       userroleadapter.CommandRepository
	logger                logger.LoggerInterface
	hashing               hash.HashPassword
	observability         observability.TraceLoggerObservability
}

func NewUserCommandService(params *userCommandDeps) UserCommandService {
	return &userCommandService{
		mencache:              params.Cache,
		userQueryRepository:   params.UserQuery,
		userCommandRepository: params.UserCommand,
		roleAdapter:           params.RoleAdapter,
		userRoleAdapter:       params.UserRoleAdapter,
		logger:                params.Logger,
		hashing:               params.Hashing,
		observability:         params.Observability,
	}
}

func (s *userCommandService) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*db.User, error) {
	const method = "CreateUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	existingUser, err := s.userQueryRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrInternalServerError.WithInternal(err),
			method,
			span,
			zap.String("email", request.Email),
			zap.Error(err),
		)
	}
	if existingUser != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrUserEmailAlready,
			method,
			span,
			zap.String("email", request.Email),
			zap.Int32("user.id", existingUser.UserID),
		)
	}

	hashedPassword, err := s.hashing.HashPassword(request.Password)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrUserPassword.WithInternal(err),
			method,
			span,
			zap.String("email", request.Email),
			zap.Error(err),
		)
	}
	request.Password = hashedPassword

	// The default role is resolved through the role service before the user is
	// created so a missing role fails fast without leaving an orphan user.
	const defaultRoleName = "ROLE_ADMIN"
	role, err := s.roleAdapter.FindByName(ctx, defaultRoleName)
	if err != nil || role == nil {
		status = "error"
		if err == nil {
			err = role_errors.ErrRoleNotFoundRes
		}
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			role_errors.ErrRoleNotFoundRes.WithInternal(err),
			method,
			span,
			zap.String("name", defaultRoleName),
			zap.Error(err),
		)
	}

	res, err := s.userCommandRepository.CreateUser(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrFailedCreateUser.WithInternal(err),
			method,
			span,
			zap.String("email", request.Email),
			zap.Error(err),
		)
	}

	// Assign the default role through the role service, which owns user_roles.
	// A failure here fails the whole request hard: a user without a role cannot
	// be authorized, so silently succeeding would be worse than an error.
	if _, err := s.userRoleAdapter.AssignRoleToUser(ctx, &requests.CreateUserRoleRequest{
		UserId: int(res.UserID),
		RoleId: int(role.RoleID),
	}); err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			err,
			method,
			span,
			zap.Int32("user.id", res.UserID),
			zap.String("name", defaultRoleName),
			zap.Error(err),
		)
	}

	logSuccess("Successfully created user", zap.Int32("user.id", res.UserID), zap.Bool("success", true))
	return res, nil
}

// CreateUserRecord persists a registration coming from the auth service. Unlike
// CreateUser it does not resolve or assign a role: the auth service owns that
// step and calls the role service separately. The password arrives already
// hashed.
func (s *userCommandService) CreateUserRecord(ctx context.Context, request *requests.RegisterRequest) (*db.User, error) {
	const method = "CreateUserRecord"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("email", request.Email),
	)
	defer func() {
		end(status)
	}()

	existingUser, err := s.userQueryRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrInternalServerError.WithInternal(err),
			method,
			span,
			zap.String("email", request.Email),
			zap.Error(err),
		)
	}
	if existingUser != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrUserEmailAlready,
			method,
			span,
			zap.String("email", request.Email),
			zap.Int32("user.id", existingUser.UserID),
		)
	}

	res, err := s.userCommandRepository.CreateUserRecord(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrFailedCreateUser.WithInternal(err),
			method,
			span,
			zap.String("email", request.Email),
			zap.Error(err),
		)
	}

	logSuccess("Successfully created user record", zap.Int32("user.id", res.UserID), zap.Bool("success", true))
	return res, nil
}

// UpdateUserIsVerified flips the verification flag, backing auth's VerifyCode.
func (s *userCommandService) UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*db.User, error) {
	const method = "UpdateUserIsVerified"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user.id", user_id),
		attribute.Bool("is_verified", is_verified),
	)
	defer func() {
		end(status)
	}()

	res, err := s.userCommandRepository.UpdateUserIsVerified(ctx, user_id, is_verified)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrUpdateUserVerificationCode.WithInternal(err),
			method,
			span,
			zap.Int("user.id", user_id),
			zap.Error(err),
		)
	}

	s.mencache.DeleteUserCache(ctx, int(res.UserID))
	logSuccess("Successfully updated user verification", zap.Int32("user.id", res.UserID), zap.Bool("success", true))
	return res, nil
}

// UpdateUserPassword stores a password that the auth service has already hashed.
func (s *userCommandService) UpdateUserPassword(ctx context.Context, user_id int, password string) (*db.User, error) {
	const method = "UpdateUserPassword"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user.id", user_id),
	)
	defer func() {
		end(status)
	}()

	res, err := s.userCommandRepository.UpdateUserPassword(ctx, user_id, password)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrUpdateUserPassword.WithInternal(err),
			method,
			span,
			zap.Int("user.id", user_id),
			zap.Error(err),
		)
	}

	s.mencache.DeleteUserCache(ctx, int(res.UserID))
	logSuccess("Successfully updated user password", zap.Int32("user.id", res.UserID), zap.Bool("success", true))
	return res, nil
}

func (s *userCommandService) UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*db.User, error) {
	const method = "UpdateUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	existingUser, err := s.userQueryRepository.FindById(ctx, *request.UserID)
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
			zap.Int("user.id", *request.UserID),
			zap.Error(err),
		)
	}

	if request.Email != "" && request.Email != existingUser.Email {
		duplicateUser, _ := s.userQueryRepository.FindByEmail(ctx, request.Email)
		if duplicateUser != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.User](
				s.logger,
				user_errors.ErrUserEmailAlready,
				method,
				span,
				zap.String("email", request.Email),
				zap.Int32("user.id", duplicateUser.UserID),
			)
		}
	}

	if request.Password != "" {
		hashedPassword, err := s.hashing.HashPassword(request.Password)
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.User](
				s.logger,
				user_errors.ErrUserPassword.WithInternal(err),
				method,
				span,
				zap.Int("user.id", *request.UserID),
				zap.Error(err),
			)
		}
		request.Password = hashedPassword
	}

	const defaultRoleName = "ROLE_ADMIN"
	role, err := s.roleAdapter.FindByName(ctx, defaultRoleName)
	if err != nil || role == nil {
		status = "error"
		if err == nil {
			err = role_errors.ErrRoleNotFoundRes
		}
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			role_errors.ErrRoleNotFoundRes.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	res, err := s.userCommandRepository.UpdateUser(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrFailedUpdateUser.WithInternal(err),
			method,
			span,
			zap.Int("user.id", *request.UserID),
			zap.Error(err),
		)
	}

	s.mencache.DeleteUserCache(ctx, int(res.UserID))
	logSuccess("Successfully updated user", zap.Int32("user.id", res.UserID), zap.Bool("success", true))
	return res, nil
}

func (s *userCommandService) TrashedUser(ctx context.Context, user_id int) (*db.User, error) {
	const method = "TrashedUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	res, err := s.userCommandRepository.TrashedUser(ctx, user_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrFailedTrashedUser.WithInternal(err),
			method,
			span,
			zap.Int("user.id", user_id),
			zap.Error(err),
		)
	}

	s.mencache.DeleteUserCache(ctx, user_id)
	logSuccess("Successfully trashed user", zap.Int32("user.id", res.UserID), zap.Bool("success", true))
	return res, nil
}

func (s *userCommandService) RestoreUser(ctx context.Context, user_id int) (*db.User, error) {
	const method = "RestoreUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	res, err := s.userCommandRepository.RestoreUser(ctx, user_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.User](
			s.logger,
			user_errors.ErrFailedRestoreUser.WithInternal(err),
			method,
			span,
			zap.Int("user.id", user_id),
			zap.Error(err),
		)
	}

	s.mencache.DeleteUserCache(ctx, user_id)
	logSuccess("Successfully restored user", zap.Int32("user.id", res.UserID), zap.Bool("success", true))
	return res, nil
}

func (s *userCommandService) DeleteUserPermanent(ctx context.Context, user_id int) (bool, error) {
	const method = "DeleteUserPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.userCommandRepository.DeleteUserPermanent(ctx, user_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			user_errors.ErrFailedDeletePermanent.WithInternal(err),
			method,
			span,
			zap.Int("user.id", user_id),
			zap.Error(err),
		)
	}

	s.mencache.DeleteUserCache(ctx, user_id)
	logSuccess("Successfully deleted user permanently", zap.Int("user.id", user_id), zap.Bool("success", success))
	return success, nil
}

func (s *userCommandService) RestoreAllUser(ctx context.Context) (bool, error) {
	const method = "RestoreAllUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.userCommandRepository.RestoreAllUser(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			user_errors.ErrFailedRestoreAll.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteUserAllCache(ctx)
	logSuccess("Successfully restored all users", zap.Bool("success", success))
	return success, nil
}

func (s *userCommandService) DeleteAllUserPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllUserPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.userCommandRepository.DeleteAllUserPermanent(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			user_errors.ErrFailedDeleteAll.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteUserAllCache(ctx)
	logSuccess("Successfully deleted all users permanently", zap.Bool("success", success))
	return success, nil
}
