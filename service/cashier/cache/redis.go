package mencache

import (
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

type Mencache interface {
	CashierQueryCache
	CashierCommandCache
	CashierStatsCache
	CashierStatsByIdCache
	CashierStatsByMerchantCache
}

type mencache struct {
	CashierQueryCache
	CashierCommandCache
	CashierStatsCache
	CashierStatsByIdCache
	CashierStatsByMerchantCache
}

func NewMencache(cacheStore *cache.CacheStore) Mencache {
	return &mencache{
		CashierQueryCache:           NewCashierQueryCache(cacheStore),
		CashierCommandCache:         NewCashierCommandCache(cacheStore),
		CashierStatsCache:           NewCashierStatsCache(cacheStore),
		CashierStatsByIdCache:       NewCashierStatsByIdCache(cacheStore),
		CashierStatsByMerchantCache: NewCashierStatsByMerchantCache(cacheStore),
	}
}
