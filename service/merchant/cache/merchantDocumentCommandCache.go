package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-pointofsale-shared/cache"
)

type merchantDocumentCommandCache struct {
	store *cache.CacheStore
}

func NewMerchantDocumentCommandCache(store *cache.CacheStore) MerchantDocumentCommandCache {
	return &merchantDocumentCommandCache{store: store}
}

func (s *merchantDocumentCommandCache) DeleteCachedMerchantDocuments(ctx context.Context, id int) {
	key := fmt.Sprintf(merchantDocumentByIdCacheKey, id)
	cache.DeleteFromCache(ctx, s.store, key)
}

func (s *merchantDocumentCommandCache) DeleteCachedMerchantDocumentsAllCache(ctx context.Context) {
	s.store.InvalidateCache(ctx, "merchant_document:*")
}
