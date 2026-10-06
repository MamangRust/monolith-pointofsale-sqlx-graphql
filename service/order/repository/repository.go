package repository

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb"
	pbcashier "github.com/MamangRust/monolith-pointofsale-pb/cashier"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pbproduct "github.com/MamangRust/monolith-pointofsale-pb/product"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	cashieradapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/cashier"
	merchantadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/merchant"
	orderitemadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/orderitem"
	productadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/product"
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	CashierQuery         cashieradapter.QueryRepository
	MerchantQuery        merchantadapter.Repository
	ProductQuery         productadapter.QueryRepository
	ProductCommand       productadapter.CommandRepository
	OrderQuery           OrderQueryRepository
	OrderCommand         OrderCommandRepository
	OrderItemQuery       orderitemadapter.QueryRepository
	OrderItemCommand     orderitemadapter.CommandRepository
	OrderStats           OrderStatsRepository
	OrderStatsByMerchant OrderStatByMerchantRepository
}

type GuardOptions struct {
	Cashier   []adapter.GuardOption
	Merchant  []adapter.GuardOption
	Product   []adapter.GuardOption
	OrderItem []adapter.GuardOption
}

type Deps struct {
	Db *sqlx.DB

	CashierQueryClient     pbcashier.CashierQueryServiceClient
	MerchantQueryClient    pbmerchant.MerchantQueryServiceClient
	ProductQueryClient     pbproduct.ProductQueryServiceClient
	ProductCommandClient   pbproduct.ProductCommandServiceClient
	OrderItemQueryClient   pb.OrderItemQueryServiceClient
	OrderItemCommandClient pb.OrderItemCommandServiceClient

	Guards GuardOptions
}

func NewRepositories(deps *Deps) *Repositories {
	g := deps.Guards

	productAdapter := productadapter.New(deps.ProductQueryClient, deps.ProductCommandClient, g.Product...)
	orderItemAdapter := orderitemadapter.New(deps.OrderItemQueryClient, deps.OrderItemCommandClient, g.OrderItem...)

	return &Repositories{
		CashierQuery:         cashieradapter.New(deps.CashierQueryClient, g.Cashier...),
		MerchantQuery:        merchantadapter.New(deps.MerchantQueryClient, g.Merchant...),
		ProductQuery:         productAdapter,
		ProductCommand:       productAdapter,
		OrderQuery:           NewOrderQueryRepository(deps.Db),
		OrderCommand:         NewOrderCommandRepository(deps.Db),
		OrderItemQuery:       orderItemAdapter,
		OrderItemCommand:     orderItemAdapter,
		OrderStats:           NewOrderStatsRepository(deps.Db),
		OrderStatsByMerchant: NewOrderStatsByMerchantRepository(deps.Db),
	}
}
