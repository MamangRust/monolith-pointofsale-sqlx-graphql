package tests

import (
	"google.golang.org/grpc"

	pb "github.com/MamangRust/monolith-pointofsale-pb"
	pbcashier "github.com/MamangRust/monolith-pointofsale-pb/cashier"
	pbcategory "github.com/MamangRust/monolith-pointofsale-pb/category"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pbmerchant_document "github.com/MamangRust/monolith-pointofsale-pb/merchant_document"
	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	pbproduct "github.com/MamangRust/monolith-pointofsale-pb/product"
	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
	pbtransaction "github.com/MamangRust/monolith-pointofsale-pb/transaction"
	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
)

// Wrapper clients: satu variabel bisa memanggil method query, command, dan stats.

type UserClient struct {
	pbuser.UserQueryServiceClient
	pbuser.UserCommandServiceClient
}

func NewUserClient(conn *grpc.ClientConn) *UserClient {
	return &UserClient{
		UserQueryServiceClient:   pbuser.NewUserQueryServiceClient(conn),
		UserCommandServiceClient: pbuser.NewUserCommandServiceClient(conn),
	}
}

type CategoryClient struct {
	pbcategory.CategoryQueryServiceClient
	pbcategory.CategoryCommandServiceClient
	pbcategory.CategoryStatsServiceClient
}

func NewCategoryClient(conn *grpc.ClientConn) *CategoryClient {
	return &CategoryClient{
		CategoryQueryServiceClient:   pbcategory.NewCategoryQueryServiceClient(conn),
		CategoryCommandServiceClient: pbcategory.NewCategoryCommandServiceClient(conn),
		CategoryStatsServiceClient:   pbcategory.NewCategoryStatsServiceClient(conn),
	}
}

type MerchantClient struct {
	pbmerchant.MerchantQueryServiceClient
	pbmerchant.MerchantCommandServiceClient
}

func NewMerchantClient(conn *grpc.ClientConn) *MerchantClient {
	return &MerchantClient{
		MerchantQueryServiceClient:   pbmerchant.NewMerchantQueryServiceClient(conn),
		MerchantCommandServiceClient: pbmerchant.NewMerchantCommandServiceClient(conn),
	}
}

type MerchantDocumentClient struct {
	pbmerchant_document.MerchantDocumentQueryServiceClient
	pbmerchant_document.MerchantDocumentCommandServiceClient
}

func NewMerchantDocumentClient(conn *grpc.ClientConn) *MerchantDocumentClient {
	return &MerchantDocumentClient{
		MerchantDocumentQueryServiceClient:   pbmerchant_document.NewMerchantDocumentQueryServiceClient(conn),
		MerchantDocumentCommandServiceClient: pbmerchant_document.NewMerchantDocumentCommandServiceClient(conn),
	}
}

type OrderClient struct {
	pborder.OrderQueryServiceClient
	pborder.OrderCommandServiceClient
	pborder.OrderStatsServiceClient
}

func NewOrderClient(conn *grpc.ClientConn) *OrderClient {
	return &OrderClient{
		OrderQueryServiceClient:   pborder.NewOrderQueryServiceClient(conn),
		OrderCommandServiceClient: pborder.NewOrderCommandServiceClient(conn),
		OrderStatsServiceClient:   pborder.NewOrderStatsServiceClient(conn),
	}
}

type ProductClient struct {
	pbproduct.ProductQueryServiceClient
	pbproduct.ProductCommandServiceClient
}

func NewProductClient(conn *grpc.ClientConn) *ProductClient {
	return &ProductClient{
		ProductQueryServiceClient:   pbproduct.NewProductQueryServiceClient(conn),
		ProductCommandServiceClient: pbproduct.NewProductCommandServiceClient(conn),
	}
}

type RoleClient struct {
	pbrole.RoleQueryServiceClient
	pbrole.RoleCommandServiceClient
}

func NewRoleClient(conn *grpc.ClientConn) *RoleClient {
	return &RoleClient{
		RoleQueryServiceClient:   pbrole.NewRoleQueryServiceClient(conn),
		RoleCommandServiceClient: pbrole.NewRoleCommandServiceClient(conn),
	}
}

type TransactionClient struct {
	pbtransaction.TransactionQueryServiceClient
	pbtransaction.TransactionCommandServiceClient
	pbtransaction.TransactionStatsServiceClient
}

func NewTransactionClient(conn *grpc.ClientConn) *TransactionClient {
	return &TransactionClient{
		TransactionQueryServiceClient:   pbtransaction.NewTransactionQueryServiceClient(conn),
		TransactionCommandServiceClient: pbtransaction.NewTransactionCommandServiceClient(conn),
		TransactionStatsServiceClient:   pbtransaction.NewTransactionStatsServiceClient(conn),
	}
}

type OrderItemClient struct {
	pb.OrderItemQueryServiceClient
	pb.OrderItemCommandServiceClient
}

func NewOrderItemClient(conn *grpc.ClientConn) *OrderItemClient {
	return &OrderItemClient{
		OrderItemQueryServiceClient:   pb.NewOrderItemQueryServiceClient(conn),
		OrderItemCommandServiceClient: pb.NewOrderItemCommandServiceClient(conn),
	}
}

type CashierClient struct {
	pbcashier.CashierQueryServiceClient
	pbcashier.CashierCommandServiceClient
	pbcashier.CashierStatsServiceClient
}

func NewCashierClient(conn *grpc.ClientConn) *CashierClient {
	return &CashierClient{
		CashierQueryServiceClient:   pbcashier.NewCashierQueryServiceClient(conn),
		CashierCommandServiceClient: pbcashier.NewCashierCommandServiceClient(conn),
		CashierStatsServiceClient:   pbcashier.NewCashierStatsServiceClient(conn),
	}
}
