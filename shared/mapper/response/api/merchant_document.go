package response_api

import (
	pbmerchant_document "github.com/MamangRust/monolith-pointofsale-pb/merchant_document"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type merchantDocumentResponse struct{}

func NewMerchantDocumentResponseMapper() *merchantDocumentResponse {
	return &merchantDocumentResponse{}
}

func (m *merchantDocumentResponse) ToApiResponseMerchantDocument(doc *pbmerchant_document.ApiResponseMerchantDocument) *response.ApiResponseMerchantDocument {
	return &response.ApiResponseMerchantDocument{
		Status:  doc.Status,
		Message: doc.Message,
		Data:    m.mapMerchantDocument(doc.Data),
	}
}

func (m *merchantDocumentResponse) ToApiResponsesMerchantDocument(docs *pbmerchant_document.ApiResponsesMerchantDocument) *response.ApiResponsesMerchantDocument {
	return &response.ApiResponsesMerchantDocument{
		Status:  docs.Status,
		Message: docs.Message,
		Data:    m.mapMerchantDocuments(docs.Data),
	}
}

func (m *merchantDocumentResponse) ToApiResponsePaginationMerchantDocument(docs *pbmerchant_document.ApiResponsePaginationMerchantDocument) *response.ApiResponsePaginationMerchantDocument {
	return &response.ApiResponsePaginationMerchantDocument{
		Status:     docs.Status,
		Message:    docs.Message,
		Data:       m.mapMerchantDocuments(docs.Data),
		Pagination: mapPaginationMeta(docs.Pagination),
	}
}

func (m *merchantDocumentResponse) ToApiResponsePaginationMerchantDocumentDeleteAt(docs *pbmerchant_document.ApiResponsePaginationMerchantDocumentAt) *response.ApiResponsePaginationMerchantDocumentDeleteAt {
	return &response.ApiResponsePaginationMerchantDocumentDeleteAt{
		Status:     docs.Status,
		Message:    docs.Message,
		Data:       m.mapMerchantDocumentsDeletedAt(docs.Data),
		Pagination: mapPaginationMeta(docs.Pagination),
	}
}

func (m *merchantDocumentResponse) ToApiResponseMerchantDocumentAll(resp *pbmerchant_document.ApiResponseMerchantDocumentAll) *response.ApiResponseMerchantDocumentAll {
	return &response.ApiResponseMerchantDocumentAll{
		Status:  resp.Status,
		Message: resp.Message,
	}
}

func (m *merchantDocumentResponse) ToApiResponseMerchantDocumentDeleteAt(resp *pbmerchant_document.ApiResponseMerchantDocumentDelete) *response.ApiResponseMerchantDocumentDelete {
	return &response.ApiResponseMerchantDocumentDelete{
		Status:  resp.Status,
		Message: resp.Message,
	}
}

func (m *merchantDocumentResponse) mapMerchantDocument(doc *pbmerchant_document.MerchantDocument) *response.MerchantDocumentResponse {
	if doc == nil {
		return nil
	}
	return &response.MerchantDocumentResponse{
		ID:           int(doc.DocumentId),
		MerchantID:   int(doc.MerchantId),
		DocumentType: doc.DocumentType,
		DocumentURL:  doc.DocumentUrl,
		Status:       doc.Status,
		Note:         doc.Note,
		CreatedAt:    doc.UploadedAt,
		UpdatedAt:    doc.UpdatedAt,
	}
}

func (m *merchantDocumentResponse) mapMerchantDocuments(docs []*pbmerchant_document.MerchantDocument) []*response.MerchantDocumentResponse {
	var responses []*response.MerchantDocumentResponse
	for _, doc := range docs {
		responses = append(responses, m.mapMerchantDocument(doc))
	}
	return responses
}

func (m *merchantDocumentResponse) mapMerchantDocumentDeletedAt(doc *pbmerchant_document.MerchantDocumentDeleteAt) *response.MerchantDocumentResponseDeleteAt {
	if doc == nil {
		return nil
	}
	var deletedAt *string

	if doc.DeletedAt != nil {
		deletedAt = &doc.DeletedAt.Value
	}

	return &response.MerchantDocumentResponseDeleteAt{
		ID:           int(doc.DocumentId),
		MerchantID:   int(doc.MerchantId),
		DocumentType: doc.DocumentType,
		DocumentURL:  doc.DocumentUrl,
		Status:       doc.Status,
		Note:         doc.Note,
		CreatedAt:    doc.UploadedAt,
		UpdatedAt:    doc.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func (m *merchantDocumentResponse) mapMerchantDocumentsDeletedAt(docs []*pbmerchant_document.MerchantDocumentDeleteAt) []*response.MerchantDocumentResponseDeleteAt {
	var responses []*response.MerchantDocumentResponseDeleteAt
	for _, doc := range docs {
		responses = append(responses, m.mapMerchantDocumentDeletedAt(doc))
	}
	return responses
}
