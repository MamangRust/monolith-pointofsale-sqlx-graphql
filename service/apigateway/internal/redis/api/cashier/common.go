package cashier_cache

import "time"

const (
	cashierAllCacheKey     = "cashier:all:page:%d:pageSize:%d:search:%s"
	cashierByIdCacheKey    = "cashier:id:%d"
	cashierActiveCacheKey  = "cashier:active:page:%d:pageSize:%d:search:%s"
	cashierTrashedCacheKey = "cashier:trashed:page:%d:pageSize:%d:search:%s"

	cashierByMerchantCacheKey = "cashier:merchant:%d:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)
