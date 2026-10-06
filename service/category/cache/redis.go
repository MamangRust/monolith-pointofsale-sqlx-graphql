package mencache

import (
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

type Mencache interface {
	CategoryQueryCache
	CategoryCommandCache
	CategoryStatsCache
	CategoryStatsByIdCache
	CategoryStatsByMerchantCache
}

type mencache struct {
	CategoryQueryCache
	CategoryCommandCache
	CategoryStatsCache
	CategoryStatsByIdCache
	CategoryStatsByMerchantCache
}

func NewMencache(cacheStore *cache.CacheStore) Mencache {
	return &mencache{
		CategoryQueryCache:           NewCategoryQueryCache(cacheStore),
		CategoryCommandCache:         NewCategoryCommandCache(cacheStore),
		CategoryStatsCache:           NewCategoryStatsCache(cacheStore),
		CategoryStatsByIdCache:       NewCategoryStatsByIdCache(cacheStore),
		CategoryStatsByMerchantCache: NewCategoryStatsByMerchantCache(cacheStore),
	}
}
