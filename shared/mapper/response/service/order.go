package response_service

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/record"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type orderResponseMapper struct {
}

func NewOrderResponseMapper() *orderResponseMapper {
	return &orderResponseMapper{}
}

func (s *orderResponseMapper) ToOrderResponse(order *record.OrderRecord) *response.OrderResponse {
	return &response.OrderResponse{
		ID:         order.ID,
		MerchantID: order.MerchantID,
		CashierID:  order.CashierID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (s *orderResponseMapper) ToOrdersResponse(orders []*record.OrderRecord) []*response.OrderResponse {
	var responses []*response.OrderResponse

	for _, order := range orders {
		responses = append(responses, s.ToOrderResponse(order))
	}

	return responses
}

func (s *orderResponseMapper) ToOrderResponseDeleteAt(order *record.OrderRecord) *response.OrderResponseDeleteAt {
	return &response.OrderResponseDeleteAt{
		ID:         order.ID,
		MerchantID: order.MerchantID,
		CashierID:  order.CashierID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		DeleteAt:   order.DeletedAt,
	}
}

func (s *orderResponseMapper) ToOrdersResponseDeleteAt(orders []*record.OrderRecord) []*response.OrderResponseDeleteAt {
	var responses []*response.OrderResponseDeleteAt

	for _, order := range orders {
		responses = append(responses, s.ToOrderResponseDeleteAt(order))
	}

	return responses
}

func (s *orderResponseMapper) ToOrderMonthlyPrice(category *record.OrderMonthlyRecord) *response.OrderMonthlyResponse {
	return &response.OrderMonthlyResponse{
		Month:          category.Month,
		OrderCount:     int(category.OrderCount),
		TotalRevenue:   int(category.TotalRevenue),
		TotalItemsSold: int(category.TotalItemsSold),
	}
}

func (s *orderResponseMapper) ToOrderMonthlyPrices(c []*record.OrderMonthlyRecord) []*response.OrderMonthlyResponse {
	var categoryRecords []*response.OrderMonthlyResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToOrderMonthlyPrice(category))
	}

	return categoryRecords
}

func (s *orderResponseMapper) ToOrderYearlyPrice(category *record.OrderYearlyRecord) *response.OrderYearlyResponse {
	return &response.OrderYearlyResponse{
		Year:               category.Year,
		OrderCount:         int(category.OrderCount),
		TotalRevenue:       int(category.TotalRevenue),
		TotalItemsSold:     int(category.TotalItemsSold),
		ActiveCashiers:     int(category.ActiveCashiers),
		UniqueProductsSold: int(category.UniqueProductsSold),
	}
}

func (s *orderResponseMapper) ToOrderYearlyPrices(c []*record.OrderYearlyRecord) []*response.OrderYearlyResponse {
	var categoryRecords []*response.OrderYearlyResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToOrderYearlyPrice(category))
	}

	return categoryRecords
}

func (s *orderResponseMapper) ToOrderMonthlyTotalRevenue(c *record.OrderMonthlyTotalRevenueRecord) *response.OrderMonthlyTotalRevenueResponse {
	return &response.OrderMonthlyTotalRevenueResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *orderResponseMapper) ToOrderMonthlyTotalRevenues(c []*record.OrderMonthlyTotalRevenueRecord) []*response.OrderMonthlyTotalRevenueResponse {
	var orderRecords []*response.OrderMonthlyTotalRevenueResponse

	for _, row := range c {
		orderRecords = append(orderRecords, s.ToOrderMonthlyTotalRevenue(row))
	}

	return orderRecords
}

func (s *orderResponseMapper) ToOrderYearlyTotalRevenue(c *record.OrderYearlyTotalRevenueRecord) *response.OrderYearlyTotalRevenueResponse {
	return &response.OrderYearlyTotalRevenueResponse{
		Year:         c.Year,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *orderResponseMapper) ToOrderYearlyTotalRevenues(c []*record.OrderYearlyTotalRevenueRecord) []*response.OrderYearlyTotalRevenueResponse {
	var orderRecords []*response.OrderYearlyTotalRevenueResponse

	for _, row := range c {
		orderRecords = append(orderRecords, s.ToOrderYearlyTotalRevenue(row))
	}

	return orderRecords
}
