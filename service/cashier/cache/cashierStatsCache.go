package mencache

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
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

func NewCashierStatsCache(store *cache.CacheStore) CashierStatsCache {
	return &cashierStatsCache{store: store}
}

func (s *cashierStatsCache) GetMonthlyTotalSalesCache(ctx context.Context, req *requests.MonthTotalSales) ([]*db.GetMonthlyTotalSalesCashierRow, bool) {
	key := fmt.Sprintf(cashierStatsMonthTotalSalesCacheKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyTotalSalesCashierRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsCache) SetMonthlyTotalSalesCache(ctx context.Context, req *requests.MonthTotalSales, res []*db.GetMonthlyTotalSalesCashierRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthTotalSalesCacheKey, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsCache) GetYearlyTotalSalesCache(ctx context.Context, year int) ([]*db.GetYearlyTotalSalesCashierRow, bool) {
	key := fmt.Sprintf(cashierStatsYearTotalSalesCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTotalSalesCashierRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsCache) SetYearlyTotalSalesCache(ctx context.Context, year int, res []*db.GetYearlyTotalSalesCashierRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearTotalSalesCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsCache) GetMonthlySalesCache(ctx context.Context, year int) ([]*db.GetMonthlyCashierRow, bool) {
	key := fmt.Sprintf(cashierStatsMonthSalesCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyCashierRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsCache) SetMonthlySalesCache(ctx context.Context, year int, res []*db.GetMonthlyCashierRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthSalesCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsCache) GetYearlySalesCache(ctx context.Context, year int) ([]*db.GetYearlyCashierRow, bool) {
	key := fmt.Sprintf(cashierStatsYearSalesCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyCashierRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsCache) SetYearlySalesCache(ctx context.Context, year int, res []*db.GetYearlyCashierRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearSalesCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}
