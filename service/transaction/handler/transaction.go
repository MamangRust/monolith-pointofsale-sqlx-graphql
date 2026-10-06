package handler

import (
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbtransaction "github.com/MamangRust/monolith-pointofsale-pb/transaction"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func mapPaginationMeta(meta *pbcommon.PaginationMeta) *pbcommon.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &pbcommon.PaginationMeta{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		TotalPages:   meta.TotalPages,
		TotalRecords: meta.TotalRecords,
	}
}

func mapResponseTransaction(transaction *db.Transaction) *pbtransaction.TransactionResponse {
	if transaction == nil {
		return nil
	}
	var createdAtStr string
	if transaction.CreatedAt.Valid {
		createdAtStr = transaction.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var updatedAtStr string
	if transaction.UpdatedAt.Valid {
		updatedAtStr = transaction.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var changeAmount int32
	if transaction.ChangeAmount != nil {
		changeAmount = *transaction.ChangeAmount
	}
	return &pbtransaction.TransactionResponse{
		Id:            transaction.TransactionID,
		OrderId:       transaction.OrderID,
		MerchantId:    transaction.MerchantID,
		PaymentMethod: transaction.PaymentMethod,
		Amount:        transaction.Amount,
		ChangeAmount:  changeAmount,
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     createdAtStr,
		UpdatedAt:     updatedAtStr,
	}
}

func mapResponsesTransaction(transactions []*db.GetTransactionsRow) []*pbtransaction.TransactionResponse {
	var mappedTransactions []*pbtransaction.TransactionResponse
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if t.CreatedAt.Valid {
			createdAtStr = t.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var updatedAtStr string
		if t.UpdatedAt.Valid {
			updatedAtStr = t.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pbtransaction.TransactionResponse{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
		})
	}
	return mappedTransactions
}

func mapResponsesTransactionByMerchant(transactions []*db.GetTransactionByMerchantRow) []*pbtransaction.TransactionResponse {
	var mappedTransactions []*pbtransaction.TransactionResponse
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if t.CreatedAt.Valid {
			createdAtStr = t.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var updatedAtStr string
		if t.UpdatedAt.Valid {
			updatedAtStr = t.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pbtransaction.TransactionResponse{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
		})
	}
	return mappedTransactions
}

func mapResponseTransactionDeleteAt(transaction *db.Transaction) *pbtransaction.TransactionResponseDeleteAt {
	if transaction == nil {
		return nil
	}
	var createdAtStr string
	if transaction.CreatedAt.Valid {
		createdAtStr = transaction.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var updatedAtStr string
	if transaction.UpdatedAt.Valid {
		updatedAtStr = transaction.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var deletedAt *wrapperspb.StringValue
	if transaction.DeletedAt.Valid {
		deletedAt = wrapperspb.String(transaction.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}
	var changeAmount int32
	if transaction.ChangeAmount != nil {
		changeAmount = *transaction.ChangeAmount
	}

	return &pbtransaction.TransactionResponseDeleteAt{
		Id:            transaction.TransactionID,
		OrderId:       transaction.OrderID,
		MerchantId:    transaction.MerchantID,
		PaymentMethod: transaction.PaymentMethod,
		Amount:        transaction.Amount,
		ChangeAmount:  changeAmount,
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     createdAtStr,
		UpdatedAt:     updatedAtStr,
		DeletedAt:     deletedAt,
	}
}

func mapResponsesTransactionActive(transactions []*db.GetTransactionsActiveRow) []*pbtransaction.TransactionResponseDeleteAt {
	var mappedTransactions []*pbtransaction.TransactionResponseDeleteAt
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if t.CreatedAt.Valid {
			createdAtStr = t.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var updatedAtStr string
		if t.UpdatedAt.Valid {
			updatedAtStr = t.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var deletedAt *wrapperspb.StringValue
		if t.DeletedAt.Valid {
			deletedAt = wrapperspb.String(t.DeletedAt.Time.Format("2006-01-02 15:04:05"))
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pbtransaction.TransactionResponseDeleteAt{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
			DeletedAt:     deletedAt,
		})
	}
	return mappedTransactions
}

func mapResponsesTransactionTrashed(transactions []*db.GetTransactionsTrashedRow) []*pbtransaction.TransactionResponseDeleteAt {
	var mappedTransactions []*pbtransaction.TransactionResponseDeleteAt
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if t.CreatedAt.Valid {
			createdAtStr = t.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var updatedAtStr string
		if t.UpdatedAt.Valid {
			updatedAtStr = t.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var deletedAt *wrapperspb.StringValue
		if t.DeletedAt.Valid {
			deletedAt = wrapperspb.String(t.DeletedAt.Time.Format("2006-01-02 15:04:05"))
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pbtransaction.TransactionResponseDeleteAt{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
			DeletedAt:     deletedAt,
		})
	}
	return mappedTransactions
}

func mapResponseTransactionMonthAmountSuccess(row *db.GetMonthlyAmountTransactionSuccessRow) *pbtransaction.TransactionMonthlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionMonthlyAmountSuccess{
		Year:         row.Year,
		Month:        row.Month,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountSuccess(rows []*db.GetMonthlyAmountTransactionSuccessRow) []*pbtransaction.TransactionMonthlyAmountSuccess {
	var transaction []*pbtransaction.TransactionMonthlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountSuccess(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountSuccess(row *db.GetYearlyAmountTransactionSuccessRow) *pbtransaction.TransactionYearlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionYearlyAmountSuccess{
		Year:         row.Year,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountSuccess(rows []*db.GetYearlyAmountTransactionSuccessRow) []*pbtransaction.TransactionYearlyAmountSuccess {
	var transaction []*pbtransaction.TransactionYearlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountSuccess(row))
	}
	return transaction
}

func mapResponseTransactionMonthAmountFailed(row *db.GetMonthlyAmountTransactionFailedRow) *pbtransaction.TransactionMonthlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionMonthlyAmountFailed{
		Year:        row.Year,
		Month:       row.Month,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountFailed(rows []*db.GetMonthlyAmountTransactionFailedRow) []*pbtransaction.TransactionMonthlyAmountFailed {
	var transaction []*pbtransaction.TransactionMonthlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountFailed(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountFailed(row *db.GetYearlyAmountTransactionFailedRow) *pbtransaction.TransactionYearlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionYearlyAmountFailed{
		Year:        row.Year,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountFailed(rows []*db.GetYearlyAmountTransactionFailedRow) []*pbtransaction.TransactionYearlyAmountFailed {
	var transaction []*pbtransaction.TransactionYearlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountFailed(row))
	}
	return transaction
}

func mapResponseTransactionMonthAmountSuccessByMerchant(row *db.GetMonthlyAmountTransactionSuccessByMerchantRow) *pbtransaction.TransactionMonthlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionMonthlyAmountSuccess{
		Year:         row.Year,
		Month:        row.Month,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountSuccessByMerchant(rows []*db.GetMonthlyAmountTransactionSuccessByMerchantRow) []*pbtransaction.TransactionMonthlyAmountSuccess {
	var transaction []*pbtransaction.TransactionMonthlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountSuccessByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountSuccessByMerchant(row *db.GetYearlyAmountTransactionSuccessByMerchantRow) *pbtransaction.TransactionYearlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionYearlyAmountSuccess{
		Year:         row.Year,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountSuccessByMerchant(rows []*db.GetYearlyAmountTransactionSuccessByMerchantRow) []*pbtransaction.TransactionYearlyAmountSuccess {
	var transaction []*pbtransaction.TransactionYearlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountSuccessByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionMonthAmountFailedByMerchant(row *db.GetMonthlyAmountTransactionFailedByMerchantRow) *pbtransaction.TransactionMonthlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionMonthlyAmountFailed{
		Year:        row.Year,
		Month:       row.Month,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountFailedByMerchant(rows []*db.GetMonthlyAmountTransactionFailedByMerchantRow) []*pbtransaction.TransactionMonthlyAmountFailed {
	var transaction []*pbtransaction.TransactionMonthlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountFailedByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountFailedByMerchant(row *db.GetYearlyAmountTransactionFailedByMerchantRow) *pbtransaction.TransactionYearlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionYearlyAmountFailed{
		Year:        row.Year,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountFailedByMerchant(rows []*db.GetYearlyAmountTransactionFailedByMerchantRow) []*pbtransaction.TransactionYearlyAmountFailed {
	var transaction []*pbtransaction.TransactionYearlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountFailedByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodSuccess(row *db.GetMonthlyTransactionMethodsSuccessRow) *pbtransaction.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodSuccess(rows []*db.GetMonthlyTransactionMethodsSuccessRow) []*pbtransaction.TransactionMonthlyMethod {
	var transaction []*pbtransaction.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodSuccess(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodSuccess(row *db.GetYearlyTransactionMethodsSuccessRow) *pbtransaction.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodSuccess(rows []*db.GetYearlyTransactionMethodsSuccessRow) []*pbtransaction.TransactionYearlyMethod {
	var transaction []*pbtransaction.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodSuccess(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodFailed(row *db.GetMonthlyTransactionMethodsFailedRow) *pbtransaction.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodFailed(rows []*db.GetMonthlyTransactionMethodsFailedRow) []*pbtransaction.TransactionMonthlyMethod {
	var transaction []*pbtransaction.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodFailed(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodFailed(row *db.GetYearlyTransactionMethodsFailedRow) *pbtransaction.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodFailed(rows []*db.GetYearlyTransactionMethodsFailedRow) []*pbtransaction.TransactionYearlyMethod {
	var transaction []*pbtransaction.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodFailed(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodByMerchantSuccess(row *db.GetMonthlyTransactionMethodsByMerchantSuccessRow) *pbtransaction.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodByMerchantSuccess(rows []*db.GetMonthlyTransactionMethodsByMerchantSuccessRow) []*pbtransaction.TransactionMonthlyMethod {
	var transaction []*pbtransaction.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodByMerchantSuccess(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodByMerchantSuccess(row *db.GetYearlyTransactionMethodsByMerchantSuccessRow) *pbtransaction.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodByMerchantSuccess(rows []*db.GetYearlyTransactionMethodsByMerchantSuccessRow) []*pbtransaction.TransactionYearlyMethod {
	var transaction []*pbtransaction.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodByMerchantSuccess(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodByMerchantFailed(row *db.GetMonthlyTransactionMethodsByMerchantFailedRow) *pbtransaction.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodByMerchantFailed(rows []*db.GetMonthlyTransactionMethodsByMerchantFailedRow) []*pbtransaction.TransactionMonthlyMethod {
	var transaction []*pbtransaction.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodByMerchantFailed(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodByMerchantFailed(row *db.GetYearlyTransactionMethodsByMerchantFailedRow) *pbtransaction.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransaction.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodByMerchantFailed(rows []*db.GetYearlyTransactionMethodsByMerchantFailedRow) []*pbtransaction.TransactionYearlyMethod {
	var transaction []*pbtransaction.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodByMerchantFailed(row))
	}
	return transaction
}
