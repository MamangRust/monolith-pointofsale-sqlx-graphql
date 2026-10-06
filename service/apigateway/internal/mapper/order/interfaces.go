package ordergraphqlmapper

import (
	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-pointofsale-pb/order"
)

type OrderGraphqlMapper interface {
	ToGraphqlResponseOrder(res *pb.ApiResponseOrder) *model.APIResponseOrder
	ToGraphqlResponsesOrder(res *pb.ApiResponsesOrder) *model.APIResponsesOrder
	ToGraphqlResponseOrderDeleteAt(res *pb.ApiResponseOrderDeleteAt) *model.APIResponseOrderDeleteAt
	ToGraphqlResponseOrderDelete(res *pb.ApiResponseOrderDelete) *model.APIResponseOrderDelete
	ToGraphqlResponseOrderAll(res *pb.ApiResponseOrderAll) *model.APIResponseOrderAll
	ToGraphqlResponsePaginationOrder(res *pb.ApiResponsePaginationOrder) *model.APIResponsePaginationOrder
	ToGraphqlResponsePaginationOrderDeleteAt(res *pb.ApiResponsePaginationOrderDeleteAt) *model.APIResponsePaginationOrderDeleteAt
	ToGraphqlResponseMonthlyRevenue(res *pb.ApiResponseOrderMonthly) *model.APIResponseOrderMonthly
	ToGraphqlResponseYearlyRevenue(res *pb.ApiResponseOrderYearly) *model.APIResponseOrderYearly
	ToGraphqlResponseMonthlyTotalRevenue(res *pb.ApiResponseOrderMonthlyTotalRevenue) *model.APIResponseOrderMonthlyTotalRevenue
	ToGraphqlResponseYearlyTotalRevenue(res *pb.ApiResponseOrderYearlyTotalRevenue) *model.APIResponseOrderYearlyTotalRevenue
}
