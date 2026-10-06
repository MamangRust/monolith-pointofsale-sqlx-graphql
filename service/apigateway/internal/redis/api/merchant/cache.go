package merchant_cache

import "github.com/MamangRust/monolith-pointofsale-shared/cache"

type MerchantMenCache interface {
	MerchantQueryCache
	MerchantCommandCache
}

type merchantMencache struct {
	MerchantQueryCache
	MerchantCommandCache
}

func NewMerchantMencache(store *cache.CacheStore) MerchantMenCache {
	return &merchantMencache{
		MerchantQueryCache:   NewMerchantQueryCache(store),
		MerchantCommandCache: NewMerchantCommandCache(store),
	}
}
