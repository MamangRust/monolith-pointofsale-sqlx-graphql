package category_cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

const (
	categoryStatsByMerchantMonthTotalPriceCacheKey = "category:stats:bymerchant:%d:month:%d:year:%d"
	categoryStatsByMerchantYearTotalPriceCacheKey  = "category:stats:bymerchant:%d:year:%d"

	categoryStatsByMerchantMonthPriceCacheKey = "category:stats:bymerchant:%d:month:%d"
	categoryStatsByMerchantYearPriceCacheKey  = "category:stats:bymerchant:%d:year:%d"
)

type categoryStatsByMerchantCache struct {
	store *cache.CacheStore
}

func NewCategoryStatsByMerchantCache(store *cache.CacheStore) *categoryStatsByMerchantCache {
	return &categoryStatsByMerchantCache{store: store}
}

func (s *categoryStatsByMerchantCache) GetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *model.FindYearMonthTotalPriceByMerchant) (*model.APIResponseCategoryMonthlyTotalPrice, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantMonthTotalPriceCacheKey, req.MerchantID, req.Month, req.Year)

	result, found := cache.GetFromCache[model.APIResponseCategoryMonthlyTotalPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsByMerchantCache) SetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *model.FindYearMonthTotalPriceByMerchant, res *model.APIResponseCategoryMonthlyTotalPrice) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByMerchantMonthTotalPriceCacheKey, req.MerchantID, req.Month, req.Year)

	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *categoryStatsByMerchantCache) GetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *model.FindYearTotalPriceByMerchant) (*model.APIResponseCategoryYearlyTotalPrice, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantYearTotalPriceCacheKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[model.APIResponseCategoryYearlyTotalPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsByMerchantCache) SetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *model.FindYearTotalPriceByMerchant, res *model.APIResponseCategoryYearlyTotalPrice) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByMerchantYearTotalPriceCacheKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *categoryStatsByMerchantCache) GetCachedMonthPriceByMerchantCache(ctx context.Context, req *model.FindYearCategoryByMerchant) (*model.APIResponseCategoryMonthPrice, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantMonthPriceCacheKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[model.APIResponseCategoryMonthPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsByMerchantCache) SetCachedMonthPriceByMerchantCache(ctx context.Context, req *model.FindYearCategoryByMerchant, res *model.APIResponseCategoryMonthPrice) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByMerchantMonthPriceCacheKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *categoryStatsByMerchantCache) GetCachedYearPriceByMerchantCache(ctx context.Context, req *model.FindYearCategoryByMerchant) (*model.APIResponseCategoryYearPrice, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantYearPriceCacheKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[model.APIResponseCategoryYearPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsByMerchantCache) SetCachedYearPriceByMerchantCache(ctx context.Context, req *model.FindYearCategoryByMerchant, res *model.APIResponseCategoryYearPrice) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByMerchantYearPriceCacheKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}
