package response_api

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type authResponseMapper struct {
}

func NewAuthResponseMapper() *authResponseMapper {
	return &authResponseMapper{}
}

func (s *authResponseMapper) ToResponseVerifyCode(res *pb.ApiResponseVerifyCode) *response.ApiResponseVerifyCode {
	return &response.ApiResponseVerifyCode{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (s *authResponseMapper) ToResponseForgotPassword(res *pb.ApiResponseForgotPassword) *response.ApiResponseForgotPassword {
	return &response.ApiResponseForgotPassword{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (s *authResponseMapper) ToResponseResetPassword(res *pb.ApiResponseResetPassword) *response.ApiResponseResetPassword {
	return &response.ApiResponseResetPassword{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (s *authResponseMapper) ToResponseLogin(res *pb.ApiResponseLogin) *response.ApiResponseLogin {
	defaultErrorResponse := &response.ApiResponseLogin{
		Status:  "error",
		Message: "authentication failed",
		Data:    nil,
	}

	if res == nil {
		return defaultErrorResponse
	}

	if res.Status == "" {
		res.Status = "error"
	}

	if res.Message == "" {
		res.Message = "authentication processed with empty response"
	}

	if res.Data == nil || res.Data.AccessToken == "" {
		return &response.ApiResponseLogin{
			Status:  res.Status,
			Message: res.Message,
			Data:    nil,
		}
	}

	return &response.ApiResponseLogin{
		Status:  res.Status,
		Message: res.Message,
		Data: &response.TokenResponse{
			AccessToken:  res.Data.AccessToken,
			RefreshToken: res.Data.RefreshToken,
		},
	}
}

func (s *authResponseMapper) ToResponseRegister(res *pb.ApiResponseRegister) *response.ApiResponseRegister {
	return &response.ApiResponseRegister{
		Status:  res.Status,
		Message: res.Message,
		Data: &response.UserResponse{
			ID:        int(res.Data.Id),
			FirstName: res.Data.Firstname,
			LastName:  res.Data.Lastname,
			Email:     res.Data.Email,
			CreatedAt: res.Data.CreatedAt,
			UpdatedAt: res.Data.UpdatedAt,
		},
	}
}

func (s *authResponseMapper) ToResponseRefreshToken(res *pb.ApiResponseRefreshToken) *response.ApiResponseRefreshToken {
	return &response.ApiResponseRefreshToken{
		Status:  res.Status,
		Message: res.Message,
		Data: &response.TokenResponse{
			AccessToken:  res.Data.AccessToken,
			RefreshToken: res.Data.RefreshToken,
		},
	}
}

func (s *authResponseMapper) ToResponseGetMe(res *pb.ApiResponseGetMe) *response.ApiResponseGetMe {
	return &response.ApiResponseGetMe{
		Status:  res.Status,
		Message: res.Message,
		Data: &response.UserResponse{
			ID:        int(res.Data.Id),
			FirstName: res.Data.Firstname,
			LastName:  res.Data.Lastname,
			Email:     res.Data.Email,
			CreatedAt: res.Data.CreatedAt,
			UpdatedAt: res.Data.UpdatedAt,
		},
	}
}
