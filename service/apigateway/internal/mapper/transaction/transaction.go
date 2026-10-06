package transactiongraphqlmapper

import (
	graphqlmapper "github.com/MamangRust/monolith-pointofsale-apigateway/internal/mapper"
	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
	pb "github.com/MamangRust/monolith-pointofsale-pb/transaction"
)

type transactionGraphqlMapper struct {
}

func NewTransactionGraphqlMapper() *transactionGraphqlMapper {
	return &transactionGraphqlMapper{}
}

func (t *transactionGraphqlMapper) ToGraphqlResponseTransaction(res *pb.ApiResponseTransaction) *model.APIResponseTransaction {
	return &model.APIResponseTransaction{
		Status:  res.Status,
		Message: res.Message,
		Data:    t.mapResponseTransaction(res.Data),
	}
}

func (t *transactionGraphqlMapper) ToGraphqlResponsesTransaction(res *pb.ApiResponsesTransaction) *model.APIResponsesTransaction {
	return &model.APIResponsesTransaction{
		Status:  res.Status,
		Message: res.Message,
		Data:    t.mapResponsesTransaction(res.Data),
	}
}

func (t *transactionGraphqlMapper) ToGraphqlResponseTransactionDeleteAt(res *pb.ApiResponseTransactionDeleteAt) *model.APIResponseTransactionDeleteAt {
	return &model.APIResponseTransactionDeleteAt{
		Status:  res.Status,
		Message: res.Message,
		Data:    t.mapResponseTransactionDeleteAt(res.Data),
	}
}

func (t *transactionGraphqlMapper) ToGraphqlResponseTransactionDelete(res *pb.ApiResponseTransactionDelete) *model.APIResponseTransactionDelete {
	return &model.APIResponseTransactionDelete{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (t *transactionGraphqlMapper) ToGraphqlResponseTransactionAll(res *pb.ApiResponseTransactionAll) *model.APIResponseTransactionAll {
	return &model.APIResponseTransactionAll{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (t *transactionGraphqlMapper) ToGraphqlResponsePaginationTransaction(res *pb.ApiResponsePaginationTransaction) *model.APIResponsePaginationTransaction {
	return &model.APIResponsePaginationTransaction{
		Status:     res.Status,
		Message:    res.Message,
		Data:       t.mapResponsesTransaction(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (t *transactionGraphqlMapper) ToGraphqlResponsePaginationTransactionDeleteAt(res *pb.ApiResponsePaginationTransactionDeleteAt) *model.APIResponsePaginationTransactionDeleteAt {
	return &model.APIResponsePaginationTransactionDeleteAt{
		Status:     res.Status,
		Message:    res.Message,
		Data:       t.mapResponsesTransactionDeleteAt(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (t *transactionGraphqlMapper) mapResponseTransaction(transaction *pb.TransactionResponse) *model.TransactionResponse {
	if transaction == nil {
		return nil
	}
	return &model.TransactionResponse{
		ID:            int32(transaction.Id),
		OrderID:       int32(transaction.OrderId),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int32(transaction.Amount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}
}

func (t *transactionGraphqlMapper) mapResponsesTransaction(transactions []*pb.TransactionResponse) []*model.TransactionResponse {
	var responses []*model.TransactionResponse
	for _, transaction := range transactions {
		responses = append(responses, t.mapResponseTransaction(transaction))
	}
	return responses
}

func (t *transactionGraphqlMapper) mapResponseTransactionDeleteAt(transaction *pb.TransactionResponseDeleteAt) *model.TransactionResponseDeleteAt {
	if transaction == nil {
		return nil
	}
	var deletedAt string
	if transaction.DeletedAt != nil {
		deletedAt = transaction.DeletedAt.Value
	}

	return &model.TransactionResponseDeleteAt{
		ID:            int32(transaction.Id),
		OrderID:       int32(transaction.OrderId),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int32(transaction.Amount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
		DeletedAt:     &deletedAt,
	}
}

func (t *transactionGraphqlMapper) mapResponsesTransactionDeleteAt(transactions []*pb.TransactionResponseDeleteAt) []*model.TransactionResponseDeleteAt {
	var responses []*model.TransactionResponseDeleteAt
	for _, transaction := range transactions {
		responses = append(responses, t.mapResponseTransactionDeleteAt(transaction))
	}
	return responses
}

func (t *transactionGraphqlMapper) ToGraphqlResponseMonthAmountSuccess(res *pb.ApiResponseTransactionMonthAmountSuccess) *model.APIResponseTransactionMonthAmountSuccess {
	var responses []*model.TransactionMonthlyAmountSuccess
	for _, item := range res.Data {
		if item == nil {
			continue
		}
		responses = append(responses, &model.TransactionMonthlyAmountSuccess{
			Year:         item.Year,
			Month:        item.Month,
			TotalSuccess: item.TotalSuccess,
			TotalAmount:  item.TotalAmount,
		})
	}
	return &model.APIResponseTransactionMonthAmountSuccess{
		Status:  res.Status,
		Message: res.Message,
		Data:    responses,
	}
}

func (t *transactionGraphqlMapper) ToGraphqlResponseYearAmountSuccess(res *pb.ApiResponseTransactionYearAmountSuccess) *model.APIResponseTransactionYearAmountSuccess {
	var responses []*model.TransactionYearlyAmountSuccess
	for _, item := range res.Data {
		if item == nil {
			continue
		}
		responses = append(responses, &model.TransactionYearlyAmountSuccess{
			Year:         item.Year,
			TotalSuccess: item.TotalSuccess,
			TotalAmount:  item.TotalAmount,
		})
	}
	return &model.APIResponseTransactionYearAmountSuccess{
		Status:  res.Status,
		Message: res.Message,
		Data:    responses,
	}
}

func (t *transactionGraphqlMapper) ToGraphqlResponseMonthAmountFailed(res *pb.ApiResponseTransactionMonthAmountFailed) *model.APIResponseTransactionMonthAmountFailed {
	var responses []*model.TransactionMonthlyAmountFailed
	for _, item := range res.Data {
		if item == nil {
			continue
		}
		responses = append(responses, &model.TransactionMonthlyAmountFailed{
			Year:        item.Year,
			Month:       item.Month,
			TotalFailed: item.TotalFailed,
			TotalAmount: item.TotalAmount,
		})
	}
	return &model.APIResponseTransactionMonthAmountFailed{
		Status:  res.Status,
		Message: res.Message,
		Data:    responses,
	}
}

func (t *transactionGraphqlMapper) ToGraphqlResponseYearAmountFailed(res *pb.ApiResponseTransactionYearAmountFailed) *model.APIResponseTransactionYearAmountFailed {
	var responses []*model.TransactionYearlyAmountFailed
	for _, item := range res.Data {
		if item == nil {
			continue
		}
		responses = append(responses, &model.TransactionYearlyAmountFailed{
			Year:        item.Year,
			TotalFailed: item.TotalFailed,
			TotalAmount: item.TotalAmount,
		})
	}
	return &model.APIResponseTransactionYearAmountFailed{
		Status:  res.Status,
		Message: res.Message,
		Data:    responses,
	}
}

func (t *transactionGraphqlMapper) ToGraphqlResponseMonthMethod(res *pb.ApiResponseTransactionMonthPaymentMethod) *model.APIResponseTransactionMonthPaymentMethod {
	var responses []*model.TransactionMonthlyMethod
	for _, item := range res.Data {
		if item == nil {
			continue
		}
		responses = append(responses, &model.TransactionMonthlyMethod{
			Month:             item.Month,
			PaymentMethod:     item.PaymentMethod,
			TotalTransactions: item.TotalTransactions,
			TotalAmount:       item.TotalAmount,
		})
	}
	return &model.APIResponseTransactionMonthPaymentMethod{
		Status:  res.Status,
		Message: res.Message,
		Data:    responses,
	}
}

func (t *transactionGraphqlMapper) ToGraphqlResponseYearMethod(res *pb.ApiResponseTransactionYearPaymentmethod) *model.APIResponseTransactionYearPaymentMethod {
	var responses []*model.TransactionYearlyMethod
	for _, item := range res.Data {
		if item == nil {
			continue
		}
		responses = append(responses, &model.TransactionYearlyMethod{
			Year:              item.Year,
			PaymentMethod:     item.PaymentMethod,
			TotalTransactions: item.TotalTransactions,
			TotalAmount:       item.TotalAmount,
		})
	}
	return &model.APIResponseTransactionYearPaymentMethod{
		Status:  res.Status,
		Message: res.Message,
		Data:    responses,
	}
}
