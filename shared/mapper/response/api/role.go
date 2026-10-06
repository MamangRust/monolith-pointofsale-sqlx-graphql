package response_api

import (
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type roleResponseMapper struct {
}

func NewRoleResponseMapper() *roleResponseMapper {
	return &roleResponseMapper{}
}

func (s *roleResponseMapper) ToApiResponseRoleAll(pbResponse *pbrole.ApiResponseRoleAll) *response.ApiResponseRoleAll {
	return &response.ApiResponseRoleAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *roleResponseMapper) ToApiResponseRoleDelete(pbResponse *pbrole.ApiResponseRoleDelete) *response.ApiResponseRoleDelete {
	return &response.ApiResponseRoleDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *roleResponseMapper) ToApiResponseRole(pbResponse *pbrole.ApiResponseRole) *response.ApiResponseRole {
	return &response.ApiResponseRole{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapResponseRole(pbResponse.Data),
	}
}

func (s *roleResponseMapper) ToApiResponsesRole(pbResponse *pbrole.ApiResponsesRole) *response.ApiResponsesRole {
	return &response.ApiResponsesRole{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapResponsesRole(pbResponse.Data),
	}
}

func (s *roleResponseMapper) ToApiResponsePaginationRole(pbResponse *pbrole.ApiResponsePaginationRole) *response.ApiResponsePaginationRole {
	return &response.ApiResponsePaginationRole{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.mapResponsesRole(pbResponse.Data),
		Pagination: mapPaginationMeta(pbResponse.Pagination),
	}
}

func (s *roleResponseMapper) ToApiResponsePaginationRoleDeleteAt(pbResponse *pbrole.ApiResponsePaginationRoleDeleteAt) *response.ApiResponsePaginationRoleDeleteAt {
	return &response.ApiResponsePaginationRoleDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.mapResponsesRoleDeleteAt(pbResponse.Data),
		Pagination: mapPaginationMeta(pbResponse.Pagination),
	}
}

func (s *roleResponseMapper) mapResponseRole(role *pbrole.RoleResponse) *response.RoleResponse {
	return &response.RoleResponse{
		ID:        int(role.Id),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

func (s *roleResponseMapper) mapResponsesRole(roles []*pbrole.RoleResponse) []*response.RoleResponse {
	var responseRoles []*response.RoleResponse

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRole(role))
	}

	return responseRoles
}

func (s *roleResponseMapper) mapResponseRoleDeleteAt(role *pbrole.RoleResponseDeleteAt) *response.RoleResponseDeleteAt {
	return &response.RoleResponseDeleteAt{
		ID:        int(role.Id),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: role.DeletedAt,
	}
}

func (s *roleResponseMapper) mapResponsesRoleDeleteAt(roles []*pbrole.RoleResponseDeleteAt) []*response.RoleResponseDeleteAt {
	var responseRoles []*response.RoleResponseDeleteAt

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRoleDeleteAt(role))
	}

	return responseRoles
}

func mapPaginationMeta(s *pbcommon.PaginationMeta) *response.PaginationMeta {
	return &response.PaginationMeta{
		CurrentPage:  int(s.CurrentPage),
		PageSize:     int(s.PageSize),
		TotalRecords: int(s.TotalRecords),
		TotalPages:   int(s.TotalPages),
	}
}
