package category_cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-pointofsale-shared/cache"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
)

const (
	categoryAllCacheKey     = "category:all:page:%d:pageSize:%d:search:%s"
	categoryByIdCacheKey    = "category:id:%d"
	categoryActiveCacheKey  = "category:active:page:%d:pageSize:%d:search:%s"
	categoryTrashedCacheKey = "category:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type categoryQueryCache struct {
	store *cache.CacheStore
}

func NewCategoryQueryCache(store *cache.CacheStore) *categoryQueryCache {
	return &categoryQueryCache{store: store}
}

func (s *categoryQueryCache) GetCachedCategoriesCache(ctx context.Context, req *model.FindAllCategoryRequest) (*model.APIResponsePaginationCategory, bool) {
	key := fmt.Sprintf(categoryAllCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationCategory](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryQueryCache) SetCachedCategoriesCache(ctx context.Context, req *model.FindAllCategoryRequest, res *model.APIResponsePaginationCategory) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(categoryAllCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryActiveCache(ctx context.Context, req *model.FindAllCategoryRequest) (*model.APIResponsePaginationCategoryDeleteAt, bool) {
	key := fmt.Sprintf(categoryActiveCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationCategoryDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryQueryCache) SetCachedCategoryActiveCache(ctx context.Context, req *model.FindAllCategoryRequest, res *model.APIResponsePaginationCategoryDeleteAt) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(categoryActiveCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryTrashedCache(ctx context.Context, req *model.FindAllCategoryRequest) (*model.APIResponsePaginationCategoryDeleteAt, bool) {
	key := fmt.Sprintf(categoryTrashedCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationCategoryDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryQueryCache) SetCachedCategoryTrashedCache(ctx context.Context, req *model.FindAllCategoryRequest, res *model.APIResponsePaginationCategoryDeleteAt) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(categoryTrashedCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryCache(ctx context.Context, id int) (*model.APIResponseCategory, bool) {
	key := fmt.Sprintf(categoryByIdCacheKey, id)
	result, found := cache.GetFromCache[model.APIResponseCategory](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryQueryCache) SetCachedCategoryCache(ctx context.Context, res *model.APIResponseCategory) {
	if res == nil || res.Data == nil {
		return
	}

	key := fmt.Sprintf(categoryByIdCacheKey, res.Data.ID)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
