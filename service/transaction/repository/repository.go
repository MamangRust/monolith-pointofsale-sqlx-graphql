package repository

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb"
	pbcashier "github.com/MamangRust/monolith-pointofsale-pb/cashier"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	cashieradapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/cashier"
	merchantadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/merchant"
	orderadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/order"
	orderitemadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/orderitem"
	"github.com/jmoiron/sqlx"
)

type GuardOptions struct {
	Cashier   []adapter.GuardOption
	Merchant  []adapter.GuardOption
	Order     []adapter.GuardOption
	OrderItem []adapter.GuardOption
}

type Repositories struct {
	CashierQuery                 cashieradapter.QueryRepository
	MerchantQuery                merchantadapter.Repository
	OrderItemQuery               orderitemadapter.QueryRepository
	OrderQuery                   orderadapter.QueryRepository
	TransactionCommandRepository TransactionCommandRepository
	TransactionQueryRepository   TransactionQueryRepository
	TransactionStatsRepository   TransactionStatsRepository
	TransactionStatsByMerchant   TransactionStatsByMerchantRepository
}

func NewRepositories(
	DB *sqlx.DB,
	cashierClient pbcashier.CashierQueryServiceClient,
	merchantClient pbmerchant.MerchantQueryServiceClient,
	orderClient pborder.OrderQueryServiceClient,
	orderItemClient pb.OrderItemQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		CashierQuery:                 cashieradapter.New(cashierClient, g.Cashier...),
		MerchantQuery:                merchantadapter.New(merchantClient, g.Merchant...),
		OrderItemQuery:               orderitemadapter.NewQueryAdapter(orderItemClient, g.OrderItem...),
		OrderQuery:                   orderadapter.New(orderClient, g.Order...),
		TransactionCommandRepository: NewTransactionCommandRepository(DB),
		TransactionQueryRepository:   NewTransactionQueryRepository(DB),
		TransactionStatsRepository:   NewTransactionStatsRepository(DB),
		TransactionStatsByMerchant:   NewTransactionStatsByMerchantRepository(DB),
	}
}
