package orderitem_cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

const (
	orderItemAllCacheKey     = "order_item:all:page:%d:pageSize:%d:search:%s"
	orderItemActiveCacheKey  = "order_item:active:page:%d:pageSize:%d:search:%s"
	orderItemTrashedCacheKey = "order_item:trashed:page:%d:pageSize:%d:search:%s"

	orderItemByIdCacheKey = "order_item:id:%d"

	ttlDefault = 5 * time.Minute
)

type orderItemQueryCache struct {
	store *cache.CacheStore
}

func NewOrderItemQueryCache(store *cache.CacheStore) *orderItemQueryCache {
	return &orderItemQueryCache{store: store}
}

func (o *orderItemQueryCache) GetCachedOrderItemsAll(ctx context.Context, req *model.FindAllOrderItemInput) (*model.APIResponsePaginationOrderItem, bool) {
	key := fmt.Sprintf(orderItemAllCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationOrderItem](ctx, o.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (o *orderItemQueryCache) SetCachedOrderItemsAll(ctx context.Context, req *model.FindAllOrderItemInput, res *model.APIResponsePaginationOrderItem) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(orderItemAllCacheKey, req.Page, req.PageSize, safeString(req.Search))

	cache.SetToCache(ctx, o.store, key, res, ttlDefault)
}

func (o *orderItemQueryCache) GetCachedOrderItemActive(ctx context.Context, req *model.FindAllOrderItemInput) (*model.APIResponsePaginationOrderItemDeleteAt, bool) {
	key := fmt.Sprintf(orderItemActiveCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationOrderItemDeleteAt](ctx, o.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (o *orderItemQueryCache) SetCachedOrderItemActive(ctx context.Context, req *model.FindAllOrderItemInput, res *model.APIResponsePaginationOrderItemDeleteAt) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(orderItemActiveCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, o.store, key, res, ttlDefault)
}

func (o *orderItemQueryCache) GetCachedOrderItemTrashed(ctx context.Context, req *model.FindAllOrderItemInput) (*model.APIResponsePaginationOrderItemDeleteAt, bool) {
	key := fmt.Sprintf(orderItemTrashedCacheKey, req.Page, req.PageSize, safeString(req.Search))
	result, found := cache.GetFromCache[model.APIResponsePaginationOrderItemDeleteAt](ctx, o.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (o *orderItemQueryCache) SetCachedOrderItemTrashed(ctx context.Context, req *model.FindAllOrderItemInput, res *model.APIResponsePaginationOrderItemDeleteAt) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(orderItemTrashedCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, o.store, key, res, ttlDefault)
}

func (o *orderItemQueryCache) GetCachedOrderItems(ctx context.Context, orderID int) (*model.APIResponsesOrderItem, bool) {
	key := fmt.Sprintf(orderItemByIdCacheKey, orderID)
	result, found := cache.GetFromCache[model.APIResponsesOrderItem](ctx, o.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (o *orderItemQueryCache) SetCachedOrderItems(ctx context.Context, res *model.APIResponsesOrderItem) {
	if res == nil || len(res.Data) == 0 {
		return
	}

	key := fmt.Sprintf(orderItemByIdCacheKey, res.Data[0].OrderID)
	cache.SetToCache(ctx, o.store, key, res, ttlDefault)
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
