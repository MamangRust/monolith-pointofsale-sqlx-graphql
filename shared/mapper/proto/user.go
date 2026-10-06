package protomapper

import (
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type userProtoMapper struct {
}

func NewUserProtoMapper() *userProtoMapper {
	return &userProtoMapper{}
}

func (u *userProtoMapper) ToProtoResponseUser(status string, message string, pbResponse *response.UserResponse) *pbuser.ApiResponseUser {
	return &pbuser.ApiResponseUser{
		Status:  status,
		Message: message,
		Data:    u.mapResponseUser(pbResponse),
	}
}

func (u *userProtoMapper) ToProtoResponseUserDeleteAt(status string, message string, pbResponse *response.UserResponseDeleteAt) *pbuser.ApiResponseUserDeleteAt {
	return &pbuser.ApiResponseUserDeleteAt{
		Status:  status,
		Message: message,
		Data:    u.mapResponseUserDeleteAt(pbResponse),
	}
}

func (u *userProtoMapper) ToProtoResponsesUser(status string, message string, pbResponse []*response.UserResponse) *pbuser.ApiResponsesUser {
	return &pbuser.ApiResponsesUser{
		Status:  status,
		Message: message,
		Data:    u.mapResponsesUser(pbResponse),
	}
}

func (u *userProtoMapper) ToProtoResponseUserDelete(status string, message string) *pbuser.ApiResponseUserDelete {
	return &pbuser.ApiResponseUserDelete{
		Status:  status,
		Message: message,
	}
}

func (u *userProtoMapper) ToProtoResponseUserAll(status string, message string) *pbuser.ApiResponseUserAll {
	return &pbuser.ApiResponseUserAll{
		Status:  status,
		Message: message,
	}
}

func (u *userProtoMapper) ToProtoResponsePaginationUserDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.UserResponseDeleteAt) *pbuser.ApiResponsePaginationUserDeleteAt {
	return &pbuser.ApiResponsePaginationUserDeleteAt{
		Status:     status,
		Message:    message,
		Data:       u.mapResponsesUserDeleteAt(users),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (u *userProtoMapper) ToProtoResponsePaginationUser(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.UserResponse) *pbuser.ApiResponsePaginationUser {
	return &pbuser.ApiResponsePaginationUser{
		Status:     status,
		Message:    message,
		Data:       u.mapResponsesUser(users),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (u *userProtoMapper) mapResponseUserDeleteAt(user *response.UserResponseDeleteAt) *pbuser.UserResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if user.DeletedAt != nil {
		deletedAt = wrapperspb.String(*user.DeletedAt)
	}

	return &pbuser.UserResponseDeleteAt{
		Id:        int32(user.ID),
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (u *userProtoMapper) mapResponseUser(user *response.UserResponse) *pbuser.UserResponse {
	return &pbuser.UserResponse{
		Id:        int32(user.ID),
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *userProtoMapper) mapResponsesUser(users []*response.UserResponse) []*pbuser.UserResponse {
	var mappedUsers []*pbuser.UserResponse

	for _, user := range users {
		mappedUsers = append(mappedUsers, u.mapResponseUser(user))
	}

	return mappedUsers
}

func (u *userProtoMapper) mapResponseUserDelete(user *response.UserResponseDeleteAt) *pbuser.UserResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if user.DeletedAt != nil {
		deletedAt = wrapperspb.String(*user.DeletedAt)
	}

	return &pbuser.UserResponseDeleteAt{
		Id:        int32(user.ID),
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (u *userProtoMapper) mapResponsesUserDeleteAt(users []*response.UserResponseDeleteAt) []*pbuser.UserResponseDeleteAt {
	var mappedUsers []*pbuser.UserResponseDeleteAt

	for _, user := range users {
		mappedUsers = append(mappedUsers, u.mapResponseUserDelete(user))
	}

	return mappedUsers
}

func mapPaginationMeta(s *pbcommon.PaginationMeta) *pbcommon.PaginationMeta {
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(s.CurrentPage),
		PageSize:     int32(s.PageSize),
		TotalPages:   int32(s.TotalPages),
		TotalRecords: int32(s.TotalRecords),
	}
}
