package transaction_cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

const (
	transactionMonthAmountSuccessKey = "transaction:month:amount:success:month:%d:year:%d"
	transactionMonthAmountFailedKey  = "transaction:month:amount:failed:month:%d:year:%d"

	transactionYearAmountSuccessKey = "transaction:year:amount:success:year:%d"
	transactionYearAmountFailedKey  = "transaction:year:amount:failed:year:%d"

	transactionMonthMethodSuccessKey = "transaction:month:method:success:month:%d:year:%d"
	transactionMonthMethodFailedKey  = "transaction:month:method:failed:month:%d:year:%d"

	transactionYearMethodSuccessKey = "transaction:year:method:success:year:%d"
	transactionYearMethodFailedKey  = "transaction:year:method:failed:year:%d"
)

type transactionStatsCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsCache(store *cache.CacheStore) *transactionStatsCache {
	return &transactionStatsCache{store: store}
}

func (t *transactionStatsCache) GetCachedMonthAmountSuccessCached(ctx context.Context, req *model.FindMonthlyTransactionStatusInput) (*model.APIResponseTransactionMonthAmountSuccess, bool) {
	key := fmt.Sprintf(transactionMonthAmountSuccessKey, req.Month, req.Year)

	result, found := cache.GetFromCache[model.APIResponseTransactionMonthAmountSuccess](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsCache) SetCachedMonthAmountSuccessCached(ctx context.Context, req *model.FindMonthlyTransactionStatusInput, res *model.APIResponseTransactionMonthAmountSuccess) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthAmountSuccessKey, req.Month, req.Year)

	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearAmountSuccessCached(ctx context.Context, year int) (*model.APIResponseTransactionYearAmountSuccess, bool) {
	key := fmt.Sprintf(transactionYearAmountSuccessKey, year)

	result, found := cache.GetFromCache[model.APIResponseTransactionYearAmountSuccess](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsCache) SetCachedYearAmountSuccessCached(ctx context.Context, year int, res *model.APIResponseTransactionYearAmountSuccess) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearAmountSuccessKey, year)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthAmountFailedCached(ctx context.Context, req *model.FindMonthlyTransactionStatusInput) (*model.APIResponseTransactionMonthAmountFailed, bool) {
	key := fmt.Sprintf(transactionMonthAmountFailedKey, req.Month, req.Year)

	result, found := cache.GetFromCache[model.APIResponseTransactionMonthAmountFailed](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsCache) SetCachedMonthAmountFailedCached(ctx context.Context, req *model.FindMonthlyTransactionStatusInput, res *model.APIResponseTransactionMonthAmountFailed) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthAmountFailedKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearAmountFailedCached(ctx context.Context, year int) (*model.APIResponseTransactionYearAmountFailed, bool) {
	key := fmt.Sprintf(transactionYearAmountFailedKey, year)

	result, found := cache.GetFromCache[model.APIResponseTransactionYearAmountFailed](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsCache) SetCachedYearAmountFailedCached(ctx context.Context, year int, res *model.APIResponseTransactionYearAmountFailed) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearAmountFailedKey, year)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthMethodSuccessCached(ctx context.Context, req *model.MonthTransactionMethodInput) (*model.APIResponseTransactionMonthPaymentMethod, bool) {
	key := fmt.Sprintf(transactionMonthMethodSuccessKey, req.Month, req.Year)

	result, found := cache.GetFromCache[model.APIResponseTransactionMonthPaymentMethod](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsCache) SetCachedMonthMethodSuccessCached(ctx context.Context, req *model.MonthTransactionMethodInput, res *model.APIResponseTransactionMonthPaymentMethod) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthMethodSuccessKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearMethodSuccessCached(ctx context.Context, year int) (*model.APIResponseTransactionYearPaymentMethod, bool) {
	key := fmt.Sprintf(transactionYearMethodSuccessKey, year)

	result, found := cache.GetFromCache[model.APIResponseTransactionYearPaymentMethod](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsCache) SetCachedYearMethodSuccessCached(ctx context.Context, year int, res *model.APIResponseTransactionYearPaymentMethod) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearMethodSuccessKey, year)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthMethodFailedCached(ctx context.Context, req *model.MonthTransactionMethodInput) (*model.APIResponseTransactionMonthPaymentMethod, bool) {
	key := fmt.Sprintf(transactionMonthMethodFailedKey, req.Month, req.Year)

	result, found := cache.GetFromCache[model.APIResponseTransactionMonthPaymentMethod](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsCache) SetCachedMonthMethodFailedCached(ctx context.Context, req *model.MonthTransactionMethodInput, res *model.APIResponseTransactionMonthPaymentMethod) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthMethodFailedKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearMethodFailedCached(ctx context.Context, year int) (*model.APIResponseTransactionYearPaymentMethod, bool) {
	key := fmt.Sprintf(transactionYearMethodFailedKey, year)

	result, found := cache.GetFromCache[model.APIResponseTransactionYearPaymentMethod](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsCache) SetCachedYearMethodFailedCached(ctx context.Context, year int, res *model.APIResponseTransactionYearPaymentMethod) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearMethodFailedKey, year)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}
