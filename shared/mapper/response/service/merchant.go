package response_service

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/record"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type merchantResponseMapper struct {
}

func NewMerchantResponseMapper() *merchantResponseMapper {
	return &merchantResponseMapper{}
}

func (s *merchantResponseMapper) ToMerchantResponse(merchant *record.MerchantRecord) *response.MerchantResponse {
	return &response.MerchantResponse{
		ID:           merchant.ID,
		UserID:       merchant.UserID,
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

func (s *merchantResponseMapper) ToMerchantsResponse(merchants []*record.MerchantRecord) []*response.MerchantResponse {
	var responses []*response.MerchantResponse

	for _, merchant := range merchants {
		responses = append(responses, s.ToMerchantResponse(merchant))
	}

	return responses
}

func (s *merchantResponseMapper) ToMerchantResponseDeleteAt(merchant *record.MerchantRecord) *response.MerchantResponseDeleteAt {
	return &response.MerchantResponseDeleteAt{
		ID:           merchant.ID,
		UserID:       merchant.UserID,
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

func (s *merchantResponseMapper) ToMerchantsResponseDeleteAt(merchants []*record.MerchantRecord) []*response.MerchantResponseDeleteAt {
	var responses []*response.MerchantResponseDeleteAt

	for _, merchant := range merchants {
		responses = append(responses, s.ToMerchantResponseDeleteAt(merchant))
	}

	return responses
}
