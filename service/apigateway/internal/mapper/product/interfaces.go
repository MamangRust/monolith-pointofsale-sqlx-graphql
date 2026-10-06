package productgraphqlmapper

import (
	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-pointofsale-pb/product"
)

type ProductGraphqlMapper interface {
	ToGraphqlResponseProduct(res *pb.ApiResponseProduct) *model.APIResponseProduct
	ToGraphqlResponsesProduct(res *pb.ApiResponsesProduct) *model.APIResponsesProduct
	ToGraphqlResponseProductDeleteAt(res *pb.ApiResponseProductDeleteAt) *model.APIResponseProductDeleteAt
	ToGraphqlResponseProductDelete(res *pb.ApiResponseProductDelete) *model.APIResponseProductDelete
	ToGraphqlResponseProductAll(res *pb.ApiResponseProductAll) *model.APIResponseProductAll
	ToGraphqlResponsePaginationProduct(res *pb.ApiResponsePaginationProduct) *model.APIResponsePaginationProduct
	ToGraphqlResponsePaginationProductDeleteAt(res *pb.ApiResponsePaginationProductDeleteAt) *model.APIResponsePaginationProductDeleteAt
}
