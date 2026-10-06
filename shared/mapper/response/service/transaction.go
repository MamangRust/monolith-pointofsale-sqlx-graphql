package response_service

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/record"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type transactionResponseMapper struct {
}

func NewTransactionResponseMapper() *transactionResponseMapper {
	return &transactionResponseMapper{}
}

func (s *transactionResponseMapper) ToTransactionResponse(transaction *record.TransactionRecord) *response.TransactionResponse {
	return &response.TransactionResponse{
		ID:            transaction.ID,
		OrderID:       transaction.OrderID,
		MerchantID:    transaction.MerchantID,
		PaymentMethod: transaction.PaymentMethod,
		Amount:        transaction.Amount,
		ChangeAmount:  transaction.ChangeAmount,
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}
}

func (s *transactionResponseMapper) ToTransactionsResponse(transactions []*record.TransactionRecord) []*response.TransactionResponse {
	var responses []*response.TransactionResponse

	for _, transaction := range transactions {
		responses = append(responses, s.ToTransactionResponse(transaction))
	}

	return responses
}

func (s *transactionResponseMapper) ToTransactionResponseDeleteAt(transaction *record.TransactionRecord) *response.TransactionResponseDeleteAt {
	return &response.TransactionResponseDeleteAt{
		ID:            transaction.ID,
		OrderID:       transaction.OrderID,
		MerchantID:    transaction.MerchantID,
		PaymentMethod: transaction.PaymentMethod,
		Amount:        transaction.Amount,
		ChangeAmount:  transaction.ChangeAmount,
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
		DeletedAt:     transaction.DeletedAt,
	}
}

func (s *transactionResponseMapper) ToTransactionsResponseDeleteAt(transactions []*record.TransactionRecord) []*response.TransactionResponseDeleteAt {
	var responses []*response.TransactionResponseDeleteAt

	for _, transaction := range transactions {
		responses = append(responses, s.ToTransactionResponseDeleteAt(transaction))
	}

	return responses
}

func (s *transactionResponseMapper) ToTransactionMonthAmountSuccess(row *record.TransactionMonthlyAmountSuccessRecord) *response.TransactionMonthlyAmountSuccessResponse {
	return &response.TransactionMonthlyAmountSuccessResponse{
		Year:         row.Year,
		Month:        row.Month,
		TotalSuccess: int(row.TotalSuccess),
		TotalAmount:  int(row.TotalAmount),
	}
}

func (s *transactionResponseMapper) ToTransactionMonthlyAmountSuccess(rows []*record.TransactionMonthlyAmountSuccessRecord) []*response.TransactionMonthlyAmountSuccessResponse {
	var transaction []*response.TransactionMonthlyAmountSuccessResponse

	for _, row := range rows {
		transaction = append(transaction, s.ToTransactionMonthAmountSuccess(row))
	}

	return transaction
}

func (s *transactionResponseMapper) ToTransactionYearAmountSuccess(row *record.TransactionYearlyAmountSuccessRecord) *response.TransactionYearlyAmountSuccessResponse {
	return &response.TransactionYearlyAmountSuccessResponse{
		Year:         row.Year,
		TotalSuccess: int(row.TotalSuccess),
		TotalAmount:  int(row.TotalAmount),
	}
}

func (s *transactionResponseMapper) ToTransactionYearlyAmountSuccess(rows []*record.TransactionYearlyAmountSuccessRecord) []*response.TransactionYearlyAmountSuccessResponse {
	var transaction []*response.TransactionYearlyAmountSuccessResponse

	for _, row := range rows {
		transaction = append(transaction, s.ToTransactionYearAmountSuccess(row))
	}

	return transaction
}

func (s *transactionResponseMapper) ToTransactionMonthAmountFailed(row *record.TransactionMonthlyAmountFailedRecord) *response.TransactionMonthlyAmountFailedResponse {
	return &response.TransactionMonthlyAmountFailedResponse{
		Year:        row.Year,
		Month:       row.Month,
		TotalFailed: int(row.TotalFailed),
		TotalAmount: int(row.TotalAmount),
	}
}

func (s *transactionResponseMapper) ToTransactionMonthlyAmountFailed(rows []*record.TransactionMonthlyAmountFailedRecord) []*response.TransactionMonthlyAmountFailedResponse {
	var transaction []*response.TransactionMonthlyAmountFailedResponse

	for _, row := range rows {
		transaction = append(transaction, s.ToTransactionMonthAmountFailed(row))
	}

	return transaction
}

func (s *transactionResponseMapper) ToTransactionYearAmountFailed(row *record.TransactionYearlyAmountFailedRecord) *response.TransactionYearlyAmountFailedResponse {
	return &response.TransactionYearlyAmountFailedResponse{
		Year:        row.Year,
		TotalFailed: int(row.TotalFailed),
		TotalAmount: int(row.TotalAmount),
	}
}

func (s *transactionResponseMapper) ToTransactionYearlyAmountFailed(rows []*record.TransactionYearlyAmountFailedRecord) []*response.TransactionYearlyAmountFailedResponse {
	var transaction []*response.TransactionYearlyAmountFailedResponse

	for _, row := range rows {
		transaction = append(transaction, s.ToTransactionYearAmountFailed(row))
	}

	return transaction
}

func (s *transactionResponseMapper) ToTransactionMonthMethod(row *record.TransactionMonthlyMethodRecord) *response.TransactionMonthlyMethodResponse {
	return &response.TransactionMonthlyMethodResponse{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int(row.TotalTransactions),
		TotalAmount:       int(row.TotalAmount),
	}
}

func (s *transactionResponseMapper) ToTransactionMonthlyMethod(rows []*record.TransactionMonthlyMethodRecord) []*response.TransactionMonthlyMethodResponse {
	var transaction []*response.TransactionMonthlyMethodResponse

	for _, row := range rows {
		transaction = append(transaction, s.ToTransactionMonthMethod(row))
	}

	return transaction
}

func (s *transactionResponseMapper) ToTransactionYearMethod(row *record.TransactionYearlyMethodRecord) *response.TransactionYearlyMethodResponse {
	return &response.TransactionYearlyMethodResponse{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int(row.TotalTransactions),
		TotalAmount:       int(row.TotalAmount),
	}
}

func (s *transactionResponseMapper) ToTransactionYearlyMethod(rows []*record.TransactionYearlyMethodRecord) []*response.TransactionYearlyMethodResponse {
	var transaction []*response.TransactionYearlyMethodResponse

	for _, row := range rows {
		transaction = append(transaction, s.ToTransactionYearMethod(row))
	}

	return transaction
}
