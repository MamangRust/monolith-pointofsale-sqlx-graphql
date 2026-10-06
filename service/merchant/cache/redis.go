package mencache

import (
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

type Mencache interface {
	MerchantQueryCache
	MerchantCommandCache
	MerchantDocumentQueryCache
	MerchantDocumentCommandCache
}

type mencache struct {
	MerchantQueryCache
	MerchantCommandCache
	MerchantDocumentQueryCache
	MerchantDocumentCommandCache
}

func NewMencache(cacheStore *cache.CacheStore) Mencache {
	return &mencache{
		MerchantQueryCache:           NewMerchantQueryCache(cacheStore),
		MerchantCommandCache:         NewMerchantCommandCache(cacheStore),
		MerchantDocumentQueryCache:   NewMerchantDocumentQueryCache(cacheStore),
		MerchantDocumentCommandCache: NewMerchantDocumentCommandCache(cacheStore),
	}
}
