package mencache

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

type TransactionStatsCache interface {
	GetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionSuccessRow, bool)
	SetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction, res []*db.GetMonthlyAmountTransactionSuccessRow)

	GetCachedYearAmountSuccessCached(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionSuccessRow, bool)
	SetCachedYearAmountSuccessCached(ctx context.Context, year int, res []*db.GetYearlyAmountTransactionSuccessRow)

	GetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionFailedRow, bool)
	SetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction, res []*db.GetMonthlyAmountTransactionFailedRow)

	GetCachedYearAmountFailedCached(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionFailedRow, bool)
	SetCachedYearAmountFailedCached(ctx context.Context, year int, res []*db.GetYearlyAmountTransactionFailedRow)

	GetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsSuccessRow, bool)
	SetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction, res []*db.GetMonthlyTransactionMethodsSuccessRow)

	GetCachedYearMethodSuccessCached(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsSuccessRow, bool)
	SetCachedYearMethodSuccessCached(ctx context.Context, year int, res []*db.GetYearlyTransactionMethodsSuccessRow)

	GetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsFailedRow, bool)
	SetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction, res []*db.GetMonthlyTransactionMethodsFailedRow)

	GetCachedYearMethodFailedCached(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsFailedRow, bool)
	SetCachedYearMethodFailedCached(ctx context.Context, year int, res []*db.GetYearlyTransactionMethodsFailedRow)
}

type TransactionStatsByMerchantCache interface {
	GetCachedMonthAmountSuccessByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionSuccessByMerchantRow, bool)
	SetCachedMonthAmountSuccessByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*db.GetMonthlyAmountTransactionSuccessByMerchantRow)

	GetCachedYearAmountSuccessByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionSuccessByMerchantRow, bool)
	SetCachedYearAmountSuccessByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*db.GetYearlyAmountTransactionSuccessByMerchantRow)

	GetCachedMonthAmountFailedByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionFailedByMerchantRow, bool)
	SetCachedMonthAmountFailedByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*db.GetMonthlyAmountTransactionFailedByMerchantRow)

	GetCachedYearAmountFailedByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionFailedByMerchantRow, bool)
	SetCachedYearAmountFailedByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*db.GetYearlyAmountTransactionFailedByMerchantRow)

	GetCachedMonthMethodSuccessByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow, bool)
	SetCachedMonthMethodSuccessByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*db.GetMonthlyTransactionMethodsByMerchantSuccessRow)

	GetCachedYearMethodSuccessByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantSuccessRow, bool)
	SetCachedYearMethodSuccessByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*db.GetYearlyTransactionMethodsByMerchantSuccessRow)

	GetCachedMonthMethodFailedByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantFailedRow, bool)
	SetCachedMonthMethodFailedByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*db.GetMonthlyTransactionMethodsByMerchantFailedRow)

	GetCachedYearMethodFailedByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantFailedRow, bool)
	SetCachedYearMethodFailedByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*db.GetYearlyTransactionMethodsByMerchantFailedRow)
}

type TransactionQueryCache interface {
	GetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsRow, *int, bool)
	SetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction, data []*db.GetTransactionsRow, total *int)

	GetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*db.GetTransactionByMerchantRow, *int, bool)
	SetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant, data []*db.GetTransactionByMerchantRow, total *int)

	GetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsActiveRow, *int, bool)
	SetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction, data []*db.GetTransactionsActiveRow, total *int)

	GetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsTrashedRow, *int, bool)
	SetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction, data []*db.GetTransactionsTrashedRow, total *int)

	GetCachedTransactionCache(ctx context.Context, id int) (*db.Transaction, bool)
	SetCachedTransactionCache(ctx context.Context, data *db.Transaction)

	GetCachedTransactionByOrderId(ctx context.Context, orderID int) (*db.Transaction, bool)
	SetCachedTransactionByOrderId(ctx context.Context, orderID int, data *db.Transaction)
}

type TransactionCommandCache interface {
	DeleteTransactionCache(ctx context.Context, transactionID int)
	DeleteTransactionAllCache(ctx context.Context)
}
