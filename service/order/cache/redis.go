package mencache

import (
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

type Mencache interface {
	OrderQueryCache
	OrderCommandCache
	OrderStatsCache
	OrderStatsByMerchantCache
}

type mencache struct {
	OrderQueryCache
	OrderCommandCache
	OrderStatsCache
	OrderStatsByMerchantCache
}

func NewMencache(cacheStore *cache.CacheStore) Mencache {
	return &mencache{
		OrderQueryCache:           NewOrderQueryCache(cacheStore),
		OrderCommandCache:         NewOrderCommandCache(cacheStore),
		OrderStatsCache:           NewOrderStatsCache(cacheStore),
		OrderStatsByMerchantCache: NewOrderStatsByMerchantCache(cacheStore),
	}
}
