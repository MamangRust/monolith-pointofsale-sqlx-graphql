package protomapper

import (
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type merchantProtoMapper struct{}

func NewMerchantProtoMaper() *merchantProtoMapper {
	return &merchantProtoMapper{}
}

func (m *merchantProtoMapper) ToProtoResponseMerchant(status string, message string, pbResponse *response.MerchantResponse) *pbmerchant.ApiResponseMerchant {
	return &pbmerchant.ApiResponseMerchant{
		Status:  status,
		Message: message,
		Data:    m.mapResponseMerchant(pbResponse),
	}
}

func (m *merchantProtoMapper) ToProtoResponsesMerchant(status string, message string, pbResponse []*response.MerchantResponse) *pbmerchant.ApiResponsesMerchant {
	return &pbmerchant.ApiResponsesMerchant{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesMerchant(pbResponse),
	}
}

func (m *merchantProtoMapper) ToProtoResponseMerchantDeleteAt(status string, message string, pbResponse *response.MerchantResponseDeleteAt) *pbmerchant.ApiResponseMerchantDeleteAt {
	return &pbmerchant.ApiResponseMerchantDeleteAt{
		Status:  status,
		Message: message,
		Data:    m.mapResponseMerchantDeleteAt(pbResponse),
	}
}

func (m *merchantProtoMapper) ToProtoResponseMerchantDelete(status string, message string) *pbmerchant.ApiResponseMerchantDelete {
	return &pbmerchant.ApiResponseMerchantDelete{
		Status:  status,
		Message: message,
	}
}

func (m *merchantProtoMapper) ToProtoResponseMerchantAll(status string, message string) *pbmerchant.ApiResponseMerchantAll {
	return &pbmerchant.ApiResponseMerchantAll{
		Status:  status,
		Message: message,
	}
}

func (m *merchantProtoMapper) ToProtoResponsePaginationMerchantDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, merchants []*response.MerchantResponseDeleteAt) *pbmerchant.ApiResponsePaginationMerchantDeleteAt {
	return &pbmerchant.ApiResponsePaginationMerchantDeleteAt{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesMerchantDeleteAt(merchants),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantProtoMapper) ToProtoResponsePaginationMerchant(pagination *pbcommon.PaginationMeta, status string, message string, merchants []*response.MerchantResponse) *pbmerchant.ApiResponsePaginationMerchant {
	return &pbmerchant.ApiResponsePaginationMerchant{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesMerchant(merchants),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantProtoMapper) mapResponseMerchant(merchant *response.MerchantResponse) *pbmerchant.MerchantResponse {
	return &pbmerchant.MerchantResponse{
		Id:           int32(merchant.ID),
		UserId:       int32(merchant.UserID),
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

func (m *merchantProtoMapper) mapResponsesMerchant(merchants []*response.MerchantResponse) []*pbmerchant.MerchantResponse {
	var mappedMerchants []*pbmerchant.MerchantResponse

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.mapResponseMerchant(merchant))
	}

	return mappedMerchants
}

func (m *merchantProtoMapper) mapResponseMerchantDeleteAt(merchant *response.MerchantResponseDeleteAt) *pbmerchant.MerchantResponseDeleteAt {
	return &pbmerchant.MerchantResponseDeleteAt{
		Id:           int32(merchant.ID),
		UserId:       int32(merchant.UserID),
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

func (m *merchantProtoMapper) mapResponsesMerchantDeleteAt(merchants []*response.MerchantResponseDeleteAt) []*pbmerchant.MerchantResponseDeleteAt {
	var mappedMerchants []*pbmerchant.MerchantResponseDeleteAt

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.mapResponseMerchantDeleteAt(merchant))
	}

	return mappedMerchants
}
