package usergraphqlmapper

import (
	graphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper"
	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-pointofsale-pb/user"
)

type userGraphqlMapper struct {
}

func NewUserGraphqlMapper() *userGraphqlMapper {
	return &userGraphqlMapper{}
}

func (u *userGraphqlMapper) ToGraphqlResponseUserDelete(res *pb.ApiResponseUserDelete) *model.APIResponseUserDelete {
	return &model.APIResponseUserDelete{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (u *userGraphqlMapper) ToGraphqlResponseUserAll(res *pb.ApiResponseUserAll) *model.APIResponseUserAll {
	return &model.APIResponseUserAll{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (u *userGraphqlMapper) ToGraphqlResponseUser(res *pb.ApiResponseUser) *model.APIResponseUserResponse {
	return &model.APIResponseUserResponse{
		Status:  res.Status,
		Message: res.Message,
		Data:    u.mapUserResponse(res.Data),
	}
}

func (u *userGraphqlMapper) ToGraphqlResponseUserDeleteAt(res *pb.ApiResponseUserDeleteAt) *model.APIResponseUserResponseDeleteAt {
	return &model.APIResponseUserResponseDeleteAt{
		Status:  res.Status,
		Message: res.Message,
		Data:    u.mapUserResponseDeleteAt(res.Data),
	}
}

func (u *userGraphqlMapper) ToGraphqlResponseUsers(res *pb.ApiResponsesUser) *model.APIResponsesUser {
	return &model.APIResponsesUser{
		Status:  res.Status,
		Message: res.Message,
		Data:    u.mapUserResponses(res.Data),
	}
}

func (u *userGraphqlMapper) ToGraphqlResponsePaginationUser(res *pb.ApiResponsePaginationUser) *model.APIResponsePaginationUser {
	return &model.APIResponsePaginationUser{
		Status:     res.Status,
		Message:    res.Message,
		Data:       u.mapUserResponses(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (u *userGraphqlMapper) ToGraphqlResponsePaginationUserDeleteAt(res *pb.ApiResponsePaginationUserDeleteAt) *model.APIResponsePaginationUserDeleteAt {
	return &model.APIResponsePaginationUserDeleteAt{
		Status:     res.Status,
		Message:    res.Message,
		Data:       u.mapUserResponsesDeleteAt(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (u *userGraphqlMapper) mapUserResponse(user *pb.UserResponse) *model.UserResponse {
	if user == nil {
		return nil
	}
	return &model.UserResponse{
		ID:        int32(user.Id),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *userGraphqlMapper) mapUserResponses(users []*pb.UserResponse) []*model.UserResponse {
	var responses []*model.UserResponse
	for _, user := range users {
		responses = append(responses, u.mapUserResponse(user))
	}
	return responses
}

func (u *userGraphqlMapper) mapUserResponseDeleteAt(user *pb.UserResponseDeleteAt) *model.UserResponseDeleteAt {
	if user == nil {
		return nil
	}
	var deletedAt *string
	if user.DeletedAt != nil {
		deletedAt = &user.DeletedAt.Value
	}

	return &model.UserResponseDeleteAt{
		ID:        int32(user.Id),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (u *userGraphqlMapper) mapUserResponsesDeleteAt(users []*pb.UserResponseDeleteAt) []*model.UserResponseDeleteAt {
	var responses []*model.UserResponseDeleteAt
	for _, user := range users {
		responses = append(responses, u.mapUserResponseDeleteAt(user))
	}
	return responses
}
