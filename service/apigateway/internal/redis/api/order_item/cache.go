package orderitem_cache

import "github.com/MamangRust/monolith-pointofsale-shared/cache"

type OrderItemCache interface {
	OrderItemQueryCache
}

type orderItemCache struct {
	OrderItemQueryCache
}

func NewOrderItemCache(store *cache.CacheStore) OrderItemCache {
	return &orderItemCache{
		OrderItemQueryCache: NewOrderItemQueryCache(store),
	}
}
