package cashiergraphqlmapper

import (
	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-pointofsale-pb/cashier"
)

type CashierGraphqlMapper interface {
	ToGraphqlResponseCashier(res *pb.ApiResponseCashier) *model.APIResponseCashier
	ToGraphqlResponsesCashier(res *pb.ApiResponsesCashier) *model.APIResponsesCashier
	ToGraphqlResponseCashierDeleteAt(res *pb.ApiResponseCashierDeleteAt) *model.APIResponseCashierDeleteAt
	ToGraphqlResponseCashierDelete(res *pb.ApiResponseCashierDelete) *model.APIResponseCashierDelete
	ToGraphqlResponseCashierAll(res *pb.ApiResponseCashierAll) *model.APIResponseCashierAll
	ToGraphqlResponsePaginationCashier(res *pb.ApiResponsePaginationCashier) *model.APIResponsePaginationCashier
	ToGraphqlResponsePaginationCashierDeleteAt(res *pb.ApiResponsePaginationCashierDeleteAt) *model.APIResponsePaginationCashierDeleteAt
	ToGraphqlResponseMonthlyTotalSales(res *pb.ApiResponseCashierMonthlyTotalSales) *model.APIResponseCashierMonthlyTotalSales
	ToGraphqlResponseMonthlySales(res *pb.ApiResponseCashierMonthSales) *model.APIResponseCashierMonthSales
	ToGraphqlResponseYearlySales(res *pb.ApiResponseCashierYearSales) *model.APIResponseCashierYearSales
	ToGraphqlResponseYearlyTotalSales(res *pb.ApiResponseCashierYearlyTotalSales) *model.APIResponseCashierYearlyTotalSales
}
