package handler

import (
	"context"

	pb "github.com/MamangRust/monolith-pointofsale-pb"

	"github.com/MamangRust/monolith-pointofsale-auth/service"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"

	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	"go.uber.org/zap"
)

type authHandleGrpc struct {
	pb.UnimplementedAuthServiceServer
	registerService      service.RegistrationService
	loginService         service.LoginService
	passwordResetService service.PasswordResetService
	identifyService      service.IdentifyService
	logger               logger.LoggerInterface
}

func NewAuthHandleGrpc(authService *service.Service, logger logger.LoggerInterface) pb.AuthServiceServer {
	return &authHandleGrpc{
		registerService:      authService.Register,
		loginService:         authService.Login,
		passwordResetService: authService.PasswordReset,
		identifyService:      authService.Identify,
		logger:               logger,
	}
}

func (s *authHandleGrpc) VerifyCode(ctx context.Context, req *pb.VerifyCodeRequest) (*pb.ApiResponseVerifyCode, error) {
	s.logger.Info("VerifyCode called", zap.Bool("verification_code_present", req.Code != ""))

	_, err := s.passwordResetService.VerifyCode(ctx, req.Code)
	if err != nil {
		s.logger.Error("VerifyCode failed", zap.Any("error", err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("VerifyCode success", zap.Bool("verification_code_present", req.Code != ""))

	return &pb.ApiResponseVerifyCode{
		Status:  "success",
		Message: "Verification successfully",
	}, nil
}

func (s *authHandleGrpc) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ApiResponseForgotPassword, error) {
	s.logger.Info("ForgotPassword called", zap.String("email", req.Email))

	_, err := s.passwordResetService.ForgotPassword(ctx, req.Email)
	if err != nil {
		s.logger.Error("ForgotPassword failed", zap.Any("error", err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("ForgotPassword successful", zap.Bool("success", true))

	return &pb.ApiResponseForgotPassword{
		Status:  "success",
		Message: "ForgotPassword successful",
	}, nil
}

func (s *authHandleGrpc) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ApiResponseResetPassword, error) {
	s.logger.Info("ResetPassword called", zap.Bool("reset_token_present", req.ResetToken != ""))

	_, err := s.passwordResetService.ResetPassword(ctx, &requests.CreateResetPasswordRequest{
		ResetToken:      req.ResetToken,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {
		s.logger.Error("ResetPassword failed", zap.Any("error", err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("ResetPassword successful", zap.Bool("success", true))

	return &pb.ApiResponseResetPassword{
		Status:  "success",
		Message: "Reset password successful",
	}, nil
}

func (s *authHandleGrpc) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponseLogin, error) {
	s.logger.Info("LoginUser called", zap.String("email", req.Email))

	request := &requests.AuthRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := s.loginService.Login(ctx, request)
	if err != nil {
		s.logger.Error("LoginUser failed", zap.Any("error", err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("LoginUser successful", zap.Bool("success", true))

	return &pb.ApiResponseLogin{
		Status:  "success",
		Message: "LoginUser successfull",
		Data: &pb.TokenResponse{
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
		},
	}, nil
}

func (s *authHandleGrpc) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.ApiResponseRefreshToken, error) {
	s.logger.Info("RefreshToken called")

	res, err := s.identifyService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		s.logger.Error("RefreshToken failed", zap.Any("error", err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RefreshToken successful", zap.Bool("success", true))

	return &pb.ApiResponseRefreshToken{
		Status:  "success",
		Message: "Refresh token successful",
		Data: &pb.TokenResponse{
			AccessToken:  res.AccessToken,
			RefreshToken: req.RefreshToken,
		},
	}, nil
}

func (s *authHandleGrpc) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.ApiResponseGetMe, error) {
	s.logger.Info("GetMe called")

	res, err := s.identifyService.GetMe(ctx, req.AccessToken)
	if err != nil {
		s.logger.Error("GetMe failed", zap.Any("error", err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("GetMe successful", zap.Bool("success", true))

	var createdAtStr, updatedAtStr string
	if res.CreatedAt.Valid {
		createdAtStr = res.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if res.UpdatedAt.Valid {
		updatedAtStr = res.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}

	return &pb.ApiResponseGetMe{
		Status:  "success",
		Message: "Get me successfully",
		Data: &pbuser.UserResponse{
			Id:        res.UserID,
			Firstname: res.Firstname,
			Lastname:  res.Lastname,
			Email:     res.Email,
			CreatedAt: createdAtStr,
			UpdatedAt: updatedAtStr,
		},
	}, nil
}

func (s *authHandleGrpc) RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponseRegister, error) {
	s.logger.Info("RegisterUser called", zap.String("email", req.Email))

	request := &requests.RegisterRequest{
		FirstName:       req.Firstname,
		LastName:        req.Lastname,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}

	res, err := s.registerService.Register(ctx, request)
	if err != nil {
		s.logger.Error("RegisterUser failed", zap.Any("error", err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RegisterUser successful", zap.Bool("success", true))

	var createdAtStr, updatedAtStr string
	if res.CreatedAt.Valid {
		createdAtStr = res.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if res.UpdatedAt.Valid {
		updatedAtStr = res.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}

	return &pb.ApiResponseRegister{
		Status:  "success",
		Message: "RegisterUser successful",
		Data: &pbuser.UserResponse{
			Id:        res.UserID,
			Firstname: res.Firstname,
			Lastname:  res.Lastname,
			Email:     res.Email,
			CreatedAt: createdAtStr,
			UpdatedAt: updatedAtStr,
		},
	}, nil
}
