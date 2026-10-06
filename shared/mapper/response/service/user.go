package response_service

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/record"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type userResponseMapper struct {
}

func NewUserResponseMapper() *userResponseMapper {
	return &userResponseMapper{}
}

func (s *userResponseMapper) ToUserResponse(user *record.UserRecord) *response.UserResponse {
	return &response.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (s *userResponseMapper) ToUsersResponse(users []*record.UserRecord) []*response.UserResponse {
	var responses []*response.UserResponse

	for _, user := range users {
		responses = append(responses, s.ToUserResponse(user))
	}

	return responses
}

func (s *userResponseMapper) ToUserResponseDeleteAt(user *record.UserRecord) *response.UserResponseDeleteAt {
	return &response.UserResponseDeleteAt{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

func (s *userResponseMapper) ToUsersResponseDeleteAt(users []*record.UserRecord) []*response.UserResponseDeleteAt {
	var responses []*response.UserResponseDeleteAt

	for _, user := range users {
		responses = append(responses, s.ToUserResponseDeleteAt(user))
	}

	return responses
}
