package protomapper

import (
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type orderProtoMapper struct{}

func NewOrderProtoMapper() *orderProtoMapper {
	return &orderProtoMapper{}
}

func (o *orderProtoMapper) ToProtoResponseOrder(status string, message string, pbResponse *response.OrderResponse) *pborder.ApiResponseOrder {
	return &pborder.ApiResponseOrder{
		Status:  status,
		Message: message,
		Data:    o.mapResponseOrder(pbResponse),
	}
}

func (o *orderProtoMapper) ToProtoResponsesOrder(status string, message string, pbResponse []*response.OrderResponse) *pborder.ApiResponsesOrder {
	return &pborder.ApiResponsesOrder{
		Status:  status,
		Message: message,
		Data:    o.mapResponsesOrder(pbResponse),
	}
}

func (o *orderProtoMapper) ToProtoResponseOrderDeleteAt(status string, message string, pbResponse *response.OrderResponseDeleteAt) *pborder.ApiResponseOrderDeleteAt {
	return &pborder.ApiResponseOrderDeleteAt{
		Status:  status,
		Message: message,
		Data:    o.mapResponseOrderDeleteAt(pbResponse),
	}
}

func (o *orderProtoMapper) ToProtoResponseOrderDelete(status string, message string) *pborder.ApiResponseOrderDelete {
	return &pborder.ApiResponseOrderDelete{
		Status:  status,
		Message: message,
	}
}

func (o *orderProtoMapper) ToProtoResponseOrderAll(status string, message string) *pborder.ApiResponseOrderAll {
	return &pborder.ApiResponseOrderAll{
		Status:  status,
		Message: message,
	}
}

func (o *orderProtoMapper) ToProtoResponsePaginationOrderDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, orders []*response.OrderResponseDeleteAt) *pborder.ApiResponsePaginationOrderDeleteAt {
	return &pborder.ApiResponsePaginationOrderDeleteAt{
		Status:     status,
		Message:    message,
		Data:       o.mapResponsesOrderDeleteAt(orders),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (o *orderProtoMapper) ToProtoResponsePaginationOrder(pagination *pbcommon.PaginationMeta, status string, message string, orders []*response.OrderResponse) *pborder.ApiResponsePaginationOrder {
	return &pborder.ApiResponsePaginationOrder{
		Status:     status,
		Message:    message,
		Data:       o.mapResponsesOrder(orders),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (o *orderProtoMapper) ToProtoResponseMonthlyRevenue(status string, message string, row []*response.OrderMonthlyResponse) *pborder.ApiResponseOrderMonthly {
	return &pborder.ApiResponseOrderMonthly{
		Status:  status,
		Message: message,
		Data:    o.mapResponsesOrderMonthlyPrices(row),
	}
}

func (o *orderProtoMapper) ToProtoResponseYearlyRevenue(status string, message string, row []*response.OrderYearlyResponse) *pborder.ApiResponseOrderYearly {
	return &pborder.ApiResponseOrderYearly{
		Status:  status,
		Message: message,
		Data:    o.mapResponsesOrderYearlyPrices(row),
	}
}

func (o *orderProtoMapper) ToProtoResponseMonthlyTotalRevenue(status string, message string, row []*response.OrderMonthlyTotalRevenueResponse) *pborder.ApiResponseOrderMonthlyTotalRevenue {
	return &pborder.ApiResponseOrderMonthlyTotalRevenue{
		Status:  status,
		Message: message,
		Data:    o.mapResponseOrderMonthlyTotalRevenues(row),
	}
}

func (o *orderProtoMapper) ToProtoResponseYearlyTotalRevenue(status string, message string, row []*response.OrderYearlyTotalRevenueResponse) *pborder.ApiResponseOrderYearlyTotalRevenue {
	return &pborder.ApiResponseOrderYearlyTotalRevenue{
		Status:  status,
		Message: message,
		Data:    o.mapResponseOrderYearlyTotalRevenues(row),
	}
}

func (o *orderProtoMapper) mapResponseOrder(order *response.OrderResponse) *pborder.OrderResponse {
	return &pborder.OrderResponse{
		Id:         int32(order.ID),
		MerchantId: int32(order.MerchantID),
		CashierId:  int32(order.CashierID),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (o *orderProtoMapper) mapResponsesOrder(orders []*response.OrderResponse) []*pborder.OrderResponse {
	var mappedOrders []*pborder.OrderResponse

	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.mapResponseOrder(order))
	}

	return mappedOrders
}

func (o *orderProtoMapper) mapResponseOrderDeleteAt(order *response.OrderResponseDeleteAt) *pborder.OrderResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue

	if order.DeleteAt != nil {
		deletedAt = wrapperspb.String(*order.DeleteAt)
	}

	return &pborder.OrderResponseDeleteAt{
		Id:         int32(order.ID),
		MerchantId: int32(order.MerchantID),
		CashierId:  int32(order.CashierID),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		DeletedAt:  deletedAt,
	}
}

func (o *orderProtoMapper) mapResponsesOrderDeleteAt(orders []*response.OrderResponseDeleteAt) []*pborder.OrderResponseDeleteAt {
	var mappedOrders []*pborder.OrderResponseDeleteAt

	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.mapResponseOrderDeleteAt(order))
	}

	return mappedOrders
}

func (s *orderProtoMapper) mapResponseOrderMonthlyPrice(category *response.OrderMonthlyResponse) *pborder.OrderMonthlyResponse {
	return &pborder.OrderMonthlyResponse{
		Month:          category.Month,
		OrderCount:     int32(category.OrderCount),
		TotalRevenue:   int32(category.TotalRevenue),
		TotalItemsSold: int32(category.TotalItemsSold),
	}
}

func (s *orderProtoMapper) mapResponsesOrderMonthlyPrices(c []*response.OrderMonthlyResponse) []*pborder.OrderMonthlyResponse {
	var categoryRecords []*pborder.OrderMonthlyResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.mapResponseOrderMonthlyPrice(category))
	}

	return categoryRecords
}

func (s *orderProtoMapper) mapResponseOrderYearlyPrice(category *response.OrderYearlyResponse) *pborder.OrderYearlyResponse {
	return &pborder.OrderYearlyResponse{
		Year:               category.Year,
		OrderCount:         int32(category.OrderCount),
		TotalRevenue:       int32(category.TotalRevenue),
		TotalItemsSold:     int32(category.TotalItemsSold),
		ActiveCashiers:     int32(category.ActiveCashiers),
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func (s *orderProtoMapper) mapResponsesOrderYearlyPrices(c []*response.OrderYearlyResponse) []*pborder.OrderYearlyResponse {
	var categoryRecords []*pborder.OrderYearlyResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.mapResponseOrderYearlyPrice(category))
	}

	return categoryRecords
}

func (s *orderProtoMapper) mapResponseOrderMonthlyTotalRevenue(c *response.OrderMonthlyTotalRevenueResponse) *pborder.OrderMonthlyTotalRevenueResponse {
	return &pborder.OrderMonthlyTotalRevenueResponse{
		Year:           c.Year,
		Month:          c.Month,
		TotalRevenue:   int32(c.TotalRevenue),
		TotalItemsSold: int32(c.TotalItemsSold),
	}
}

func (s *orderProtoMapper) mapResponseOrderMonthlyTotalRevenues(c []*response.OrderMonthlyTotalRevenueResponse) []*pborder.OrderMonthlyTotalRevenueResponse {
	var orderRecords []*pborder.OrderMonthlyTotalRevenueResponse

	for _, row := range c {
		orderRecords = append(orderRecords, s.mapResponseOrderMonthlyTotalRevenue(row))
	}

	return orderRecords
}

func (s *orderProtoMapper) mapResponseOrderYearlyTotalRevenue(c *response.OrderYearlyTotalRevenueResponse) *pborder.OrderYearlyTotalRevenueResponse {
	return &pborder.OrderYearlyTotalRevenueResponse{
		Year:         c.Year,
		TotalRevenue: int32(c.TotalRevenue),
	}
}

func (s *orderProtoMapper) mapResponseOrderYearlyTotalRevenues(c []*response.OrderYearlyTotalRevenueResponse) []*pborder.OrderYearlyTotalRevenueResponse {
	var orderRecords []*pborder.OrderYearlyTotalRevenueResponse

	for _, row := range c {
		orderRecords = append(orderRecords, s.mapResponseOrderYearlyTotalRevenue(row))
	}

	return orderRecords
}
