package protomapper

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb"
	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type authProtoMapper struct {
}

func NewAuthProtoMapper() *authProtoMapper {
	return &authProtoMapper{}
}

func (s *authProtoMapper) ToProtoResponseVerifyCode(status string, message string) *pb.ApiResponseVerifyCode {
	return &pb.ApiResponseVerifyCode{
		Status:  status,
		Message: message,
	}
}

func (s *authProtoMapper) ToProtoResponseForgotPassword(status string, message string) *pb.ApiResponseForgotPassword {
	return &pb.ApiResponseForgotPassword{
		Status:  status,
		Message: message,
	}
}

func (s *authProtoMapper) ToProtoResponseResetPassword(status string, message string) *pb.ApiResponseResetPassword {
	return &pb.ApiResponseResetPassword{
		Status:  status,
		Message: message,
	}
}

func (s *authProtoMapper) ToProtoResponseLogin(status string, message string, response *response.TokenResponse) *pb.ApiResponseLogin {
	return &pb.ApiResponseLogin{
		Status:  status,
		Message: message,
		Data: &pb.TokenResponse{
			AccessToken:  response.AccessToken,
			RefreshToken: response.RefreshToken,
		},
	}
}

func (s *authProtoMapper) ToProtoResponseRegister(status string, message string, response *response.UserResponse) *pb.ApiResponseRegister {
	return &pb.ApiResponseRegister{
		Status:  status,
		Message: message,
		Data: &pbuser.UserResponse{
			Id:        int32(response.ID),
			Firstname: response.FirstName,
			Lastname:  response.LastName,
			Email:     response.Email,
			CreatedAt: response.CreatedAt,
			UpdatedAt: response.UpdatedAt,
		},
	}
}

func (s *authProtoMapper) ToProtoResponseRefreshToken(status string, message string, response *response.TokenResponse) *pb.ApiResponseRefreshToken {
	return &pb.ApiResponseRefreshToken{
		Status:  status,
		Message: message,
		Data: &pb.TokenResponse{
			AccessToken:  response.AccessToken,
			RefreshToken: response.RefreshToken,
		},
	}
}

func (s *authProtoMapper) ToProtoResponseGetMe(status string, message string, response *response.UserResponse) *pb.ApiResponseGetMe {
	return &pb.ApiResponseGetMe{
		Status:  status,
		Message: message,
		Data: &pbuser.UserResponse{
			Id:        int32(response.ID),
			Firstname: response.FirstName,
			Lastname:  response.LastName,
			Email:     response.Email,
			CreatedAt: response.CreatedAt,
			UpdatedAt: response.UpdatedAt,
		},
	}
}
