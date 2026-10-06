package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

type productCommandCache struct {
	store *cache.CacheStore
}

func NewProductCommandCache(store *cache.CacheStore) ProductCommandCache {
	return &productCommandCache{store: store}
}

func (c *productCommandCache) DeleteCachedProduct(ctx context.Context, productID int) {
	cache.DeleteFromCache(ctx, c.store, fmt.Sprintf(productByIdCacheKey, productID))
}

func (c *productCommandCache) DeleteCachedProductAllCache(ctx context.Context) {
	c.store.InvalidateCache(ctx, "product:*")
}
