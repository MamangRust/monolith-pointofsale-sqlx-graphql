package response_api

import (
	pbtransaction "github.com/MamangRust/monolith-pointofsale-pb/transaction"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type transactionResponseMapper struct {
}

func NewTransactionResponseMapper() *transactionResponseMapper {
	return &transactionResponseMapper{}
}

func (t *transactionResponseMapper) ToResponseTransaction(transaction *pbtransaction.TransactionResponse) *response.TransactionResponse {
	return &response.TransactionResponse{
		ID:            int(transaction.Id),
		OrderID:       int(transaction.OrderId),
		MerchantID:    int(transaction.MerchantId),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int(transaction.Amount),
		ChangeAmount:  int(transaction.ChangeAmount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}
}

func (t *transactionResponseMapper) ToResponsesTransaction(transactions []*pbtransaction.TransactionResponse) []*response.TransactionResponse {
	var mappedTransactions []*response.TransactionResponse

	for _, transaction := range transactions {
		mappedTransactions = append(mappedTransactions, t.ToResponseTransaction(transaction))
	}

	return mappedTransactions
}

func (t *transactionResponseMapper) ToResponseTransactionDeleteAt(transaction *pbtransaction.TransactionResponseDeleteAt) *response.TransactionResponseDeleteAt {
	var deletedAt string
	if transaction.DeletedAt != nil {
		deletedAt = transaction.DeletedAt.Value
	}

	return &response.TransactionResponseDeleteAt{
		ID:            int(transaction.Id),
		OrderID:       int(transaction.OrderId),
		MerchantID:    int(transaction.MerchantId),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int(transaction.Amount),
		ChangeAmount:  int(transaction.ChangeAmount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
		DeletedAt:     &deletedAt,
	}
}

func (t *transactionResponseMapper) ToResponsesTransactionDeleteAt(transactions []*pbtransaction.TransactionResponseDeleteAt) []*response.TransactionResponseDeleteAt {
	var mappedTransactions []*response.TransactionResponseDeleteAt

	for _, transaction := range transactions {
		mappedTransactions = append(mappedTransactions, t.ToResponseTransactionDeleteAt(transaction))
	}

	return mappedTransactions
}

func (s *transactionResponseMapper) ToTransactionMonthAmountSuccess(row *pbtransaction.TransactionMonthlyAmountSuccess) *response.TransactionMonthlyAmountSuccessResponse {
	return &response.TransactionMonthlyAmountSuccessResponse{
		Year:         row.Year,
		Month:        row.Month,
		TotalSuccess: int(row.TotalSuccess),
		TotalAmount:  int(row.TotalAmount),
	}
}

func (s *transactionResponseMapper) ToTransactionMonthlyAmountSuccess(rows []*pbtransaction.TransactionMonthlyAmountSuccess) []*response.TransactionMonthlyAmountSuccessResponse {
	var transaction []*response.TransactionMonthlyAmountSuccessResponse

	for _, row := range rows {
		transaction = append(transaction, s.ToTransactionMonthAmountSuccess(row))
	}

	return transaction
}

func (s *transactionResponseMapper) ToTransactionYearAmountSuccess(row *pbtransaction.TransactionYearlyAmountSuccess) *response.TransactionYearlyAmountSuccessResponse {
	return &response.TransactionYearlyAmountSuccessResponse{
		Year:         row.Year,
		TotalSuccess: int(row.TotalSuccess),
		TotalAmount:  int(row.TotalAmount),
	}
}

func (s *transactionResponseMapper) ToTransactionYearlyAmountSuccess(rows []*pbtransaction.TransactionYearlyAmountSuccess) []*response.TransactionYearlyAmountSuccessResponse {
	var transaction []*response.TransactionYearlyAmountSuccessResponse

	for _, row := range rows {
		transaction = append(transaction, s.ToTransactionYearAmountSuccess(row))
	}

	return transaction
}

func (s *transactionResponseMapper) ToTransactionMonthAmountFailed(row *pbtransaction.TransactionMonthlyAmountFailed) *response.TransactionMonthlyAmountFailedResponse {
	return &response.TransactionMonthlyAmountFailedResponse{
		Year:        row.Year,
		Month:       row.Month,
		TotalFailed: int(row.TotalFailed),
		TotalAmount: int(row.TotalAmount),
	}
}

func (s *transactionResponseMapper) ToTransactionMonthlyAmountFailed(rows []*pbtransaction.TransactionMonthlyAmountFailed) []*response.TransactionMonthlyAmountFailedResponse {
	var transaction []*response.TransactionMonthlyAmountFailedResponse

	for _, row := range rows {
		transaction = append(transaction, s.ToTransactionMonthAmountFailed(row))
	}

	return transaction
}

func (s *transactionResponseMapper) ToTransactionYearAmountFailed(row *pbtransaction.TransactionYearlyAmountFailed) *response.TransactionYearlyAmountFailedResponse {
	return &response.TransactionYearlyAmountFailedResponse{
		Year:        row.Year,
		TotalFailed: int(row.TotalFailed),
		TotalAmount: int(row.TotalAmount),
	}
}

func (s *transactionResponseMapper) ToTransactionYearlyAmountFailed(rows []*pbtransaction.TransactionYearlyAmountFailed) []*response.TransactionYearlyAmountFailedResponse {
	var transaction []*response.TransactionYearlyAmountFailedResponse

	for _, row := range rows {
		transaction = append(transaction, s.ToTransactionYearAmountFailed(row))
	}

	return transaction
}

func (s *transactionResponseMapper) ToTransactionMonthMethod(row *pbtransaction.TransactionMonthlyMethod) *response.TransactionMonthlyMethodResponse {
	return &response.TransactionMonthlyMethodResponse{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int(row.TotalTransactions),
		TotalAmount:       int(row.TotalAmount),
	}
}

func (s *transactionResponseMapper) ToTransactionMonthlyMethod(rows []*pbtransaction.TransactionMonthlyMethod) []*response.TransactionMonthlyMethodResponse {
	var transaction []*response.TransactionMonthlyMethodResponse

	for _, row := range rows {
		transaction = append(transaction, s.ToTransactionMonthMethod(row))
	}

	return transaction
}

func (s *transactionResponseMapper) ToTransactionYearMethod(row *pbtransaction.TransactionYearlyMethod) *response.TransactionYearlyMethodResponse {
	return &response.TransactionYearlyMethodResponse{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int(row.TotalTransactions),
		TotalAmount:       int(row.TotalAmount),
	}
}

func (s *transactionResponseMapper) ToTransactionYearlyMethod(rows []*pbtransaction.TransactionYearlyMethod) []*response.TransactionYearlyMethodResponse {
	var transaction []*response.TransactionYearlyMethodResponse

	for _, row := range rows {
		transaction = append(transaction, s.ToTransactionYearMethod(row))
	}

	return transaction
}

func (t *transactionResponseMapper) ToApiResponseTransaction(pbResponse *pbtransaction.ApiResponseTransaction) *response.ApiResponseTransaction {
	return &response.ApiResponseTransaction{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToResponseTransaction(pbResponse.Data),
	}
}

func (t *transactionResponseMapper) ToApiResponseTransactionDeleteAt(pbResponse *pbtransaction.ApiResponseTransactionDeleteAt) *response.ApiResponseTransactionDeleteAt {
	return &response.ApiResponseTransactionDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToResponseTransactionDeleteAt(pbResponse.Data),
	}
}

func (t *transactionResponseMapper) ToApiResponsesTransaction(pbResponse *pbtransaction.ApiResponsesTransaction) *response.ApiResponsesTransaction {
	return &response.ApiResponsesTransaction{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToResponsesTransaction(pbResponse.Data),
	}
}

func (t *transactionResponseMapper) ToApiResponseTransactionDelete(pbResponse *pbtransaction.ApiResponseTransactionDelete) *response.ApiResponseTransactionDelete {
	return &response.ApiResponseTransactionDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (t *transactionResponseMapper) ToApiResponseTransactionAll(pbResponse *pbtransaction.ApiResponseTransactionAll) *response.ApiResponseTransactionAll {
	return &response.ApiResponseTransactionAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (t *transactionResponseMapper) ToApiResponsePaginationTransactionDeleteAt(pbResponse *pbtransaction.ApiResponsePaginationTransactionDeleteAt) *response.ApiResponsePaginationTransactionDeleteAt {
	return &response.ApiResponsePaginationTransactionDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       t.ToResponsesTransactionDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (t *transactionResponseMapper) ToApiResponsePaginationTransaction(pbResponse *pbtransaction.ApiResponsePaginationTransaction) *response.ApiResponsePaginationTransaction {
	return &response.ApiResponsePaginationTransaction{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       t.ToResponsesTransaction(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (t *transactionResponseMapper) ToApiResponseTransactionMonthAmountSuccess(pbResponse *pbtransaction.ApiResponseTransactionMonthAmountSuccess) *response.ApiResponsesTransactionMonthSuccess {
	return &response.ApiResponsesTransactionMonthSuccess{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToTransactionMonthlyAmountSuccess(pbResponse.Data),
	}
}

func (t *transactionResponseMapper) ToApiResponseTransactionMonthAmountFailed(pbResponse *pbtransaction.ApiResponseTransactionMonthAmountFailed) *response.ApiResponsesTransactionMonthFailed {
	return &response.ApiResponsesTransactionMonthFailed{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToTransactionMonthlyAmountFailed(pbResponse.Data),
	}
}

func (t *transactionResponseMapper) ToApiResponseTransactionYearAmountSuccess(pbResponse *pbtransaction.ApiResponseTransactionYearAmountSuccess) *response.ApiResponsesTransactionYearSuccess {
	return &response.ApiResponsesTransactionYearSuccess{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToTransactionYearlyAmountSuccess(pbResponse.Data),
	}
}

func (t *transactionResponseMapper) ToApiResponseTransactionYearAmountFailed(pbResponse *pbtransaction.ApiResponseTransactionYearAmountFailed) *response.ApiResponsesTransactionYearFailed {
	return &response.ApiResponsesTransactionYearFailed{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToTransactionYearlyAmountFailed(pbResponse.Data),
	}
}

func (t *transactionResponseMapper) ToApiResponseTransactionMonthMethod(pbResponse *pbtransaction.ApiResponseTransactionMonthPaymentMethod) *response.ApiResponsesTransactionMonthMethod {
	return &response.ApiResponsesTransactionMonthMethod{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToTransactionMonthlyMethod(pbResponse.Data),
	}
}

func (t *transactionResponseMapper) ToApiResponseTransactionYearMethod(pbResponse *pbtransaction.ApiResponseTransactionYearPaymentmethod) *response.ApiResponsesTransactionYearMethod {
	return &response.ApiResponsesTransactionYearMethod{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToTransactionYearlyMethod(pbResponse.Data),
	}
}
