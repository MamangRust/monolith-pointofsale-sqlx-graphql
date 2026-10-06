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

type categoryStatsByIdHandleGrpc struct {
	pbcategory.UnimplementedCategoryStatsByIdServiceServer
	categoryStatsById service.CategoryStatsByIdService
	logger            logger.LoggerInterface
}

func NewCategoryStatsByIdHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) *categoryStatsByIdHandleGrpc {
	return &categoryStatsByIdHandleGrpc{
		categoryStatsById: service.CategoryStatsById,
		logger:            logger,
	}
}

func (s *categoryStatsByIdHandleGrpc) FindMonthlyTotalPricesById(ctx context.Context, req *pbcategory.FindYearMonthTotalPriceById) (*pbcategory.ApiResponseCategoryMonthlyTotalPrice, error) {
	s.logger.Info("FindMonthlyTotalPricesById categories called", zap.Int32("id", req.GetCategoryId()))

	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcFailedInvalidMonth
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.MonthTotalPriceCategory{
		Year:       year,
		Month:      month,
		CategoryID: id,
	}

	methods, err := s.categoryStatsById.FindMonthlyTotalPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalPricesById categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalPricesById categories success")

	return &pbcategory.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCategoryMonthlyTotalPricesById(methods),
	}, nil
}

func (s *categoryStatsByIdHandleGrpc) FindYearlyTotalPricesById(ctx context.Context, req *pbcategory.FindYearTotalPriceById) (*pbcategory.ApiResponseCategoryYearlyTotalPrice, error) {
	s.logger.Info("FindYearlyTotalPricesById categories called", zap.Int32("id", req.GetCategoryId()))

	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearTotalPriceCategory{
		Year:       year,
		CategoryID: id,
	}

	methods, err := s.categoryStatsById.FindYearlyTotalPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTotalPricesById categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalPricesById categories success")

	return &pbcategory.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCategoryYearlyTotalPricesById(methods),
	}, nil
}

func (s *categoryStatsByIdHandleGrpc) FindMonthPriceById(ctx context.Context, req *pbcategory.FindYearCategoryById) (*pbcategory.ApiResponseCategoryMonthPrice, error) {
	s.logger.Info("FindMonthPriceById categories called", zap.Int32("id", req.GetCategoryId()))

	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.MonthPriceId{
		Year:       year,
		CategoryID: id,
	}

	methods, err := s.categoryStatsById.FindMonthPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthPriceById categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthPriceById categories success")

	return &pbcategory.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Merchant monthly payment methods retrieved successfully",
		Data:    mapResponsesCategoryMonthlyPricesById(methods),
	}, nil
}

func (s *categoryStatsByIdHandleGrpc) FindYearPriceById(ctx context.Context, req *pbcategory.FindYearCategoryById) (*pbcategory.ApiResponseCategoryYearPrice, error) {
	s.logger.Info("FindYearPriceById categories called", zap.Int32("id", req.GetCategoryId()))

	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearPriceId{
		Year:       year,
		CategoryID: id,
	}

	methods, err := s.categoryStatsById.FindYearPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearPriceById categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearPriceById categories success")

	return &pbcategory.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    mapResponsesCategoryYearlyPricesById(methods),
	}, nil
}
