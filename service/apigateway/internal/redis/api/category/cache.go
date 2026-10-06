package category_cache

import "github.com/MamangRust/monolith-pointofsale-shared/cache"

type CategoryMencache interface {
	CategoryQueryCache
	CategoryCommandCache
	CategoryStatsByIdCache
	CategoryStatsCache
	CategoryStatsByMerchantCache
}

type categoryMencache struct {
	CategoryQueryCache
	CategoryCommandCache
	CategoryStatsCache
	CategoryStatsByIdCache
	CategoryStatsByMerchantCache
}

func NewCategoryMencache(store *cache.CacheStore) CategoryMencache {
	return &categoryMencache{
		CategoryQueryCache:           NewCategoryQueryCache(store),
		CategoryCommandCache:         NewCategoryCommandCache(store),
		CategoryStatsCache:           NewCategoryStatsCache(store),
		CategoryStatsByIdCache:       NewCategoryStatsByIdCache(store),
		CategoryStatsByMerchantCache: NewCategoryStatsByMerchantCache(store),
	}
}
