package transaction_cache

import "github.com/MamangRust/monolith-pointofsale-shared/cache"

type TransactionMencache interface {
	TransactionQueryCache
	TransactionCommandCache
	TransactionStatsCache
	TransactionStatsByMerchantCache
}

type transactionMencache struct {
	TransactionQueryCache
	TransactionCommandCache
	TransactionStatsCache
	TransactionStatsByMerchantCache
}

func NewTransactionMencache(cacheStore *cache.CacheStore) TransactionMencache {
	return &transactionMencache{
		TransactionQueryCache:           NewTransactionQueryCache(cacheStore),
		TransactionCommandCache:         NewTransactionCommandCache(cacheStore),
		TransactionStatsCache:           NewTransactionStatsCache(cacheStore),
		TransactionStatsByMerchantCache: NewTransactionStatsByMerchantCache(cacheStore),
	}
}
