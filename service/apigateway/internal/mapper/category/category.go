package categorygraphqlmapper

import (
	graphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper"
	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-pointofsale-pb/category"
)

type categoryGraphqlMapper struct {
}

func NewCategoryGraphqlMapper() *categoryGraphqlMapper {
	return &categoryGraphqlMapper{}
}

func (c *categoryGraphqlMapper) ToGraphqlResponseCategory(res *pb.ApiResponseCategory) *model.APIResponseCategory {
	return &model.APIResponseCategory{
		Status:  res.Status,
		Message: res.Message,
		Data:    c.mapResponseCategory(res.Data),
	}
}

func (c *categoryGraphqlMapper) ToGraphqlResponsesCategory(res *pb.ApiResponsesCategory) *model.APIResponsesCategory {
	return &model.APIResponsesCategory{
		Status:  res.Status,
		Message: res.Message,
		Data:    c.mapResponsesCategory(res.Data),
	}
}

func (c *categoryGraphqlMapper) ToGraphqlResponseCategoryDeleteAt(res *pb.ApiResponseCategoryDeleteAt) *model.APIResponseCategoryDeleteAt {
	return &model.APIResponseCategoryDeleteAt{
		Status:  res.Status,
		Message: res.Message,
		Data:    c.mapResponseCategoryDeleteAt(res.Data),
	}
}

func (c *categoryGraphqlMapper) ToGraphqlResponseCategoryDelete(res *pb.ApiResponseCategoryDelete) *model.APIResponseCategoryDelete {
	return &model.APIResponseCategoryDelete{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (c *categoryGraphqlMapper) ToGraphqlResponseCategoryAll(res *pb.ApiResponseCategoryAll) *model.APIResponseCategoryAll {
	return &model.APIResponseCategoryAll{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (c *categoryGraphqlMapper) ToGraphqlResponsePaginationCategory(res *pb.ApiResponsePaginationCategory) *model.APIResponsePaginationCategory {
	return &model.APIResponsePaginationCategory{
		Status:     res.Status,
		Message:    res.Message,
		Data:       c.mapResponsesCategory(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (c *categoryGraphqlMapper) ToGraphqlResponsePaginationCategoryDeleteAt(res *pb.ApiResponsePaginationCategoryDeleteAt) *model.APIResponsePaginationCategoryDeleteAt {
	return &model.APIResponsePaginationCategoryDeleteAt{
		Status:     res.Status,
		Message:    res.Message,
		Data:       c.mapResponsesCategoryDeleteAt(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (c *categoryGraphqlMapper) mapResponseCategory(category *pb.CategoryResponse) *model.CategoryResponse {
	if category == nil {
		return nil
	}
	return &model.CategoryResponse{
		ID:            int32(category.Id),
		Name:          category.Name,
		Description:   &category.Description,
		SlugCategory:  &category.SlugCategory,
		ImageCategory: &category.ImageCategory,
		CreatedAt:     &category.CreatedAt,
		UpdatedAt:     &category.UpdatedAt,
	}
}

func (c *categoryGraphqlMapper) mapResponsesCategory(categories []*pb.CategoryResponse) []*model.CategoryResponse {
	var responses []*model.CategoryResponse
	for _, category := range categories {
		responses = append(responses, c.mapResponseCategory(category))
	}
	return responses
}

func (c *categoryGraphqlMapper) mapResponseCategoryDeleteAt(category *pb.CategoryResponseDeleteAt) *model.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	var deletedAt *string
	if category.DeletedAt != nil {
		deletedAt = &category.DeletedAt.Value
	}

	return &model.CategoryResponseDeleteAt{
		ID:            int32(category.Id),
		Name:          category.Name,
		Description:   &category.Description,
		SlugCategory:  &category.SlugCategory,
		ImageCategory: &category.ImageCategory,
		CreatedAt:     &category.CreatedAt,
		UpdatedAt:     &category.UpdatedAt,
		DeletedAt:     deletedAt,
	}
}

func (c *categoryGraphqlMapper) mapResponsesCategoryDeleteAt(categories []*pb.CategoryResponseDeleteAt) []*model.CategoryResponseDeleteAt {
	var responses []*model.CategoryResponseDeleteAt
	for _, category := range categories {
		responses = append(responses, c.mapResponseCategoryDeleteAt(category))
	}
	return responses
}

func (c *categoryGraphqlMapper) ToGraphqlResponseCategoryMonthlyTotalPrice(res *pb.ApiResponseCategoryMonthlyTotalPrice) *model.APIResponseCategoryMonthlyTotalPrice {
	return &model.APIResponseCategoryMonthlyTotalPrice{
		Status:  res.Status,
		Message: res.Message,
		Data:    c.mapResponsesCategoryMonthlyTotalPrice(res.Data),
	}
}

func (c *categoryGraphqlMapper) ToGraphqlResponseCategoryYearlyTotalPrice(res *pb.ApiResponseCategoryYearlyTotalPrice) *model.APIResponseCategoryYearlyTotalPrice {
	return &model.APIResponseCategoryYearlyTotalPrice{
		Status:  res.Status,
		Message: res.Message,
		Data:    c.mapResponsesCategoryYearlyTotalPrice(res.Data),
	}
}

func (c *categoryGraphqlMapper) ToGraphqlResponseCategoryMonthlyPrice(res *pb.ApiResponseCategoryMonthPrice) *model.APIResponseCategoryMonthPrice {
	return &model.APIResponseCategoryMonthPrice{
		Status:  res.Status,
		Message: res.Message,
		Data:    c.mapResponsesCategoryMonthPrice(res.Data),
	}
}

func (c *categoryGraphqlMapper) ToGraphqlResponseCategoryYearlyPrice(res *pb.ApiResponseCategoryYearPrice) *model.APIResponseCategoryYearPrice {
	return &model.APIResponseCategoryYearPrice{
		Status:  res.Status,
		Message: res.Message,
		Data:    c.mapResponsesCategoryYearPrice(res.Data),
	}
}

func (c *categoryGraphqlMapper) mapResponsesCategoryMonthlyTotalPrice(data []*pb.CategoriesMonthlyTotalPriceResponse) []*model.CategoriesMonthlyTotalPriceResponse {
	var responses []*model.CategoriesMonthlyTotalPriceResponse
	for _, item := range data {
		if item == nil {
			continue
		}
		responses = append(responses, &model.CategoriesMonthlyTotalPriceResponse{
			Year:         item.Year,
			Month:        item.Month,
			TotalRevenue: int32(item.TotalRevenue),
		})
	}
	return responses
}

func (c *categoryGraphqlMapper) mapResponsesCategoryYearlyTotalPrice(data []*pb.CategoriesYearlyTotalPriceResponse) []*model.CategoriesYearlyTotalPriceResponse {
	var responses []*model.CategoriesYearlyTotalPriceResponse
	for _, item := range data {
		if item == nil {
			continue
		}
		responses = append(responses, &model.CategoriesYearlyTotalPriceResponse{
			Year:         item.Year,
			TotalRevenue: int32(item.TotalRevenue),
		})
	}
	return responses
}

func (c *categoryGraphqlMapper) mapResponsesCategoryMonthPrice(data []*pb.CategoryMonthPriceResponse) []*model.CategoryMonthPriceResponse {
	var responses []*model.CategoryMonthPriceResponse
	for _, item := range data {
		if item == nil {
			continue
		}
		responses = append(responses, &model.CategoryMonthPriceResponse{
			Month:        item.Month,
			CategoryID:   int32(item.CategoryId),
			CategoryName: item.CategoryName,
			OrderCount:   int32(item.OrderCount),
			ItemsSold:    int32(item.ItemsSold),
			TotalRevenue: int32(item.TotalRevenue),
		})
	}
	return responses
}

func (c *categoryGraphqlMapper) mapResponsesCategoryYearPrice(data []*pb.CategoryYearPriceResponse) []*model.CategoryYearPriceResponse {
	var responses []*model.CategoryYearPriceResponse
	for _, item := range data {
		if item == nil {
			continue
		}
		responses = append(responses, &model.CategoryYearPriceResponse{
			Year:               item.Year,
			CategoryID:         int32(item.CategoryId),
			CategoryName:       item.CategoryName,
			OrderCount:         int32(item.OrderCount),
			ItemsSold:          int32(item.ItemsSold),
			TotalRevenue:       int32(item.TotalRevenue),
			UniqueProductsSold: int32(item.UniqueProductsSold),
		})
	}
	return responses
}
