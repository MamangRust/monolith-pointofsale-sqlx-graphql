package productgraphqlmapper

import (
	graphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper"
	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-pointofsale-pb/product"
)

type productGraphqlMapper struct {
}

func NewProductGraphqlMapper() *productGraphqlMapper {
	return &productGraphqlMapper{}
}

func (p *productGraphqlMapper) ToGraphqlResponseProduct(res *pb.ApiResponseProduct) *model.APIResponseProduct {
	return &model.APIResponseProduct{
		Status:  res.Status,
		Message: res.Message,
		Data:    p.mapResponseProduct(res.Data),
	}
}

func (p *productGraphqlMapper) ToGraphqlResponsesProduct(res *pb.ApiResponsesProduct) *model.APIResponsesProduct {
	return &model.APIResponsesProduct{
		Status:  res.Status,
		Message: res.Message,
		Data:    p.mapResponsesProduct(res.Data),
	}
}

func (p *productGraphqlMapper) ToGraphqlResponseProductDeleteAt(res *pb.ApiResponseProductDeleteAt) *model.APIResponseProductDeleteAt {
	return &model.APIResponseProductDeleteAt{
		Status:  res.Status,
		Message: res.Message,
		Data:    p.mapResponseProductDeleteAt(res.Data),
	}
}

func (p *productGraphqlMapper) ToGraphqlResponseProductDelete(res *pb.ApiResponseProductDelete) *model.APIResponseProductDelete {
	return &model.APIResponseProductDelete{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (p *productGraphqlMapper) ToGraphqlResponseProductAll(res *pb.ApiResponseProductAll) *model.APIResponseProductAll {
	return &model.APIResponseProductAll{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (p *productGraphqlMapper) ToGraphqlResponsePaginationProduct(res *pb.ApiResponsePaginationProduct) *model.APIResponsePaginationProduct {
	return &model.APIResponsePaginationProduct{
		Status:     res.Status,
		Message:    res.Message,
		Data:       p.mapResponsesProduct(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (p *productGraphqlMapper) ToGraphqlResponsePaginationProductDeleteAt(res *pb.ApiResponsePaginationProductDeleteAt) *model.APIResponsePaginationProductDeleteAt {
	return &model.APIResponsePaginationProductDeleteAt{
		Status:     res.Status,
		Message:    res.Message,
		Data:       p.mapResponsesProductDeleteAt(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (p *productGraphqlMapper) mapResponseProduct(product *pb.ProductResponse) *model.ProductResponse {
	if product == nil {
		return nil
	}
	rating := float64(product.Rating)
	return &model.ProductResponse{
		ID:           int32(product.Id),
		CategoryID:   int32(product.CategoryId),
		MerchantID:   int32(product.MerchantId),
		Name:         product.Name,
		Description:  &product.Description,
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        &product.Brand,
		Weight:       &product.Weight,
		Rating:       &rating,
		SlugProduct:  &product.SlugProduct,
		ImageProduct: &product.ImageProduct,
		Barcode:      &product.Barcode,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}

func (p *productGraphqlMapper) mapResponsesProduct(products []*pb.ProductResponse) []*model.ProductResponse {
	var responses []*model.ProductResponse
	for _, product := range products {
		responses = append(responses, p.mapResponseProduct(product))
	}
	return responses
}

func (p *productGraphqlMapper) mapResponseProductDeleteAt(product *pb.ProductResponseDeleteAt) *model.ProductResponseDeleteAt {
	if product == nil {
		return nil
	}
	var deletedAt *string
	if product.DeletedAt != nil {
		deletedAt = &product.DeletedAt.Value
	}
	rating := float64(product.Rating)

	return &model.ProductResponseDeleteAt{
		ID:           int32(product.Id),
		CategoryID:   int32(product.CategoryId),
		MerchantID:   int32(product.MerchantId),
		Name:         product.Name,
		Description:  &product.Description,
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        &product.Brand,
		Weight:       &product.Weight,
		Rating:       &rating,
		SlugProduct:  &product.SlugProduct,
		ImageProduct: &product.ImageProduct,
		Barcode:      &product.Barcode,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func (p *productGraphqlMapper) mapResponsesProductDeleteAt(products []*pb.ProductResponseDeleteAt) []*model.ProductResponseDeleteAt {
	var responses []*model.ProductResponseDeleteAt
	for _, product := range products {
		responses = append(responses, p.mapResponseProductDeleteAt(product))
	}
	return responses
}
