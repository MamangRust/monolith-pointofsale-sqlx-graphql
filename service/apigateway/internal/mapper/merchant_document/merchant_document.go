package merchantdocumentgraphqlmapper

import (
	graphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper"
	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-pointofsale-pb/merchant_document"
)

type merchantDocumentGraphqlMapper struct {
}

func NewMerchantDocumentGraphqlMapper() *merchantDocumentGraphqlMapper {
	return &merchantDocumentGraphqlMapper{}
}

func (m *merchantDocumentGraphqlMapper) ToGraphqlResponseMerchantDocument(res *pb.ApiResponseMerchantDocument) *model.APIResponseMerchantDocument {
	return &model.APIResponseMerchantDocument{
		Status:  res.Status,
		Message: res.Message,
		Data:    m.mapResponseMerchantDocument(res.Data),
	}
}

func (m *merchantDocumentGraphqlMapper) ToGraphqlResponseMerchantDocumentDeleteAt(res *pb.ApiResponseMerchantDocument) *model.APIResponseMerchantDocumentDeleteAt {
	var data *model.MerchantDocumentResponseDeleteAt
	if res.Data != nil {
		var deletedAt string = ""
		data = &model.MerchantDocumentResponseDeleteAt{
			DocumentID:   res.Data.DocumentId,
			MerchantID:   res.Data.MerchantId,
			DocumentType: res.Data.DocumentType,
			DocumentURL:  res.Data.DocumentUrl,
			Status:       res.Data.Status,
			Note:         res.Data.Note,
			UploadedAt:   res.Data.UploadedAt,
			UpdatedAt:    res.Data.UpdatedAt,
			DeletedAt:    &deletedAt,
		}
	}
	return &model.APIResponseMerchantDocumentDeleteAt{
		Status:  res.Status,
		Message: res.Message,
		Data:    data,
	}
}

func (m *merchantDocumentGraphqlMapper) ToGraphqlResponseDelete(res *pb.ApiResponseMerchantDocumentDelete) *model.APIResponseMerchantDocumentDelete {
	return &model.APIResponseMerchantDocumentDelete{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (m *merchantDocumentGraphqlMapper) ToGraphqlResponseAll(res *pb.ApiResponseMerchantDocumentAll) *model.APIResponseMerchantDocumentAll {
	return &model.APIResponseMerchantDocumentAll{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (m *merchantDocumentGraphqlMapper) ToGraphqlResponsePaginationMerchantDocument(res *pb.ApiResponsePaginationMerchantDocument) *model.APIResponsePaginationMerchantDocument {
	return &model.APIResponsePaginationMerchantDocument{
		Status:     res.Status,
		Message:    res.Message,
		Data:       m.mapResponsesMerchantDocument(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (m *merchantDocumentGraphqlMapper) ToGraphqlResponsePaginationMerchantDocumentDeleteAt(res *pb.ApiResponsePaginationMerchantDocumentAt) *model.APIResponsePaginationMerchantDocumentAt {
	return &model.APIResponsePaginationMerchantDocumentAt{
		Status:     res.Status,
		Message:    res.Message,
		Data:       m.mapResponsesMerchantDocumentDeleteAt(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (m *merchantDocumentGraphqlMapper) mapResponseMerchantDocument(doc *pb.MerchantDocument) *model.MerchantDocumentResponse {
	if doc == nil {
		return nil
	}
	return &model.MerchantDocumentResponse{
		DocumentID:   doc.DocumentId,
		MerchantID:   doc.MerchantId,
		DocumentType: doc.DocumentType,
		DocumentURL:  doc.DocumentUrl,
		Status:       doc.Status,
		Note:         doc.Note,
		UploadedAt:   doc.UploadedAt,
		UpdatedAt:    doc.UpdatedAt,
	}
}

func (m *merchantDocumentGraphqlMapper) mapResponsesMerchantDocument(docs []*pb.MerchantDocument) []*model.MerchantDocumentResponse {
	var responses []*model.MerchantDocumentResponse
	for _, doc := range docs {
		responses = append(responses, m.mapResponseMerchantDocument(doc))
	}
	return responses
}

func (m *merchantDocumentGraphqlMapper) mapResponseMerchantDocumentDeleteAt(doc *pb.MerchantDocumentDeleteAt) *model.MerchantDocumentResponseDeleteAt {
	if doc == nil {
		return nil
	}
	var deletedAt string = ""
	if doc.DeletedAt != nil {
		deletedAt = doc.DeletedAt.Value
	}
	return &model.MerchantDocumentResponseDeleteAt{
		DocumentID:   doc.DocumentId,
		MerchantID:   doc.MerchantId,
		DocumentType: doc.DocumentType,
		DocumentURL:  doc.DocumentUrl,
		Status:       doc.Status,
		Note:         doc.Note,
		UploadedAt:   doc.UploadedAt,
		UpdatedAt:    doc.UpdatedAt,
		DeletedAt:    &deletedAt,
	}
}

func (m *merchantDocumentGraphqlMapper) mapResponsesMerchantDocumentDeleteAt(docs []*pb.MerchantDocumentDeleteAt) []*model.MerchantDocumentResponseDeleteAt {
	var responses []*model.MerchantDocumentResponseDeleteAt
	for _, doc := range docs {
		responses = append(responses, m.mapResponseMerchantDocumentDeleteAt(doc))
	}
	return responses
}

func (m *merchantDocumentGraphqlMapper) ToGraphqlResponsePaginationMerchantDocumentActive(res *pb.ApiResponsePaginationMerchantDocument) *model.APIResponsePaginationMerchantDocumentAt {
	var responses []*model.MerchantDocumentResponseDeleteAt
	for _, doc := range res.Data {
		if doc == nil {
			continue
		}
		responses = append(responses, &model.MerchantDocumentResponseDeleteAt{
			DocumentID:   doc.DocumentId,
			MerchantID:   doc.MerchantId,
			DocumentType: doc.DocumentType,
			DocumentURL:  doc.DocumentUrl,
			Status:       doc.Status,
			Note:         doc.Note,
			UploadedAt:   doc.UploadedAt,
			UpdatedAt:    doc.UpdatedAt,
			DeletedAt:    nil,
		})
	}
	return &model.APIResponsePaginationMerchantDocumentAt{
		Status:     res.Status,
		Message:    res.Message,
		Data:       responses,
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}
