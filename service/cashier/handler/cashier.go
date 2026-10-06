package handler

import (
	pbcashier "github.com/MamangRust/monolith-pointofsale-pb/cashier"
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
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

func mapResponseCashier(cashier *db.Cashier) *pbcashier.CashierResponse {
	if cashier == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if cashier.CreatedAt.Valid {
		createdAtStr = cashier.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if cashier.UpdatedAt.Valid {
		updatedAtStr = cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return &pbcashier.CashierResponse{
		Id:         int32(cashier.CashierID),
		MerchantId: int32(cashier.MerchantID),
		Name:       cashier.Name,
		CreatedAt:  createdAtStr,
		UpdatedAt:  updatedAtStr,
	}
}

func mapResponsesCashier(cashiers []*db.GetCashiersRow) []*pbcashier.CashierResponse {
	var mappedCashiers []*pbcashier.CashierResponse
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if cashier.CreatedAt.Valid {
			createdAtStr = cashier.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		if cashier.UpdatedAt.Valid {
			updatedAtStr = cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		mappedCashiers = append(mappedCashiers, &pbcashier.CashierResponse{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
		})
	}
	return mappedCashiers
}

func mapResponsesCashierByMerchant(cashiers []*db.GetCashiersByMerchantRow) []*pbcashier.CashierResponse {
	var mappedCashiers []*pbcashier.CashierResponse
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if cashier.CreatedAt.Valid {
			createdAtStr = cashier.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		if cashier.UpdatedAt.Valid {
			updatedAtStr = cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		mappedCashiers = append(mappedCashiers, &pbcashier.CashierResponse{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
		})
	}
	return mappedCashiers
}

func mapResponseCashierDeleteAt(cashier *db.Cashier) *pbcashier.CashierResponseDeleteAt {
	if cashier == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if cashier.CreatedAt.Valid {
		createdAtStr = cashier.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if cashier.UpdatedAt.Valid {
		updatedAtStr = cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var deletedAt *wrapperspb.StringValue
	if cashier.DeletedAt.Valid {
		deletedAt = wrapperspb.String(cashier.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}

	return &pbcashier.CashierResponseDeleteAt{
		Id:         int32(cashier.CashierID),
		MerchantId: int32(cashier.MerchantID),
		Name:       cashier.Name,
		CreatedAt:  createdAtStr,
		UpdatedAt:  updatedAtStr,
		DeletedAt:  deletedAt,
	}
}

func mapResponsesCashierActive(cashiers []*db.GetCashiersActiveRow) []*pbcashier.CashierResponseDeleteAt {
	var mappedCashiers []*pbcashier.CashierResponseDeleteAt
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if cashier.CreatedAt.Valid {
			createdAtStr = cashier.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		if cashier.UpdatedAt.Valid {
			updatedAtStr = cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var deletedAt *wrapperspb.StringValue
		if cashier.DeletedAt.Valid {
			deletedAt = wrapperspb.String(cashier.DeletedAt.Time.Format("2006-01-02 15:04:05"))
		}
		mappedCashiers = append(mappedCashiers, &pbcashier.CashierResponseDeleteAt{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
			DeletedAt:  deletedAt,
		})
	}
	return mappedCashiers
}

func mapResponsesCashierTrashed(cashiers []*db.GetCashiersTrashedRow) []*pbcashier.CashierResponseDeleteAt {
	var mappedCashiers []*pbcashier.CashierResponseDeleteAt
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if cashier.CreatedAt.Valid {
			createdAtStr = cashier.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		if cashier.UpdatedAt.Valid {
			updatedAtStr = cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var deletedAt *wrapperspb.StringValue
		if cashier.DeletedAt.Valid {
			deletedAt = wrapperspb.String(cashier.DeletedAt.Time.Format("2006-01-02 15:04:05"))
		}
		mappedCashiers = append(mappedCashiers, &pbcashier.CashierResponseDeleteAt{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
			DeletedAt:  deletedAt,
		})
	}
	return mappedCashiers
}

func mapResponseCashierMonthlySale(cashier *db.GetMonthlyCashierRow) *pbcashier.CashierResponseMonthSales {
	if cashier == nil {
		return nil
	}
	return &pbcashier.CashierResponseMonthSales{
		Month:       cashier.Month,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierMonthlySales(c []*db.GetMonthlyCashierRow) []*pbcashier.CashierResponseMonthSales {
	var cashierRecords []*pbcashier.CashierResponseMonthSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierMonthlySale(cashier))
	}
	return cashierRecords
}

func mapResponseCashierMonthlySaleById(cashier *db.GetMonthlyCashierByCashierIdRow) *pbcashier.CashierResponseMonthSales {
	if cashier == nil {
		return nil
	}
	return &pbcashier.CashierResponseMonthSales{
		Month:       cashier.Month,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierMonthlySalesById(c []*db.GetMonthlyCashierByCashierIdRow) []*pbcashier.CashierResponseMonthSales {
	var cashierRecords []*pbcashier.CashierResponseMonthSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierMonthlySaleById(cashier))
	}
	return cashierRecords
}

func mapResponseCashierMonthlySaleByMerchant(cashier *db.GetMonthlyCashierByMerchantRow) *pbcashier.CashierResponseMonthSales {
	if cashier == nil {
		return nil
	}
	return &pbcashier.CashierResponseMonthSales{
		Month:       cashier.Month,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierMonthlySalesByMerchant(c []*db.GetMonthlyCashierByMerchantRow) []*pbcashier.CashierResponseMonthSales {
	var cashierRecords []*pbcashier.CashierResponseMonthSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierMonthlySaleByMerchant(cashier))
	}
	return cashierRecords
}

func mapResponseCashierYearlySale(cashier *db.GetYearlyCashierRow) *pbcashier.CashierResponseYearSales {
	if cashier == nil {
		return nil
	}
	return &pbcashier.CashierResponseYearSales{
		Year:        cashier.Year,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierYearlySales(c []*db.GetYearlyCashierRow) []*pbcashier.CashierResponseYearSales {
	var cashierRecords []*pbcashier.CashierResponseYearSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierYearlySale(cashier))
	}
	return cashierRecords
}

func mapResponseCashierYearlySaleById(cashier *db.GetYearlyCashierByCashierIdRow) *pbcashier.CashierResponseYearSales {
	if cashier == nil {
		return nil
	}
	return &pbcashier.CashierResponseYearSales{
		Year:        cashier.Year,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierYearlySalesById(c []*db.GetYearlyCashierByCashierIdRow) []*pbcashier.CashierResponseYearSales {
	var cashierRecords []*pbcashier.CashierResponseYearSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierYearlySaleById(cashier))
	}
	return cashierRecords
}

func mapResponseCashierYearlySaleByMerchant(cashier *db.GetYearlyCashierByMerchantRow) *pbcashier.CashierResponseYearSales {
	if cashier == nil {
		return nil
	}
	return &pbcashier.CashierResponseYearSales{
		Year:        cashier.Year,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierYearlySalesByMerchant(c []*db.GetYearlyCashierByMerchantRow) []*pbcashier.CashierResponseYearSales {
	var cashierRecords []*pbcashier.CashierResponseYearSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierYearlySaleByMerchant(cashier))
	}
	return cashierRecords
}

func mapResponseCashierMonthlyTotalSales(c []*db.GetMonthlyTotalSalesCashierRow) []*pbcashier.CashierResponseMonthTotalSales {
	var cashierRecords []*pbcashier.CashierResponseMonthTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashier.CashierResponseMonthTotalSales{
			Year:       cashier.Year,
			Month:      cashier.Month,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierMonthlyTotalSalesById(c []*db.GetMonthlyTotalSalesByIdRow) []*pbcashier.CashierResponseMonthTotalSales {
	var cashierRecords []*pbcashier.CashierResponseMonthTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashier.CashierResponseMonthTotalSales{
			Year:       cashier.Year,
			Month:      cashier.Month,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierMonthlyTotalSalesByMerchant(c []*db.GetMonthlyTotalSalesByMerchantRow) []*pbcashier.CashierResponseMonthTotalSales {
	var cashierRecords []*pbcashier.CashierResponseMonthTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashier.CashierResponseMonthTotalSales{
			Year:       cashier.Year,
			Month:      cashier.Month,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierYearlyTotalSales(c []*db.GetYearlyTotalSalesCashierRow) []*pbcashier.CashierResponseYearTotalSales {
	var cashierRecords []*pbcashier.CashierResponseYearTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashier.CashierResponseYearTotalSales{
			Year:       cashier.Year,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierYearlyTotalSalesById(c []*db.GetYearlyTotalSalesByIdRow) []*pbcashier.CashierResponseYearTotalSales {
	var cashierRecords []*pbcashier.CashierResponseYearTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashier.CashierResponseYearTotalSales{
			Year:       cashier.Year,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierYearlyTotalSalesByMerchant(c []*db.GetYearlyTotalSalesByMerchantRow) []*pbcashier.CashierResponseYearTotalSales {
	var cashierRecords []*pbcashier.CashierResponseYearTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashier.CashierResponseYearTotalSales{
			Year:       cashier.Year,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}
