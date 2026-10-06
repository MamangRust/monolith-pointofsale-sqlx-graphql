package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	CategoryQuery           CategoryQueryRepository
	CategoryCommand         CategoryCommandRepository
	CategoryStats           CategoryStatsRepository
	CategoryStatsById       CategoryStatsByIdRepository
	CategoryStatsByMerchant CategoryStatsByMerchantRepository
}

func NewRepositories(DB *sqlx.DB) *Repositories {
	return &Repositories{
		CategoryQuery:           NewCategoryQueryRepository(DB),
		CategoryCommand:         NewCategoryCommandRepository(DB),
		CategoryStats:           NewCategoryStatsRepository(DB),
		CategoryStatsById:       NewCategoryStatsByIdRepository(DB),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantRepository(DB),
	}
}
