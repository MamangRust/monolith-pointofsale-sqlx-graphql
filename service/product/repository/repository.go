package repository

import (
	pbcategory "github.com/MamangRust/monolith-pointofsale-pb/category"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	categoryadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/category"
	merchantadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/merchant"
	"github.com/jmoiron/sqlx"
)

type GuardOptions struct {
	Category []adapter.GuardOption
	Merchant []adapter.GuardOption
}

type Repositories struct {
	ProductQuery   ProductQueryRepository
	ProductCommand ProductCommandRepository
	CategoryQuery  categoryadapter.Repository
	MerchantQuery  merchantadapter.Repository
}

func NewRepositories(
	DB *sqlx.DB,
	categoryQueryClient pbcategory.CategoryQueryServiceClient,
	merchantQueryClient pbmerchant.MerchantQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		ProductQuery:   NewProductQueryRepository(DB),
		ProductCommand: NewProductCommandRepository(DB),
		CategoryQuery:  categoryadapter.New(categoryQueryClient, g.Category...),
		MerchantQuery:  merchantadapter.New(merchantQueryClient, g.Merchant...),
	}
}
