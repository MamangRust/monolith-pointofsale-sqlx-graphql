package protomapper

import (
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbproduct "github.com/MamangRust/monolith-pointofsale-pb/product"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type productProtoMapper struct{}

func NewProductProtoMapper() *productProtoMapper {
	return &productProtoMapper{}
}

func (p *productProtoMapper) ToProtoResponseProduct(status string, message string, pbResponse *response.ProductResponse) *pbproduct.ApiResponseProduct {
	return &pbproduct.ApiResponseProduct{
		Status:  status,
		Message: message,
		Data:    p.mapResponseProduct(pbResponse),
	}
}

func (p *productProtoMapper) ToProtoResponsesProduct(status string, message string, pbResponse []*response.ProductResponse) *pbproduct.ApiResponsesProduct {
	return &pbproduct.ApiResponsesProduct{
		Status:  status,
		Message: message,
		Data:    p.mapResponsesProduct(pbResponse),
	}
}

func (p *productProtoMapper) ToProtoResponseProductDeleteAt(status string, message string, pbResponse *response.ProductResponseDeleteAt) *pbproduct.ApiResponseProductDeleteAt {
	return &pbproduct.ApiResponseProductDeleteAt{
		Status:  status,
		Message: message,
		Data:    p.mapResponseProductDeleteAt(pbResponse),
	}
}

func (p *productProtoMapper) ToProtoResponseProductDelete(status string, message string) *pbproduct.ApiResponseProductDelete {
	return &pbproduct.ApiResponseProductDelete{
		Status:  status,
		Message: message,
	}
}

func (p *productProtoMapper) ToProtoResponseProductAll(status string, message string) *pbproduct.ApiResponseProductAll {
	return &pbproduct.ApiResponseProductAll{
		Status:  status,
		Message: message,
	}
}

func (p *productProtoMapper) ToProtoResponsePaginationProductDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, products []*response.ProductResponseDeleteAt) *pbproduct.ApiResponsePaginationProductDeleteAt {
	return &pbproduct.ApiResponsePaginationProductDeleteAt{
		Status:     status,
		Message:    message,
		Data:       p.mapResponsesProductDeleteAt(products),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (p *productProtoMapper) ToProtoResponsePaginationProduct(pagination *pbcommon.PaginationMeta, status string, message string, products []*response.ProductResponse) *pbproduct.ApiResponsePaginationProduct {
	return &pbproduct.ApiResponsePaginationProduct{
		Status:     status,
		Message:    message,
		Data:       p.mapResponsesProduct(products),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (p *productProtoMapper) mapResponseProduct(product *response.ProductResponse) *pbproduct.ProductResponse {
	return &pbproduct.ProductResponse{
		Id:           int32(product.ID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  product.Description,
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        product.Brand,
		Weight:       int32(product.Weight),
		SlugProduct:  product.SlugProduct,
		ImageProduct: product.ImageProduct,
		Barcode:      product.Barcode,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}

func (p *productProtoMapper) mapResponsesProduct(products []*response.ProductResponse) []*pbproduct.ProductResponse {
	var mappedProducts []*pbproduct.ProductResponse

	for _, product := range products {
		mappedProducts = append(mappedProducts, p.mapResponseProduct(product))
	}

	return mappedProducts
}

func (p *productProtoMapper) mapResponseProductDeleteAt(product *response.ProductResponseDeleteAt) *pbproduct.ProductResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if product.DeleteAt != nil {
		deletedAt = wrapperspb.String(*product.DeleteAt)
	}

	return &pbproduct.ProductResponseDeleteAt{
		Id:           int32(product.ID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  product.Description,
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        product.Brand,
		Weight:       int32(product.Weight),
		SlugProduct:  product.SlugProduct,
		ImageProduct: product.ImageProduct,
		Barcode:      product.Barcode,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func (p *productProtoMapper) mapResponsesProductDeleteAt(products []*response.ProductResponseDeleteAt) []*pbproduct.ProductResponseDeleteAt {
	var mappedProducts []*pbproduct.ProductResponseDeleteAt

	for _, product := range products {
		mappedProducts = append(mappedProducts, p.mapResponseProductDeleteAt(product))
	}

	return mappedProducts
}
