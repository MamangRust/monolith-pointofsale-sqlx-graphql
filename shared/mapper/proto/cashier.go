package protomapper

import (
	pbcashier "github.com/MamangRust/monolith-pointofsale-pb/cashier"
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type cashierProtoMapper struct {
}

func NewCashierProtoMapper() *cashierProtoMapper {
	return &cashierProtoMapper{}
}

func (c *cashierProtoMapper) ToProtoResponseCashier(status string, message string, pbResponse *response.CashierResponse) *pbcashier.ApiResponseCashier {
	return &pbcashier.ApiResponseCashier{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCashier(pbResponse),
	}
}

func (c *cashierProtoMapper) ToProtoResponsesCashier(status string, message string, pbResponse []*response.CashierResponse) *pbcashier.ApiResponsesCashier {
	return &pbcashier.ApiResponsesCashier{
		Status:  status,
		Message: message,
		Data:    c.mapResponsesCashier(pbResponse),
	}
}

func (c *cashierProtoMapper) ToProtoResponseCashierDeleteAt(status string, message string, pbResponse *response.CashierResponseDeleteAt) *pbcashier.ApiResponseCashierDeleteAt {
	return &pbcashier.ApiResponseCashierDeleteAt{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCashierDeleteAt(pbResponse),
	}
}

func (c *cashierProtoMapper) ToProtoResponseCashierDelete(status string, message string) *pbcashier.ApiResponseCashierDelete {
	return &pbcashier.ApiResponseCashierDelete{
		Status:  status,
		Message: message,
	}
}

func (u *cashierProtoMapper) ToProtoResponseCashierAll(status string, message string) *pbcashier.ApiResponseCashierAll {
	return &pbcashier.ApiResponseCashierAll{
		Status:  status,
		Message: message,
	}
}

func (u *cashierProtoMapper) ToProtoResponsePaginationCashierDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.CashierResponseDeleteAt) *pbcashier.ApiResponsePaginationCashierDeleteAt {
	return &pbcashier.ApiResponsePaginationCashierDeleteAt{
		Status:     status,
		Message:    message,
		Data:       u.mapResponsesCashierDeleteAt(users),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (u *cashierProtoMapper) ToProtoResponsePaginationCashier(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.CashierResponse) *pbcashier.ApiResponsePaginationCashier {
	return &pbcashier.ApiResponsePaginationCashier{
		Status:     status,
		Message:    message,
		Data:       u.mapResponsesCashier(users),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (u *cashierProtoMapper) ToProtoResponseMonthlyTotalSales(status, message string, row []*response.CashierResponseMonthSales) *pbcashier.ApiResponseCashierMonthSales {
	return &pbcashier.ApiResponseCashierMonthSales{
		Status:  status,
		Message: message,
		Data:    u.mapResponsesCashierMonthlySales(row),
	}
}

func (u *cashierProtoMapper) ToProtoResponseYearlyTotalSales(status, message string, row []*response.CashierResponseYearSales) *pbcashier.ApiResponseCashierYearSales {
	return &pbcashier.ApiResponseCashierYearSales{
		Status:  status,
		Message: message,
		Data:    u.mapResponsesCashierYearlySales(row),
	}
}

func (u *cashierProtoMapper) ToProtoMonthlyTotalSales(status, message string, row []*response.CashierResponseMonthTotalSales) *pbcashier.ApiResponseCashierMonthlyTotalSales {
	return &pbcashier.ApiResponseCashierMonthlyTotalSales{
		Status:  status,
		Message: message,
		Data:    u.mapResponseCashierMonthlyTotalSales(row),
	}
}

func (u *cashierProtoMapper) ToProtoYearlyTotalSales(status, message string, row []*response.CashierResponseYearTotalSales) *pbcashier.ApiResponseCashierYearlyTotalSales {
	return &pbcashier.ApiResponseCashierYearlyTotalSales{
		Status:  status,
		Message: message,
		Data:    u.mapResponseCashierYearlyTotalSales(row),
	}
}

func (c *cashierProtoMapper) mapResponseCashier(cashier *response.CashierResponse) *pbcashier.CashierResponse {
	return &pbcashier.CashierResponse{
		Id:         int32(cashier.ID),
		MerchantId: int32(cashier.MerchantID),
		Name:       cashier.Name,
		CreatedAt:  cashier.CreatedAt,
		UpdatedAt:  cashier.UpdatedAt,
	}
}

func (c *cashierProtoMapper) mapResponsesCashier(cashiers []*response.CashierResponse) []*pbcashier.CashierResponse {
	var mappedCashiers []*pbcashier.CashierResponse

	for _, cashier := range cashiers {
		mappedCashiers = append(mappedCashiers, c.mapResponseCashier(cashier))
	}

	return mappedCashiers
}

func (c *cashierProtoMapper) mapResponseCashierDeleteAt(cashier *response.CashierResponseDeleteAt) *pbcashier.CashierResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if cashier.DeletedAt != nil {
		deletedAt = wrapperspb.String(*cashier.DeletedAt)
	}

	return &pbcashier.CashierResponseDeleteAt{
		Id:         int32(cashier.ID),
		MerchantId: int32(cashier.MerchantID),
		Name:       cashier.Name,
		CreatedAt:  cashier.CreatedAt,
		UpdatedAt:  cashier.UpdatedAt,
		DeletedAt:  deletedAt,
	}
}

func (c *cashierProtoMapper) mapResponsesCashierDeleteAt(cashiers []*response.CashierResponseDeleteAt) []*pbcashier.CashierResponseDeleteAt {
	var mappedCashiers []*pbcashier.CashierResponseDeleteAt

	for _, cashier := range cashiers {
		mappedCashiers = append(mappedCashiers, c.mapResponseCashierDeleteAt(cashier))
	}

	return mappedCashiers
}

func (s *cashierProtoMapper) mapResponseCashierMonthlySale(cashier *response.CashierResponseMonthSales) *pbcashier.CashierResponseMonthSales {
	return &pbcashier.CashierResponseMonthSales{
		Month:       cashier.Month,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func (s *cashierProtoMapper) mapResponsesCashierMonthlySales(c []*response.CashierResponseMonthSales) []*pbcashier.CashierResponseMonthSales {
	var cashierRecords []*pbcashier.CashierResponseMonthSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.mapResponseCashierMonthlySale(cashier))
	}

	return cashierRecords
}

func (s *cashierProtoMapper) mapResponseCashierYearlySale(cashier *response.CashierResponseYearSales) *pbcashier.CashierResponseYearSales {
	return &pbcashier.CashierResponseYearSales{
		Year:        cashier.Year,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func (s *cashierProtoMapper) mapResponsesCashierYearlySales(c []*response.CashierResponseYearSales) []*pbcashier.CashierResponseYearSales {
	var cashierRecords []*pbcashier.CashierResponseYearSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.mapResponseCashierYearlySale(cashier))
	}

	return cashierRecords
}

func (s *cashierProtoMapper) mapResponseCashierMonthlyTotalSale(c *response.CashierResponseMonthTotalSales) *pbcashier.CashierResponseMonthTotalSales {
	return &pbcashier.CashierResponseMonthTotalSales{
		Year:       c.Year,
		Month:      c.Month,
		TotalSales: int32(c.TotalSales),
	}
}

func (s *cashierProtoMapper) mapResponseCashierMonthlyTotalSales(c []*response.CashierResponseMonthTotalSales) []*pbcashier.CashierResponseMonthTotalSales {
	var cashierRecords []*pbcashier.CashierResponseMonthTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.mapResponseCashierMonthlyTotalSale(cashier))
	}

	return cashierRecords
}

func (s *cashierProtoMapper) mapResponseCashierYearlyTotalSale(c *response.CashierResponseYearTotalSales) *pbcashier.CashierResponseYearTotalSales {
	return &pbcashier.CashierResponseYearTotalSales{
		Year:       c.Year,
		TotalSales: int32(c.TotalSales),
	}
}

func (s *cashierProtoMapper) mapResponseCashierYearlyTotalSales(c []*response.CashierResponseYearTotalSales) []*pbcashier.CashierResponseYearTotalSales {
	var cashierRecords []*pbcashier.CashierResponseYearTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.mapResponseCashierYearlyTotalSale(cashier))
	}

	return cashierRecords
}
