package response_service

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/record"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type UserResponseMapper interface {
	ToUserResponse(user *record.UserRecord) *response.UserResponse
	ToUsersResponse(users []*record.UserRecord) []*response.UserResponse

	ToUserResponseDeleteAt(user *record.UserRecord) *response.UserResponseDeleteAt
	ToUsersResponseDeleteAt(users []*record.UserRecord) []*response.UserResponseDeleteAt
}

type RoleResponseMapper interface {
	ToRoleResponse(role *record.RoleRecord) *response.RoleResponse
	ToRolesResponse(roles []*record.RoleRecord) []*response.RoleResponse

	ToRoleResponseDeleteAt(role *record.RoleRecord) *response.RoleResponseDeleteAt
	ToRolesResponseDeleteAt(roles []*record.RoleRecord) []*response.RoleResponseDeleteAt
}

type RefreshTokenResponseMapper interface {
	ToRefreshTokenResponse(refresh *record.RefreshTokenRecord) *response.RefreshTokenResponse
	ToRefreshTokenResponses(refreshs []*record.RefreshTokenRecord) []*response.RefreshTokenResponse
}

type CategoryResponseMapper interface {
	ToCategoryMonthlyTotalPrice(c *record.CategoriesMonthlyTotalPriceRecord) *response.CategoriesMonthlyTotalPriceResponse
	ToCategoryMonthlyTotalPrices(c []*record.CategoriesMonthlyTotalPriceRecord) []*response.CategoriesMonthlyTotalPriceResponse
	ToCategoryYearlyTotalPrice(c *record.CategoriesYearlyTotalPriceRecord) *response.CategoriesYearlyTotalPriceResponse
	ToCategoryYearlyTotalPrices(c []*record.CategoriesYearlyTotalPriceRecord) []*response.CategoriesYearlyTotalPriceResponse

	ToCategoryMonthlyPrice(category *record.CategoriesMonthPriceRecord) *response.CategoryMonthPriceResponse
	ToCategoryMonthlyPrices(c []*record.CategoriesMonthPriceRecord) []*response.CategoryMonthPriceResponse
	ToCategoryYearlyPrice(category *record.CategoriesYearPriceRecord) *response.CategoryYearPriceResponse
	ToCategoryYearlyPrices(c []*record.CategoriesYearPriceRecord) []*response.CategoryYearPriceResponse

	ToCategoryResponse(category *record.CategoriesRecord) *response.CategoryResponse
	ToCategorysResponse(categories []*record.CategoriesRecord) []*response.CategoryResponse
	ToCategoryResponseDeleteAt(category *record.CategoriesRecord) *response.CategoryResponseDeleteAt
	ToCategoryResponsesDeleteAt(categories []*record.CategoriesRecord) []*response.CategoryResponseDeleteAt
}

type CashierResponseMapper interface {
	ToCashierMonthlyTotalSale(c *record.CashierRecordMonthTotalSales) *response.CashierResponseMonthTotalSales
	ToCashierMonthlyTotalSales(c []*record.CashierRecordMonthTotalSales) []*response.CashierResponseMonthTotalSales
	ToCashierYearlyTotalSale(c *record.CashierRecordYearTotalSales) *response.CashierResponseYearTotalSales
	ToCashierYearlyTotalSales(c []*record.CashierRecordYearTotalSales) []*response.CashierResponseYearTotalSales

	ToCashierMonthlySale(cashier *record.CashierRecordMonthSales) *response.CashierResponseMonthSales
	ToCashierMonthlySales(c []*record.CashierRecordMonthSales) []*response.CashierResponseMonthSales
	ToCashierYearlySale(cashier *record.CashierRecordYearSales) *response.CashierResponseYearSales
	ToCashierYearlySales(c []*record.CashierRecordYearSales) []*response.CashierResponseYearSales

	ToCashierResponse(cashier *record.CashierRecord) *response.CashierResponse
	ToCashiersResponse(cashiers []*record.CashierRecord) []*response.CashierResponse
	ToCashierResponseDeleteAt(cashier *record.CashierRecord) *response.CashierResponseDeleteAt
	ToCashiersResponseDeleteAt(cashiers []*record.CashierRecord) []*response.CashierResponseDeleteAt
}

type MerchantResponseMapper interface {
	ToMerchantResponse(merchant *record.MerchantRecord) *response.MerchantResponse
	ToMerchantsResponse(merchants []*record.MerchantRecord) []*response.MerchantResponse
	ToMerchantResponseDeleteAt(merchant *record.MerchantRecord) *response.MerchantResponseDeleteAt
	ToMerchantsResponseDeleteAt(merchants []*record.MerchantRecord) []*response.MerchantResponseDeleteAt
}

type MerchantDocumentResponseMapper interface {
	ToMerchantDocumentResponse(doc *record.MerchantDocumentRecord) *response.MerchantDocumentResponse
	ToMerchantDocumentsResponse(docs []*record.MerchantDocumentRecord) []*response.MerchantDocumentResponse

	ToMerchantDocumentResponseDeleteAt(doc *record.MerchantDocumentRecord) *response.MerchantDocumentResponseDeleteAt
	ToMerchantDocumentsResponseDeleteAt(docs []*record.MerchantDocumentRecord) []*response.MerchantDocumentResponseDeleteAt
}

type OrderResponseMapper interface {
	ToOrderMonthlyTotalRevenue(c *record.OrderMonthlyTotalRevenueRecord) *response.OrderMonthlyTotalRevenueResponse
	ToOrderMonthlyTotalRevenues(c []*record.OrderMonthlyTotalRevenueRecord) []*response.OrderMonthlyTotalRevenueResponse
	ToOrderYearlyTotalRevenue(c *record.OrderYearlyTotalRevenueRecord) *response.OrderYearlyTotalRevenueResponse
	ToOrderYearlyTotalRevenues(c []*record.OrderYearlyTotalRevenueRecord) []*response.OrderYearlyTotalRevenueResponse

	ToOrderMonthlyPrice(category *record.OrderMonthlyRecord) *response.OrderMonthlyResponse
	ToOrderMonthlyPrices(c []*record.OrderMonthlyRecord) []*response.OrderMonthlyResponse
	ToOrderYearlyPrice(category *record.OrderYearlyRecord) *response.OrderYearlyResponse
	ToOrderYearlyPrices(c []*record.OrderYearlyRecord) []*response.OrderYearlyResponse

	ToOrderResponse(order *record.OrderRecord) *response.OrderResponse
	ToOrdersResponse(orders []*record.OrderRecord) []*response.OrderResponse
	ToOrderResponseDeleteAt(order *record.OrderRecord) *response.OrderResponseDeleteAt
	ToOrdersResponseDeleteAt(orders []*record.OrderRecord) []*response.OrderResponseDeleteAt
}

type OrderItemResponseMapper interface {
	ToOrderItemResponse(order *record.OrderItemRecord) *response.OrderItemResponse
	ToOrderItemsResponse(orders []*record.OrderItemRecord) []*response.OrderItemResponse
	ToOrderItemResponseDeleteAt(order *record.OrderItemRecord) *response.OrderItemResponseDeleteAt
	ToOrderItemsResponseDeleteAt(orders []*record.OrderItemRecord) []*response.OrderItemResponseDeleteAt
}

type ProductResponseMapper interface {
	ToProductResponse(product *record.ProductRecord) *response.ProductResponse
	ToProductsResponse(products []*record.ProductRecord) []*response.ProductResponse
	ToProductResponseDeleteAt(product *record.ProductRecord) *response.ProductResponseDeleteAt
	ToProductsResponseDeleteAt(products []*record.ProductRecord) []*response.ProductResponseDeleteAt
}

type TransactionResponseMapper interface {
	ToTransactionMonthAmountSuccess(row *record.TransactionMonthlyAmountSuccessRecord) *response.TransactionMonthlyAmountSuccessResponse
	ToTransactionMonthlyAmountSuccess(rows []*record.TransactionMonthlyAmountSuccessRecord) []*response.TransactionMonthlyAmountSuccessResponse
	ToTransactionYearAmountSuccess(row *record.TransactionYearlyAmountSuccessRecord) *response.TransactionYearlyAmountSuccessResponse
	ToTransactionYearlyAmountSuccess(rows []*record.TransactionYearlyAmountSuccessRecord) []*response.TransactionYearlyAmountSuccessResponse
	ToTransactionMonthAmountFailed(row *record.TransactionMonthlyAmountFailedRecord) *response.TransactionMonthlyAmountFailedResponse
	ToTransactionMonthlyAmountFailed(rows []*record.TransactionMonthlyAmountFailedRecord) []*response.TransactionMonthlyAmountFailedResponse
	ToTransactionYearAmountFailed(row *record.TransactionYearlyAmountFailedRecord) *response.TransactionYearlyAmountFailedResponse
	ToTransactionYearlyAmountFailed(rows []*record.TransactionYearlyAmountFailedRecord) []*response.TransactionYearlyAmountFailedResponse
	ToTransactionMonthMethod(row *record.TransactionMonthlyMethodRecord) *response.TransactionMonthlyMethodResponse
	ToTransactionMonthlyMethod(rows []*record.TransactionMonthlyMethodRecord) []*response.TransactionMonthlyMethodResponse
	ToTransactionYearMethod(row *record.TransactionYearlyMethodRecord) *response.TransactionYearlyMethodResponse
	ToTransactionYearlyMethod(rows []*record.TransactionYearlyMethodRecord) []*response.TransactionYearlyMethodResponse

	ToTransactionResponse(transaction *record.TransactionRecord) *response.TransactionResponse
	ToTransactionsResponse(transactions []*record.TransactionRecord) []*response.TransactionResponse
	ToTransactionResponseDeleteAt(transaction *record.TransactionRecord) *response.TransactionResponseDeleteAt
	ToTransactionsResponseDeleteAt(transactions []*record.TransactionRecord) []*response.TransactionResponseDeleteAt
}
