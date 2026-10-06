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

type cashierStatsByMerchantHandleGrpc struct {
	pbcashier.UnimplementedCashierStatsByMerchantServiceServer
	cashierStatsByMerchant service.CashierStatsByMerchant
	logger                 logger.LoggerInterface
}

func NewCashierStatsByMerchantHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) *cashierStatsByMerchantHandleGrpc {
	return &cashierStatsByMerchantHandleGrpc{
		cashierStatsByMerchant: service.CashierStatsByMerchant,
		logger:                 logger,
	}
}

func (s *cashierStatsByMerchantHandleGrpc) FindMonthlyTotalSalesByMerchant(ctx context.Context, req *pbcashier.FindYearMonthTotalSalesByMerchant) (*pbcashier.ApiResponseCashierMonthlyTotalSales, error) {
	s.logger.Info("FindMonthlyTotalSalesByMerchant cashier called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	month := int(req.GetMonth())
	merchantId := int(req.GetMerchantId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMonth
	}
	if merchantId <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthTotalSalesMerchant{
		Year:       year,
		Month:      month,
		MerchantID: merchantId,
	}

	methods, err := s.cashierStatsByMerchant.FindMonthlyTotalSalesByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalSalesByMerchant cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalSalesByMerchant cashier success")

	return &pbcashier.ApiResponseCashierMonthlyTotalSales{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCashierMonthlyTotalSalesByMerchant(methods),
	}, nil
}

func (s *cashierStatsByMerchantHandleGrpc) FindYearlyTotalSalesByMerchant(ctx context.Context, req *pbcashier.FindYearTotalSalesByMerchant) (*pbcashier.ApiResponseCashierYearlyTotalSales, error) {
	s.logger.Info("FindYearlyTotalSalesByMerchant cashier called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	merchantId := int(req.GetMerchantId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if merchantId <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.YearTotalSalesMerchant{
		Year:       year,
		MerchantID: merchantId,
	}

	methods, err := s.cashierStatsByMerchant.FindYearlyTotalSalesByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTotalSalesByMerchant cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalSalesByMerchant cashier success")

	return &pbcashier.ApiResponseCashierYearlyTotalSales{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCashierYearlyTotalSalesByMerchant(methods),
	}, nil
}

func (s *cashierStatsByMerchantHandleGrpc) FindMonthSalesByMerchant(ctx context.Context, req *pbcashier.FindYearCashierByMerchant) (*pbcashier.ApiResponseCashierMonthSales, error) {
	s.logger.Info("FindMonthSalesByMerchant cashier called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	merchantId := int(req.GetMerchantId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if merchantId <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthCashierMerchant{
		Year:       year,
		MerchantID: merchantId,
	}

	methods, err := s.cashierStatsByMerchant.FindMonthlyCashierByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthSalesByMerchant cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthSalesByMerchant cashier success")

	return &pbcashier.ApiResponseCashierMonthSales{
		Status:  "success",
		Message: "Merchant monthly revenue retrieved successfully",
		Data:    mapResponsesCashierMonthlySalesByMerchant(methods),
	}, nil
}

func (s *cashierStatsByMerchantHandleGrpc) FindYearSalesByMerchant(ctx context.Context, req *pbcashier.FindYearCashierByMerchant) (*pbcashier.ApiResponseCashierYearSales, error) {
	s.logger.Info("FindYearSalesByMerchant cashier called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	merchantId := int(req.GetMerchantId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if merchantId <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.YearCashierMerchant{
		Year:       year,
		MerchantID: merchantId,
	}

	methods, err := s.cashierStatsByMerchant.FindYearlyCashierByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearSalesByMerchant cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearSalesByMerchant cashier success")

	return &pbcashier.ApiResponseCashierYearSales{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    mapResponsesCashierYearlySalesByMerchant(methods),
	}, nil
}
