package category_cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

const (
	categoryStatsByIdMonthTotalPriceCacheKey = "category:stats:byid:%d:month:%d:year:%d"
	categoryStatsByIdYearTotalPriceCacheKey  = "category:stats:byid:%d:year:%d"

	categoryStatsByIdMonthPriceCacheKey = "category:stats:byid:%d:month:%d"
	categoryStatsByIdYearPriceCacheKey  = "category:stats:byid:%d:year:%d"
)

type categoryStatsByIdCache struct {
	store *cache.CacheStore
}

func NewCategoryStatsByIdCache(store *cache.CacheStore) *categoryStatsByIdCache {
	return &categoryStatsByIdCache{store: store}
}

func (s *categoryStatsByIdCache) GetCachedMonthTotalPriceByIdCache(ctx context.Context, req *model.FindYearMonthTotalPriceByID) (*model.APIResponseCategoryMonthlyTotalPrice, bool) {
	key := fmt.Sprintf(categoryStatsByIdMonthTotalPriceCacheKey, req.CategoryID, req.Month, req.Year)

	result, found := cache.GetFromCache[model.APIResponseCategoryMonthlyTotalPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsByIdCache) SetCachedMonthTotalPriceByIdCache(ctx context.Context, req *model.FindYearMonthTotalPriceByID, res *model.APIResponseCategoryMonthlyTotalPrice) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByIdMonthTotalPriceCacheKey, req.CategoryID, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedYearTotalPriceByIdCache(ctx context.Context, req *model.FindYearTotalPriceByID) (*model.APIResponseCategoryYearlyTotalPrice, bool) {
	key := fmt.Sprintf(categoryStatsByIdYearTotalPriceCacheKey, req.CategoryID, req.Year)

	result, found := cache.GetFromCache[model.APIResponseCategoryYearlyTotalPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsByIdCache) SetCachedYearTotalPriceByIdCache(ctx context.Context, req *model.FindYearTotalPriceByID, res *model.APIResponseCategoryYearlyTotalPrice) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByIdYearTotalPriceCacheKey, req.CategoryID, req.Year)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedMonthPriceByIdCache(ctx context.Context, req *model.FindYearCategoryByID) (*model.APIResponseCategoryMonthPrice, bool) {
	key := fmt.Sprintf(categoryStatsByIdMonthPriceCacheKey, req.CategoryID, req.Year)

	result, found := cache.GetFromCache[model.APIResponseCategoryMonthPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsByIdCache) SetCachedMonthPriceByIdCache(ctx context.Context, req *model.FindYearCategoryByID, res *model.APIResponseCategoryMonthPrice) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByIdMonthPriceCacheKey, req.CategoryID, req.Year)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedYearPriceByIdCache(ctx context.Context, req *model.FindYearCategoryByID) (*model.APIResponseCategoryYearPrice, bool) {
	key := fmt.Sprintf(categoryStatsByIdYearPriceCacheKey, req.CategoryID, req.Year)

	result, found := cache.GetFromCache[model.APIResponseCategoryYearPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsByIdCache) SetCachedYearPriceByIdCache(ctx context.Context, req *model.FindYearCategoryByID, res *model.APIResponseCategoryYearPrice) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByIdYearPriceCacheKey, req.CategoryID, req.Year)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}
