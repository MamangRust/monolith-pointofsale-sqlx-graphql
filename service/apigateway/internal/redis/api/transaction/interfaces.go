package transaction_cache

import (
	"context"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
)

type TransactionStatsCache interface {
	GetCachedMonthAmountSuccessCached(ctx context.Context, req *model.FindMonthlyTransactionStatusInput) (*model.APIResponseTransactionMonthAmountSuccess, bool)
	SetCachedMonthAmountSuccessCached(ctx context.Context, req *model.FindMonthlyTransactionStatusInput, res *model.APIResponseTransactionMonthAmountSuccess)

	GetCachedYearAmountSuccessCached(ctx context.Context, year int) (*model.APIResponseTransactionYearAmountSuccess, bool)
	SetCachedYearAmountSuccessCached(ctx context.Context, year int, res *model.APIResponseTransactionYearAmountSuccess)

	GetCachedMonthAmountFailedCached(ctx context.Context, req *model.FindMonthlyTransactionStatusInput) (*model.APIResponseTransactionMonthAmountFailed, bool)
	SetCachedMonthAmountFailedCached(ctx context.Context, req *model.FindMonthlyTransactionStatusInput, res *model.APIResponseTransactionMonthAmountFailed)

	GetCachedYearAmountFailedCached(ctx context.Context, year int) (*model.APIResponseTransactionYearAmountFailed, bool)
	SetCachedYearAmountFailedCached(ctx context.Context, year int, res *model.APIResponseTransactionYearAmountFailed)

	GetCachedMonthMethodSuccessCached(ctx context.Context, req *model.MonthTransactionMethodInput) (*model.APIResponseTransactionMonthPaymentMethod, bool)
	SetCachedMonthMethodSuccessCached(ctx context.Context, req *model.MonthTransactionMethodInput, res *model.APIResponseTransactionMonthPaymentMethod)

	GetCachedYearMethodSuccessCached(ctx context.Context, year int) (*model.APIResponseTransactionYearPaymentMethod, bool)
	SetCachedYearMethodSuccessCached(ctx context.Context, year int, res *model.APIResponseTransactionYearPaymentMethod)

	GetCachedMonthMethodFailedCached(ctx context.Context, req *model.MonthTransactionMethodInput) (*model.APIResponseTransactionMonthPaymentMethod, bool)
	SetCachedMonthMethodFailedCached(ctx context.Context, req *model.MonthTransactionMethodInput, res *model.APIResponseTransactionMonthPaymentMethod)

	GetCachedYearMethodFailedCached(ctx context.Context, year int) (*model.APIResponseTransactionYearPaymentMethod, bool)
	SetCachedYearMethodFailedCached(ctx context.Context, year int, res *model.APIResponseTransactionYearPaymentMethod)
}

type TransactionStatsByMerchantCache interface {
	GetCachedMonthAmountSuccessByMerchantCached(ctx context.Context, req *model.FindMonthlyTransactionStatusByMerchantInput) (*model.APIResponseTransactionMonthAmountSuccess, bool)
	SetCachedMonthAmountSuccessByMerchantCached(ctx context.Context, req *model.FindMonthlyTransactionStatusByMerchantInput, res *model.APIResponseTransactionMonthAmountSuccess)

	GetCachedYearAmountSuccessByMerchantCached(ctx context.Context, req *model.FindYearlyTransactionStatusByMerchantInput) (*model.APIResponseTransactionYearAmountSuccess, bool)
	SetCachedYearAmountSuccessByMerchantCached(ctx context.Context, req *model.FindYearlyTransactionStatusByMerchantInput, res *model.APIResponseTransactionYearAmountSuccess)

	GetCachedMonthAmountFailedByMerchantCached(ctx context.Context, req *model.FindMonthlyTransactionStatusByMerchantInput) (*model.APIResponseTransactionMonthAmountFailed, bool)
	SetCachedMonthAmountFailedByMerchantCached(ctx context.Context, req *model.FindMonthlyTransactionStatusByMerchantInput, res *model.APIResponseTransactionMonthAmountFailed)

	GetCachedYearAmountFailedByMerchantCached(ctx context.Context, req *model.FindYearlyTransactionStatusByMerchantInput) (*model.APIResponseTransactionYearAmountFailed, bool)
	SetCachedYearAmountFailedByMerchantCached(ctx context.Context, req *model.FindYearlyTransactionStatusByMerchantInput, res *model.APIResponseTransactionYearAmountFailed)

	GetCachedMonthMethodSuccessByMerchantCached(ctx context.Context, req *model.MonthTransactionMethodByMerchantInput) (*model.APIResponseTransactionMonthPaymentMethod, bool)
	SetCachedMonthMethodSuccessByMerchantCached(ctx context.Context, req *model.MonthTransactionMethodByMerchantInput, res *model.APIResponseTransactionMonthPaymentMethod)

	GetCachedYearMethodSuccessByMerchantCached(ctx context.Context, req *model.YearTransactionMethodByMerchantInput) (*model.APIResponseTransactionYearPaymentMethod, bool)
	SetCachedYearMethodSuccessByMerchantCached(ctx context.Context, req *model.YearTransactionMethodByMerchantInput, res *model.APIResponseTransactionYearPaymentMethod)

	GetCachedMonthMethodFailedByMerchantCached(ctx context.Context, req *model.MonthTransactionMethodByMerchantInput) (*model.APIResponseTransactionMonthPaymentMethod, bool)
	SetCachedMonthMethodFailedByMerchantCached(ctx context.Context, req *model.MonthTransactionMethodByMerchantInput, res *model.APIResponseTransactionMonthPaymentMethod)

	GetCachedYearMethodFailedByMerchantCached(ctx context.Context, req *model.YearTransactionMethodByMerchantInput) (*model.APIResponseTransactionYearPaymentMethod, bool)
	SetCachedYearMethodFailedByMerchantCached(ctx context.Context, req *model.YearTransactionMethodByMerchantInput, res *model.APIResponseTransactionYearPaymentMethod)
}

type TransactionQueryCache interface {
	GetCachedTransactionsCache(ctx context.Context, req *model.FindAllTransactionInput) (*model.APIResponsePaginationTransaction, bool)
	SetCachedTransactionsCache(ctx context.Context, req *model.FindAllTransactionInput, res *model.APIResponsePaginationTransaction)

	GetCachedTransactionByMerchant(ctx context.Context, req *model.FindAllTransactionMerchantInput) (*model.APIResponsePaginationTransaction, bool)
	SetCachedTransactionByMerchant(ctx context.Context, req *model.FindAllTransactionMerchantInput, res *model.APIResponsePaginationTransaction)

	GetCachedTransactionActiveCache(ctx context.Context, req *model.FindAllTransactionInput) (*model.APIResponsePaginationTransactionDeleteAt, bool)
	SetCachedTransactionActiveCache(ctx context.Context, req *model.FindAllTransactionInput, res *model.APIResponsePaginationTransactionDeleteAt)

	GetCachedTransactionTrashedCache(ctx context.Context, req *model.FindAllTransactionInput) (*model.APIResponsePaginationTransactionDeleteAt, bool)
	SetCachedTransactionTrashedCache(ctx context.Context, req *model.FindAllTransactionInput, res *model.APIResponsePaginationTransactionDeleteAt)

	GetCachedTransactionCache(ctx context.Context, id int) (*model.APIResponseTransaction, bool)
	SetCachedTransactionCache(ctx context.Context, res *model.APIResponseTransaction)

	GetCachedTransactionByOrderId(ctx context.Context, orderID int) (*model.APIResponseTransaction, bool)
	SetCachedTransactionByOrderId(ctx context.Context, orderID int, res *model.APIResponseTransaction)
}

type TransactionCommandCache interface {
	DeleteTransactionCache(ctx context.Context, transactionID int)
}
