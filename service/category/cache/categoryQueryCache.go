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
	categoryAllCacheKey     = "category:all:page:%d:pageSize:%d:search:%s"
	categoryByIdCacheKey    = "category:id:%d"
	categoryActiveCacheKey  = "category:active:page:%d:pageSize:%d:search:%s"
	categoryTrashedCacheKey = "category:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type categoryCacheResponse struct {
	Data         []*db.GetCategoriesRow `json:"data"`
	TotalRecords *int                   `json:"totalRecords"`
}

type categoryCacheResponseActive struct {
	Data         []*db.GetCategoriesActiveRow `json:"data"`
	TotalRecords *int                         `json:"totalRecords"`
}

type categoryCacheResponseTrashed struct {
	Data         []*db.GetCategoriesTrashedRow `json:"data"`
	TotalRecords *int                          `json:"totalRecords"`
}

type categoryQueryCache struct {
	store *cache.CacheStore
}

func NewCategoryQueryCache(store *cache.CacheStore) CategoryQueryCache {
	return &categoryQueryCache{store: store}
}

func (s *categoryQueryCache) GetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesRow, *int, bool) {
	key := fmt.Sprintf(categoryAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[categoryCacheResponse](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *categoryQueryCache) SetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory, data []*db.GetCategoriesRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetCategoriesRow{}
	}
	key := fmt.Sprintf(categoryAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &categoryCacheResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesActiveRow, *int, bool) {
	key := fmt.Sprintf(categoryActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[categoryCacheResponseActive](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *categoryQueryCache) SetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory, data []*db.GetCategoriesActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetCategoriesActiveRow{}
	}
	key := fmt.Sprintf(categoryActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &categoryCacheResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesTrashedRow, *int, bool) {
	key := fmt.Sprintf(categoryTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[categoryCacheResponseTrashed](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *categoryQueryCache) SetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory, data []*db.GetCategoriesTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetCategoriesTrashedRow{}
	}
	key := fmt.Sprintf(categoryTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &categoryCacheResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryCache(ctx context.Context, id int) (*db.Category, bool) {
	key := fmt.Sprintf(categoryByIdCacheKey, id)
	result, found := cache.GetFromCache[db.Category](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *categoryQueryCache) SetCachedCategoryCache(ctx context.Context, data *db.Category) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryByIdCacheKey, data.CategoryID)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
