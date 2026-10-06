package response_api

import (
	pbproduct "github.com/MamangRust/monolith-pointofsale-pb/product"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type productResponseMapper struct {
}

func NewProductResponseMapper() *productResponseMapper {
	return &productResponseMapper{}
}

func (p *productResponseMapper) ToResponseProduct(product *pbproduct.ProductResponse) *response.ProductResponse {
	return &response.ProductResponse{
		ID:           int(product.Id),
		MerchantID:   int(product.MerchantId),
		CategoryID:   int(product.CategoryId),
		Name:         product.Name,
		Description:  product.Description,
		Price:        int(product.Price),
		CountInStock: int(product.CountInStock),
		Brand:        product.Brand,
		Weight:       int(product.Weight),
		SlugProduct:  product.SlugProduct,
		ImageProduct: product.ImageProduct,
		Barcode:      product.Barcode,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}

func (p *productResponseMapper) ToResponsesProduct(products []*pbproduct.ProductResponse) []*response.ProductResponse {
	var mappedProducts []*response.ProductResponse

	for _, product := range products {
		mappedProducts = append(mappedProducts, p.ToResponseProduct(product))
	}

	return mappedProducts
}

func (p *productResponseMapper) ToResponseProductDeleteAt(product *pbproduct.ProductResponseDeleteAt) *response.ProductResponseDeleteAt {
	var deletedAt string

	if product.DeletedAt != nil {
		deletedAt = product.DeletedAt.Value
	}

	return &response.ProductResponseDeleteAt{
		ID:           int(product.Id),
		MerchantID:   int(product.MerchantId),
		CategoryID:   int(product.CategoryId),
		Name:         product.Name,
		Description:  product.Description,
		Price:        int(product.Price),
		CountInStock: int(product.CountInStock),
		Brand:        product.Brand,
		Weight:       int(product.Weight),
		SlugProduct:  product.SlugProduct,
		ImageProduct: product.ImageProduct,
		Barcode:      product.Barcode,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
		DeleteAt:     &deletedAt,
	}
}

func (p *productResponseMapper) ToResponsesProductDeleteAt(products []*pbproduct.ProductResponseDeleteAt) []*response.ProductResponseDeleteAt {
	var mappedProducts []*response.ProductResponseDeleteAt

	for _, product := range products {
		mappedProducts = append(mappedProducts, p.ToResponseProductDeleteAt(product))
	}

	return mappedProducts
}

func (p *productResponseMapper) ToApiResponseProduct(pbResponse *pbproduct.ApiResponseProduct) *response.ApiResponseProduct {
	return &response.ApiResponseProduct{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    p.ToResponseProduct(pbResponse.Data),
	}
}

func (p *productResponseMapper) ToApiResponsesProduct(pbResponse *pbproduct.ApiResponsesProduct) *response.ApiResponsesProduct {
	return &response.ApiResponsesProduct{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    p.ToResponsesProduct(pbResponse.Data),
	}
}

func (p *productResponseMapper) ToApiResponsesProductDeleteAt(pbResponse *pbproduct.ApiResponseProductDeleteAt) *response.ApiResponseProductDeleteAt {
	return &response.ApiResponseProductDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    p.ToResponseProductDeleteAt(pbResponse.Data),
	}
}

func (p *productResponseMapper) ToApiResponseProductDelete(pbResponse *pbproduct.ApiResponseProductDelete) *response.ApiResponseProductDelete {
	return &response.ApiResponseProductDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (p *productResponseMapper) ToApiResponseProductAll(pbResponse *pbproduct.ApiResponseProductAll) *response.ApiResponseProductAll {
	return &response.ApiResponseProductAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (p *productResponseMapper) ToApiResponsePaginationProductDeleteAt(pbResponse *pbproduct.ApiResponsePaginationProductDeleteAt) *response.ApiResponsePaginationProductDeleteAt {
	return &response.ApiResponsePaginationProductDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       p.ToResponsesProductDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (p *productResponseMapper) ToApiResponsePaginationProduct(pbResponse *pbproduct.ApiResponsePaginationProduct) *response.ApiResponsePaginationProduct {
	return &response.ApiResponsePaginationProduct{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       p.ToResponsesProduct(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}
