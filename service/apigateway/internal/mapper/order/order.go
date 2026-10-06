package ordergraphqlmapper

import (
	graphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper"
	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-pointofsale-pb/order"
)

type orderGraphqlMapper struct {
}

func NewOrderGraphqlMapper() *orderGraphqlMapper {
	return &orderGraphqlMapper{}
}

func (o *orderGraphqlMapper) ToGraphqlResponseOrder(res *pb.ApiResponseOrder) *model.APIResponseOrder {
	return &model.APIResponseOrder{
		Status:  res.Status,
		Message: res.Message,
		Data:    o.mapResponseOrder(res.Data),
	}
}

func (o *orderGraphqlMapper) ToGraphqlResponsesOrder(res *pb.ApiResponsesOrder) *model.APIResponsesOrder {
	return &model.APIResponsesOrder{
		Status:  res.Status,
		Message: res.Message,
		Data:    o.mapResponsesOrder(res.Data),
	}
}

func (o *orderGraphqlMapper) ToGraphqlResponseOrderDeleteAt(res *pb.ApiResponseOrderDeleteAt) *model.APIResponseOrderDeleteAt {
	return &model.APIResponseOrderDeleteAt{
		Status:  res.Status,
		Message: res.Message,
		Data:    o.mapResponseOrderDeleteAt(res.Data),
	}
}

func (o *orderGraphqlMapper) ToGraphqlResponseOrderDelete(res *pb.ApiResponseOrderDelete) *model.APIResponseOrderDelete {
	return &model.APIResponseOrderDelete{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (o *orderGraphqlMapper) ToGraphqlResponseOrderAll(res *pb.ApiResponseOrderAll) *model.APIResponseOrderAll {
	return &model.APIResponseOrderAll{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (o *orderGraphqlMapper) ToGraphqlResponsePaginationOrder(res *pb.ApiResponsePaginationOrder) *model.APIResponsePaginationOrder {
	return &model.APIResponsePaginationOrder{
		Status:     res.Status,
		Message:    res.Message,
		Data:       o.mapResponsesOrder(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (o *orderGraphqlMapper) ToGraphqlResponsePaginationOrderDeleteAt(res *pb.ApiResponsePaginationOrderDeleteAt) *model.APIResponsePaginationOrderDeleteAt {
	return &model.APIResponsePaginationOrderDeleteAt{
		Status:     res.Status,
		Message:    res.Message,
		Data:       o.mapResponsesOrderDeleteAt(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (o *orderGraphqlMapper) ToGraphqlResponseMonthlyRevenue(res *pb.ApiResponseOrderMonthly) *model.APIResponseOrderMonthly {
	return &model.APIResponseOrderMonthly{
		Status:  res.Status,
		Message: res.Message,
		Data:    o.mapResponsesOrderMonthlyPrice(res.Data),
	}
}

func (o *orderGraphqlMapper) ToGraphqlResponseYearlyRevenue(res *pb.ApiResponseOrderYearly) *model.APIResponseOrderYearly {
	return &model.APIResponseOrderYearly{
		Status:  res.Status,
		Message: res.Message,
		Data:    o.mapResponsesOrderYearlyPrice(res.Data),
	}
}

func (o *orderGraphqlMapper) ToGraphqlResponseMonthlyTotalRevenue(res *pb.ApiResponseOrderMonthlyTotalRevenue) *model.APIResponseOrderMonthlyTotalRevenue {
	return &model.APIResponseOrderMonthlyTotalRevenue{
		Status:  res.Status,
		Message: res.Message,
		Data:    o.mapResponsesOrderMonthlyTotalRevenue(res.Data),
	}
}

func (o *orderGraphqlMapper) ToGraphqlResponseYearlyTotalRevenue(res *pb.ApiResponseOrderYearlyTotalRevenue) *model.APIResponseOrderYearlyTotalRevenue {
	return &model.APIResponseOrderYearlyTotalRevenue{
		Status:  res.Status,
		Message: res.Message,
		Data:    o.mapResponsesOrderYearlyTotalRevenue(res.Data),
	}
}

func (o *orderGraphqlMapper) mapResponseOrder(order *pb.OrderResponse) *model.OrderResponse {
	if order == nil {
		return nil
	}
	return &model.OrderResponse{
		ID:         int32(order.Id),
		MerchantID: int32(order.MerchantId),
		CashierID:  int32(order.CashierId),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  &order.CreatedAt,
		UpdatedAt:  &order.UpdatedAt,
	}
}

func (o *orderGraphqlMapper) mapResponsesOrder(orders []*pb.OrderResponse) []*model.OrderResponse {
	var responses []*model.OrderResponse
	for _, order := range orders {
		responses = append(responses, o.mapResponseOrder(order))
	}
	return responses
}

func (o *orderGraphqlMapper) mapResponseOrderDeleteAt(order *pb.OrderResponseDeleteAt) *model.OrderResponseDeleteAt {
	if order == nil {
		return nil
	}
	var deletedAt string
	if order.DeletedAt != nil {
		deletedAt = order.DeletedAt.Value
	}

	return &model.OrderResponseDeleteAt{
		ID:         int32(order.Id),
		MerchantID: int32(order.MerchantId),
		CashierID:  int32(order.CashierId),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  &order.CreatedAt,
		UpdatedAt:  &order.UpdatedAt,
		DeletedAt:  &deletedAt,
	}
}

func (o *orderGraphqlMapper) mapResponsesOrderDeleteAt(orders []*pb.OrderResponseDeleteAt) []*model.OrderResponseDeleteAt {
	var responses []*model.OrderResponseDeleteAt
	for _, order := range orders {
		responses = append(responses, o.mapResponseOrderDeleteAt(order))
	}
	return responses
}

func (o *orderGraphqlMapper) mapResponseOrderMonthlyPrice(order *pb.OrderMonthlyResponse) *model.OrderMonthlyResponse {
	if order == nil {
		return nil
	}
	orderCount := int32(order.OrderCount)
	totalRevenue := int32(order.TotalRevenue)
	totalItemSold := int32(order.TotalItemsSold)

	return &model.OrderMonthlyResponse{
		Month:          &order.Month,
		OrderCount:     &orderCount,
		TotalRevenue:   &totalRevenue,
		TotalItemsSold: &totalItemSold,
	}
}

func (o *orderGraphqlMapper) mapResponsesOrderMonthlyPrice(orders []*pb.OrderMonthlyResponse) []*model.OrderMonthlyResponse {
	var responses []*model.OrderMonthlyResponse
	for _, order := range orders {
		responses = append(responses, o.mapResponseOrderMonthlyPrice(order))
	}
	return responses
}

func (o *orderGraphqlMapper) mapResponseOrderYearlyPrice(order *pb.OrderYearlyResponse) *model.OrderYearlyResponse {
	if order == nil {
		return nil
	}
	orderCount := int32(order.OrderCount)
	totalRevenue := int32(order.TotalRevenue)
	totalItemSold := int32(order.TotalItemsSold)
	activeCashiers := int32(order.ActiveCashiers)
	uniqueProductsSold := int32(order.UniqueProductsSold)

	return &model.OrderYearlyResponse{
		Year:               &order.Year,
		OrderCount:         &orderCount,
		TotalRevenue:       &totalRevenue,
		TotalItemsSold:     &totalItemSold,
		ActiveCashiers:     &activeCashiers,
		UniqueProductsSold: &uniqueProductsSold,
	}
}

func (o *orderGraphqlMapper) mapResponsesOrderYearlyPrice(orders []*pb.OrderYearlyResponse) []*model.OrderYearlyResponse {
	var responses []*model.OrderYearlyResponse
	for _, order := range orders {
		responses = append(responses, o.mapResponseOrderYearlyPrice(order))
	}
	return responses
}

func (o *orderGraphqlMapper) mapResponseOrderMonthlyTotalRevenue(c *pb.OrderMonthlyTotalRevenueResponse) *model.OrderMonthlyTotalRevenue {
	if c == nil {
		return nil
	}
	totalRevenue := int32(c.TotalRevenue)
	totalItemsSold := int32(c.TotalItemsSold)

	return &model.OrderMonthlyTotalRevenue{
		Year:           &c.Year,
		Month:          &c.Month,
		TotalRevenue:   &totalRevenue,
		TotalItemsSold: &totalItemsSold,
	}
}

func (o *orderGraphqlMapper) mapResponsesOrderMonthlyTotalRevenue(c []*pb.OrderMonthlyTotalRevenueResponse) []*model.OrderMonthlyTotalRevenue {
	var responses []*model.OrderMonthlyTotalRevenue
	for _, row := range c {
		responses = append(responses, o.mapResponseOrderMonthlyTotalRevenue(row))
	}
	return responses
}

func (o *orderGraphqlMapper) mapResponseOrderYearlyTotalRevenue(c *pb.OrderYearlyTotalRevenueResponse) *model.OrderYearlyTotalRevenue {
	if c == nil {
		return nil
	}
	totalRevenue := int32(c.TotalRevenue)

	return &model.OrderYearlyTotalRevenue{
		Year:         &c.Year,
		TotalRevenue: &totalRevenue,
	}
}

func (o *orderGraphqlMapper) mapResponsesOrderYearlyTotalRevenue(c []*pb.OrderYearlyTotalRevenueResponse) []*model.OrderYearlyTotalRevenue {
	var responses []*model.OrderYearlyTotalRevenue
	for _, row := range c {
		responses = append(responses, o.mapResponseOrderYearlyTotalRevenue(row))
	}
	return responses
}
