package cashier_cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-shared/cache"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
)

const (
	cashierStatsMonthTotalSalesCacheKey = "cashier:stats:month:%d:year:%d"
	cashierStatsYearTotalSalesCacheKey  = "cashier:stats:year:%d"

	cashierStatsMonthSalesCacheKey = "cashier:stats:month:%d"
	cashierStatsYearSalesCacheKey  = "cashier:stats:year:%d"
)

type cashierStatsCache struct {
	store *cache.CacheStore
}

func NewCashierStatsCache(store *cache.CacheStore) *cashierStatsCache {
	return &cashierStatsCache{store: store}
}

func (s *cashierStatsCache) GetMonthlyTotalSalesCache(ctx context.Context, req *model.FindYearMonthTotalSales) (*model.APIResponseCashierMonthlyTotalSales, bool) {
	key := fmt.Sprintf(cashierStatsMonthTotalSalesCacheKey, req.Month, req.Year)
	result, found := cache.GetFromCache[model.APIResponseCashierMonthlyTotalSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *cashierStatsCache) SetMonthlyTotalSalesCache(ctx context.Context, req *model.FindYearMonthTotalSales, res *model.APIResponseCashierMonthlyTotalSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthTotalSalesCacheKey, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *cashierStatsCache) GetYearlyTotalSalesCache(ctx context.Context, year int) (*model.APIResponseCashierYearlyTotalSales, bool) {
	key := fmt.Sprintf(cashierStatsYearTotalSalesCacheKey, year)
	result, found := cache.GetFromCache[model.APIResponseCashierYearlyTotalSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *cashierStatsCache) SetYearlyTotalSalesCache(ctx context.Context, year int, res *model.APIResponseCashierYearlyTotalSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearTotalSalesCacheKey, year)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *cashierStatsCache) GetMonthlySalesCache(ctx context.Context, year int) (*model.APIResponseCashierMonthSales, bool) {
	key := fmt.Sprintf(cashierStatsMonthSalesCacheKey, year)
	result, found := cache.GetFromCache[model.APIResponseCashierMonthSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *cashierStatsCache) SetMonthlySalesCache(ctx context.Context, year int, res *model.APIResponseCashierMonthSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthSalesCacheKey, year)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *cashierStatsCache) GetYearlySalesCache(ctx context.Context, year int) (*model.APIResponseCashierYearSales, bool) {
	key := fmt.Sprintf(cashierStatsYearSalesCacheKey, year)
	result, found := cache.GetFromCache[model.APIResponseCashierYearSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *cashierStatsCache) SetYearlySalesCache(ctx context.Context, year int, res *model.APIResponseCashierYearSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearSalesCacheKey, year)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}
