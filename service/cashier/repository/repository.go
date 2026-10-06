package repository

import (
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	merchantadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/merchant"
	useradapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/user"
	"github.com/jmoiron/sqlx"
)

type GuardOptions struct {
	User     []adapter.GuardOption
	Merchant []adapter.GuardOption
}

type Repositories struct {
	UserQuery              useradapter.QueryRepository
	MerchantQuery          merchantadapter.Repository
	CashierQuery           CashierQueryRepository
	CashierCommand         CashierCommandRepository
	CashierStats           CashierStatsRepository
	CashierStatsByMerchant CashierStatByMerchantRepository
	CashierStatsById       CashierStatByIdRepository
}

func NewRepositories(
	DB *sqlx.DB,
	userClient pbuser.UserQueryServiceClient,
	merchantClient pbmerchant.MerchantQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		UserQuery:              useradapter.NewQueryAdapter(userClient, g.User...),
		MerchantQuery:          merchantadapter.New(merchantClient, g.Merchant...),
		CashierQuery:           NewCashierQueryRepository(DB),
		CashierCommand:         NewCashierCommandRepository(DB),
		CashierStats:           NewCashierStatsRepository(DB),
		CashierStatsByMerchant: NewCashierStatsByMerchantRepository(DB),
		CashierStatsById:       NewCashierStatsByIdRepository(DB),
	}
}
