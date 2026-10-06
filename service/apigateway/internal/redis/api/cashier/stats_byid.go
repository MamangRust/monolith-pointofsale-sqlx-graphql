package cashier_cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-shared/cache"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
)

const (
	cashierStatsMonthTotalSalesByIdCacheKey = "cashier:stats:month:%d:year:%d:id:%d"
	cashierStatsYearTotalSalesByIdCacheKey  = "cashier:stats:year:%d:id:%d"

	cashierStatsMonthSalesByIdCacheKey = "cashier:stats:month:%d:id:%d"
	cashierStatsYearSalesByIdCacheKey  = "cashier:stats:year:%d:id:%d"
)

type cashierStatsByIdCache struct {
	store *cache.CacheStore
}

func NewCashierStatsByIdCache(store *cache.CacheStore) *cashierStatsByIdCache {
	return &cashierStatsByIdCache{store: store}
}

func (s *cashierStatsByIdCache) GetMonthlyTotalSalesByIdCache(ctx context.Context, req *model.FindYearMonthTotalSalesByID) (*model.APIResponseCashierMonthlyTotalSales, bool) {
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByIdCacheKey, req.Month, req.Year, req.CashierID)
	result, found := cache.GetFromCache[model.APIResponseCashierMonthlyTotalSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *cashierStatsByIdCache) SetMonthlyTotalSalesByIdCache(ctx context.Context, req *model.FindYearMonthTotalSalesByID, res *model.APIResponseCashierMonthlyTotalSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByIdCacheKey, req.Month, req.Year, req.CashierID)

	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *cashierStatsByIdCache) GetYearlyTotalSalesByIdCache(ctx context.Context, req *model.FindYearTotalSalesByID) (*model.APIResponseCashierYearlyTotalSales, bool) {
	key := fmt.Sprintf(cashierStatsYearTotalSalesByIdCacheKey, req.Year, req.CashierID)
	result, found := cache.GetFromCache[model.APIResponseCashierYearlyTotalSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *cashierStatsByIdCache) SetYearlyTotalSalesByIdCache(ctx context.Context, req *model.FindYearTotalSalesByID, res *model.APIResponseCashierYearlyTotalSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearTotalSalesByIdCacheKey, req.Year, req.CashierID)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *cashierStatsByIdCache) GetMonthlyCashierByIdCache(ctx context.Context, req *model.FindYearCashierByID) (*model.APIResponseCashierMonthSales, bool) {
	key := fmt.Sprintf(cashierStatsMonthSalesByIdCacheKey, req.Year, req.CashierID)
	result, found := cache.GetFromCache[model.APIResponseCashierMonthSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *cashierStatsByIdCache) SetMonthlyCashierByIdCache(ctx context.Context, req *model.FindYearCashierByID, res *model.APIResponseCashierMonthSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthSalesByIdCacheKey, req.Year, req.CashierID)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *cashierStatsByIdCache) GetYearlyCashierByIdCache(ctx context.Context, req *model.FindYearCashierByID) (*model.APIResponseCashierYearSales, bool) {
	key := fmt.Sprintf(cashierStatsYearSalesByIdCacheKey, req.Year, req.CashierID)
	result, found := cache.GetFromCache[model.APIResponseCashierYearSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *cashierStatsByIdCache) SetYearlyCashierByIdCache(ctx context.Context, req *model.FindYearCashierByID, res *model.APIResponseCashierYearSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearSalesByIdCacheKey, req.Year, req.CashierID)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}
