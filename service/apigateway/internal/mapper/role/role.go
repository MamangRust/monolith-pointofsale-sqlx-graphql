package rolegraphqlmapper

import (
	graphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper"
	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-pointofsale-pb/role"
)

type roleGraphqlMapper struct {
}

func NewRoleGraphqlMapper() *roleGraphqlMapper {
	return &roleGraphqlMapper{}
}

func (s *roleGraphqlMapper) ToGraphqlResponseAll(res *pb.ApiResponseRoleAll) *model.APIResponseRoleAll {
	return &model.APIResponseRoleAll{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (s *roleGraphqlMapper) ToGraphqlResponseDelete(res *pb.ApiResponseRoleDelete) *model.APIResponseRoleDelete {
	return &model.APIResponseRoleDelete{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (s *roleGraphqlMapper) ToGraphqlResponseRole(res *pb.ApiResponseRole) *model.APIResponseRole {
	return &model.APIResponseRole{
		Status:  res.Status,
		Message: res.Message,
		Data:    s.mapResponseRole(res.Data),
	}
}

func (s *roleGraphqlMapper) ToGraphqlResponseRoleDeleteAt(res *pb.ApiResponseRole) *model.APIResponseRoleDeleteAt {
	return &model.APIResponseRoleDeleteAt{
		Status:  res.Status,
		Message: res.Message,
		Data:    s.mapResponseRoleToDeletedAt(res.Data),
	}
}

func (s *roleGraphqlMapper) mapResponseRoleToDeletedAt(role *pb.RoleResponse) *model.RoleResponseDeleteAt {
	if role == nil {
		return nil
	}
	return &model.RoleResponseDeleteAt{
		ID:        int32(role.Id),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: nil,
	}
}

func (s *roleGraphqlMapper) ToGraphqlResponsesRole(res *pb.ApiResponsesRole) *model.APIResponsesRole {
	return &model.APIResponsesRole{
		Status:  res.Status,
		Message: res.Message,
		Data:    s.mapResponsesRole(res.Data),
	}
}

func (s *roleGraphqlMapper) ToGraphqlResponsePaginationRole(res *pb.ApiResponsePaginationRole) *model.APIResponsePaginationRole {
	return &model.APIResponsePaginationRole{
		Status:     res.Status,
		Message:    res.Message,
		Data:       s.mapResponsesRole(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (s *roleGraphqlMapper) ToGraphqlResponsePaginationRoleDeleteAt(res *pb.ApiResponsePaginationRoleDeleteAt) *model.APIResponsePaginationRoleDeleteAt {
	return &model.APIResponsePaginationRoleDeleteAt{
		Status:     res.Status,
		Message:    res.Message,
		Data:       s.mapResponsesRoleDeleteAt(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (s *roleGraphqlMapper) mapResponseRole(role *pb.RoleResponse) *model.RoleResponse {
	return &model.RoleResponse{
		ID:        int32(role.Id),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

func (s *roleGraphqlMapper) mapResponsesRole(roles []*pb.RoleResponse) []*model.RoleResponse {
	var responseRoles []*model.RoleResponse

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRole(role))
	}

	return responseRoles
}

func (s *roleGraphqlMapper) mapResponseRoleDeleteAt(role *pb.RoleResponseDeleteAt) *model.RoleResponseDeleteAt {
	var deletedAt *string
	if role.DeletedAt != "" {
		deletedAt = &role.DeletedAt
	}

	return &model.RoleResponseDeleteAt{
		ID:        int32(role.Id),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (s *roleGraphqlMapper) mapResponsesRoleDeleteAt(roles []*pb.RoleResponseDeleteAt) []*model.RoleResponseDeleteAt {
	var responseRoles []*model.RoleResponseDeleteAt

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRoleDeleteAt(role))
	}

	return responseRoles
}
