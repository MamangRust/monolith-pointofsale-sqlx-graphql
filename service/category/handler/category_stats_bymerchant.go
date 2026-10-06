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

type categoryStatsByMerchantHandleGrpc struct {
	pbcategory.UnimplementedCategoryStatsByMerchantServiceServer
	categoryStatsByMerchant service.CategoryStatsByMerchantService
	logger                  logger.LoggerInterface
}

func NewCategoryStatsByMerchantHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) *categoryStatsByMerchantHandleGrpc {
	return &categoryStatsByMerchantHandleGrpc{
		categoryStatsByMerchant: service.CategoryStatsByMerchant,
		logger:                  logger,
	}
}

func (s *categoryStatsByMerchantHandleGrpc) FindMonthlyTotalPricesByMerchant(ctx context.Context, req *pbcategory.FindYearMonthTotalPriceByMerchant) (*pbcategory.ApiResponseCategoryMonthlyTotalPrice, error) {
	s.logger.Info("FindMonthlyTotalPricesByMerchant categories called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcFailedInvalidMonth
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthTotalPriceMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	methods, err := s.categoryStatsByMerchant.FindMonthlyTotalPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalPricesByMerchant categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalPricesByMerchant categories success")

	return &pbcategory.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCategoryMonthlyTotalPricesByMerchant(methods),
	}, nil
}

func (s *categoryStatsByMerchantHandleGrpc) FindYearlyTotalPricesByMerchant(ctx context.Context, req *pbcategory.FindYearTotalPriceByMerchant) (*pbcategory.ApiResponseCategoryYearlyTotalPrice, error) {
	s.logger.Info("FindYearlyTotalPricesByMerchant categories called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.YearTotalPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.categoryStatsByMerchant.FindYearlyTotalPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTotalPricesByMerchant categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalPricesByMerchant categories success")

	return &pbcategory.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCategoryYearlyTotalPricesByMerchant(methods),
	}, nil
}

func (s *categoryStatsByMerchantHandleGrpc) FindMonthPriceByMerchant(ctx context.Context, req *pbcategory.FindYearCategoryByMerchant) (*pbcategory.ApiResponseCategoryMonthPrice, error) {
	s.logger.Info("FindMonthPriceByMerchant categories called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.categoryStatsByMerchant.FindMonthPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthPriceByMerchant categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthPriceByMerchant categories success")

	return &pbcategory.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Merchant monthly payment methods retrieved successfully",
		Data:    mapResponsesCategoryMonthlyPricesByMerchant(methods),
	}, nil
}

func (s *categoryStatsByMerchantHandleGrpc) FindYearPriceByMerchant(ctx context.Context, req *pbcategory.FindYearCategoryByMerchant) (*pbcategory.ApiResponseCategoryYearPrice, error) {
	s.logger.Info("FindYearPriceByMerchant categories called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.YearPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.categoryStatsByMerchant.FindYearPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearPriceByMerchant categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearPriceByMerchant categories success")

	return &pbcategory.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    mapResponsesCategoryYearlyPricesByMerchant(methods),
	}, nil
}
