package protomapper

import (
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type roleProtoMapper struct {
}

func NewRoleProtoMapper() *roleProtoMapper {
	return &roleProtoMapper{}
}

func (s *roleProtoMapper) ToProtoResponseRoleAll(status string, message string) *pbrole.ApiResponseRoleAll {
	return &pbrole.ApiResponseRoleAll{
		Status:  status,
		Message: message,
	}
}

func (s *roleProtoMapper) ToProtoResponseRoleDelete(status string, message string) *pbrole.ApiResponseRoleDelete {
	return &pbrole.ApiResponseRoleDelete{
		Status:  status,
		Message: message,
	}
}

func (s *roleProtoMapper) ToProtoResponseRole(status string, message string, pbResponse *response.RoleResponse) *pbrole.ApiResponseRole {
	return &pbrole.ApiResponseRole{
		Status:  status,
		Message: message,
		Data:    s.mapResponseRole(pbResponse),
	}
}

func (s *roleProtoMapper) ToProtoResponsesRole(status string, message string, pbResponse []*response.RoleResponse) *pbrole.ApiResponsesRole {
	return &pbrole.ApiResponsesRole{
		Status:  status,
		Message: message,
		Data:    s.mapResponsesRole(pbResponse),
	}
}

func (s *roleProtoMapper) ToProtoResponsePaginationRole(pagination *pbcommon.PaginationMeta, status string, message string, pbResponse []*response.RoleResponse) *pbrole.ApiResponsePaginationRole {
	return &pbrole.ApiResponsePaginationRole{
		Status:     status,
		Message:    message,
		Data:       s.mapResponsesRole(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (s *roleProtoMapper) ToProtoResponsePaginationRoleDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, pbResponse []*response.RoleResponseDeleteAt) *pbrole.ApiResponsePaginationRoleDeleteAt {
	return &pbrole.ApiResponsePaginationRoleDeleteAt{
		Status:     status,
		Message:    message,
		Data:       s.mapResponsesRoleDeleteAt(pbResponse),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (s *roleProtoMapper) mapResponseRole(role *response.RoleResponse) *pbrole.RoleResponse {
	return &pbrole.RoleResponse{
		Id:        int32(role.ID),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

func (s *roleProtoMapper) mapResponsesRole(roles []*response.RoleResponse) []*pbrole.RoleResponse {
	var responseRoles []*pbrole.RoleResponse

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRole(role))
	}

	return responseRoles
}

func (s *roleProtoMapper) mapResponseRoleDeleteAt(role *response.RoleResponseDeleteAt) *pbrole.RoleResponseDeleteAt {
	return &pbrole.RoleResponseDeleteAt{
		Id:        int32(role.ID),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: role.DeletedAt,
	}
}

func (s *roleProtoMapper) mapResponsesRoleDeleteAt(roles []*response.RoleResponseDeleteAt) []*pbrole.RoleResponseDeleteAt {
	var responseRoles []*pbrole.RoleResponseDeleteAt

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRoleDeleteAt(role))
	}

	return responseRoles
}
