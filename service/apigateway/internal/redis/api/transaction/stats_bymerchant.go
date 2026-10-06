package transaction_cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

const (
	transactonMonthAmountSuccessByMerchantKey = "transaction:month:amount:success:merchant:%d:month:%d:year:%d"
	transactonMonthAmountFailedByMerchantKey  = "transaction:month:amount:failed:merchant:%d:month:%d:year:%d"

	transactonYearAmountSuccessByMerchantKey = "transaction:year:amount:success:merchant:%d:year:%d"
	transactonYearAmountFailedByMerchantKey  = "transaction:year:amount:failed:merchant:%d:year:%d"

	transactonMonthMethodSuccessByMerchantKey = "transaction:month:method:success:merchant:%d:month:%d:year:%d"
	transactonMonthMethodFailedByMerchantKey  = "transaction:month:method:failed:merchant:%d:month:%d:year:%d"

	transactonYearMethodSuccessByMerchantKey = "transaction:year:method:success:merchant:%d:year:%d"
	transactonYearMethodFailedByMerchantKey  = "transaction:year:method:failed:merchant:%d:year:%d"
)

type transactionStatsByMerchantCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsByMerchantCache(store *cache.CacheStore) *transactionStatsByMerchantCache {
	return &transactionStatsByMerchantCache{store: store}
}

func (t *transactionStatsByMerchantCache) GetCachedMonthAmountSuccessByMerchantCached(ctx context.Context, req *model.FindMonthlyTransactionStatusByMerchantInput) (*model.APIResponseTransactionMonthAmountSuccess, bool) {
	key := fmt.Sprintf(transactonMonthAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := cache.GetFromCache[model.APIResponseTransactionMonthAmountSuccess](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthAmountSuccessByMerchantCached(ctx context.Context, req *model.FindMonthlyTransactionStatusByMerchantInput, res *model.APIResponseTransactionMonthAmountSuccess) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)

	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthAmountFailedByMerchantCached(ctx context.Context, req *model.FindMonthlyTransactionStatusByMerchantInput) (*model.APIResponseTransactionMonthAmountFailed, bool) {
	key := fmt.Sprintf(transactonMonthAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := cache.GetFromCache[model.APIResponseTransactionMonthAmountFailed](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthAmountFailedByMerchantCached(ctx context.Context, req *model.FindMonthlyTransactionStatusByMerchantInput, res *model.APIResponseTransactionMonthAmountFailed) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearAmountFailedByMerchantCached(ctx context.Context, req *model.FindYearlyTransactionStatusByMerchantInput) (*model.APIResponseTransactionYearAmountFailed, bool) {
	key := fmt.Sprintf(transactonYearAmountFailedByMerchantKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[model.APIResponseTransactionYearAmountFailed](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearAmountFailedByMerchantCached(ctx context.Context, req *model.FindYearlyTransactionStatusByMerchantInput, res *model.APIResponseTransactionYearAmountFailed) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearAmountFailedByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearAmountSuccessByMerchantCached(ctx context.Context, req *model.FindYearlyTransactionStatusByMerchantInput) (*model.APIResponseTransactionYearAmountSuccess, bool) {
	key := fmt.Sprintf(transactonYearAmountSuccessByMerchantKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[model.APIResponseTransactionYearAmountSuccess](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearAmountSuccessByMerchantCached(ctx context.Context, req *model.FindYearlyTransactionStatusByMerchantInput, res *model.APIResponseTransactionYearAmountSuccess) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearAmountSuccessByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthMethodSuccessByMerchantCached(ctx context.Context, req *model.MonthTransactionMethodByMerchantInput) (*model.APIResponseTransactionMonthPaymentMethod, bool) {
	key := fmt.Sprintf(transactonMonthMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := cache.GetFromCache[model.APIResponseTransactionMonthPaymentMethod](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthMethodSuccessByMerchantCached(ctx context.Context, req *model.MonthTransactionMethodByMerchantInput, res *model.APIResponseTransactionMonthPaymentMethod) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearMethodSuccessByMerchantCached(ctx context.Context, req *model.YearTransactionMethodByMerchantInput) (*model.APIResponseTransactionYearPaymentMethod, bool) {
	key := fmt.Sprintf(transactonYearMethodSuccessByMerchantKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[model.APIResponseTransactionYearPaymentMethod](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearMethodSuccessByMerchantCached(ctx context.Context, req *model.YearTransactionMethodByMerchantInput, res *model.APIResponseTransactionYearPaymentMethod) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearMethodSuccessByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthMethodFailedByMerchantCached(ctx context.Context, req *model.MonthTransactionMethodByMerchantInput) (*model.APIResponseTransactionMonthPaymentMethod, bool) {
	key := fmt.Sprintf(transactonMonthMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := cache.GetFromCache[model.APIResponseTransactionMonthPaymentMethod](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthMethodFailedByMerchantCached(ctx context.Context, req *model.MonthTransactionMethodByMerchantInput, res *model.APIResponseTransactionMonthPaymentMethod) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearMethodFailedByMerchantCached(ctx context.Context, req *model.YearTransactionMethodByMerchantInput) (*model.APIResponseTransactionYearPaymentMethod, bool) {
	key := fmt.Sprintf(transactonYearMethodFailedByMerchantKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[model.APIResponseTransactionYearPaymentMethod](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearMethodFailedByMerchantCached(ctx context.Context, req *model.YearTransactionMethodByMerchantInput, res *model.APIResponseTransactionYearPaymentMethod) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearMethodFailedByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, res, ttlDefault)
}
