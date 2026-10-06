package product_cache

import "github.com/MamangRust/monolith-pointofsale-shared/cache"

type productMencache struct {
	ProductQueryCache
	ProductCommandCache
}

type ProductMencache interface {
	ProductQueryCache
	ProductCommandCache
}

func NewProductMencache(store *cache.CacheStore) ProductMencache {
	return &productMencache{
		ProductQueryCache:   NewProductQueryCache(store),
		ProductCommandCache: NewProductCommandCache(store),
	}
}
