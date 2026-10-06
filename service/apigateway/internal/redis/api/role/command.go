package role_cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

type roleCommandCache struct {
	store *cache.CacheStore
}

func NewRoleCommandCache(store *cache.CacheStore) *roleCommandCache {
	return &roleCommandCache{store: store}
}

func (s *roleCommandCache) DeleteCachedRole(ctx context.Context, id int) {
	key := fmt.Sprintf(roleByIdCacheKey, id)

	cache.DeleteFromCache(ctx, s.store, key)
}
