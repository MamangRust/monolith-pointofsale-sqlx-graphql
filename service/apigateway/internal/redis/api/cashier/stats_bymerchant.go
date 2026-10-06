package cashier_cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-shared/cache"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
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

func NewCashierStatsByMerchantCache(store *cache.CacheStore) *cashierStatsByMerchantCache {
	return &cashierStatsByMerchantCache{store: store}
}

func (s *cashierStatsByMerchantCache) GetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *model.FindYearMonthTotalSalesByMerchant) (*model.APIResponseCashierMonthlyTotalSales, bool) {
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByMerchantCacheKey, req.Month, req.Year, req.MerchantID)
	result, found := cache.GetFromCache[model.APIResponseCashierMonthlyTotalSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *cashierStatsByMerchantCache) SetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *model.FindYearMonthTotalSalesByMerchant, res *model.APIResponseCashierMonthlyTotalSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByMerchantCacheKey, req.Month, req.Year, req.MerchantID)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *cashierStatsByMerchantCache) GetYearlyTotalSalesByMerchantCache(ctx context.Context, req *model.FindYearTotalSalesByMerchant) (*model.APIResponseCashierYearlyTotalSales, bool) {
	key := fmt.Sprintf(cashierStatsYearTotalSalesByMerchantCacheKey, req.Year, req.MerchantID)
	result, found := cache.GetFromCache[model.APIResponseCashierYearlyTotalSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *cashierStatsByMerchantCache) SetYearlyTotalSalesByMerchantCache(ctx context.Context, req *model.FindYearTotalSalesByMerchant, res *model.APIResponseCashierYearlyTotalSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearTotalSalesByMerchantCacheKey, req.Year, req.MerchantID)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *cashierStatsByMerchantCache) GetMonthlyCashierByMerchantCache(ctx context.Context, req *model.FindYearCashierByMerchant) (*model.APIResponseCashierMonthSales, bool) {
	key := fmt.Sprintf(cashierStatsMonthSalesByMerchantCacheKey, req.Year, req.MerchantID)
	result, found := cache.GetFromCache[model.APIResponseCashierMonthSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *cashierStatsByMerchantCache) SetMonthlyCashierByMerchantCache(ctx context.Context, req *model.FindYearCashierByMerchant, res *model.APIResponseCashierMonthSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthSalesByMerchantCacheKey, req.Year, req.MerchantID)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}

func (s *cashierStatsByMerchantCache) GetYearlyCashierByMerchantCache(ctx context.Context, req *model.FindYearCashierByMerchant) (*model.APIResponseCashierYearSales, bool) {
	key := fmt.Sprintf(cashierStatsYearSalesByMerchantCacheKey, req.Year, req.MerchantID)
	result, found := cache.GetFromCache[model.APIResponseCashierYearSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *cashierStatsByMerchantCache) SetYearlyCashierByMerchantCache(ctx context.Context, req *model.FindYearCashierByMerchant, res *model.APIResponseCashierYearSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearSalesByMerchantCacheKey, req.Year, req.MerchantID)
	cache.SetToCache(ctx, s.store, key, res, ttlDefault)
}
