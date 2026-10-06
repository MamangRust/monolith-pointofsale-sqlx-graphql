package mencache

import (
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

type Mencache interface {
	IdentityCache
	LoginCache
	PasswordResetCache
	RegisterCache
}

type mencache struct {
	IdentityCache
	LoginCache
	PasswordResetCache
	RegisterCache
}

func NewMencache(cacheStore *cache.CacheStore) Mencache {
	return &mencache{
		IdentityCache:      NewidentityCache(cacheStore),
		LoginCache:         NewLoginCache(cacheStore),
		PasswordResetCache: NewPasswordResetCache(cacheStore),
		RegisterCache:      NewRegisterCache(cacheStore),
	}
}
