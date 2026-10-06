package transaction_cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

const (
	transactionAllCacheKey  = "transaction:all:page:%d:pageSize:%d:search:%s"
	transactionByIdCacheKey = "transaction:id:%d"

	transactionByMerchantCacheKey = "transaction:merchant:%d:page:%d:pageSize:%d:search:%s"

	transactionActiveCacheKey  = "transaction:active:page:%d:pageSize:%d:search:%s"
	transactionTrashedCacheKey = "transaction:trashed:page:%d:pageSize:%d:search:%s"

	transactionByOrderCacheKey = "transaction:order:%d"

	ttlDefault = 5 * time.Minute
)

type transactionQueryCache struct {
	store *cache.CacheStore
}

func NewTransactionQueryCache(store *cache.CacheStore) *transactionQueryCache {
	return &transactionQueryCache{store: store}
}

func (t *transactionQueryCache) GetCachedTransactionsCache(ctx context.Context, req *model.FindAllTransactionInput) (*model.APIResponsePaginationTransaction, bool) {
	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationTransaction](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionsCache(ctx context.Context, req *model.FindAllTransactionInput, res *model.APIResponsePaginationTransaction) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, safeString(req.Search))

	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionByMerchant(ctx context.Context, req *model.FindAllTransactionMerchantInput) (*model.APIResponsePaginationTransaction, bool) {
	key := fmt.Sprintf(transactionByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationTransaction](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionByMerchant(ctx context.Context, req *model.FindAllTransactionMerchantInput, res *model.APIResponsePaginationTransaction) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionActiveCache(ctx context.Context, req *model.FindAllTransactionInput) (*model.APIResponsePaginationTransactionDeleteAt, bool) {
	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationTransactionDeleteAt](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionActiveCache(ctx context.Context, req *model.FindAllTransactionInput, res *model.APIResponsePaginationTransactionDeleteAt) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionTrashedCache(ctx context.Context, req *model.FindAllTransactionInput) (*model.APIResponsePaginationTransactionDeleteAt, bool) {
	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, safeString(req.Search))

	result, found := cache.GetFromCache[model.APIResponsePaginationTransactionDeleteAt](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionTrashedCache(ctx context.Context, req *model.FindAllTransactionInput, res *model.APIResponsePaginationTransactionDeleteAt) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, safeString(req.Search))
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionCache(ctx context.Context, id int) (*model.APIResponseTransaction, bool) {
	key := fmt.Sprintf(transactionByIdCacheKey, id)

	result, found := cache.GetFromCache[model.APIResponseTransaction](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionCache(ctx context.Context, res *model.APIResponseTransaction) {
	if res == nil || res.Data == nil {
		return
	}

	key := fmt.Sprintf(transactionByIdCacheKey, res.Data.ID)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionByOrderId(ctx context.Context, orderID int) (*model.APIResponseTransaction, bool) {
	key := fmt.Sprintf(transactionByOrderCacheKey, orderID)

	result, found := cache.GetFromCache[model.APIResponseTransaction](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionByOrderId(ctx context.Context, orderID int, res *model.APIResponseTransaction) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionByOrderCacheKey, orderID)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
