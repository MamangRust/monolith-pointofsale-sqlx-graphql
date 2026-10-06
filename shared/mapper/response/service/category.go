package response_service

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/record"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type categoryResponseMapper struct {
}

func NewCategoryResponseMapper() *categoryResponseMapper {
	return &categoryResponseMapper{}
}

func (s *categoryResponseMapper) ToCategoryResponse(category *record.CategoriesRecord) *response.CategoryResponse {
	return &response.CategoryResponse{
		ID:            category.ID,
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
	}
}

func (s *categoryResponseMapper) ToCategorysResponse(categories []*record.CategoriesRecord) []*response.CategoryResponse {
	var responses []*response.CategoryResponse

	for _, category := range categories {
		responses = append(responses, s.ToCategoryResponse(category))
	}

	return responses
}

func (s *categoryResponseMapper) ToCategoryResponseDeleteAt(category *record.CategoriesRecord) *response.CategoryResponseDeleteAt {
	return &response.CategoryResponseDeleteAt{
		ID:            category.ID,
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
		DeletedAt:     category.DeletedAt,
	}
}

func (s *categoryResponseMapper) ToCategoryResponsesDeleteAt(categories []*record.CategoriesRecord) []*response.CategoryResponseDeleteAt {
	var responses []*response.CategoryResponseDeleteAt

	for _, category := range categories {
		responses = append(responses, s.ToCategoryResponseDeleteAt(category))
	}

	return responses
}

func (s *categoryResponseMapper) ToCategoryMonthlyPrice(category *record.CategoriesMonthPriceRecord) *response.CategoryMonthPriceResponse {
	return &response.CategoryMonthPriceResponse{
		Month:        category.Month,
		CategoryID:   int(category.CategoryID),
		CategoryName: category.CategoryName,
		OrderCount:   int(category.OrderCount),
		ItemsSold:    int(category.ItemsSold),
		TotalRevenue: int(category.TotalRevenue),
	}
}

func (s *categoryResponseMapper) ToCategoryMonthlyPrices(c []*record.CategoriesMonthPriceRecord) []*response.CategoryMonthPriceResponse {
	var categoryRecords []*response.CategoryMonthPriceResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryMonthlyPrice(category))
	}

	return categoryRecords
}

func (s *categoryResponseMapper) ToCategoryYearlyPrice(category *record.CategoriesYearPriceRecord) *response.CategoryYearPriceResponse {
	return &response.CategoryYearPriceResponse{
		Year:               category.Year,
		CategoryID:         int(category.CategoryID),
		CategoryName:       category.CategoryName,
		OrderCount:         int(category.OrderCount),
		ItemsSold:          int(category.ItemsSold),
		TotalRevenue:       int(category.TotalRevenue),
		UniqueProductsSold: int(category.UniqueProductsSold),
	}
}

func (s *categoryResponseMapper) ToCategoryYearlyPrices(c []*record.CategoriesYearPriceRecord) []*response.CategoryYearPriceResponse {
	var categoryRecords []*response.CategoryYearPriceResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryYearlyPrice(category))
	}

	return categoryRecords
}

func (s *categoryResponseMapper) ToCategoryMonthlyTotalPrice(c *record.CategoriesMonthlyTotalPriceRecord) *response.CategoriesMonthlyTotalPriceResponse {
	return &response.CategoriesMonthlyTotalPriceResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *categoryResponseMapper) ToCategoryMonthlyTotalPrices(c []*record.CategoriesMonthlyTotalPriceRecord) []*response.CategoriesMonthlyTotalPriceResponse {
	var categoryRecords []*response.CategoriesMonthlyTotalPriceResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryMonthlyTotalPrice(category))
	}

	return categoryRecords
}

func (s *categoryResponseMapper) ToCategoryYearlyTotalPrice(c *record.CategoriesYearlyTotalPriceRecord) *response.CategoriesYearlyTotalPriceResponse {
	return &response.CategoriesYearlyTotalPriceResponse{
		Year:         c.Year,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *categoryResponseMapper) ToCategoryYearlyTotalPrices(c []*record.CategoriesYearlyTotalPriceRecord) []*response.CategoriesYearlyTotalPriceResponse {
	var categoryRecords []*response.CategoriesYearlyTotalPriceResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryYearlyTotalPrice(category))
	}

	return categoryRecords
}
