package handler

import (
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func mapPaginationMeta(meta *pbcommon.PaginationMeta) *pbcommon.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &pbcommon.PaginationMeta{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		TotalPages:   meta.TotalPages,
		TotalRecords: meta.TotalRecords,
	}
}

func mapResponseOrder(order *db.Order) *pborder.OrderResponse {
	if order == nil {
		return nil
	}
	return &pborder.OrderResponse{
		Id:         int32(order.OrderID),
		MerchantId: int32(order.MerchantID),
		CashierId:  int32(order.CashierID),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}
}

func mapResponsesOrderByMerchant(orders []*db.GetOrdersByMerchantRow) []*pborder.OrderResponse {
	var mappedOrders []*pborder.OrderResponse
	for _, order := range orders {
		if order == nil {
			continue
		}
		mappedOrders = append(mappedOrders, &pborder.OrderResponse{
			Id:         order.OrderID,
			MerchantId: order.MerchantID,
			CashierId:  order.CashierID,
			TotalPrice: int32(order.TotalPrice),
			CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return mappedOrders
}

func mapResponsesOrder(orders []*db.GetOrdersRow) []*pborder.OrderResponse {
	var mappedOrders []*pborder.OrderResponse
	for _, order := range orders {
		if order == nil {
			continue
		}
		mappedOrders = append(mappedOrders, &pborder.OrderResponse{
			Id:         int32(order.OrderID),
			MerchantId: int32(order.MerchantID),
			CashierId:  int32(order.CashierID),
			TotalPrice: int32(order.TotalPrice),
			CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return mappedOrders
}

func mapResponseOrderDeleteAt(order *db.Order) *pborder.OrderResponseDeleteAt {
	if order == nil {
		return nil
	}
	var deletedAt *wrapperspb.StringValue
	if order.DeletedAt.Valid {
		deletedAt = wrapperspb.String(order.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}

	return &pborder.OrderResponseDeleteAt{
		Id:         int32(order.OrderID),
		MerchantId: int32(order.MerchantID),
		CashierId:  int32(order.CashierID),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:  deletedAt,
	}
}

func mapResponsesOrderActive(orders []*db.GetOrdersActiveRow) []*pborder.OrderResponseDeleteAt {
	var mappedOrders []*pborder.OrderResponseDeleteAt
	for _, order := range orders {
		if order == nil {
			continue
		}
		var deletedAt *wrapperspb.StringValue
		if order.DeletedAt.Valid {
			deletedAt = wrapperspb.String(order.DeletedAt.Time.Format("2006-01-02 15:04:05"))
		}
		mappedOrders = append(mappedOrders, &pborder.OrderResponseDeleteAt{
			Id:         int32(order.OrderID),
			MerchantId: int32(order.MerchantID),
			CashierId:  int32(order.CashierID),
			TotalPrice: int32(order.TotalPrice),
			CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt:  deletedAt,
		})
	}
	return mappedOrders
}

func mapResponsesOrderTrashed(orders []*db.GetOrdersTrashedRow) []*pborder.OrderResponseDeleteAt {
	var mappedOrders []*pborder.OrderResponseDeleteAt
	for _, order := range orders {
		if order == nil {
			continue
		}
		var deletedAt *wrapperspb.StringValue
		if order.DeletedAt.Valid {
			deletedAt = wrapperspb.String(order.DeletedAt.Time.Format("2006-01-02 15:04:05"))
		}
		mappedOrders = append(mappedOrders, &pborder.OrderResponseDeleteAt{
			Id:         int32(order.OrderID),
			MerchantId: int32(order.MerchantID),
			CashierId:  int32(order.CashierID),
			TotalPrice: int32(order.TotalPrice),
			CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt:  deletedAt,
		})
	}
	return mappedOrders
}

func mapResponseOrderMonthlyTotalRevenue(row *db.GetMonthlyTotalRevenueRow) *pborder.OrderMonthlyTotalRevenueResponse {
	if row == nil {
		return nil
	}
	return &pborder.OrderMonthlyTotalRevenueResponse{
		Year:           row.Year,
		Month:          row.Month,
		TotalRevenue:   int32(row.TotalRevenue),
		TotalItemsSold: 0,
	}
}

func mapResponseOrderMonthlyTotalRevenues(c []*db.GetMonthlyTotalRevenueRow) []*pborder.OrderMonthlyTotalRevenueResponse {
	var orderRecords []*pborder.OrderMonthlyTotalRevenueResponse
	for _, row := range c {
		orderRecords = append(orderRecords, mapResponseOrderMonthlyTotalRevenue(row))
	}
	return orderRecords
}

func mapResponseOrderMonthlyTotalRevenueByMerchant(row *db.GetMonthlyTotalRevenueByMerchantRow) *pborder.OrderMonthlyTotalRevenueResponse {
	if row == nil {
		return nil
	}
	return &pborder.OrderMonthlyTotalRevenueResponse{
		Year:           row.Year,
		Month:          row.Month,
		TotalRevenue:   int32(row.TotalRevenue),
		TotalItemsSold: 0,
	}
}

func mapResponseOrderMonthlyTotalRevenuesByMerchant(c []*db.GetMonthlyTotalRevenueByMerchantRow) []*pborder.OrderMonthlyTotalRevenueResponse {
	var orderRecords []*pborder.OrderMonthlyTotalRevenueResponse
	for _, row := range c {
		orderRecords = append(orderRecords, mapResponseOrderMonthlyTotalRevenueByMerchant(row))
	}
	return orderRecords
}

func mapResponseOrderYearlyTotalRevenue(row *db.GetYearlyTotalRevenueRow) *pborder.OrderYearlyTotalRevenueResponse {
	if row == nil {
		return nil
	}
	return &pborder.OrderYearlyTotalRevenueResponse{
		Year:         row.Year,
		TotalRevenue: int32(row.TotalRevenue),
	}
}

func mapResponseOrderYearlyTotalRevenues(c []*db.GetYearlyTotalRevenueRow) []*pborder.OrderYearlyTotalRevenueResponse {
	var orderRecords []*pborder.OrderYearlyTotalRevenueResponse
	for _, row := range c {
		orderRecords = append(orderRecords, mapResponseOrderYearlyTotalRevenue(row))
	}
	return orderRecords
}

func mapResponseOrderYearlyTotalRevenueByMerchant(row *db.GetYearlyTotalRevenueByMerchantRow) *pborder.OrderYearlyTotalRevenueResponse {
	if row == nil {
		return nil
	}
	return &pborder.OrderYearlyTotalRevenueResponse{
		Year:         row.Year,
		TotalRevenue: int32(row.TotalRevenue),
	}
}

func mapResponseOrderYearlyTotalRevenuesByMerchant(c []*db.GetYearlyTotalRevenueByMerchantRow) []*pborder.OrderYearlyTotalRevenueResponse {
	var orderRecords []*pborder.OrderYearlyTotalRevenueResponse
	for _, row := range c {
		orderRecords = append(orderRecords, mapResponseOrderYearlyTotalRevenueByMerchant(row))
	}
	return orderRecords
}

func mapResponsesOrderMonthlyPrices(c []*db.GetMonthlyOrderRow) []*pborder.OrderMonthlyResponse {
	var categoryRecords []*pborder.OrderMonthlyResponse
	for _, category := range c {
		if category == nil {
			continue
		}
		categoryRecords = append(categoryRecords, &pborder.OrderMonthlyResponse{
			Month:          category.Month,
			OrderCount:     int32(category.OrderCount),
			TotalRevenue:   int32(category.TotalRevenue),
			TotalItemsSold: int32(category.TotalItemsSold),
		})
	}
	return categoryRecords
}

func mapResponsesOrderMonthlyPricesByMerchant(c []*db.GetMonthlyOrderByMerchantRow) []*pborder.OrderMonthlyResponse {
	var categoryRecords []*pborder.OrderMonthlyResponse
	for _, category := range c {
		if category == nil {
			continue
		}
		categoryRecords = append(categoryRecords, &pborder.OrderMonthlyResponse{
			Month:          category.Month,
			OrderCount:     int32(category.OrderCount),
			TotalRevenue:   int32(category.TotalRevenue),
			TotalItemsSold: int32(category.TotalItemsSold),
		})
	}
	return categoryRecords
}

func mapResponsesOrderYearlyPrices(c []*db.GetYearlyOrderRow) []*pborder.OrderYearlyResponse {
	var categoryRecords []*pborder.OrderYearlyResponse
	for _, category := range c {
		if category == nil {
			continue
		}
		categoryRecords = append(categoryRecords, &pborder.OrderYearlyResponse{
			Year:               category.Year,
			OrderCount:         int32(category.OrderCount),
			TotalRevenue:       int32(category.TotalRevenue),
			TotalItemsSold:     int32(category.TotalItemsSold),
			ActiveCashiers:     int32(category.ActiveCashiers),
			UniqueProductsSold: int32(category.UniqueProductsSold),
		})
	}
	return categoryRecords
}

func mapResponsesOrderYearlyPricesByMerchant(c []*db.GetYearlyOrderByMerchantRow) []*pborder.OrderYearlyResponse {
	var categoryRecords []*pborder.OrderYearlyResponse
	for _, category := range c {
		if category == nil {
			continue
		}
		categoryRecords = append(categoryRecords, &pborder.OrderYearlyResponse{
			Year:               category.Year,
			OrderCount:         int32(category.OrderCount),
			TotalRevenue:       int32(category.TotalRevenue),
			TotalItemsSold:     int32(category.TotalItemsSold),
			ActiveCashiers:     int32(category.ActiveCashiers),
			UniqueProductsSold: int32(category.UniqueProductsSold),
		})
	}
	return categoryRecords
}
