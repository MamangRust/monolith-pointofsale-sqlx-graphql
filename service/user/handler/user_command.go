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
	"google.golang.org/protobuf/types/known/emptypb"
)

type userCommandHandleGrpc struct {
	pbuser.UnimplementedUserCommandServiceServer
	userQueryService   service.UserQueryService
	userCommandService service.UserCommandService
	logger             logger.LoggerInterface
}

func NewUserCommandHandleGrpc(
	user *service.Service,
	logger logger.LoggerInterface,
) pbuser.UserCommandServiceServer {
	return &userCommandHandleGrpc{
		userQueryService:   user.UserQuery,
		userCommandService: user.UserCommand,
		logger:             logger,
	}
}

func (s *userCommandHandleGrpc) FindByTrashed(ctx context.Context, request *pbuser.FindAllUserRequest) (*pbuser.ApiResponsePaginationUserDeleteAt, error) {
	s.logger.Info("FindByTrashed users called")

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

	users, totalRecords, err := s.userQueryService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed users failed", zap.Error(err))
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

	s.logger.Info("FindByTrashed users success")

	return &pbuser.ApiResponsePaginationUserDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed users",
		Data:       mapTrashedUsersToProto(users),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *userCommandHandleGrpc) Create(ctx context.Context, request *pbuser.CreateUserRequest) (*pbuser.ApiResponseUser, error) {
	s.logger.Info("Create user called", zap.String("email", request.GetEmail()))

	req := &requests.CreateUserRequest{
		FirstName:       request.GetFirstname(),
		LastName:        request.GetLastname(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		ConfirmPassword: request.GetConfirmPassword(),
	}

	if err := req.Validate(); err != nil {
		return nil, user_errors.ErrGrpcValidateCreateUser
	}

	user, err := s.userCommandService.CreateUser(ctx, req)
	if err != nil {
		s.logger.Error("Create user failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create user success")

	return &pbuser.ApiResponseUser{
		Status:  "success",
		Message: "Successfully created user",
		Data:    mapUserToProto(user),
	}, nil
}

func (s *userCommandHandleGrpc) Update(ctx context.Context, request *pbuser.UpdateUserRequest) (*pbuser.ApiResponseUser, error) {
	s.logger.Info("Update user called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	req := &requests.UpdateUserRequest{
		UserID:          &id,
		FirstName:       request.GetFirstname(),
		LastName:        request.GetLastname(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		ConfirmPassword: request.GetConfirmPassword(),
	}

	if err := req.Validate(); err != nil {
		return nil, user_errors.ErrGrpcValidateCreateUser
	}

	user, err := s.userCommandService.UpdateUser(ctx, req)
	if err != nil {
		s.logger.Error("Update user failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update user success")

	return &pbuser.ApiResponseUser{
		Status:  "success",
		Message: "Successfully updated user",
		Data:    mapUserToProto(user),
	}, nil
}

func (s *userCommandHandleGrpc) TrashedUser(ctx context.Context, request *pbuser.FindByIdUserRequest) (*pbuser.ApiResponseUserDeleteAt, error) {
	s.logger.Info("TrashedUser called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.userCommandService.TrashedUser(ctx, id)
	if err != nil {
		s.logger.Error("TrashedUser failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedUser success")

	return &pbuser.ApiResponseUserDeleteAt{
		Status:  "success",
		Message: "Successfully trashed user",
		Data:    mapUserDeleteAtToProto(user),
	}, nil
}

func (s *userCommandHandleGrpc) RestoreUser(ctx context.Context, request *pbuser.FindByIdUserRequest) (*pbuser.ApiResponseUserDeleteAt, error) {
	s.logger.Info("RestoreUser called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.userCommandService.RestoreUser(ctx, id)
	if err != nil {
		s.logger.Error("RestoreUser failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreUser success")

	return &pbuser.ApiResponseUserDeleteAt{
		Status:  "success",
		Message: "Successfully restored user",
		Data:    mapUserDeleteAtToProto(user),
	}, nil
}

func (s *userCommandHandleGrpc) DeleteUserPermanent(ctx context.Context, request *pbuser.FindByIdUserRequest) (*pbuser.ApiResponseUserDelete, error) {
	s.logger.Info("DeleteUserPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	_, err := s.userCommandService.DeleteUserPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteUserPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteUserPermanent success")

	return &pbuser.ApiResponseUserDelete{
		Status:  "success",
		Message: "Successfully deleted user permanently",
	}, nil
}

func (s *userCommandHandleGrpc) RestoreAllUser(ctx context.Context, _ *emptypb.Empty) (*pbuser.ApiResponseUserAll, error) {
	s.logger.Info("RestoreAllUser called")

	_, err := s.userCommandService.RestoreAllUser(ctx)
	if err != nil {
		s.logger.Error("RestoreAllUser failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllUser success")

	return &pbuser.ApiResponseUserAll{
		Status:  "success",
		Message: "Successfully restore all user",
	}, nil
}

func (s *userCommandHandleGrpc) DeleteAllUserPermanent(ctx context.Context, _ *emptypb.Empty) (*pbuser.ApiResponseUserAll, error) {
	s.logger.Info("DeleteAllUserPermanent called")

	_, err := s.userCommandService.DeleteAllUserPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllUserPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllUserPermanent success")

	return &pbuser.ApiResponseUserAll{
		Status:  "success",
		Message: "Successfully delete user permanen",
	}, nil
}

// CreateUserRecord persists a registration coming from the auth service. The
// password is already hashed and the role is assigned by the auth service, so
// this handler deliberately does neither.
func (s *userCommandHandleGrpc) CreateUserRecord(ctx context.Context, request *pbuser.CreateUserRecordRequest) (*pbuser.ApiResponseUser, error) {
	s.logger.Info("CreateUserRecord called", zap.String("email", request.GetEmail()))

	req := &requests.RegisterRequest{
		FirstName:    request.GetFirstname(),
		LastName:     request.GetLastname(),
		Email:        request.GetEmail(),
		Password:     request.GetPassword(),
		VerifiedCode: request.GetVerificationCode(),
		IsVerified:   request.GetIsVerified(),
	}

	user, err := s.userCommandService.CreateUserRecord(ctx, req)
	if err != nil {
		s.logger.Error("CreateUserRecord failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("CreateUserRecord success")

	return &pbuser.ApiResponseUser{
		Status:  "success",
		Message: "Successfully created user",
		Data:    mapUserToProto(user),
	}, nil
}

func (s *userCommandHandleGrpc) UpdateUserIsVerified(ctx context.Context, request *pbuser.UpdateUserIsVerifiedUserRequest) (*pbuser.ApiResponseUser, error) {
	s.logger.Info("UpdateUserIsVerified called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.userCommandService.UpdateUserIsVerified(ctx, id, request.GetIsVerified())
	if err != nil {
		s.logger.Error("UpdateUserIsVerified failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateUserIsVerified success")

	return &pbuser.ApiResponseUser{
		Status:  "success",
		Message: "Successfully updated user verification",
		Data:    mapUserToProto(user),
	}, nil
}

func (s *userCommandHandleGrpc) UpdateUserPassword(ctx context.Context, request *pbuser.UpdateUserPasswordUserRequest) (*pbuser.ApiResponseUser, error) {
	s.logger.Info("UpdateUserPassword called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.userCommandService.UpdateUserPassword(ctx, id, request.GetPassword())
	if err != nil {
		s.logger.Error("UpdateUserPassword failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateUserPassword success")

	return &pbuser.ApiResponseUser{
		Status:  "success",
		Message: "Successfully updated user password",
		Data:    mapUserToProto(user),
	}, nil
}
