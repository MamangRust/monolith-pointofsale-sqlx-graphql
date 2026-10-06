package response_api

import (
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type merchantResponseMapper struct {
}

func NewMerchantResponseMapper() *merchantResponseMapper {
	return &merchantResponseMapper{}
}

func (m *merchantResponseMapper) ToResponseMerchant(merchant *pbmerchant.MerchantResponse) *response.MerchantResponse {
	return &response.MerchantResponse{
		ID:           int(merchant.Id),
		UserID:       int(merchant.UserId),
		Name:         merchant.Name,
		Description:  merchant.Description,
		Address:      merchant.Address,
		ContactEmail: merchant.ContactEmail,
		ContactPhone: merchant.ContactPhone,
		Status:       merchant.Status,
		CreatedAt:    merchant.CreatedAt,
		UpdatedAt:    merchant.UpdatedAt,
	}
}

func (m *merchantResponseMapper) ToResponsesMerchant(merchants []*pbmerchant.MerchantResponse) []*response.MerchantResponse {
	var mappedMerchants []*response.MerchantResponse

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.ToResponseMerchant(merchant))
	}

	return mappedMerchants
}

func (m *merchantResponseMapper) ToResponseMerchantDeleteAt(merchant *pbmerchant.MerchantResponseDeleteAt) *response.MerchantResponseDeleteAt {
	return &response.MerchantResponseDeleteAt{
		ID:           int(merchant.Id),
		UserID:       int(merchant.UserId),
		Name:         merchant.Name,
		Description:  merchant.Description,
		Address:      merchant.Address,
		ContactEmail: merchant.ContactEmail,
		ContactPhone: merchant.ContactPhone,
		Status:       merchant.Status,
		CreatedAt:    merchant.CreatedAt,
		UpdatedAt:    merchant.UpdatedAt,
		DeletedAt:    merchant.DeletedAt,
	}
}

func (m *merchantResponseMapper) ToResponsesMerchantDeleteAt(merchants []*pbmerchant.MerchantResponseDeleteAt) []*response.MerchantResponseDeleteAt {
	var mappedMerchants []*response.MerchantResponseDeleteAt

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.ToResponseMerchantDeleteAt(merchant))
	}

	return mappedMerchants
}

func (m *merchantResponseMapper) ToApiResponseMerchant(pbResponse *pbmerchant.ApiResponseMerchant) *response.ApiResponseMerchant {
	return &response.ApiResponseMerchant{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchant(pbResponse.Data),
	}
}

func (m *merchantResponseMapper) ToApiResponseMerchantDeleteAt(pbResponse *pbmerchant.ApiResponseMerchantDeleteAt) *response.ApiResponseMerchantDeleteAt {
	return &response.ApiResponseMerchantDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantDeleteAt(pbResponse.Data),
	}
}

func (m *merchantResponseMapper) ToApiResponsesMerchant(pbResponse *pbmerchant.ApiResponsesMerchant) *response.ApiResponsesMerchant {
	return &response.ApiResponsesMerchant{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesMerchant(pbResponse.Data),
	}
}

func (m *merchantResponseMapper) ToApiResponseMerchantDelete(pbResponse *pbmerchant.ApiResponseMerchantDelete) *response.ApiResponseMerchantDelete {
	return &response.ApiResponseMerchantDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (m *merchantResponseMapper) ToApiResponseMerchantAll(pbResponse *pbmerchant.ApiResponseMerchantAll) *response.ApiResponseMerchantAll {
	return &response.ApiResponseMerchantAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (m *merchantResponseMapper) ToApiResponsePaginationMerchantDeleteAt(pbResponse *pbmerchant.ApiResponsePaginationMerchantDeleteAt) *response.ApiResponsePaginationMerchantDeleteAt {
	return &response.ApiResponsePaginationMerchantDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchantDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *merchantResponseMapper) ToApiResponsePaginationMerchant(pbResponse *pbmerchant.ApiResponsePaginationMerchant) *response.ApiResponsePaginationMerchant {
	return &response.ApiResponsePaginationMerchant{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchant(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}
