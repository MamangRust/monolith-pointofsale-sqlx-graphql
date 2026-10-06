package handler

import (
	pbcategory "github.com/MamangRust/monolith-pointofsale-pb/category"
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func mapPaginationMeta(meta *pbcommon.PaginationMeta) *pbcommon.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &pbcommon.PaginationMeta{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		TotalPages:   meta.TotalPages,
		TotalRecords: meta.TotalRecords,
	}
}

func mapResponseCategory(category *db.Category) *pbcategory.CategoryResponse {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if category.CreatedAt.Valid {
		createdAtStr = category.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if category.UpdatedAt.Valid {
		updatedAtStr = category.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return &pbcategory.CategoryResponse{
		Id:           int32(category.CategoryID),
		Name:         category.Name,
		Description:  description,
		SlugCategory: slugCategory,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
	}
}

func mapResponseGetCategory(category *db.GetCategoriesRow) *pbcategory.CategoryResponse {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if category.CreatedAt.Valid {
		createdAtStr = category.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if category.UpdatedAt.Valid {
		updatedAtStr = category.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return &pbcategory.CategoryResponse{
		Id:           int32(category.CategoryID),
		Name:         category.Name,
		Description:  description,
		SlugCategory: slugCategory,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
	}
}

func mapResponsesCategory(categories []*db.GetCategoriesRow) []*pbcategory.CategoryResponse {
	var mappedCategories []*pbcategory.CategoryResponse
	for _, category := range categories {
		mappedCategories = append(mappedCategories, mapResponseGetCategory(category))
	}
	return mappedCategories
}

func mapResponseCategoryDeleteAt(category *db.Category) *pbcategory.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if category.CreatedAt.Valid {
		createdAtStr = category.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if category.UpdatedAt.Valid {
		updatedAtStr = category.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var deletedAt *wrapperspb.StringValue
	if category.DeletedAt.Valid {
		deletedAt = wrapperspb.String(category.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}

	return &pbcategory.CategoryResponseDeleteAt{
		Id:           int32(category.CategoryID),
		Name:         category.Name,
		Description:  description,
		SlugCategory: slugCategory,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
		DeletedAt:    deletedAt,
	}
}

func mapResponseGetCategoryActive(category *db.GetCategoriesActiveRow) *pbcategory.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if category.CreatedAt.Valid {
		createdAtStr = category.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if category.UpdatedAt.Valid {
		updatedAtStr = category.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var deletedAt *wrapperspb.StringValue
	if category.DeletedAt.Valid {
		deletedAt = wrapperspb.String(category.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}

	return &pbcategory.CategoryResponseDeleteAt{
		Id:           int32(category.CategoryID),
		Name:         category.Name,
		Description:  description,
		SlugCategory: slugCategory,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
		DeletedAt:    deletedAt,
	}
}

func mapResponsesCategoryActive(categories []*db.GetCategoriesActiveRow) []*pbcategory.CategoryResponseDeleteAt {
	var mappedCategories []*pbcategory.CategoryResponseDeleteAt
	for _, category := range categories {
		mappedCategories = append(mappedCategories, mapResponseGetCategoryActive(category))
	}
	return mappedCategories
}

func mapResponseGetCategoryTrashed(category *db.GetCategoriesTrashedRow) *pbcategory.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if category.CreatedAt.Valid {
		createdAtStr = category.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if category.UpdatedAt.Valid {
		updatedAtStr = category.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var deletedAt *wrapperspb.StringValue
	if category.DeletedAt.Valid {
		deletedAt = wrapperspb.String(category.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}

	return &pbcategory.CategoryResponseDeleteAt{
		Id:           int32(category.CategoryID),
		Name:         category.Name,
		Description:  description,
		SlugCategory: slugCategory,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
		DeletedAt:    deletedAt,
	}
}

func mapResponsesCategoryTrashed(categories []*db.GetCategoriesTrashedRow) []*pbcategory.CategoryResponseDeleteAt {
	var mappedCategories []*pbcategory.CategoryResponseDeleteAt
	for _, category := range categories {
		mappedCategories = append(mappedCategories, mapResponseGetCategoryTrashed(category))
	}
	return mappedCategories
}

func mapResponseCategoryMonthlyPrice(category *db.GetMonthlyCategoryRow) *pbcategory.CategoryMonthPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategory.CategoryMonthPriceResponse{
		Month:        category.Month,
		CategoryId:   category.CategoryID,
		CategoryName: category.CategoryName,
		OrderCount:   int32(category.OrderCount),
		ItemsSold:    int32(category.ItemsSold),
		TotalRevenue: category.TotalRevenue,
	}
}

func mapResponsesCategoryMonthlyPrices(c []*db.GetMonthlyCategoryRow) []*pbcategory.CategoryMonthPriceResponse {
	var categoryRecords []*pbcategory.CategoryMonthPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryMonthlyPrice(category))
	}
	return categoryRecords
}

func mapResponseCategoryMonthlyPriceById(category *db.GetMonthlyCategoryByIdRow) *pbcategory.CategoryMonthPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategory.CategoryMonthPriceResponse{
		Month:        category.Month,
		CategoryId:   category.CategoryID,
		CategoryName: category.CategoryName,
		OrderCount:   int32(category.OrderCount),
		ItemsSold:    int32(category.ItemsSold),
		TotalRevenue: category.TotalRevenue,
	}
}

func mapResponsesCategoryMonthlyPricesById(c []*db.GetMonthlyCategoryByIdRow) []*pbcategory.CategoryMonthPriceResponse {
	var categoryRecords []*pbcategory.CategoryMonthPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryMonthlyPriceById(category))
	}
	return categoryRecords
}

func mapResponseCategoryMonthlyPriceByMerchant(category *db.GetMonthlyCategoryByMerchantRow) *pbcategory.CategoryMonthPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategory.CategoryMonthPriceResponse{
		Month:        category.Month,
		CategoryId:   category.CategoryID,
		CategoryName: category.CategoryName,
		OrderCount:   int32(category.OrderCount),
		ItemsSold:    int32(category.ItemsSold),
		TotalRevenue: category.TotalRevenue,
	}
}

func mapResponsesCategoryMonthlyPricesByMerchant(c []*db.GetMonthlyCategoryByMerchantRow) []*pbcategory.CategoryMonthPriceResponse {
	var categoryRecords []*pbcategory.CategoryMonthPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryMonthlyPriceByMerchant(category))
	}
	return categoryRecords
}

func mapResponseCategoryYearlyPrice(category *db.GetYearlyCategoryRow) *pbcategory.CategoryYearPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategory.CategoryYearPriceResponse{
		Year:               category.Year,
		CategoryId:         category.CategoryID,
		CategoryName:       category.CategoryName,
		OrderCount:         int32(category.OrderCount),
		ItemsSold:          int32(category.ItemsSold),
		TotalRevenue:       category.TotalRevenue,
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func mapResponsesCategoryYearlyPrices(c []*db.GetYearlyCategoryRow) []*pbcategory.CategoryYearPriceResponse {
	var categoryRecords []*pbcategory.CategoryYearPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryYearlyPrice(category))
	}
	return categoryRecords
}

func mapResponseCategoryYearlyPriceById(category *db.GetYearlyCategoryByIdRow) *pbcategory.CategoryYearPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategory.CategoryYearPriceResponse{
		Year:               category.Year,
		CategoryId:         category.CategoryID,
		CategoryName:       category.CategoryName,
		OrderCount:         int32(category.OrderCount),
		ItemsSold:          int32(category.ItemsSold),
		TotalRevenue:       category.TotalRevenue,
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func mapResponsesCategoryYearlyPricesById(c []*db.GetYearlyCategoryByIdRow) []*pbcategory.CategoryYearPriceResponse {
	var categoryRecords []*pbcategory.CategoryYearPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryYearlyPriceById(category))
	}
	return categoryRecords
}

func mapResponseCategoryYearlyPriceByMerchant(category *db.GetYearlyCategoryByMerchantRow) *pbcategory.CategoryYearPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategory.CategoryYearPriceResponse{
		Year:               category.Year,
		CategoryId:         category.CategoryID,
		CategoryName:       category.CategoryName,
		OrderCount:         int32(category.OrderCount),
		ItemsSold:          int32(category.ItemsSold),
		TotalRevenue:       category.TotalRevenue,
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func mapResponsesCategoryYearlyPricesByMerchant(c []*db.GetYearlyCategoryByMerchantRow) []*pbcategory.CategoryYearPriceResponse {
	var categoryRecords []*pbcategory.CategoryYearPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryYearlyPriceByMerchant(category))
	}
	return categoryRecords
}

func mapResponseCashierMonthlyTotalPrice(c *db.GetMonthlyTotalPriceRow) *pbcategory.CategoriesMonthlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategory.CategoriesMonthlyTotalPriceResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryMonthlyTotalPrices(c []*db.GetMonthlyTotalPriceRow) []*pbcategory.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*pbcategory.CategoriesMonthlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCashierMonthlyTotalPrice(Category))
	}
	return CategoryRecords
}

func mapResponseCashierMonthlyTotalPriceById(c *db.GetMonthlyTotalPriceByIdRow) *pbcategory.CategoriesMonthlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategory.CategoriesMonthlyTotalPriceResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryMonthlyTotalPricesById(c []*db.GetMonthlyTotalPriceByIdRow) []*pbcategory.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*pbcategory.CategoriesMonthlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCashierMonthlyTotalPriceById(Category))
	}
	return CategoryRecords
}

func mapResponseCashierMonthlyTotalPriceByMerchant(c *db.GetMonthlyTotalPriceByMerchantRow) *pbcategory.CategoriesMonthlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategory.CategoriesMonthlyTotalPriceResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryMonthlyTotalPricesByMerchant(c []*db.GetMonthlyTotalPriceByMerchantRow) []*pbcategory.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*pbcategory.CategoriesMonthlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCashierMonthlyTotalPriceByMerchant(Category))
	}
	return CategoryRecords
}

func mapResponseCategoryYearlyTotalSale(c *db.GetYearlyTotalPriceRow) *pbcategory.CategoriesYearlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategory.CategoriesYearlyTotalPriceResponse{
		Year:         c.Year,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryYearlyTotalPrices(c []*db.GetYearlyTotalPriceRow) []*pbcategory.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*pbcategory.CategoriesYearlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCategoryYearlyTotalSale(Category))
	}
	return CategoryRecords
}

func mapResponseCategoryYearlyTotalSaleById(c *db.GetYearlyTotalPriceByIdRow) *pbcategory.CategoriesYearlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategory.CategoriesYearlyTotalPriceResponse{
		Year:         c.Year,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryYearlyTotalPricesById(c []*db.GetYearlyTotalPriceByIdRow) []*pbcategory.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*pbcategory.CategoriesYearlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCategoryYearlyTotalSaleById(Category))
	}
	return CategoryRecords
}

func mapResponseCategoryYearlyTotalSaleByMerchant(c *db.GetYearlyTotalPriceByMerchantRow) *pbcategory.CategoriesYearlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategory.CategoriesYearlyTotalPriceResponse{
		Year:         c.Year,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryYearlyTotalPricesByMerchant(c []*db.GetYearlyTotalPriceByMerchantRow) []*pbcategory.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*pbcategory.CategoriesYearlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCategoryYearlyTotalSaleByMerchant(Category))
	}
	return CategoryRecords
}
