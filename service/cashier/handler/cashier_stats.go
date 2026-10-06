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

type cashierStatsHandleGrpc struct {
	*cashierStatsCoreHandleGrpc
	*cashierStatsByIdHandleGrpc
	*cashierStatsByMerchantHandleGrpc
}

func NewCashierStatsHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbcashier.CashierStatsServiceServer {
	return &cashierStatsHandleGrpc{
		cashierStatsCoreHandleGrpc:       NewCashierStatsCoreHandleGrpc(service, logger),
		cashierStatsByIdHandleGrpc:       NewCashierStatsByIdHandleGrpc(service, logger),
		cashierStatsByMerchantHandleGrpc: NewCashierStatsByMerchantHandleGrpc(service, logger),
	}
}

type cashierStatsCoreHandleGrpc struct {
	pbcashier.UnimplementedCashierStatsServiceServer
	cashierStats service.CashierStatsService
	logger       logger.LoggerInterface
}

func NewCashierStatsCoreHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) *cashierStatsCoreHandleGrpc {
	return &cashierStatsCoreHandleGrpc{
		cashierStats: service.CashierStats,
		logger:       logger,
	}
}

func (s *cashierStatsCoreHandleGrpc) FindMonthlyTotalSales(ctx context.Context, req *pbcashier.FindYearMonthTotalSales) (*pbcashier.ApiResponseCashierMonthlyTotalSales, error) {
	s.logger.Info("FindMonthlyTotalSales cashier called", zap.Int32("year", req.GetYear()), zap.Int32("month", req.GetMonth()))

	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMonth
	}

	reqService := requests.MonthTotalSales{
		Year:  year,
		Month: month,
	}

	methods, err := s.cashierStats.FindMonthlyTotalSales(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalSales cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalSales cashier success")

	return &pbcashier.ApiResponseCashierMonthlyTotalSales{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCashierMonthlyTotalSales(methods),
	}, nil
}

func (s *cashierStatsCoreHandleGrpc) FindYearlyTotalSales(ctx context.Context, req *pbcashier.FindYearTotalSales) (*pbcashier.ApiResponseCashierYearlyTotalSales, error) {
	s.logger.Info("FindYearlyTotalSales cashier called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.cashierStats.FindYearlyTotalSales(ctx, year)
	if err != nil {
		s.logger.Error("FindYearlyTotalSales cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalSales cashier success")

	return &pbcashier.ApiResponseCashierYearlyTotalSales{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCashierYearlyTotalSales(methods),
	}, nil
}

func (s *cashierStatsCoreHandleGrpc) FindMonthSales(ctx context.Context, req *pbcashier.FindYearCashier) (*pbcashier.ApiResponseCashierMonthSales, error) {
	s.logger.Info("FindMonthSales cashier called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.cashierStats.FindMonthlySales(ctx, year)
	if err != nil {
		s.logger.Error("FindMonthSales cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthSales cashier success")

	return &pbcashier.ApiResponseCashierMonthSales{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponsesCashierMonthlySales(methods),
	}, nil
}

func (s *cashierStatsCoreHandleGrpc) FindYearSales(ctx context.Context, req *pbcashier.FindYearCashier) (*pbcashier.ApiResponseCashierYearSales, error) {
	s.logger.Info("FindYearSales cashier called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.cashierStats.FindYearlySales(ctx, year)
	if err != nil {
		s.logger.Error("FindYearSales cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearSales cashier success")

	return &pbcashier.ApiResponseCashierYearSales{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponsesCashierYearlySales(methods),
	}, nil
}
