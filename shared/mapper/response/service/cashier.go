package response_service

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/record"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type cashierResponseMapper struct {
}

func NewCashierResponseMapper() *cashierResponseMapper {
	return &cashierResponseMapper{}
}

func (s *cashierResponseMapper) ToCashierResponse(cashier *record.CashierRecord) *response.CashierResponse {
	return &response.CashierResponse{
		ID:         cashier.ID,
		MerchantID: cashier.MerchantID,
		Name:       cashier.Name,
		CreatedAt:  cashier.CreatedAt,
		UpdatedAt:  cashier.UpdatedAt,
	}
}

func (s *cashierResponseMapper) ToCashiersResponse(cashiers []*record.CashierRecord) []*response.CashierResponse {
	var responses []*response.CashierResponse

	for _, cashier := range cashiers {
		responses = append(responses, s.ToCashierResponse(cashier))
	}

	return responses
}

func (s *cashierResponseMapper) ToCashierResponseDeleteAt(cashier *record.CashierRecord) *response.CashierResponseDeleteAt {
	return &response.CashierResponseDeleteAt{
		ID:         cashier.ID,
		MerchantID: cashier.MerchantID,
		Name:       cashier.Name,
		CreatedAt:  cashier.CreatedAt,
		UpdatedAt:  cashier.UpdatedAt,
		DeletedAt:  cashier.DeletedAt,
	}
}

func (s *cashierResponseMapper) ToCashiersResponseDeleteAt(cashiers []*record.CashierRecord) []*response.CashierResponseDeleteAt {
	var responses []*response.CashierResponseDeleteAt

	for _, cashier := range cashiers {
		responses = append(responses, s.ToCashierResponseDeleteAt(cashier))
	}

	return responses
}

func (s *cashierResponseMapper) ToCashierMonthlySale(cashier *record.CashierRecordMonthSales) *response.CashierResponseMonthSales {
	return &response.CashierResponseMonthSales{
		Month:       cashier.Month,
		CashierID:   int(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int(cashier.OrderCount),
		TotalSales:  int(cashier.TotalSales),
	}
}

func (s *cashierResponseMapper) ToCashierMonthlySales(c []*record.CashierRecordMonthSales) []*response.CashierResponseMonthSales {
	var cashierRecords []*response.CashierResponseMonthSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierMonthlySale(cashier))
	}

	return cashierRecords
}

func (s *cashierResponseMapper) ToCashierYearlySale(cashier *record.CashierRecordYearSales) *response.CashierResponseYearSales {
	return &response.CashierResponseYearSales{
		Year:        cashier.Year,
		CashierID:   int(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int(cashier.OrderCount),
		TotalSales:  int(cashier.TotalSales),
	}
}

func (s *cashierResponseMapper) ToCashierYearlySales(c []*record.CashierRecordYearSales) []*response.CashierResponseYearSales {
	var cashierRecords []*response.CashierResponseYearSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierYearlySale(cashier))
	}

	return cashierRecords
}

func (s *cashierResponseMapper) ToCashierMonthlyTotalSale(c *record.CashierRecordMonthTotalSales) *response.CashierResponseMonthTotalSales {
	return &response.CashierResponseMonthTotalSales{
		Year:       c.Year,
		Month:      c.Month,
		TotalSales: int(c.TotalSales),
	}
}

func (s *cashierResponseMapper) ToCashierMonthlyTotalSales(c []*record.CashierRecordMonthTotalSales) []*response.CashierResponseMonthTotalSales {
	var cashierRecords []*response.CashierResponseMonthTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierMonthlyTotalSale(cashier))
	}

	return cashierRecords
}

func (s *cashierResponseMapper) ToCashierYearlyTotalSale(c *record.CashierRecordYearTotalSales) *response.CashierResponseYearTotalSales {
	return &response.CashierResponseYearTotalSales{
		Year:       c.Year,
		TotalSales: int(c.TotalSales),
	}
}

func (s *cashierResponseMapper) ToCashierYearlyTotalSales(c []*record.CashierRecordYearTotalSales) []*response.CashierResponseYearTotalSales {
	var cashierRecords []*response.CashierResponseYearTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierYearlyTotalSale(cashier))
	}

	return cashierRecords
}
