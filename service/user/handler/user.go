package handler

import (
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func mapUserToProto(user *db.User) *pbuser.UserResponse {
	if user == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if user.CreatedAt.Valid {
		createdAtStr = user.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if user.UpdatedAt.Valid {
		updatedAtStr = user.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return &pbuser.UserResponse{
		Id:        int32(user.UserID),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

func mapGetUsersRowToProto(user *db.GetUsersRow) *pbuser.UserResponse {
	if user == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if user.CreatedAt.Valid {
		createdAtStr = user.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if user.UpdatedAt.Valid {
		updatedAtStr = user.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return &pbuser.UserResponse{
		Id:        int32(user.UserID),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

// mapUserWithPasswordToProto carries the password hash and must only be used for
// the auth-facing RPCs (FindByEmail/FindByEmailAndVerify); never expose it on
// listing endpoints.
func mapUserWithPasswordToProto(user *db.User) *pbuser.UserResponseWithPassword {
	if user == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if user.CreatedAt.Valid {
		createdAtStr = user.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if user.UpdatedAt.Valid {
		updatedAtStr = user.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return &pbuser.UserResponseWithPassword{
		Id:        int32(user.UserID),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

func mapGetUsersRowsToProto(users []*db.GetUsersRow) []*pbuser.UserResponse {
	var responseUsers []*pbuser.UserResponse
	for _, u := range users {
		responseUsers = append(responseUsers, mapGetUsersRowToProto(u))
	}
	return responseUsers
}

func mapUserDeleteAtToProto(user *db.User) *pbuser.UserResponseDeleteAt {
	if user == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	var deletedAt *wrapperspb.StringValue
	if user.CreatedAt.Valid {
		createdAtStr = user.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if user.UpdatedAt.Valid {
		updatedAtStr = user.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if user.DeletedAt.Valid {
		deletedAt = wrapperspb.String(user.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}
	return &pbuser.UserResponseDeleteAt{
		Id:        int32(user.UserID),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
		DeletedAt: deletedAt,
	}
}

func mapActiveUserToProto(user *db.GetUsersActiveRow) *pbuser.UserResponseDeleteAt {
	if user == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	var deletedAt *wrapperspb.StringValue
	if user.CreatedAt.Valid {
		createdAtStr = user.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if user.UpdatedAt.Valid {
		updatedAtStr = user.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if user.DeletedAt.Valid {
		deletedAt = wrapperspb.String(user.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}
	return &pbuser.UserResponseDeleteAt{
		Id:        int32(user.UserID),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
		DeletedAt: deletedAt,
	}
}

func mapActiveUsersToProto(users []*db.GetUsersActiveRow) []*pbuser.UserResponseDeleteAt {
	var responseUsers []*pbuser.UserResponseDeleteAt
	for _, u := range users {
		responseUsers = append(responseUsers, mapActiveUserToProto(u))
	}
	return responseUsers
}

func mapTrashedUserToProto(user *db.GetUserTrashedRow) *pbuser.UserResponseDeleteAt {
	if user == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	var deletedAt *wrapperspb.StringValue
	if user.CreatedAt.Valid {
		createdAtStr = user.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if user.UpdatedAt.Valid {
		updatedAtStr = user.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if user.DeletedAt.Valid {
		deletedAt = wrapperspb.String(user.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}
	return &pbuser.UserResponseDeleteAt{
		Id:        int32(user.UserID),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
		DeletedAt: deletedAt,
	}
}

func mapTrashedUsersToProto(users []*db.GetUserTrashedRow) []*pbuser.UserResponseDeleteAt {
	var responseUsers []*pbuser.UserResponseDeleteAt
	for _, u := range users {
		responseUsers = append(responseUsers, mapTrashedUserToProto(u))
	}
	return responseUsers
}
