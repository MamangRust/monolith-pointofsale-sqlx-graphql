package order_cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

const (
	orderAllCacheKey      = "order:all:page:%d:pageSize:%d:search:%s"
	orderByIdCacheKey     = "order:id:%d"
	orderActiveCacheKey   = "order:active:page:%d:pageSize:%d:search:%s"
	orderTrashedCacheKey  = "order:trashed:page:%d:pageSize:%d:search:%s"
	orderMerchantCacheKey = "order:merchant:%d:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type orderQueryCache struct {
	store *cache.CacheStore
}

func NewOrderQueryCache(store *cache.CacheStore) *orderQueryCache {
	return &orderQueryCache{store: store}
}

func (s *orderQueryCache) GetOrderAllCache(ctx context.Context, req *model.FindAllOrderInput) (*model.APIResponsePaginationOrder, bool) {
	key := fmt.Sprintf(orderAllCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationOrder](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderQueryCache) SetOrderAllCache(ctx context.Context, req *model.FindAllOrderInput, res *model.APIResponsePaginationOrder) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(orderAllCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *orderQueryCache) GetCachedOrderCache(ctx context.Context, orderID int) (*model.APIResponseOrder, bool) {
	key := fmt.Sprintf(orderByIdCacheKey, orderID)

	result, found := cache.GetFromCache[model.APIResponseOrder](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderQueryCache) SetCachedOrderCache(ctx context.Context, res *model.APIResponseOrder) {
	if res == nil || res.Data == nil {
		return
	}

	key := fmt.Sprintf(orderByIdCacheKey, res.Data.ID)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *orderQueryCache) GetCachedOrderMerchant(ctx context.Context, req *model.FindAllOrderMerchantInput) (*model.APIResponsePaginationOrder, bool) {
	key := fmt.Sprintf(orderMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationOrder](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderQueryCache) SetCachedOrderMerchant(ctx context.Context, req *model.FindAllOrderMerchantInput, res *model.APIResponsePaginationOrder) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(orderMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *orderQueryCache) GetOrderActiveCache(ctx context.Context, req *model.FindAllOrderInput) (*model.APIResponsePaginationOrderDeleteAt, bool) {
	key := fmt.Sprintf(orderActiveCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationOrderDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderQueryCache) SetOrderActiveCache(ctx context.Context, req *model.FindAllOrderInput, res *model.APIResponsePaginationOrderDeleteAt) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(orderActiveCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *orderQueryCache) GetOrderTrashedCache(ctx context.Context, req *model.FindAllOrderInput) (*model.APIResponsePaginationOrderDeleteAt, bool) {
	key := fmt.Sprintf(orderTrashedCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationOrderDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderQueryCache) SetOrderTrashedCache(ctx context.Context, req *model.FindAllOrderInput, res *model.APIResponsePaginationOrderDeleteAt) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(orderTrashedCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
