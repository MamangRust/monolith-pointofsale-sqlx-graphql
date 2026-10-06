package mencache

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
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

func NewTransactionStatsCache(store *cache.CacheStore) TransactionStatsCache {
	return &transactionStatsCache{store: store}
}

func (t *transactionStatsCache) GetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionSuccessRow, bool) {
	key := fmt.Sprintf(transactionMonthAmountSuccessKey, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyAmountTransactionSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction, res []*db.GetMonthlyAmountTransactionSuccessRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthAmountSuccessKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearAmountSuccessCached(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionSuccessRow, bool) {
	key := fmt.Sprintf(transactionYearAmountSuccessKey, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyAmountTransactionSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearAmountSuccessCached(ctx context.Context, year int, res []*db.GetYearlyAmountTransactionSuccessRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearAmountSuccessKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionFailedRow, bool) {
	key := fmt.Sprintf(transactionMonthAmountFailedKey, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyAmountTransactionFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction, res []*db.GetMonthlyAmountTransactionFailedRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthAmountFailedKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearAmountFailedCached(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionFailedRow, bool) {
	key := fmt.Sprintf(transactionYearAmountFailedKey, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyAmountTransactionFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearAmountFailedCached(ctx context.Context, year int, res []*db.GetYearlyAmountTransactionFailedRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearAmountFailedKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsSuccessRow, bool) {
	key := fmt.Sprintf(transactionMonthMethodSuccessKey, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTransactionMethodsSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction, res []*db.GetMonthlyTransactionMethodsSuccessRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthMethodSuccessKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearMethodSuccessCached(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsSuccessRow, bool) {
	key := fmt.Sprintf(transactionYearMethodSuccessKey, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionMethodsSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearMethodSuccessCached(ctx context.Context, year int, res []*db.GetYearlyTransactionMethodsSuccessRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearMethodSuccessKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsFailedRow, bool) {
	key := fmt.Sprintf(transactionMonthMethodFailedKey, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTransactionMethodsFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction, res []*db.GetMonthlyTransactionMethodsFailedRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthMethodFailedKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearMethodFailedCached(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsFailedRow, bool) {
	key := fmt.Sprintf(transactionYearMethodFailedKey, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionMethodsFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearMethodFailedCached(ctx context.Context, year int, res []*db.GetYearlyTransactionMethodsFailedRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearMethodFailedKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}
