package merchantgraphqlmapper

import (
	graphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper"
	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-pointofsale-pb/merchant"
)

type merchantGraphqlMapper struct {
}

func NewMerchantGraphqlMapper() *merchantGraphqlMapper {
	return &merchantGraphqlMapper{}
}

func (m *merchantGraphqlMapper) ToGraphqlResponseMerchant(res *pb.ApiResponseMerchant) *model.APIResponseMerchant {
	return &model.APIResponseMerchant{
		Status:  res.Status,
		Message: res.Message,
		Data:    m.mapResponseMerchant(res.Data),
	}
}

func (m *merchantGraphqlMapper) ToGraphqlResponsesMerchant(res *pb.ApiResponsesMerchant) *model.APIResponsesMerchant {
	return &model.APIResponsesMerchant{
		Status:  res.Status,
		Message: res.Message,
		Data:    m.mapResponsesMerchant(res.Data),
	}
}

func (m *merchantGraphqlMapper) ToGraphqlResponseMerchantDeleteAt(res *pb.ApiResponseMerchantDeleteAt) *model.APIResponseMerchantDeleteAt {
	return &model.APIResponseMerchantDeleteAt{
		Status:  res.Status,
		Message: res.Message,
		Data:    m.mapResponseMerchantDeleteAt(res.Data),
	}
}

func (m *merchantGraphqlMapper) ToGraphqlResponseMerchantDelete(res *pb.ApiResponseMerchantDelete) *model.APIResponseMerchantDelete {
	return &model.APIResponseMerchantDelete{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (m *merchantGraphqlMapper) ToGraphqlResponseMerchantAll(res *pb.ApiResponseMerchantAll) *model.APIResponseMerchantAll {
	return &model.APIResponseMerchantAll{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (m *merchantGraphqlMapper) ToGraphqlResponsePaginationMerchant(res *pb.ApiResponsePaginationMerchant) *model.APIResponsePaginationMerchant {
	return &model.APIResponsePaginationMerchant{
		Status:     res.Status,
		Message:    res.Message,
		Data:       m.mapResponsesMerchant(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (m *merchantGraphqlMapper) ToGraphqlResponsePaginationMerchantDeleteAt(res *pb.ApiResponsePaginationMerchantDeleteAt) *model.APIResponsePaginationMerchantDeleteAt {
	return &model.APIResponsePaginationMerchantDeleteAt{
		Status:     res.Status,
		Message:    res.Message,
		Data:       m.mapResponsesMerchantDeleteAt(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (m *merchantGraphqlMapper) mapResponseMerchant(merchant *pb.MerchantResponse) *model.MerchantResponse {
	if merchant == nil {
		return nil
	}
	return &model.MerchantResponse{
		ID:           int32(merchant.Id),
		UserID:       int32(merchant.UserId),
		Name:         merchant.Name,
		Description:  &merchant.Description,
		Address:      &merchant.Address,
		ContactEmail: &merchant.ContactEmail,
		ContactPhone: &merchant.ContactPhone,
		Status:       &merchant.Status,
		CreatedAt:    &merchant.CreatedAt,
		UpdatedAt:    &merchant.UpdatedAt,
	}
}

func (m *merchantGraphqlMapper) mapResponsesMerchant(merchants []*pb.MerchantResponse) []*model.MerchantResponse {
	var responses []*model.MerchantResponse
	for _, merchant := range merchants {
		responses = append(responses, m.mapResponseMerchant(merchant))
	}
	return responses
}

func (m *merchantGraphqlMapper) mapResponseMerchantDeleteAt(merchant *pb.MerchantResponseDeleteAt) *model.MerchantResponseDeleteAt {
	if merchant == nil {
		return nil
	}
	return &model.MerchantResponseDeleteAt{
		ID:           int32(merchant.Id),
		UserID:       int32(merchant.UserId),
		Name:         merchant.Name,
		Description:  &merchant.Description,
		Address:      &merchant.Address,
		ContactEmail: &merchant.ContactEmail,
		ContactPhone: &merchant.ContactPhone,
		Status:       &merchant.Status,
		CreatedAt:    &merchant.CreatedAt,
		UpdatedAt:    &merchant.UpdatedAt,
		DeletedAt:    &merchant.DeletedAt,
	}
}

func (m *merchantGraphqlMapper) mapResponsesMerchantDeleteAt(merchants []*pb.MerchantResponseDeleteAt) []*model.MerchantResponseDeleteAt {
	var responses []*model.MerchantResponseDeleteAt
	for _, merchant := range merchants {
		responses = append(responses, m.mapResponseMerchantDeleteAt(merchant))
	}
	return responses
}

func (m *merchantGraphqlMapper) ToGraphqlResponseMerchantRestore(res *pb.ApiResponseMerchant) *model.APIResponseMerchantDeleteAt {
	var deletedAt string = ""
	var data *model.MerchantResponseDeleteAt
	if res.Data != nil {
		data = &model.MerchantResponseDeleteAt{
			ID:           int32(res.Data.Id),
			UserID:       int32(res.Data.UserId),
			Name:         res.Data.Name,
			Description:  &res.Data.Description,
			Address:      &res.Data.Address,
			ContactEmail: &res.Data.ContactEmail,
			ContactPhone: &res.Data.ContactPhone,
			Status:       &res.Data.Status,
			CreatedAt:    &res.Data.CreatedAt,
			UpdatedAt:    &res.Data.UpdatedAt,
			DeletedAt:    &deletedAt,
		}
	}
	return &model.APIResponseMerchantDeleteAt{
		Status:  res.Status,
		Message: res.Message,
		Data:    data,
	}
}
