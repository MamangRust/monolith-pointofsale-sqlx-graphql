package product_cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
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

type productQueryCache struct {
	store *cache.CacheStore
}

func NewProductQueryCache(store *cache.CacheStore) *productQueryCache {
	return &productQueryCache{store: store}
}

func (p *productQueryCache) GetCachedProducts(ctx context.Context, req *model.FindAllProductInput) (*model.APIResponsePaginationProduct, bool) {
	key := fmt.Sprintf(productAllCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationProduct](ctx, p.store, key)

	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (p *productQueryCache) SetCachedProducts(ctx context.Context, req *model.FindAllProductInput, res *model.APIResponsePaginationProduct) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(productAllCacheKey, req.Page, req.PageSize, safeString(req.Search))

	cache.SetToCache(ctx, p.store, key, res, ttlDefault)
}

func (p *productQueryCache) GetCachedProductsByMerchant(ctx context.Context, req *model.FindAllProductMerchantInput) (*model.APIResponsePaginationProduct, bool) {
	key := fmt.Sprintf(productMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationProduct](ctx, p.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (p *productQueryCache) SetCachedProductsByMerchant(ctx context.Context, req *model.FindAllProductMerchantInput, res *model.APIResponsePaginationProduct) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(productMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, p.store, key, res, ttlDefault)
}

func (p *productQueryCache) GetCachedProductsByCategory(ctx context.Context, req *model.FindAllProductCategoryInput) (*model.APIResponsePaginationProduct, bool) {
	key := fmt.Sprintf(productCategoryCacheKey, safeString(req.CategoryName), req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationProduct](ctx, p.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (p *productQueryCache) SetCachedProductsByCategory(ctx context.Context, req *model.FindAllProductCategoryInput, res *model.APIResponsePaginationProduct) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(productCategoryCacheKey, safeString(req.CategoryName), req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, p.store, key, res, ttlDefault)
}

func (p *productQueryCache) GetCachedProductActive(ctx context.Context, req *model.FindAllProductInput) (*model.APIResponsePaginationProductDeleteAt, bool) {
	key := fmt.Sprintf(productActiveCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationProductDeleteAt](ctx, p.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (p *productQueryCache) SetCachedProductActive(ctx context.Context, req *model.FindAllProductInput, res *model.APIResponsePaginationProductDeleteAt) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(productActiveCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, p.store, key, res, ttlDefault)
}

func (p *productQueryCache) GetCachedProductTrashed(ctx context.Context, req *model.FindAllProductInput) (*model.APIResponsePaginationProductDeleteAt, bool) {
	key := fmt.Sprintf(productTrashedCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationProductDeleteAt](ctx, p.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (p *productQueryCache) SetCachedProductTrashed(ctx context.Context, req *model.FindAllProductInput, res *model.APIResponsePaginationProductDeleteAt) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(productTrashedCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, p.store, key, res, ttlDefault)
}

func (p *productQueryCache) GetCachedProduct(ctx context.Context, productID int) (*model.APIResponseProduct, bool) {
	key := fmt.Sprintf(productByIdCacheKey, productID)

	result, found := cache.GetFromCache[model.APIResponseProduct](ctx, p.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (p *productQueryCache) SetCachedProduct(ctx context.Context, res *model.APIResponseProduct) {
	if res == nil || res.Data == nil {
		return
	}

	key := fmt.Sprintf(productByIdCacheKey, res.Data.ID)
	cache.SetToCache(ctx, p.store, key, res, ttlDefault)
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
