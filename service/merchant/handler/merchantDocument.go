package handler

import (
	pbmerchant_document "github.com/MamangRust/monolith-pointofsale-pb/merchant_document"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func mapMerchantDocument(doc *db.MerchantDocument) *pbmerchant_document.MerchantDocument {
	if doc == nil {
		return nil
	}
	return &pbmerchant_document.MerchantDocument{
		DocumentId:   int32(doc.DocumentID),
		MerchantId:   int32(doc.MerchantID),
		DocumentType: doc.DocumentType,
		DocumentUrl:  doc.DocumentUrl,
		Status:       doc.Status,
		Note:         mapSqlNullString(doc.Note),
		UploadedAt:   mapSqlNullTime(doc.CreatedAt),
		UpdatedAt:    mapSqlNullTime(doc.UpdatedAt),
	}
}

func mapResponsesGetMerchantDocumentsRow(docs []*db.GetMerchantDocumentsRow) []*pbmerchant_document.MerchantDocument {
	var res []*pbmerchant_document.MerchantDocument
	for _, doc := range docs {
		res = append(res, &pbmerchant_document.MerchantDocument{
			DocumentId:   int32(doc.DocumentID),
			MerchantId:   int32(doc.MerchantID),
			DocumentType: doc.DocumentType,
			DocumentUrl:  doc.DocumentUrl,
			Status:       doc.Status,
			Note:         mapSqlNullString(doc.Note),
			UploadedAt:   mapSqlNullTime(doc.CreatedAt),
			UpdatedAt:    mapSqlNullTime(doc.UpdatedAt),
		})
	}
	return res
}

func mapResponsesGetActiveMerchantDocumentsRow(docs []*db.GetActiveMerchantDocumentsRow) []*pbmerchant_document.MerchantDocument {
	var res []*pbmerchant_document.MerchantDocument
	for _, doc := range docs {
		res = append(res, &pbmerchant_document.MerchantDocument{
			DocumentId:   int32(doc.DocumentID),
			MerchantId:   int32(doc.MerchantID),
			DocumentType: doc.DocumentType,
			DocumentUrl:  doc.DocumentUrl,
			Status:       doc.Status,
			Note:         mapSqlNullString(doc.Note),
			UploadedAt:   mapSqlNullTime(doc.CreatedAt),
			UpdatedAt:    mapSqlNullTime(doc.UpdatedAt),
		})
	}
	return res
}

func mapMerchantDocumentDeleteAt(doc *db.MerchantDocument) *pbmerchant_document.MerchantDocumentDeleteAt {
	if doc == nil {
		return nil
	}
	var deletedAt *wrapperspb.StringValue
	if doc.DeletedAt.Valid {
		deletedAt = wrapperspb.String(doc.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}

	return &pbmerchant_document.MerchantDocumentDeleteAt{
		DocumentId:   int32(doc.DocumentID),
		MerchantId:   int32(doc.MerchantID),
		DocumentType: doc.DocumentType,
		DocumentUrl:  doc.DocumentUrl,
		Status:       doc.Status,
		Note:         mapSqlNullString(doc.Note),
		UploadedAt:   mapSqlNullTime(doc.CreatedAt),
		UpdatedAt:    mapSqlNullTime(doc.UpdatedAt),
		DeletedAt:    deletedAt,
	}
}

func mapResponsesGetTrashedMerchantDocumentsRow(docs []*db.GetTrashedMerchantDocumentsRow) []*pbmerchant_document.MerchantDocumentDeleteAt {
	var res []*pbmerchant_document.MerchantDocumentDeleteAt
	for _, doc := range docs {
		var deletedAt *wrapperspb.StringValue
		if doc.DeletedAt.Valid {
			deletedAt = wrapperspb.String(doc.DeletedAt.Time.Format("2006-01-02 15:04:05"))
		}
		res = append(res, &pbmerchant_document.MerchantDocumentDeleteAt{
			DocumentId:   int32(doc.DocumentID),
			MerchantId:   int32(doc.MerchantID),
			DocumentType: doc.DocumentType,
			DocumentUrl:  doc.DocumentUrl,
			Status:       doc.Status,
			Note:         mapSqlNullString(doc.Note),
			UploadedAt:   mapSqlNullTime(doc.CreatedAt),
			UpdatedAt:    mapSqlNullTime(doc.UpdatedAt),
			DeletedAt:    deletedAt,
		})
	}
	return res
}
