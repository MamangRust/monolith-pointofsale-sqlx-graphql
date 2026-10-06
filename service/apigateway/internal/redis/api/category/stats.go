package category_cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

const (
	categoryStatsMonthTotalPriceCacheKey = "category:stats:month:%d:year:%d"
	categoryStatsYearTotalPriceCacheKey  = "category:stats:year:%d"

	categoryStatsMonthPriceCacheKey = "category:stats:month:%d"
	categoryStatsYearPriceCacheKey  = "category:stats:year:%d"
)

type categoryStatsCache struct {
	store *cache.CacheStore
}

func NewCategoryStatsCache(store *cache.CacheStore) *categoryStatsCache {
	return &categoryStatsCache{store: store}
}

func (s *categoryStatsCache) GetCachedMonthTotalPriceCache(ctx context.Context, req *model.FindYearMonthTotalPrices) (*model.APIResponseCategoryMonthlyTotalPrice, bool) {
	key := fmt.Sprintf(categoryStatsMonthTotalPriceCacheKey, req.Month, req.Year)

	result, found := cache.GetFromCache[model.APIResponseCategoryMonthlyTotalPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsCache) SetCachedMonthTotalPriceCache(ctx context.Context, req *model.FindYearMonthTotalPrices, res *model.APIResponseCategoryMonthlyTotalPrice) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsMonthTotalPriceCacheKey, req.Month, req.Year)

	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *categoryStatsCache) GetCachedYearTotalPriceCache(ctx context.Context, year int) (*model.APIResponseCategoryYearlyTotalPrice, bool) {
	key := fmt.Sprintf(categoryStatsYearTotalPriceCacheKey, year)
	result, found := cache.GetFromCache[model.APIResponseCategoryYearlyTotalPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsCache) SetCachedYearTotalPriceCache(ctx context.Context, year int, res *model.APIResponseCategoryYearlyTotalPrice) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsYearTotalPriceCacheKey, year)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *categoryStatsCache) GetCachedMonthPriceCache(ctx context.Context, year int) (*model.APIResponseCategoryMonthPrice, bool) {
	key := fmt.Sprintf(categoryStatsMonthPriceCacheKey, year)
	result, found := cache.GetFromCache[model.APIResponseCategoryMonthPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsCache) SetCachedMonthPriceCache(ctx context.Context, year int, res *model.APIResponseCategoryMonthPrice) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsMonthPriceCacheKey, year)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *categoryStatsCache) GetCachedYearPriceCache(ctx context.Context, year int) (*model.APIResponseCategoryYearPrice, bool) {
	key := fmt.Sprintf(categoryStatsYearPriceCacheKey, year)
	result, found := cache.GetFromCache[model.APIResponseCategoryYearPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsCache) SetCachedYearPriceCache(ctx context.Context, year int, res *model.APIResponseCategoryYearPrice) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsYearPriceCacheKey, year)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}
