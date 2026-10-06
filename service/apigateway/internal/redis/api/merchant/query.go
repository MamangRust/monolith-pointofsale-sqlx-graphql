package merchant_cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

const (
	merchantAllCacheKey      = "merchant:all:page:%d:pageSize:%d:search:%s"
	merchantByIdCacheKey     = "merchant:id:%d"
	merchantActiveCacheKey   = "merchant:active:page:%d:pageSize:%d:search:%s"
	merchantTrashedCacheKey  = "merchant:trashed:page:%d:pageSize:%d:search:%s"
	merchantByUserIdCacheKey = "merchant:user_id:%d"

	ttlDefault = 5 * time.Minute
)

type merchantQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantQueryCache(store *cache.CacheStore) *merchantQueryCache {
	return &merchantQueryCache{store: store}
}

func (m *merchantQueryCache) GetCachedMerchants(ctx context.Context, req *model.FindAllMerchantInput) (*model.APIResponsePaginationMerchant, bool) {
	key := fmt.Sprintf(merchantAllCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationMerchant](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantQueryCache) SetCachedMerchants(ctx context.Context, req *model.FindAllMerchantInput, res *model.APIResponsePaginationMerchant) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(merchantAllCacheKey, req.Page, req.PageSize, safeString(req.Search))

	cache.SetToCache(ctx, m.store, key, res, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantActive(ctx context.Context, req *model.FindAllMerchantInput) (*model.APIResponsePaginationMerchantDeleteAt, bool) {
	key := fmt.Sprintf(merchantActiveCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationMerchantDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantQueryCache) SetCachedMerchantActive(ctx context.Context, req *model.FindAllMerchantInput, res *model.APIResponsePaginationMerchantDeleteAt) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(merchantActiveCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, m.store, key, res, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantTrashed(ctx context.Context, req *model.FindAllMerchantInput) (*model.APIResponsePaginationMerchantDeleteAt, bool) {
	key := fmt.Sprintf(merchantTrashedCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationMerchantDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantQueryCache) SetCachedMerchantTrashed(ctx context.Context, req *model.FindAllMerchantInput, res *model.APIResponsePaginationMerchantDeleteAt) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(merchantTrashedCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, m.store, key, res, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchant(ctx context.Context, id int) (*model.APIResponseMerchant, bool) {
	key := fmt.Sprintf(merchantByIdCacheKey, id)

	result, found := cache.GetFromCache[model.APIResponseMerchant](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantQueryCache) SetCachedMerchant(ctx context.Context, res *model.APIResponseMerchant) {
	if res == nil || res.Data == nil {
		return
	}

	key := fmt.Sprintf(merchantByIdCacheKey, res.Data.ID)
	cache.SetToCache(ctx, m.store, key, res, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantsByUserId(ctx context.Context, id int) (*model.APIResponsesMerchant, bool) {
	key := fmt.Sprintf(merchantByUserIdCacheKey, id)

	result, found := cache.GetFromCache[model.APIResponsesMerchant](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantQueryCache) SetCachedMerchantsByUserId(ctx context.Context, userId int, res *model.APIResponsesMerchant) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(merchantByUserIdCacheKey, userId)
	cache.SetToCache(ctx, m.store, key, res, ttlDefault)
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
