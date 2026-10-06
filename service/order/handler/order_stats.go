package handler

import (
	"context"

	"github.com/MamangRust/monolith-pointofsale-order/service"
	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/order_errors"
	"go.uber.org/zap"
)

type orderStatsHandleGrpc struct {
	*orderStatsCoreHandleGrpc
	*orderStatsByMerchantHandleGrpc
}

func NewOrderStatsHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pborder.OrderStatsServiceServer {
	return &orderStatsHandleGrpc{
		orderStatsCoreHandleGrpc:       NewOrderStatsCoreHandleGrpc(service, logger),
		orderStatsByMerchantHandleGrpc: NewOrderStatsByMerchantHandleGrpc(service, logger),
	}
}

type orderStatsCoreHandleGrpc struct {
	pborder.UnimplementedOrderStatsServiceServer
	orderStats service.OrderStatsService
	logger     logger.LoggerInterface
}

func NewOrderStatsCoreHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) *orderStatsCoreHandleGrpc {
	return &orderStatsCoreHandleGrpc{
		orderStats: service.OrderStats,
		logger:     logger,
	}
}

func (s *orderStatsCoreHandleGrpc) FindMonthlyTotalRevenue(ctx context.Context, req *pborder.FindYearMonthTotalRevenue) (*pborder.ApiResponseOrderMonthlyTotalRevenue, error) {
	s.logger.Info("FindMonthlyTotalRevenue orders called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, order_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthTotalRevenue{
		Year:  year,
		Month: month,
	}

	methods, err := s.orderStats.FindMonthlyTotalRevenue(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalRevenue orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalRevenue orders success")

	return &pborder.ApiResponseOrderMonthlyTotalRevenue{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseOrderMonthlyTotalRevenues(methods),
	}, nil
}

func (s *orderStatsCoreHandleGrpc) FindYearlyTotalRevenue(ctx context.Context, req *pborder.FindYearTotalRevenue) (*pborder.ApiResponseOrderYearlyTotalRevenue, error) {
	s.logger.Info("FindYearlyTotalRevenue orders called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}

	methods, err := s.orderStats.FindYearlyTotalRevenue(ctx, year)
	if err != nil {
		s.logger.Error("FindYearlyTotalRevenue orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalRevenue orders success")

	return &pborder.ApiResponseOrderYearlyTotalRevenue{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseOrderYearlyTotalRevenues(methods),
	}, nil
}

func (s *orderStatsCoreHandleGrpc) FindMonthlyRevenue(ctx context.Context, request *pborder.FindYearOrder) (*pborder.ApiResponseOrderMonthly, error) {
	s.logger.Info("FindMonthlyRevenue orders called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	if year <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	res, err := s.orderStats.FindMonthlyOrder(ctx, year)
	if err != nil {
		s.logger.Error("FindMonthlyRevenue orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyRevenue orders success")

	return &pborder.ApiResponseOrderMonthly{
		Status:  "success",
		Message: "Monthly revenue data retrieved",
		Data:    mapResponsesOrderMonthlyPrices(res),
	}, nil
}

func (s *orderStatsCoreHandleGrpc) FindYearlyRevenue(ctx context.Context, request *pborder.FindYearOrder) (*pborder.ApiResponseOrderYearly, error) {
	s.logger.Info("FindYearlyRevenue orders called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	if year <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	res, err := s.orderStats.FindYearlyOrder(ctx, year)
	if err != nil {
		s.logger.Error("FindYearlyRevenue orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyRevenue orders success")

	return &pborder.ApiResponseOrderYearly{
		Status:  "success",
		Message: "Yearly revenue data retrieved",
		Data:    mapResponsesOrderYearlyPrices(res),
	}, nil
}
