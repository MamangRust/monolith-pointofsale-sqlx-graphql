package response_service

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/record"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type roleResponseMapper struct {
}

func NewRoleResponseMapper() *roleResponseMapper {
	return &roleResponseMapper{}
}

func (s *roleResponseMapper) ToRoleResponse(role *record.RoleRecord) *response.RoleResponse {
	return &response.RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

func (s *roleResponseMapper) ToRolesResponse(roles []*record.RoleRecord) []*response.RoleResponse {
	var responseRoles []*response.RoleResponse

	for _, role := range roles {
		responseRoles = append(responseRoles, s.ToRoleResponse(role))
	}

	return responseRoles
}

func (s *roleResponseMapper) ToRoleResponseDeleteAt(role *record.RoleRecord) *response.RoleResponseDeleteAt {
	return &response.RoleResponseDeleteAt{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: *role.DeletedAt,
	}
}

func (s *roleResponseMapper) ToRolesResponseDeleteAt(roles []*record.RoleRecord) []*response.RoleResponseDeleteAt {
	var responseRoles []*response.RoleResponseDeleteAt

	for _, role := range roles {
		responseRoles = append(responseRoles, s.ToRoleResponseDeleteAt(role))
	}

	return responseRoles
}
