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
	orderAllCacheKey      = "order:all:page:%d:pageSize:%d:search:%s"
	orderByIdCacheKey     = "order:id:%d"
	orderMerchantCacheKey = "order:merchant:%d:page:%d:pageSize:%d:search:%s"
	orderActiveCacheKey   = "order:active:page:%d:pageSize:%d:search:%s"
	orderTrashedCacheKey  = "order:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type orderCacheResponse struct {
	Data         []*db.GetOrdersRow `json:"data"`
	TotalRecords *int               `json:"total_records"`
}

type orderCacheResponseByMerchant struct {
	Data         []*db.GetOrdersByMerchantRow `json:"data"`
	TotalRecords *int                         `json:"total_records"`
}

type orderCacheResponseActive struct {
	Data         []*db.GetOrdersActiveRow `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type orderCacheResponseTrashed struct {
	Data         []*db.GetOrdersTrashedRow `json:"data"`
	TotalRecords *int                      `json:"total_records"`
}

type orderQueryCache struct {
	store *cache.CacheStore
}

func NewOrderQueryCache(store *cache.CacheStore) OrderQueryCache {
	return &orderQueryCache{store: store}
}

func (s *orderQueryCache) GetOrderAllCache(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersRow, *int, bool) {
	key := fmt.Sprintf(orderAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[orderCacheResponse](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *orderQueryCache) SetOrderAllCache(ctx context.Context, req *requests.FindAllOrders, data []*db.GetOrdersRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetOrdersRow{}
	}

	key := fmt.Sprintf(orderAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderCacheResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *orderQueryCache) GetCachedOrderMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*db.GetOrdersByMerchantRow, *int, bool) {
	key := fmt.Sprintf(orderMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[orderCacheResponseByMerchant](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *orderQueryCache) SetCachedOrderMerchant(ctx context.Context, req *requests.FindAllOrderMerchant, res []*db.GetOrdersByMerchantRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if res == nil {
		res = []*db.GetOrdersByMerchantRow{}
	}

	key := fmt.Sprintf(orderMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	payload := &orderCacheResponseByMerchant{Data: res, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *orderQueryCache) GetOrderActiveCache(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersActiveRow, *int, bool) {
	key := fmt.Sprintf(orderActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[orderCacheResponseActive](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *orderQueryCache) SetOrderActiveCache(ctx context.Context, req *requests.FindAllOrders, data []*db.GetOrdersActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetOrdersActiveRow{}
	}

	key := fmt.Sprintf(orderActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderCacheResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *orderQueryCache) GetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersTrashedRow, *int, bool) {
	key := fmt.Sprintf(orderTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[orderCacheResponseTrashed](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *orderQueryCache) SetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrders, data []*db.GetOrdersTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetOrdersTrashedRow{}
	}

	key := fmt.Sprintf(orderTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderCacheResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *orderQueryCache) GetCachedOrderCache(ctx context.Context, orderID int) (*db.Order, bool) {
	key := fmt.Sprintf(orderByIdCacheKey, orderID)
	result, found := cache.GetFromCache[db.Order](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *orderQueryCache) SetCachedOrderCache(ctx context.Context, data *db.Order) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(orderByIdCacheKey, data.OrderID)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
