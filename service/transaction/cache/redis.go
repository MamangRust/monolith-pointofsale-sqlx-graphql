package mencache

import (
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

type Mencache interface {
	TransactionQueryCache
	TransactionCommandCache
	TransactionStatsCache
	TransactionStatsByMerchantCache
}

type mencache struct {
	TransactionQueryCache
	TransactionCommandCache
	TransactionStatsCache
	TransactionStatsByMerchantCache
}

func NewMencache(cacheStore *cache.CacheStore) Mencache {
	return &mencache{
		TransactionQueryCache:           NewTransactionQueryCache(cacheStore),
		TransactionCommandCache:         NewTransactionCommandCache(cacheStore),
		TransactionStatsCache:           NewTransactionStatsCache(cacheStore),
		TransactionStatsByMerchantCache: NewTransactionStatsByMerchantCache(cacheStore),
	}
}
