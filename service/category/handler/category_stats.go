package handler

import (
	"context"

	"github.com/MamangRust/monolith-pointofsale-category/service"
	pbcategory "github.com/MamangRust/monolith-pointofsale-pb/category"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/category_errors"
	"go.uber.org/zap"
)

type categoryStatsHandleGrpc struct {
	*categoryStatsCoreHandleGrpc
	*categoryStatsByIdHandleGrpc
	*categoryStatsByMerchantHandleGrpc
}

func NewCategoryStatsHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbcategory.CategoryStatsServiceServer {
	return &categoryStatsHandleGrpc{
		categoryStatsCoreHandleGrpc:       NewCategoryStatsCoreHandleGrpc(service, logger),
		categoryStatsByIdHandleGrpc:       NewCategoryStatsByIdHandleGrpc(service, logger),
		categoryStatsByMerchantHandleGrpc: NewCategoryStatsByMerchantHandleGrpc(service, logger),
	}
}

type categoryStatsCoreHandleGrpc struct {
	pbcategory.UnimplementedCategoryStatsServiceServer
	categoryStats service.CategoryStatsService
	logger        logger.LoggerInterface
}

func NewCategoryStatsCoreHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) *categoryStatsCoreHandleGrpc {
	return &categoryStatsCoreHandleGrpc{
		categoryStats: service.CategoryStats,
		logger:        logger,
	}
}

func (s *categoryStatsCoreHandleGrpc) FindMonthlyTotalPrices(ctx context.Context, req *pbcategory.FindYearMonthTotalPrices) (*pbcategory.ApiResponseCategoryMonthlyTotalPrice, error) {
	s.logger.Info("FindMonthlyTotalPrices categories called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcFailedInvalidMonth
	}

	reqService := requests.MonthTotalPrice{
		Year:  year,
		Month: month,
	}

	methods, err := s.categoryStats.FindMonthlyTotalPrice(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalPrices categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalPrices categories success")

	return &pbcategory.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCategoryMonthlyTotalPrices(methods),
	}, nil
}

func (s *categoryStatsCoreHandleGrpc) FindYearlyTotalPrices(ctx context.Context, req *pbcategory.FindYearTotalPrices) (*pbcategory.ApiResponseCategoryYearlyTotalPrice, error) {
	s.logger.Info("FindYearlyTotalPrices categories called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.categoryStats.FindYearlyTotalPrice(ctx, year)
	if err != nil {
		s.logger.Error("FindYearlyTotalPrices categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalPrices categories success")

	return &pbcategory.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCategoryYearlyTotalPrices(methods),
	}, nil
}

func (s *categoryStatsCoreHandleGrpc) FindMonthPrice(ctx context.Context, req *pbcategory.FindYearCategory) (*pbcategory.ApiResponseCategoryMonthPrice, error) {
	s.logger.Info("FindMonthPrice categories called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.categoryStats.FindMonthPrice(ctx, year)
	if err != nil {
		s.logger.Error("FindMonthPrice categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthPrice categories success")

	return &pbcategory.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Monthly payment methods retrieved successfully",
		Data:    mapResponsesCategoryMonthlyPrices(methods),
	}, nil
}

func (s *categoryStatsCoreHandleGrpc) FindYearPrice(ctx context.Context, req *pbcategory.FindYearCategory) (*pbcategory.ApiResponseCategoryYearPrice, error) {
	s.logger.Info("FindYearPrice categories called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.categoryStats.FindYearPrice(ctx, year)
	if err != nil {
		s.logger.Error("FindYearPrice categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearPrice categories success")

	return &pbcategory.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponsesCategoryYearlyPrices(methods),
	}, nil
}
