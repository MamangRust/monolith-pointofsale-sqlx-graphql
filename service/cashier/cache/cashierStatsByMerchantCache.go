package mencache

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

const (
	cashierStatsMonthTotalSalesByMerchantCacheKey = "cashier:stats:month:%d:year:%d:id:%d"
	cashierStatsYearTotalSalesByMerchantCacheKey  = "cashier:stats:year:%d:merchant:%d"

	cashierStatsMonthSalesByMerchantCacheKey = "cashier:stats:month:%d:merchant:%d"
	cashierStatsYearSalesByMerchantCacheKey  = "cashier:stats:year:%d:merchant:%d"
)

type cashierStatsByMerchantCache struct {
	store *cache.CacheStore
}

func NewCashierStatsByMerchantCache(store *cache.CacheStore) CashierStatsByMerchantCache {
	return &cashierStatsByMerchantCache{store: store}
}

func (s *cashierStatsByMerchantCache) GetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*db.GetMonthlyTotalSalesByMerchantRow, bool) {
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByMerchantCacheKey, req.Month, req.Year, req.MerchantID)
	result, found := cache.GetFromCache[[]*db.GetMonthlyTotalSalesByMerchantRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByMerchantCache) SetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *requests.MonthTotalSalesMerchant, res []*db.GetMonthlyTotalSalesByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByMerchantCacheKey, req.Month, req.Year, req.MerchantID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByMerchantCache) GetYearlyTotalSalesByMerchantCache(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*db.GetYearlyTotalSalesByMerchantRow, bool) {
	key := fmt.Sprintf(cashierStatsYearTotalSalesByMerchantCacheKey, req.Year, req.MerchantID)
	result, found := cache.GetFromCache[[]*db.GetYearlyTotalSalesByMerchantRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByMerchantCache) SetYearlyTotalSalesByMerchantCache(ctx context.Context, req *requests.YearTotalSalesMerchant, res []*db.GetYearlyTotalSalesByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearTotalSalesByMerchantCacheKey, req.Year, req.MerchantID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByMerchantCache) GetMonthlyCashierByMerchantCache(ctx context.Context, req *requests.MonthCashierMerchant) ([]*db.GetMonthlyCashierByMerchantRow, bool) {
	key := fmt.Sprintf(cashierStatsMonthSalesByMerchantCacheKey, req.Year, req.MerchantID)
	result, found := cache.GetFromCache[[]*db.GetMonthlyCashierByMerchantRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByMerchantCache) SetMonthlyCashierByMerchantCache(ctx context.Context, req *requests.MonthCashierMerchant, res []*db.GetMonthlyCashierByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthSalesByMerchantCacheKey, req.Year, req.MerchantID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByMerchantCache) GetYearlyCashierByMerchantCache(ctx context.Context, req *requests.YearCashierMerchant) ([]*db.GetYearlyCashierByMerchantRow, bool) {
	key := fmt.Sprintf(cashierStatsYearSalesByMerchantCacheKey, req.Year, req.MerchantID)
	result, found := cache.GetFromCache[[]*db.GetYearlyCashierByMerchantRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByMerchantCache) SetYearlyCashierByMerchantCache(ctx context.Context, req *requests.YearCashierMerchant, res []*db.GetYearlyCashierByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearSalesByMerchantCacheKey, req.Year, req.MerchantID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}
