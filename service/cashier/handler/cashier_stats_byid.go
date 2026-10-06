package handler

import (
	"context"

	"github.com/MamangRust/monolith-pointofsale-cashier/service"
	pbcashier "github.com/MamangRust/monolith-pointofsale-pb/cashier"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/cashier_errors"
	"go.uber.org/zap"
)

type cashierStatsByIdHandleGrpc struct {
	pbcashier.UnimplementedCashierStatsByIdServiceServer
	cashierStatsById service.CashierStatsByIdService
	logger           logger.LoggerInterface
}

func NewCashierStatsByIdHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) *cashierStatsByIdHandleGrpc {
	return &cashierStatsByIdHandleGrpc{
		cashierStatsById: service.CashierStatsById,
		logger:           logger,
	}
}

func (s *cashierStatsByIdHandleGrpc) FindMonthlyTotalSalesById(ctx context.Context, req *pbcashier.FindYearMonthTotalSalesById) (*pbcashier.ApiResponseCashierMonthlyTotalSales, error) {
	s.logger.Info("FindMonthlyTotalSalesById cashier called", zap.Int32("id", req.GetCashierId()))

	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetCashierId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMonth
	}
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.MonthTotalSalesCashier{
		Year:      year,
		Month:     month,
		CashierID: id,
	}

	methods, err := s.cashierStatsById.FindMonthlyTotalSalesById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalSalesById cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalSalesById cashier success")

	return &pbcashier.ApiResponseCashierMonthlyTotalSales{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCashierMonthlyTotalSalesById(methods),
	}, nil
}

func (s *cashierStatsByIdHandleGrpc) FindYearlyTotalSalesById(ctx context.Context, req *pbcashier.FindYearTotalSalesById) (*pbcashier.ApiResponseCashierYearlyTotalSales, error) {
	s.logger.Info("FindYearlyTotalSalesById cashier called", zap.Int32("id", req.GetCashierId()))

	year := int(req.GetYear())
	id := int(req.GetCashierId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearTotalSalesCashier{
		Year:      year,
		CashierID: id,
	}

	methods, err := s.cashierStatsById.FindYearlyTotalSalesById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTotalSalesById cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalSalesById cashier success")

	return &pbcashier.ApiResponseCashierYearlyTotalSales{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCashierYearlyTotalSalesById(methods),
	}, nil
}

func (s *cashierStatsByIdHandleGrpc) FindMonthSalesById(ctx context.Context, req *pbcashier.FindYearCashierById) (*pbcashier.ApiResponseCashierMonthSales, error) {
	s.logger.Info("FindMonthSalesById cashier called", zap.Int32("id", req.GetCashierId()))

	year := int(req.GetYear())
	cashierId := int(req.GetCashierId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if cashierId <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.MonthCashierId{
		Year:      year,
		CashierID: cashierId,
	}

	methods, err := s.cashierStatsById.FindMonthlyCashierById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthSalesById cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthSalesById cashier success")

	return &pbcashier.ApiResponseCashierMonthSales{
		Status:  "success",
		Message: "Cashier monthly sales retrieved successfully",
		Data:    mapResponsesCashierMonthlySalesById(methods),
	}, nil
}

func (s *cashierStatsByIdHandleGrpc) FindYearSalesById(ctx context.Context, req *pbcashier.FindYearCashierById) (*pbcashier.ApiResponseCashierYearSales, error) {
	s.logger.Info("FindYearSalesById cashier called", zap.Int32("id", req.GetCashierId()))

	year := int(req.GetYear())
	cashierId := int(req.GetCashierId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if cashierId <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearCashierId{
		Year:      year,
		CashierID: cashierId,
	}

	methods, err := s.cashierStatsById.FindYearlyCashierById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearSalesById cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearSalesById cashier success")

	return &pbcashier.ApiResponseCashierYearSales{
		Status:  "success",
		Message: "Cashier yearly sales retrieved successfully",
		Data:    mapResponsesCashierYearlySalesById(methods),
	}, nil
}
