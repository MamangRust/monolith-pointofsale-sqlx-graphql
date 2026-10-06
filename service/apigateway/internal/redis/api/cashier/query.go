package cashier_cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-shared/cache"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
)

type cashierQueryCache struct {
	store *cache.CacheStore
}

func NewCashierQueryCache(store *cache.CacheStore) *cashierQueryCache {
	return &cashierQueryCache{store: store}
}

func (s *cashierQueryCache) GetCachedCashiersCache(ctx context.Context, req *model.FindAllCashierRequest) (*model.APIResponsePaginationCashier, bool) {
	key := fmt.Sprintf(cashierAllCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationCashier](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *cashierQueryCache) SetCachedCashiersCache(ctx context.Context, req *model.FindAllCashierRequest, res *model.APIResponsePaginationCashier) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(cashierAllCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *cashierQueryCache) GetCachedCashier(ctx context.Context, cashierID int) (*model.APIResponseCashier, bool) {
	key := fmt.Sprintf(cashierByIdCacheKey, cashierID)

	result, found := cache.GetFromCache[model.APIResponseCashier](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *cashierQueryCache) SetCachedCashier(ctx context.Context, res *model.APIResponseCashier) {
	if res == nil || res.Data == nil {
		return
	}

	key := fmt.Sprintf(cashierByIdCacheKey, res.Data.ID)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *cashierQueryCache) GetCachedCashiersActive(ctx context.Context, req *model.FindAllCashierRequest) (*model.APIResponsePaginationCashierDeleteAt, bool) {
	key := fmt.Sprintf(cashierActiveCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationCashierDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *cashierQueryCache) SetCachedCashiersActive(ctx context.Context, req *model.FindAllCashierRequest, res *model.APIResponsePaginationCashierDeleteAt) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(cashierActiveCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *cashierQueryCache) GetCachedCashiersTrashed(ctx context.Context, req *model.FindAllCashierRequest) (*model.APIResponsePaginationCashierDeleteAt, bool) {
	key := fmt.Sprintf(cashierTrashedCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationCashierDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *cashierQueryCache) SetCachedCashiersTrashed(ctx context.Context, req *model.FindAllCashierRequest, res *model.APIResponsePaginationCashierDeleteAt) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(cashierTrashedCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *cashierQueryCache) GetCachedCashiersByMerchant(ctx context.Context, req *model.FindByMerchantCashierRequest) (*model.APIResponsePaginationCashier, bool) {
	key := fmt.Sprintf(cashierByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationCashier](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *cashierQueryCache) SetCachedCashiersByMerchant(ctx context.Context, req *model.FindByMerchantCashierRequest, res *model.APIResponsePaginationCashier) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(cashierByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
