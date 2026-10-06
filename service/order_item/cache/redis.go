package mencache

import (
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

type Mencache interface {
	OrderItemQueryCache
	OrderItemCommandCache
}

type mencache struct {
	OrderItemQueryCache
	OrderItemCommandCache
}

func NewMencache(cacheStore *cache.CacheStore) Mencache {
	return &mencache{
		OrderItemQueryCache:   NewOrderItemQueryCache(cacheStore),
		OrderItemCommandCache: NewOrderItemCommandCache(cacheStore),
	}
}
