package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

type categoryCommandCache struct {
	store *cache.CacheStore
}

func NewCategoryCommandCache(store *cache.CacheStore) CategoryCommandCache {
	return &categoryCommandCache{store: store}
}

func (c *categoryCommandCache) DeleteCachedCategoryCache(ctx context.Context, id int) {
	key := fmt.Sprintf(categoryByIdCacheKey, id)
	cache.DeleteFromCache(ctx, c.store, key)
}

func (c *categoryCommandCache) DeleteCachedCategoryAllCache(ctx context.Context) {
	c.store.InvalidateCache(ctx, "category:*")
}
