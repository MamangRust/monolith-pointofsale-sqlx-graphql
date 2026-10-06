package mencache

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
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

func NewCashierStatsByIdCache(store *cache.CacheStore) CashierStatsByIdCache {
	return &cashierStatsByIdCache{store: store}
}

func (s *cashierStatsByIdCache) GetMonthlyTotalSalesByIdCache(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*db.GetMonthlyTotalSalesByIdRow, bool) {
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByIdCacheKey, req.Month, req.Year, req.CashierID)
	result, found := cache.GetFromCache[[]*db.GetMonthlyTotalSalesByIdRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByIdCache) SetMonthlyTotalSalesByIdCache(ctx context.Context, req *requests.MonthTotalSalesCashier, res []*db.GetMonthlyTotalSalesByIdRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByIdCacheKey, req.Month, req.Year, req.CashierID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByIdCache) GetYearlyTotalSalesByIdCache(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*db.GetYearlyTotalSalesByIdRow, bool) {
	key := fmt.Sprintf(cashierStatsYearTotalSalesByIdCacheKey, req.Year, req.CashierID)
	result, found := cache.GetFromCache[[]*db.GetYearlyTotalSalesByIdRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByIdCache) SetYearlyTotalSalesByIdCache(ctx context.Context, req *requests.YearTotalSalesCashier, res []*db.GetYearlyTotalSalesByIdRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearTotalSalesByIdCacheKey, req.Year, req.CashierID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByIdCache) GetMonthlyCashierByIdCache(ctx context.Context, req *requests.MonthCashierId) ([]*db.GetMonthlyCashierByCashierIdRow, bool) {
	key := fmt.Sprintf(cashierStatsMonthSalesByIdCacheKey, req.Year, req.CashierID)
	result, found := cache.GetFromCache[[]*db.GetMonthlyCashierByCashierIdRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByIdCache) SetMonthlyCashierByIdCache(ctx context.Context, req *requests.MonthCashierId, res []*db.GetMonthlyCashierByCashierIdRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthSalesByIdCacheKey, req.Year, req.CashierID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByIdCache) GetYearlyCashierByIdCache(ctx context.Context, req *requests.YearCashierId) ([]*db.GetYearlyCashierByCashierIdRow, bool) {
	key := fmt.Sprintf(cashierStatsYearSalesByIdCacheKey, req.Year, req.CashierID)
	result, found := cache.GetFromCache[[]*db.GetYearlyCashierByCashierIdRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByIdCache) SetYearlyCashierByIdCache(ctx context.Context, req *requests.YearCashierId, res []*db.GetYearlyCashierByCashierIdRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearSalesByIdCacheKey, req.Year, req.CashierID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}
