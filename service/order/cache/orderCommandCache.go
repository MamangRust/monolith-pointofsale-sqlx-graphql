package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

type orderCommandCache struct {
	store *cache.CacheStore
}

func NewOrderCommandCache(store *cache.CacheStore) OrderCommandCache {
	return &orderCommandCache{store: store}
}

func (s *orderCommandCache) DeleteOrderCache(ctx context.Context, orderID int) {
	cache.DeleteFromCache(ctx, s.store, fmt.Sprintf(orderByIdCacheKey, orderID))
}

func (s *orderCommandCache) DeleteOrderAllCache(ctx context.Context) {
	s.store.InvalidateCache(ctx, "order:*")
}
