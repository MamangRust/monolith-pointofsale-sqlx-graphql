package response_service

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/record"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type merchantDocumentResponseMapper struct{}

func NewMerchantDocumentResponseMapper() *merchantDocumentResponseMapper {
	return &merchantDocumentResponseMapper{}
}

func (s *merchantDocumentResponseMapper) ToMerchantDocumentResponse(doc *record.MerchantDocumentRecord) *response.MerchantDocumentResponse {
	return &response.MerchantDocumentResponse{
		ID:           doc.ID,
		MerchantID:   doc.MerchantID,
		DocumentType: doc.DocumentType,
		DocumentURL:  doc.DocumentURL,
		Status:       doc.Status,
		Note:         doc.Note,
		CreatedAt:    doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
	}
}

func (s *merchantDocumentResponseMapper) ToMerchantDocumentsResponse(docs []*record.MerchantDocumentRecord) []*response.MerchantDocumentResponse {
	var responses []*response.MerchantDocumentResponse
	for _, doc := range docs {
		responses = append(responses, s.ToMerchantDocumentResponse(doc))
	}
	return responses
}

func (s *merchantDocumentResponseMapper) ToMerchantDocumentResponseDeleteAt(doc *record.MerchantDocumentRecord) *response.MerchantDocumentResponseDeleteAt {
	return &response.MerchantDocumentResponseDeleteAt{
		ID:           doc.ID,
		MerchantID:   doc.MerchantID,
		DocumentType: doc.DocumentType,
		DocumentURL:  doc.DocumentURL,
		Status:       doc.Status,
		Note:         doc.Note,
		CreatedAt:    doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
		DeletedAt:    doc.DeletedAt,
	}
}

func (s *merchantDocumentResponseMapper) ToMerchantDocumentsResponseDeleteAt(docs []*record.MerchantDocumentRecord) []*response.MerchantDocumentResponseDeleteAt {
	var responses []*response.MerchantDocumentResponseDeleteAt
	for _, doc := range docs {
		responses = append(responses, s.ToMerchantDocumentResponseDeleteAt(doc))
	}
	return responses
}
