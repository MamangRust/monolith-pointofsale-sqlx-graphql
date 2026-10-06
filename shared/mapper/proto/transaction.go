package protomapper

import (
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbtransaction "github.com/MamangRust/monolith-pointofsale-pb/transaction"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type transactionProtoMapper struct{}

func NewTransactionProtoMapper() *transactionProtoMapper {
	return &transactionProtoMapper{}
}

func (t *transactionProtoMapper) ToProtoResponseTransaction(status string, message string, trans *response.TransactionResponse) *pbtransaction.ApiResponseTransaction {
	return &pbtransaction.ApiResponseTransaction{
		Status:  status,
		Message: message,
		Data:    t.mapResponseTransaction(trans),
	}
}

func (t *transactionProtoMapper) ToProtoResponsesTransaction(status string, message string, transList []*response.TransactionResponse) *pbtransaction.ApiResponsesTransaction {
	return &pbtransaction.ApiResponsesTransaction{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransaction(transList),
	}
}

func (t *transactionProtoMapper) ToProtoResponseTransactionDeleteAt(status string, message string, trans *response.TransactionResponseDeleteAt) *pbtransaction.ApiResponseTransactionDeleteAt {
	return &pbtransaction.ApiResponseTransactionDeleteAt{
		Status:  status,
		Message: message,
		Data:    t.mapResponseTransactionDeleteAt(trans),
	}
}

func (t *transactionProtoMapper) ToProtoResponseTransactionDelete(status string, message string) *pbtransaction.ApiResponseTransactionDelete {
	return &pbtransaction.ApiResponseTransactionDelete{
		Status:  status,
		Message: message,
	}
}

func (t *transactionProtoMapper) ToProtoResponseTransactionAll(status string, message string) *pbtransaction.ApiResponseTransactionAll {
	return &pbtransaction.ApiResponseTransactionAll{
		Status:  status,
		Message: message,
	}
}

func (t *transactionProtoMapper) ToProtoResponsePaginationTransactionDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, transactions []*response.TransactionResponseDeleteAt) *pbtransaction.ApiResponsePaginationTransactionDeleteAt {
	return &pbtransaction.ApiResponsePaginationTransactionDeleteAt{
		Status:     status,
		Message:    message,
		Data:       t.mapResponsesTransactionDeleteAt(transactions),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (t *transactionProtoMapper) ToProtoResponsePaginationTransaction(pagination *pbcommon.PaginationMeta, status string, message string, transactions []*response.TransactionResponse) *pbtransaction.ApiResponsePaginationTransaction {
	return &pbtransaction.ApiResponsePaginationTransaction{
		Status:     status,
		Message:    message,
		Data:       t.mapResponsesTransaction(transactions),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (t *transactionProtoMapper) ToProtoResponseMonthAmountSuccess(status string, message string, row []*response.TransactionMonthlyAmountSuccessResponse) *pbtransaction.ApiResponseTransactionMonthAmountSuccess {
	return &pbtransaction.ApiResponseTransactionMonthAmountSuccess{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionMonthlyAmountSuccess(row),
	}
}

func (t *transactionProtoMapper) ToProtoResponseYearAmountSuccess(status string, message string, row []*response.TransactionYearlyAmountSuccessResponse) *pbtransaction.ApiResponseTransactionYearAmountSuccess {
	return &pbtransaction.ApiResponseTransactionYearAmountSuccess{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionYearlyAmountSuccess(row),
	}
}

func (t *transactionProtoMapper) ToProtoResponseMonthAmountFailed(status string, message string, row []*response.TransactionMonthlyAmountFailedResponse) *pbtransaction.ApiResponseTransactionMonthAmountFailed {
	return &pbtransaction.ApiResponseTransactionMonthAmountFailed{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionMonthlyAmountFailed(row),
	}
}

func (t *transactionProtoMapper) ToProtoResponseYearAmountFailed(status string, message string, row []*response.TransactionYearlyAmountFailedResponse) *pbtransaction.ApiResponseTransactionYearAmountFailed {
	return &pbtransaction.ApiResponseTransactionYearAmountFailed{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionYearlyAmountFailed(row),
	}
}

func (t *transactionProtoMapper) ToProtoResponseMonthMethod(status string, message string, row []*response.TransactionMonthlyMethodResponse) *pbtransaction.ApiResponseTransactionMonthPaymentMethod {
	return &pbtransaction.ApiResponseTransactionMonthPaymentMethod{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionMonthlyMethod(row),
	}
}

func (t *transactionProtoMapper) ToProtoResponseYearMethod(status string, message string, row []*response.TransactionYearlyMethodResponse) *pbtransaction.ApiResponseTransactionYearPaymentmethod {
	return &pbtransaction.ApiResponseTransactionYearPaymentmethod{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionYearlyMethod(row),
	}
}

func (t *transactionProtoMapper) mapResponseTransaction(transaction *response.TransactionResponse) *pbtransaction.TransactionResponse {
	return &pbtransaction.TransactionResponse{
		Id:            int32(transaction.ID),
		OrderId:       int32(transaction.OrderID),
		MerchantId:    int32(transaction.MerchantID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int32(transaction.Amount),
		ChangeAmount:  int32(transaction.ChangeAmount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}
}

func (t *transactionProtoMapper) mapResponsesTransaction(transactions []*response.TransactionResponse) []*pbtransaction.TransactionResponse {
	var mappedTransactions []*pbtransaction.TransactionResponse

	for _, transaction := range transactions {
		mappedTransactions = append(mappedTransactions, t.mapResponseTransaction(transaction))
	}

	return mappedTransactions
}

func (t *transactionProtoMapper) mapResponseTransactionDeleteAt(transaction *response.TransactionResponseDeleteAt) *pbtransaction.TransactionResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if transaction.DeletedAt != nil {
		deletedAt = wrapperspb.String(*transaction.DeletedAt)
	}

	return &pbtransaction.TransactionResponseDeleteAt{
		Id:            int32(transaction.ID),
		OrderId:       int32(transaction.OrderID),
		MerchantId:    int32(transaction.MerchantID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int32(transaction.Amount),
		ChangeAmount:  int32(transaction.ChangeAmount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
		DeletedAt:     deletedAt,
	}
}

func (t *transactionProtoMapper) mapResponsesTransactionDeleteAt(transactions []*response.TransactionResponseDeleteAt) []*pbtransaction.TransactionResponseDeleteAt {
	var mappedTransactions []*pbtransaction.TransactionResponseDeleteAt

	for _, transaction := range transactions {
		mappedTransactions = append(mappedTransactions, t.mapResponseTransactionDeleteAt(transaction))
	}

	return mappedTransactions
}

func (s *transactionProtoMapper) mapResponseTransactionMonthAmountSuccess(row *response.TransactionMonthlyAmountSuccessResponse) *pbtransaction.TransactionMonthlyAmountSuccess {
	return &pbtransaction.TransactionMonthlyAmountSuccess{
		Year:         row.Year,
		Month:        row.Month,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionMonthlyAmountSuccess(rows []*response.TransactionMonthlyAmountSuccessResponse) []*pbtransaction.TransactionMonthlyAmountSuccess {
	var transaction []*pbtransaction.TransactionMonthlyAmountSuccess

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionMonthAmountSuccess(row))
	}

	return transaction
}

func (s *transactionProtoMapper) mapResponseTransactionYearAmountSuccess(row *response.TransactionYearlyAmountSuccessResponse) *pbtransaction.TransactionYearlyAmountSuccess {
	return &pbtransaction.TransactionYearlyAmountSuccess{
		Year:         row.Year,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionYearlyAmountSuccess(rows []*response.TransactionYearlyAmountSuccessResponse) []*pbtransaction.TransactionYearlyAmountSuccess {
	var transaction []*pbtransaction.TransactionYearlyAmountSuccess

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionYearAmountSuccess(row))
	}

	return transaction
}

func (s *transactionProtoMapper) mapResponseTransactionMonthAmountFailed(row *response.TransactionMonthlyAmountFailedResponse) *pbtransaction.TransactionMonthlyAmountFailed {
	return &pbtransaction.TransactionMonthlyAmountFailed{
		Year:        row.Year,
		Month:       row.Month,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionMonthlyAmountFailed(rows []*response.TransactionMonthlyAmountFailedResponse) []*pbtransaction.TransactionMonthlyAmountFailed {
	var transaction []*pbtransaction.TransactionMonthlyAmountFailed

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionMonthAmountFailed(row))
	}

	return transaction
}

func (s *transactionProtoMapper) mapResponseTransactionYearAmountFailed(row *response.TransactionYearlyAmountFailedResponse) *pbtransaction.TransactionYearlyAmountFailed {
	return &pbtransaction.TransactionYearlyAmountFailed{
		Year:        row.Year,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionYearlyAmountFailed(rows []*response.TransactionYearlyAmountFailedResponse) []*pbtransaction.TransactionYearlyAmountFailed {
	var transaction []*pbtransaction.TransactionYearlyAmountFailed

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionYearAmountFailed(row))
	}

	return transaction
}

func (s *transactionProtoMapper) mapResponseTransactionMonthMethod(row *response.TransactionMonthlyMethodResponse) *pbtransaction.TransactionMonthlyMethod {
	return &pbtransaction.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionMonthlyMethod(rows []*response.TransactionMonthlyMethodResponse) []*pbtransaction.TransactionMonthlyMethod {
	var transaction []*pbtransaction.TransactionMonthlyMethod

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionMonthMethod(row))
	}

	return transaction
}

func (s *transactionProtoMapper) mapResponseTransactionYearMethod(row *response.TransactionYearlyMethodResponse) *pbtransaction.TransactionYearlyMethod {
	return &pbtransaction.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionYearlyMethod(rows []*response.TransactionYearlyMethodResponse) []*pbtransaction.TransactionYearlyMethod {
	var transaction []*pbtransaction.TransactionYearlyMethod

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionYearMethod(row))
	}

	return transaction
}
