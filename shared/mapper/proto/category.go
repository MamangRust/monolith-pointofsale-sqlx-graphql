package protomapper

import (
	pbcategory "github.com/MamangRust/monolith-pointofsale-pb/category"
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type categoryProtoMapper struct {
}

func NewCategoryProtoMapper() *categoryProtoMapper {
	return &categoryProtoMapper{}
}

func (c *categoryProtoMapper) ToProtoResponseCategory(status string, message string, pbResponse *response.CategoryResponse) *pbcategory.ApiResponseCategory {
	return &pbcategory.ApiResponseCategory{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCategory(pbResponse),
	}
}

func (c *categoryProtoMapper) ToProtoResponseCategoryDeleteAt(status string, message string, pbResponse *response.CategoryResponseDeleteAt) *pbcategory.ApiResponseCategoryDeleteAt {
	return &pbcategory.ApiResponseCategoryDeleteAt{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCategoryDeleteAt(pbResponse),
	}
}

func (c *categoryProtoMapper) ToProtoResponsesCategory(status string, message string, pbResponse []*response.CategoryResponse) *pbcategory.ApiResponsesCategory {
	return &pbcategory.ApiResponsesCategory{
		Status:  status,
		Message: message,
		Data:    c.mapResponsesCategory(pbResponse),
	}
}

func (c *categoryProtoMapper) ToProtoResponseCategoryDelete(status string, message string) *pbcategory.ApiResponseCategoryDelete {
	return &pbcategory.ApiResponseCategoryDelete{
		Status:  status,
		Message: message,
	}
}

func (c *categoryProtoMapper) ToProtoResponseCategoryAll(status string, message string) *pbcategory.ApiResponseCategoryAll {
	return &pbcategory.ApiResponseCategoryAll{
		Status:  status,
		Message: message,
	}
}

func (c *categoryProtoMapper) ToProtoResponsePaginationCategoryDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, categories []*response.CategoryResponseDeleteAt) *pbcategory.ApiResponsePaginationCategoryDeleteAt {
	return &pbcategory.ApiResponsePaginationCategoryDeleteAt{
		Status:     status,
		Message:    message,
		Data:       c.mapResponsesCategoryDeleteAt(categories),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (c *categoryProtoMapper) ToProtoResponsePaginationCategory(pagination *pbcommon.PaginationMeta, status string, message string, categories []*response.CategoryResponse) *pbcategory.ApiResponsePaginationCategory {
	return &pbcategory.ApiResponsePaginationCategory{
		Status:     status,
		Message:    message,
		Data:       c.mapResponsesCategory(categories),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (c *categoryProtoMapper) ToProtoResponseCategoryMonthlyPrice(status string, message string, row []*response.CategoryMonthPriceResponse) *pbcategory.ApiResponseCategoryMonthPrice {
	return &pbcategory.ApiResponseCategoryMonthPrice{
		Status:  status,
		Message: message,
		Data:    c.mapResponsesCategoryMonthlyPrices(row),
	}
}

func (c *categoryProtoMapper) ToProtoResponseCategoryYearlyPrice(status string, message string, row []*response.CategoryYearPriceResponse) *pbcategory.ApiResponseCategoryYearPrice {
	return &pbcategory.ApiResponseCategoryYearPrice{
		Status:  status,
		Message: message,
		Data:    c.mapResponsesCategoryYearlyPrices(row),
	}
}

func (c *categoryProtoMapper) ToProtoResponseMonthlyTotalPrice(status string, message string, row []*response.CategoriesMonthlyTotalPriceResponse) *pbcategory.ApiResponseCategoryMonthlyTotalPrice {
	return &pbcategory.ApiResponseCategoryMonthlyTotalPrice{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCategoryMonthlyTotalPrices(row),
	}
}

func (c *categoryProtoMapper) ToProtoResponseYearlyTotalPrice(status string, message string, row []*response.CategoriesYearlyTotalPriceResponse) *pbcategory.ApiResponseCategoryYearlyTotalPrice {
	return &pbcategory.ApiResponseCategoryYearlyTotalPrice{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCategoryYearlyTotalPrices(row),
	}
}

func (c *categoryProtoMapper) mapResponseCategory(category *response.CategoryResponse) *pbcategory.CategoryResponse {
	return &pbcategory.CategoryResponse{
		Id:            int32(category.ID),
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
	}
}

func (c *categoryProtoMapper) mapResponsesCategory(categories []*response.CategoryResponse) []*pbcategory.CategoryResponse {
	var mappedCategories []*pbcategory.CategoryResponse

	for _, category := range categories {
		mappedCategories = append(mappedCategories, c.mapResponseCategory(category))
	}

	return mappedCategories
}

func (c *categoryProtoMapper) mapResponseCategoryDeleteAt(category *response.CategoryResponseDeleteAt) *pbcategory.CategoryResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if category.DeletedAt != nil {
		deletedAt = wrapperspb.String(*category.DeletedAt)
	}

	return &pbcategory.CategoryResponseDeleteAt{
		Id:            int32(category.ID),
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
		DeletedAt:     deletedAt,
	}
}

func (c *categoryProtoMapper) mapResponsesCategoryDeleteAt(categories []*response.CategoryResponseDeleteAt) []*pbcategory.CategoryResponseDeleteAt {
	var mappedCategories []*pbcategory.CategoryResponseDeleteAt

	for _, category := range categories {
		mappedCategories = append(mappedCategories, c.mapResponseCategoryDeleteAt(category))
	}

	return mappedCategories
}

func (s *categoryProtoMapper) mapResponseCategoryMonthlyPrice(category *response.CategoryMonthPriceResponse) *pbcategory.CategoryMonthPriceResponse {
	return &pbcategory.CategoryMonthPriceResponse{
		Month:        category.Month,
		CategoryId:   int32(category.CategoryID),
		CategoryName: category.CategoryName,
		OrderCount:   int32(category.OrderCount),
		ItemsSold:    int32(category.ItemsSold),
		TotalRevenue: int32(category.TotalRevenue),
	}
}

func (s *categoryProtoMapper) mapResponsesCategoryMonthlyPrices(c []*response.CategoryMonthPriceResponse) []*pbcategory.CategoryMonthPriceResponse {
	var categoryRecords []*pbcategory.CategoryMonthPriceResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.mapResponseCategoryMonthlyPrice(category))
	}

	return categoryRecords
}

func (s *categoryProtoMapper) mapResponseCategoryYearlyPrice(category *response.CategoryYearPriceResponse) *pbcategory.CategoryYearPriceResponse {
	return &pbcategory.CategoryYearPriceResponse{
		Year:               category.Year,
		CategoryId:         int32(category.CategoryID),
		CategoryName:       category.CategoryName,
		OrderCount:         int32(category.OrderCount),
		ItemsSold:          int32(category.ItemsSold),
		TotalRevenue:       int32(category.TotalRevenue),
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func (s *categoryProtoMapper) mapResponsesCategoryYearlyPrices(c []*response.CategoryYearPriceResponse) []*pbcategory.CategoryYearPriceResponse {
	var categoryRecords []*pbcategory.CategoryYearPriceResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.mapResponseCategoryYearlyPrice(category))
	}

	return categoryRecords
}

func (s *categoryProtoMapper) mapResponseCashierMonthlyTotalPrice(c *response.CategoriesMonthlyTotalPriceResponse) *pbcategory.CategoriesMonthlyTotalPriceResponse {
	return &pbcategory.CategoriesMonthlyTotalPriceResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: int32(c.TotalRevenue),
	}
}

func (s *categoryProtoMapper) mapResponseCategoryMonthlyTotalPrices(c []*response.CategoriesMonthlyTotalPriceResponse) []*pbcategory.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*pbcategory.CategoriesMonthlyTotalPriceResponse

	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, s.mapResponseCashierMonthlyTotalPrice(Category))
	}

	return CategoryRecords
}

func (s *categoryProtoMapper) mapResponseCategoryYearlyTotalSale(c *response.CategoriesYearlyTotalPriceResponse) *pbcategory.CategoriesYearlyTotalPriceResponse {
	return &pbcategory.CategoriesYearlyTotalPriceResponse{
		Year:         c.Year,
		TotalRevenue: int32(c.TotalRevenue),
	}
}

func (s *categoryProtoMapper) mapResponseCategoryYearlyTotalPrices(c []*response.CategoriesYearlyTotalPriceResponse) []*pbcategory.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*pbcategory.CategoriesYearlyTotalPriceResponse

	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, s.mapResponseCategoryYearlyTotalSale(Category))
	}

	return CategoryRecords
}
