package categorygraphqlmapper

import (
	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-pointofsale-pb/category"
)

type CategoryGraphqlMapper interface {
	ToGraphqlResponseCategory(res *pb.ApiResponseCategory) *model.APIResponseCategory
	ToGraphqlResponsesCategory(res *pb.ApiResponsesCategory) *model.APIResponsesCategory
	ToGraphqlResponseCategoryDeleteAt(res *pb.ApiResponseCategoryDeleteAt) *model.APIResponseCategoryDeleteAt
	ToGraphqlResponseCategoryDelete(res *pb.ApiResponseCategoryDelete) *model.APIResponseCategoryDelete
	ToGraphqlResponseCategoryAll(res *pb.ApiResponseCategoryAll) *model.APIResponseCategoryAll
	ToGraphqlResponsePaginationCategory(res *pb.ApiResponsePaginationCategory) *model.APIResponsePaginationCategory
	ToGraphqlResponseCategoryMonthlyTotalPrice(res *pb.ApiResponseCategoryMonthlyTotalPrice) *model.APIResponseCategoryMonthlyTotalPrice
	ToGraphqlResponseCategoryYearlyTotalPrice(res *pb.ApiResponseCategoryYearlyTotalPrice) *model.APIResponseCategoryYearlyTotalPrice
	ToGraphqlResponseCategoryMonthlyPrice(res *pb.ApiResponseCategoryMonthPrice) *model.APIResponseCategoryMonthPrice
	ToGraphqlResponseCategoryYearlyPrice(res *pb.ApiResponseCategoryYearPrice) *model.APIResponseCategoryYearPrice
	ToGraphqlResponsePaginationCategoryDeleteAt(res *pb.ApiResponsePaginationCategoryDeleteAt) *model.APIResponsePaginationCategoryDeleteAt
}
