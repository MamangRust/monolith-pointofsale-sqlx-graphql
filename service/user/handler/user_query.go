package handler

import (
	"context"
	"math"

	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/user_errors"
	"github.com/MamangRust/monolith-pointofsale-user/service"
	"go.uber.org/zap"
)

type userQueryHandleGrpc struct {
	pbuser.UnimplementedUserQueryServiceServer
	userQueryService service.UserQueryService
	logger           logger.LoggerInterface
}

func NewUserQueryHandleGrpc(
	user *service.Service,
	logger logger.LoggerInterface,
) pbuser.UserQueryServiceServer {
	return &userQueryHandleGrpc{
		userQueryService: user.UserQuery,
		logger:           logger,
	}
}

func (s *userQueryHandleGrpc) FindAll(ctx context.Context, request *pbuser.FindAllUserRequest) (*pbuser.ApiResponsePaginationUser, error) {
	s.logger.Info("FindAll users called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll users failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	var totalRecordsVal int
	if totalRecords != nil {
		totalRecordsVal = *totalRecords
	}
	totalPages := int(math.Ceil(float64(totalRecordsVal) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecordsVal),
	}

	s.logger.Info("FindAll users success")

	return &pbuser.ApiResponsePaginationUser{
		Status:     "success",
		Message:    "Successfully fetched users",
		Data:       mapGetUsersRowsToProto(users),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *userQueryHandleGrpc) FindById(ctx context.Context, request *pbuser.FindByIdUserRequest) (*pbuser.ApiResponseUser, error) {
	s.logger.Info("FindById user called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, user_errors.ErrGrpcUserNotFound
	}

	user, err := s.userQueryService.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("FindById user failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById user success")

	return &pbuser.ApiResponseUser{
		Status:  "success",
		Message: "Successfully fetched user",
		Data:    mapUserToProto(user),
	}, nil
}

func (s *userQueryHandleGrpc) FindByActive(ctx context.Context, request *pbuser.FindAllUserRequest) (*pbuser.ApiResponsePaginationUserDeleteAt, error) {
	s.logger.Info("FindByActive users called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive users failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	var totalRecordsVal int
	if totalRecords != nil {
		totalRecordsVal = *totalRecords
	}
	totalPages := int(math.Ceil(float64(totalRecordsVal) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecordsVal),
	}

	s.logger.Info("FindByActive users success")

	return &pbuser.ApiResponsePaginationUserDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active users",
		Data:       mapActiveUsersToProto(users),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

// FindByEmail returns the account with its password hash so the auth service can
// verify credentials. This RPC must never be exposed through the API gateway.
func (s *userQueryHandleGrpc) FindByEmail(ctx context.Context, request *pbuser.FindByEmailRequest) (*pbuser.ApiResponseUserWithPassword, error) {
	s.logger.Info("FindByEmail user called", zap.String("email", request.GetEmail()))

	user, err := s.userQueryService.FindByEmailWithPassword(ctx, request.GetEmail())
	if err != nil {
		s.logger.Error("FindByEmail user failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindByEmail user success")

	return &pbuser.ApiResponseUserWithPassword{
		Status:  "success",
		Message: "Successfully fetched user",
		Data:    mapUserWithPasswordToProto(user),
	}, nil
}

// FindByEmailAndVerify is the login-time lookup restricted to verified accounts.
func (s *userQueryHandleGrpc) FindByEmailAndVerify(ctx context.Context, request *pbuser.FindByEmailRequest) (*pbuser.ApiResponseUserWithPassword, error) {
	s.logger.Info("FindByEmailAndVerify user called", zap.String("email", request.GetEmail()))

	user, err := s.userQueryService.FindByEmailAndVerify(ctx, request.GetEmail())
	if err != nil {
		s.logger.Error("FindByEmailAndVerify user failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindByEmailAndVerify user success")

	return &pbuser.ApiResponseUserWithPassword{
		Status:  "success",
		Message: "Successfully fetched user",
		Data:    mapUserWithPasswordToProto(user),
	}, nil
}

func (s *userQueryHandleGrpc) FindByVerificationCode(ctx context.Context, request *pbuser.FindByVerificationCodeRequest) (*pbuser.ApiResponseUser, error) {
	s.logger.Info("FindByVerificationCode user called")

	user, err := s.userQueryService.FindByVerificationCode(ctx, request.GetVerificationCode())
	if err != nil {
		s.logger.Error("FindByVerificationCode user failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindByVerificationCode user success")

	return &pbuser.ApiResponseUser{
		Status:  "success",
		Message: "Successfully fetched user",
		Data:    mapUserToProto(user),
	}, nil
}
