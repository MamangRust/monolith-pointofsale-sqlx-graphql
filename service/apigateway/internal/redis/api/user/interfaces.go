package user_cache

import (
	"context"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
)

type UserQueryCache interface {
	GetCachedUsersCache(ctx context.Context, req *model.FindAllUserInput) (*model.APIResponsePaginationUser, bool)
	SetCachedUsersCache(ctx context.Context, req *model.FindAllUserInput, data *model.APIResponsePaginationUser)

	GetCachedUserActiveCache(ctx context.Context, req *model.FindAllUserInput) (*model.APIResponsePaginationUserDeleteAt, bool)
	SetCachedUserActiveCache(ctx context.Context, req *model.FindAllUserInput, data *model.APIResponsePaginationUserDeleteAt)

	GetCachedUserTrashedCache(ctx context.Context, req *model.FindAllUserInput) (*model.APIResponsePaginationUserDeleteAt, bool)
	SetCachedUserTrashedCache(ctx context.Context, req *model.FindAllUserInput, data *model.APIResponsePaginationUserDeleteAt)

	GetCachedUserCache(ctx context.Context, id int) (*model.APIResponseUserResponse, bool)
	SetCachedUserCache(ctx context.Context, data *model.APIResponseUserResponse)
}

type UserCommandCache interface {
	DeleteUserCache(ctx context.Context, id int)
}
