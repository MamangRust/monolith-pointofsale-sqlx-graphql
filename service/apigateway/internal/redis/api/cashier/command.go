package cashier_cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

type cashierCommandCache struct {
	store *cache.CacheStore
}

func NewCashierCommandCache(store *cache.CacheStore) *cashierCommandCache {
	return &cashierCommandCache{store: store}
}

func (c *cashierCommandCache) DeleteCashierCache(ctx context.Context, id int) {
	key := fmt.Sprintf(cashierByIdCacheKey, id)

	cache.DeleteFromCache(ctx, c.store, key)
}
