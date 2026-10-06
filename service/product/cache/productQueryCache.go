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
	productAllCacheKey      = "product:all:page:%d:pageSize:%d:search:%s"
	productCategoryCacheKey = "product:category:%s:page:%d:pageSize:%d:search:%s"
	productMerchantCacheKey = "product:merchant:%d:page:%d:pageSize:%d:search:%s"

	productActiveCacheKey  = "product:active:page:%d:pageSize:%d:search:%s"
	productTrashedCacheKey = "product:trashed:page:%d:pageSize:%d:search:%s"
	productByIdCacheKey    = "product:id:%d"

	ttlDefault = 5 * time.Minute
)

type productCacheResponse struct {
	Data         []*db.GetProductsRow `json:"data"`
	TotalRecords *int                 `json:"total_records"`
}

type productCacheResponseMerchant struct {
	Data         []*db.GetProductsByMerchantRow `json:"data"`
	TotalRecords *int                           `json:"total_records"`
}

type productCacheResponseCategory struct {
	Data         []*db.GetProductsByCategoryNameRow `json:"data"`
	TotalRecords *int                               `json:"total_records"`
}

type productCacheResponseActive struct {
	Data         []*db.GetProductsActiveRow `json:"data"`
	TotalRecords *int                       `json:"total_records"`
}

type productCacheResponseTrashed struct {
	Data         []*db.GetProductsTrashedRow `json:"data"`
	TotalRecords *int                        `json:"total_records"`
}

type productQueryCache struct {
	store *cache.CacheStore
}

func NewProductQueryCache(store *cache.CacheStore) ProductQueryCache {
	return &productQueryCache{store: store}
}

func (p *productQueryCache) GetCachedProducts(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsRow, *int, bool) {
	key := fmt.Sprintf(productAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[productCacheResponse](ctx, p.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProducts(ctx context.Context, req *requests.FindAllProducts, data []*db.GetProductsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetProductsRow{}
	}

	key := fmt.Sprintf(productAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductsByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*db.GetProductsByMerchantRow, *int, bool) {
	key := fmt.Sprintf(productMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[productCacheResponseMerchant](ctx, p.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductsByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest, data []*db.GetProductsByMerchantRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetProductsByMerchantRow{}
	}

	key := fmt.Sprintf(productMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponseMerchant{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductsByCategory(ctx context.Context, req *requests.ProductByCategoryRequest) ([]*db.GetProductsByCategoryNameRow, *int, bool) {
	key := fmt.Sprintf(productCategoryCacheKey, req.CategoryName, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[productCacheResponseCategory](ctx, p.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductsByCategory(ctx context.Context, req *requests.ProductByCategoryRequest, data []*db.GetProductsByCategoryNameRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetProductsByCategoryNameRow{}
	}

	key := fmt.Sprintf(productCategoryCacheKey, req.CategoryName, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponseCategory{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductActive(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsActiveRow, *int, bool) {
	key := fmt.Sprintf(productActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[productCacheResponseActive](ctx, p.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductActive(ctx context.Context, req *requests.FindAllProducts, data []*db.GetProductsActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetProductsActiveRow{}
	}

	key := fmt.Sprintf(productActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsTrashedRow, *int, bool) {
	key := fmt.Sprintf(productTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[productCacheResponseTrashed](ctx, p.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductTrashed(ctx context.Context, req *requests.FindAllProducts, data []*db.GetProductsTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetProductsTrashedRow{}
	}

	key := fmt.Sprintf(productTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProduct(ctx context.Context, productID int) (*db.Product, bool) {
	key := fmt.Sprintf(productByIdCacheKey, productID)

	result, found := cache.GetFromCache[db.Product](ctx, p.store, key)
	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (p *productQueryCache) SetCachedProduct(ctx context.Context, data *db.Product) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(productByIdCacheKey, data.ProductID)
	cache.SetToCache(ctx, p.store, key, data, ttlDefault)
}
