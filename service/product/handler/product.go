package handler

import (
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbproduct "github.com/MamangRust/monolith-pointofsale-pb/product"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func parseStrPointer(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func parseInt32Pointer(i *int32) int32 {
	if i != nil {
		return *i
	}
	return 0
}

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

func mapResponseProduct(product *db.Product) *pbproduct.ProductResponse {
	if product == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if product.CreatedAt.Valid {
		createdAtStr = product.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if product.UpdatedAt.Valid {
		updatedAtStr = product.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return &pbproduct.ProductResponse{
		Id:           int32(product.ProductID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  parseStrPointer(product.Description),
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        parseStrPointer(product.Brand),
		Weight:       parseInt32Pointer(product.Weight),
		SlugProduct:  parseStrPointer(product.SlugProduct),
		ImageProduct: parseStrPointer(product.ImageProduct),
		Barcode:      parseStrPointer(product.Barcode),
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
	}
}

func mapResponsesProduct(products []*db.GetProductsRow) []*pbproduct.ProductResponse {
	var mappedProducts []*pbproduct.ProductResponse
	for _, p := range products {
		if p == nil {
			continue
		}
		var createdAtStr, updatedAtStr string
		if p.CreatedAt.Valid {
			createdAtStr = p.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		if p.UpdatedAt.Valid {
			updatedAtStr = p.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		mappedProducts = append(mappedProducts, &pbproduct.ProductResponse{
			Id:           int32(p.ProductID),
			MerchantId:   int32(p.MerchantID),
			CategoryId:   int32(p.CategoryID),
			Name:         p.Name,
			Description:  parseStrPointer(p.Description),
			Price:        int32(p.Price),
			CountInStock: int32(p.CountInStock),
			Brand:        parseStrPointer(p.Brand),
			Weight:       parseInt32Pointer(p.Weight),
			SlugProduct:  parseStrPointer(p.SlugProduct),
			ImageProduct: parseStrPointer(p.ImageProduct),
			Barcode:      parseStrPointer(p.Barcode),
			CreatedAt:    createdAtStr,
			UpdatedAt:    updatedAtStr,
		})
	}
	return mappedProducts
}

func mapResponsesProductByMerchant(products []*db.GetProductsByMerchantRow, merchantId int32) []*pbproduct.ProductResponse {
	var mappedProducts []*pbproduct.ProductResponse
	for _, p := range products {
		if p == nil {
			continue
		}
		var createdAtStr string
		if p.CreatedAt.Valid {
			createdAtStr = p.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		mappedProducts = append(mappedProducts, &pbproduct.ProductResponse{
			Id:           int32(p.ProductID),
			MerchantId:   merchantId,
			Name:         p.Name,
			Description:  parseStrPointer(p.Description),
			Price:        int32(p.Price),
			CountInStock: int32(p.CountInStock),
			Brand:        parseStrPointer(p.Brand),
			ImageProduct: parseStrPointer(p.ImageProduct),
			CreatedAt:    createdAtStr,
		})
	}
	return mappedProducts
}

func mapResponsesProductByCategory(products []*db.GetProductsByCategoryNameRow) []*pbproduct.ProductResponse {
	var mappedProducts []*pbproduct.ProductResponse
	for _, p := range products {
		if p == nil {
			continue
		}
		var createdAtStr, updatedAtStr string
		if p.CreatedAt.Valid {
			createdAtStr = p.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		if p.UpdatedAt.Valid {
			updatedAtStr = p.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		mappedProducts = append(mappedProducts, &pbproduct.ProductResponse{
			Id:           int32(p.ProductID),
			MerchantId:   int32(p.MerchantID),
			CategoryId:   int32(p.CategoryID),
			Name:         p.Name,
			Description:  parseStrPointer(p.Description),
			Price:        int32(p.Price),
			CountInStock: int32(p.CountInStock),
			Brand:        parseStrPointer(p.Brand),
			Weight:       parseInt32Pointer(p.Weight),
			SlugProduct:  parseStrPointer(p.SlugProduct),
			ImageProduct: parseStrPointer(p.ImageProduct),
			Barcode:      parseStrPointer(p.Barcode),
			CreatedAt:    createdAtStr,
			UpdatedAt:    updatedAtStr,
		})
	}
	return mappedProducts
}

func mapResponseProductDeleteAt(product *db.Product) *pbproduct.ProductResponseDeleteAt {
	if product == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if product.CreatedAt.Valid {
		createdAtStr = product.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if product.UpdatedAt.Valid {
		updatedAtStr = product.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var deletedAt *wrapperspb.StringValue
	if product.DeletedAt.Valid {
		deletedAt = wrapperspb.String(product.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}

	return &pbproduct.ProductResponseDeleteAt{
		Id:           int32(product.ProductID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  parseStrPointer(product.Description),
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        parseStrPointer(product.Brand),
		Weight:       parseInt32Pointer(product.Weight),
		SlugProduct:  parseStrPointer(product.SlugProduct),
		ImageProduct: parseStrPointer(product.ImageProduct),
		Barcode:      parseStrPointer(product.Barcode),
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
		DeletedAt:    deletedAt,
	}
}

func mapResponsesProductActive(products []*db.GetProductsActiveRow) []*pbproduct.ProductResponseDeleteAt {
	var mappedProducts []*pbproduct.ProductResponseDeleteAt
	for _, p := range products {
		if p == nil {
			continue
		}
		var createdAtStr, updatedAtStr string
		if p.CreatedAt.Valid {
			createdAtStr = p.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		if p.UpdatedAt.Valid {
			updatedAtStr = p.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var deletedAt *wrapperspb.StringValue
		if p.DeletedAt.Valid {
			deletedAt = wrapperspb.String(p.DeletedAt.Time.Format("2006-01-02 15:04:05"))
		}

		mappedProducts = append(mappedProducts, &pbproduct.ProductResponseDeleteAt{
			Id:           int32(p.ProductID),
			MerchantId:   int32(p.MerchantID),
			CategoryId:   int32(p.CategoryID),
			Name:         p.Name,
			Description:  parseStrPointer(p.Description),
			Price:        int32(p.Price),
			CountInStock: int32(p.CountInStock),
			Brand:        parseStrPointer(p.Brand),
			Weight:       parseInt32Pointer(p.Weight),
			SlugProduct:  parseStrPointer(p.SlugProduct),
			ImageProduct: parseStrPointer(p.ImageProduct),
			Barcode:      parseStrPointer(p.Barcode),
			CreatedAt:    createdAtStr,
			UpdatedAt:    updatedAtStr,
			DeletedAt:    deletedAt,
		})
	}
	return mappedProducts
}

func mapResponsesProductTrashed(products []*db.GetProductsTrashedRow) []*pbproduct.ProductResponseDeleteAt {
	var mappedProducts []*pbproduct.ProductResponseDeleteAt
	for _, p := range products {
		if p == nil {
			continue
		}
		var createdAtStr, updatedAtStr string
		if p.CreatedAt.Valid {
			createdAtStr = p.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		if p.UpdatedAt.Valid {
			updatedAtStr = p.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var deletedAt *wrapperspb.StringValue
		if p.DeletedAt.Valid {
			deletedAt = wrapperspb.String(p.DeletedAt.Time.Format("2006-01-02 15:04:05"))
		}

		mappedProducts = append(mappedProducts, &pbproduct.ProductResponseDeleteAt{
			Id:           int32(p.ProductID),
			MerchantId:   int32(p.MerchantID),
			CategoryId:   int32(p.CategoryID),
			Name:         p.Name,
			Description:  parseStrPointer(p.Description),
			Price:        int32(p.Price),
			CountInStock: int32(p.CountInStock),
			Brand:        parseStrPointer(p.Brand),
			Weight:       parseInt32Pointer(p.Weight),
			SlugProduct:  parseStrPointer(p.SlugProduct),
			ImageProduct: parseStrPointer(p.ImageProduct),
			Barcode:      parseStrPointer(p.Barcode),
			CreatedAt:    createdAtStr,
			UpdatedAt:    updatedAtStr,
			DeletedAt:    deletedAt,
		})
	}
	return mappedProducts
}
