package mencache

import (
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

type Mencache interface {
	RoleQueryCache
	RoleCommandCache
}

type mencache struct {
	RoleQueryCache
	RoleCommandCache
}

func NewMencache(cacheStore *cache.CacheStore) Mencache {
	return &mencache{
		RoleQueryCache:   NewRoleQueryCache(cacheStore),
		RoleCommandCache: NewRoleCommandCache(cacheStore),
	}
}
