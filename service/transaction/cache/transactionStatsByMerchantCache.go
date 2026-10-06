package mencache

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
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

func NewTransactionStatsByMerchantCache(store *cache.CacheStore) TransactionStatsByMerchantCache {
	return &transactionStatsByMerchantCache{store: store}
}

func (t *transactionStatsByMerchantCache) GetCachedMonthAmountSuccessByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionSuccessByMerchantRow, bool) {
	key := fmt.Sprintf(transactonMonthAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyAmountTransactionSuccessByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthAmountSuccessByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*db.GetMonthlyAmountTransactionSuccessByMerchantRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthAmountFailedByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionFailedByMerchantRow, bool) {
	key := fmt.Sprintf(transactonMonthAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyAmountTransactionFailedByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthAmountFailedByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*db.GetMonthlyAmountTransactionFailedByMerchantRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearAmountFailedByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionFailedByMerchantRow, bool) {
	key := fmt.Sprintf(transactonYearAmountFailedByMerchantKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyAmountTransactionFailedByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearAmountFailedByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*db.GetYearlyAmountTransactionFailedByMerchantRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearAmountFailedByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearAmountSuccessByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionSuccessByMerchantRow, bool) {
	key := fmt.Sprintf(transactonYearAmountSuccessByMerchantKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyAmountTransactionSuccessByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearAmountSuccessByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*db.GetYearlyAmountTransactionSuccessByMerchantRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearAmountSuccessByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthMethodSuccessByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow, bool) {
	key := fmt.Sprintf(transactonMonthMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthMethodSuccessByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*db.GetMonthlyTransactionMethodsByMerchantSuccessRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearMethodSuccessByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantSuccessRow, bool) {
	key := fmt.Sprintf(transactonYearMethodSuccessByMerchantKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionMethodsByMerchantSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearMethodSuccessByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*db.GetYearlyTransactionMethodsByMerchantSuccessRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearMethodSuccessByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthMethodFailedByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantFailedRow, bool) {
	key := fmt.Sprintf(transactonMonthMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTransactionMethodsByMerchantFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthMethodFailedByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*db.GetMonthlyTransactionMethodsByMerchantFailedRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearMethodFailedByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantFailedRow, bool) {
	key := fmt.Sprintf(transactonYearMethodFailedByMerchantKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionMethodsByMerchantFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearMethodFailedByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*db.GetYearlyTransactionMethodsByMerchantFailedRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearMethodFailedByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}
