package mencache

import (
	"context"
	"fmt"
	"time"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

const (
	roleAllCacheKey      = "role:all:page:%d:pageSize:%d:search:%s"
	roleByIdCacheKey     = "role:id:%d"
	roleByUserIdCacheKey = "role:user:%d"
	roleActiveCacheKey   = "role:active:page:%d:pageSize:%d:search:%s"
	roleTrashedCacheKey  = "role:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type roleCachedResponse struct {
	Data         []*db.GetRolesRow `json:"data"`
	TotalRecords *int              `json:"total_records"`
}

type roleCachedResponseActive struct {
	Data         []*db.GetActiveRolesRow `json:"data"`
	TotalRecords *int                    `json:"total_records"`
}

type roleCachedResponseTrashed struct {
	Data         []*db.GetTrashedRolesRow `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type roleQueryCache struct {
	store *cache.CacheStore
}

func NewRoleQueryCache(store *cache.CacheStore) RoleQueryCache {
	return &roleQueryCache{store: store}
}

func (m *roleQueryCache) SetCachedRoles(ctx context.Context, req *requests.FindAllRoles, data []*db.GetRolesRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetRolesRow{}
	}

	key := fmt.Sprintf(roleAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &roleCachedResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *roleQueryCache) SetCachedRoleById(ctx context.Context, data *db.Role) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleByIdCacheKey, data.RoleID)
	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *roleQueryCache) SetCachedRoleByUserId(ctx context.Context, userId int, data []*db.Role) {
	if data == nil {
		data = []*db.Role{}
	}

	key := fmt.Sprintf(roleByUserIdCacheKey, userId)
	cache.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

func (m *roleQueryCache) SetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles, data []*db.GetActiveRolesRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetActiveRolesRow{}
	}

	key := fmt.Sprintf(roleActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &roleCachedResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *roleQueryCache) SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles, data []*db.GetTrashedRolesRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetTrashedRolesRow{}
	}

	key := fmt.Sprintf(roleTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &roleCachedResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *roleQueryCache) GetCachedRoles(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetRolesRow, *int, bool) {
	key := fmt.Sprintf(roleAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[roleCachedResponse](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *roleQueryCache) GetCachedRoleById(ctx context.Context, id int) (*db.Role, bool) {
	key := fmt.Sprintf(roleByIdCacheKey, id)

	result, found := cache.GetFromCache[db.Role](ctx, m.store, key)
	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *roleQueryCache) GetCachedRoleByUserId(ctx context.Context, userId int) ([]*db.Role, bool) {
	key := fmt.Sprintf(roleByUserIdCacheKey, userId)

	result, found := cache.GetFromCache[[]*db.Role](ctx, m.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *roleQueryCache) GetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetActiveRolesRow, *int, bool) {
	key := fmt.Sprintf(roleActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[roleCachedResponseActive](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *roleQueryCache) GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetTrashedRolesRow, *int, bool) {
	key := fmt.Sprintf(roleTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[roleCachedResponseTrashed](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}
