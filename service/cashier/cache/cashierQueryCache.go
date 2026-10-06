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
	cashierAllCacheKey     = "cashier:all:page:%d:pageSize:%d:search:%s"
	cashierByIdCacheKey    = "cashier:id:%d"
	cashierActiveCacheKey  = "cashier:active:page:%d:pageSize:%d:search:%s"
	cashierTrashedCacheKey = "cashier:trashed:page:%d:pageSize:%d:search:%s"

	cashierByMerchantCacheKey = "cashier:merchant:%d:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type cashierCacheResponse struct {
	Data         []*db.GetCashiersRow `json:"data"`
	TotalRecords *int                 `json:"totalRecords"`
}

type cashierCacheResponseActive struct {
	Data         []*db.GetCashiersActiveRow `json:"data"`
	TotalRecords *int                       `json:"totalRecords"`
}

type cashierCacheResponseTrashed struct {
	Data         []*db.GetCashiersTrashedRow `json:"data"`
	TotalRecords *int                        `json:"totalRecords"`
}

type cashierCacheResponseMerchant struct {
	Data         []*db.GetCashiersByMerchantRow `json:"data"`
	TotalRecords *int                           `json:"totalRecords"`
}

type cashierQueryCache struct {
	store *cache.CacheStore
}

func NewCashierQueryCache(store *cache.CacheStore) CashierQueryCache {
	return &cashierQueryCache{store: store}
}

func (s *cashierQueryCache) GetCachedCashiersCache(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersRow, *int, bool) {
	key := fmt.Sprintf(cashierAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[cashierCacheResponse](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *cashierQueryCache) SetCachedCashiersCache(ctx context.Context, req *requests.FindAllCashiers, data []*db.GetCashiersRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetCashiersRow{}
	}

	key := fmt.Sprintf(cashierAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &cashierCacheResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *cashierQueryCache) GetCachedCashiersByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*db.GetCashiersByMerchantRow, *int, bool) {
	key := fmt.Sprintf(cashierByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[cashierCacheResponseMerchant](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *cashierQueryCache) SetCachedCashiersByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant, data []*db.GetCashiersByMerchantRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetCashiersByMerchantRow{}
	}

	key := fmt.Sprintf(cashierByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	payload := &cashierCacheResponseMerchant{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *cashierQueryCache) GetCachedCashiersActive(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersActiveRow, *int, bool) {
	key := fmt.Sprintf(cashierActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[cashierCacheResponseActive](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *cashierQueryCache) SetCachedCashiersActive(ctx context.Context, req *requests.FindAllCashiers, data []*db.GetCashiersActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetCashiersActiveRow{}
	}

	key := fmt.Sprintf(cashierActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &cashierCacheResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *cashierQueryCache) GetCachedCashiersTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersTrashedRow, *int, bool) {
	key := fmt.Sprintf(cashierTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[cashierCacheResponseTrashed](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *cashierQueryCache) SetCachedCashiersTrashed(ctx context.Context, req *requests.FindAllCashiers, data []*db.GetCashiersTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetCashiersTrashedRow{}
	}

	key := fmt.Sprintf(cashierTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &cashierCacheResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *cashierQueryCache) GetCachedCashier(ctx context.Context, id int) (*db.Cashier, bool) {
	key := fmt.Sprintf(cashierByIdCacheKey, id)
	result, found := cache.GetFromCache[*db.Cashier](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *cashierQueryCache) SetCachedCashier(ctx context.Context, data *db.Cashier) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cashierByIdCacheKey, data.CashierID)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
