package protomapper

import (
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbmerchant_document "github.com/MamangRust/monolith-pointofsale-pb/merchant_document"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type merchantDocumentProtoMapper struct{}

func NewMerchantDocumentProtoMapper() *merchantDocumentProtoMapper {
	return &merchantDocumentProtoMapper{}
}

func (m *merchantDocumentProtoMapper) ToProtoResponseMerchantDocument(status string, message string, doc *response.MerchantDocumentResponse) *pbmerchant_document.ApiResponseMerchantDocument {
	return &pbmerchant_document.ApiResponseMerchantDocument{
		Status:  status,
		Message: message,
		Data:    m.mapMerchantDocument(doc),
	}
}

func (m *merchantDocumentProtoMapper) ToProtoResponsesMerchantDocument(status string, message string, docs []*response.MerchantDocumentResponse) *pbmerchant_document.ApiResponsesMerchantDocument {
	return &pbmerchant_document.ApiResponsesMerchantDocument{
		Status:  status,
		Message: message,
		Data:    m.mapMerchantDocuments(docs),
	}
}

func (m *merchantDocumentProtoMapper) ToProtoResponsePaginationMerchantDocument(pagination *pbcommon.PaginationMeta, status string, message string, docs []*response.MerchantDocumentResponse) *pbmerchant_document.ApiResponsePaginationMerchantDocument {
	return &pbmerchant_document.ApiResponsePaginationMerchantDocument{
		Status:     status,
		Message:    message,
		Data:       m.mapMerchantDocuments(docs),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantDocumentProtoMapper) ToProtoResponsePaginationMerchantDocumentDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, docs []*response.MerchantDocumentResponseDeleteAt) *pbmerchant_document.ApiResponsePaginationMerchantDocumentAt {
	return &pbmerchant_document.ApiResponsePaginationMerchantDocumentAt{
		Status:     status,
		Message:    message,
		Data:       m.mapMerchantDocumentsDeleteAt(docs),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantDocumentProtoMapper) ToProtoResponseMerchantDocumentDelete(status string, message string) *pbmerchant_document.ApiResponseMerchantDocumentDelete {
	return &pbmerchant_document.ApiResponseMerchantDocumentDelete{
		Status:  status,
		Message: message,
	}
}

func (m *merchantDocumentProtoMapper) ToProtoResponseMerchantDocumentAll(status string, message string) *pbmerchant_document.ApiResponseMerchantDocumentAll {
	return &pbmerchant_document.ApiResponseMerchantDocumentAll{
		Status:  status,
		Message: message,
	}
}

func (m *merchantDocumentProtoMapper) mapMerchantDocument(doc *response.MerchantDocumentResponse) *pbmerchant_document.MerchantDocument {
	return &pbmerchant_document.MerchantDocument{
		DocumentId:   int32(doc.ID),
		MerchantId:   int32(doc.MerchantID),
		DocumentType: doc.DocumentType,
		DocumentUrl:  doc.DocumentURL,
		Status:       doc.Status,
		Note:         doc.Note,
		UploadedAt:   doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
	}
}

func (m *merchantDocumentProtoMapper) mapMerchantDocuments(docs []*response.MerchantDocumentResponse) []*pbmerchant_document.MerchantDocument {
	var res []*pbmerchant_document.MerchantDocument
	for _, doc := range docs {
		res = append(res, m.mapMerchantDocument(doc))
	}
	return res
}

func (m *merchantDocumentProtoMapper) mapMerchantDocumentDeleteAt(doc *response.MerchantDocumentResponseDeleteAt) *pbmerchant_document.MerchantDocumentDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if doc.DeletedAt != nil {
		deletedAt = wrapperspb.String(*doc.DeletedAt)
	}

	return &pbmerchant_document.MerchantDocumentDeleteAt{
		DocumentId:   int32(doc.ID),
		MerchantId:   int32(doc.MerchantID),
		DocumentType: doc.DocumentType,
		DocumentUrl:  doc.DocumentURL,
		Status:       doc.Status,
		Note:         doc.Note,
		UploadedAt:   doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func (m *merchantDocumentProtoMapper) mapMerchantDocumentsDeleteAt(docs []*response.MerchantDocumentResponseDeleteAt) []*pbmerchant_document.MerchantDocumentDeleteAt {
	var res []*pbmerchant_document.MerchantDocumentDeleteAt
	for _, doc := range docs {
		res = append(res, m.mapMerchantDocumentDeleteAt(doc))
	}
	return res
}
