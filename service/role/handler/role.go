package handler

import (
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
	pbuserrole "github.com/MamangRust/monolith-pointofsale-pb/user_role"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
)

func mapPaginationMeta(meta *pbcommon.PaginationMeta) *pbcommon.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &pbcommon.PaginationMeta{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		TotalPages:   meta.TotalPages,
		TotalRecords: meta.TotalRecords,
	}
}

func mapResponseRole(role *db.Role) *pbrole.RoleResponse {
	if role == nil {
		return nil
	}
	var createdAtStr string
	if role.CreatedAt.Valid {
		createdAtStr = role.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var updatedAtStr string
	if role.UpdatedAt.Valid {
		updatedAtStr = role.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return &pbrole.RoleResponse{
		Id:        role.RoleID,
		Name:      role.RoleName,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

func mapResponsesRole(roles []*db.GetRolesRow) []*pbrole.RoleResponse {
	var responseRoles []*pbrole.RoleResponse
	for _, role := range roles {
		if role == nil {
			continue
		}
		var createdAtStr string
		if role.CreatedAt.Valid {
			createdAtStr = role.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var updatedAtStr string
		if role.UpdatedAt.Valid {
			updatedAtStr = role.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		responseRoles = append(responseRoles, &pbrole.RoleResponse{
			Id:        role.RoleID,
			Name:      role.RoleName,
			CreatedAt: createdAtStr,
			UpdatedAt: updatedAtStr,
		})
	}
	return responseRoles
}

func mapResponseUserRole(userRole *db.UserRole) *pbuserrole.UserRoleResponse {
	if userRole == nil {
		return nil
	}
	return &pbuserrole.UserRoleResponse{
		UserRoleId: userRole.UserRoleID,
		UserId:     userRole.UserID,
		RoleId:     userRole.RoleID,
	}
}

func mapResponsesRoleFromDB(roles []*db.Role) []*pbrole.RoleResponse {
	var responseRoles []*pbrole.RoleResponse
	for _, role := range roles {
		responseRoles = append(responseRoles, mapResponseRole(role))
	}
	return responseRoles
}

func mapResponsesRoleFromActive(roles []*db.GetActiveRolesRow) []*pbrole.RoleResponseDeleteAt {
	var responseRoles []*pbrole.RoleResponseDeleteAt
	for _, role := range roles {
		if role == nil {
			continue
		}
		var createdAtStr string
		if role.CreatedAt.Valid {
			createdAtStr = role.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var updatedAtStr string
		if role.UpdatedAt.Valid {
			updatedAtStr = role.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var deletedAtStr string
		if role.DeletedAt.Valid {
			deletedAtStr = role.DeletedAt.Time.Format("2006-01-02 15:04:05")
		}
		responseRoles = append(responseRoles, &pbrole.RoleResponseDeleteAt{
			Id:        role.RoleID,
			Name:      role.RoleName,
			CreatedAt: createdAtStr,
			UpdatedAt: updatedAtStr,
			DeletedAt: deletedAtStr,
		})
	}
	return responseRoles
}

func mapResponsesRoleFromTrashed(roles []*db.GetTrashedRolesRow) []*pbrole.RoleResponseDeleteAt {
	var responseRoles []*pbrole.RoleResponseDeleteAt
	for _, role := range roles {
		if role == nil {
			continue
		}
		var createdAtStr string
		if role.CreatedAt.Valid {
			createdAtStr = role.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var updatedAtStr string
		if role.UpdatedAt.Valid {
			updatedAtStr = role.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var deletedAtStr string
		if role.DeletedAt.Valid {
			deletedAtStr = role.DeletedAt.Time.Format("2006-01-02 15:04:05")
		}
		responseRoles = append(responseRoles, &pbrole.RoleResponseDeleteAt{
			Id:        role.RoleID,
			Name:      role.RoleName,
			CreatedAt: createdAtStr,
			UpdatedAt: updatedAtStr,
			DeletedAt: deletedAtStr,
		})
	}
	return responseRoles
}
