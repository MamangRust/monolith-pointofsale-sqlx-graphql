package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

type orderItemCommandCache struct {
	store *cache.CacheStore
}

func NewOrderItemCommandCache(store *cache.CacheStore) OrderItemCommandCache {
	return &orderItemCommandCache{store: store}
}

func (c *orderItemCommandCache) DeleteCachedOrderItems(ctx context.Context, orderID int) {
	cache.DeleteFromCache(ctx, c.store, fmt.Sprintf(orderItemByIdCacheKey, orderID))
}

func (c *orderItemCommandCache) DeleteCachedOrderItemsAll(ctx context.Context) {
	c.store.InvalidateCache(ctx, "order_item:*")
}
