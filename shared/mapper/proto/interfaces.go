package protomapper

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb"
	pbcashier "github.com/MamangRust/monolith-pointofsale-pb/cashier"
	pbcategory "github.com/MamangRust/monolith-pointofsale-pb/category"
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pbmerchant_document "github.com/MamangRust/monolith-pointofsale-pb/merchant_document"
	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	pbproduct "github.com/MamangRust/monolith-pointofsale-pb/product"
	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
	pbtransaction "github.com/MamangRust/monolith-pointofsale-pb/transaction"
	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type AuthProtoMapper interface {
	ToProtoResponseVerifyCode(status string, message string) *pb.ApiResponseVerifyCode
	ToProtoResponseForgotPassword(status string, message string) *pb.ApiResponseForgotPassword
	ToProtoResponseResetPassword(status string, message string) *pb.ApiResponseResetPassword
	ToProtoResponseLogin(status string, message string, response *response.TokenResponse) *pb.ApiResponseLogin
	ToProtoResponseRegister(status string, message string, response *response.UserResponse) *pb.ApiResponseRegister
	ToProtoResponseRefreshToken(status string, message string, response *response.TokenResponse) *pb.ApiResponseRefreshToken
	ToProtoResponseGetMe(status string, message string, response *response.UserResponse) *pb.ApiResponseGetMe
}

type UserProtoMapper interface {
	ToProtoResponseUserDeleteAt(status string, message string, pbResponse *response.UserResponseDeleteAt) *pbuser.ApiResponseUserDeleteAt
	ToProtoResponsesUser(status string, message string, pbResponse []*response.UserResponse) *pbuser.ApiResponsesUser
	ToProtoResponseUser(status string, message string, pbResponse *response.UserResponse) *pbuser.ApiResponseUser
	ToProtoResponseUserDelete(status string, message string) *pbuser.ApiResponseUserDelete
	ToProtoResponseUserAll(status string, message string) *pbuser.ApiResponseUserAll
	ToProtoResponsePaginationUserDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.UserResponseDeleteAt) *pbuser.ApiResponsePaginationUserDeleteAt
	ToProtoResponsePaginationUser(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.UserResponse) *pbuser.ApiResponsePaginationUser
}

type RoleProtoMapper interface {
	ToProtoResponseRoleAll(status string, message string) *pbrole.ApiResponseRoleAll
	ToProtoResponseRoleDelete(status string, message string) *pbrole.ApiResponseRoleDelete
	ToProtoResponseRole(status string, message string, pbResponse *response.RoleResponse) *pbrole.ApiResponseRole
	ToProtoResponsesRole(status string, message string, pbResponse []*response.RoleResponse) *pbrole.ApiResponsesRole
	ToProtoResponsePaginationRole(pagination *pbcommon.PaginationMeta, status string, message string, pbResponse []*response.RoleResponse) *pbrole.ApiResponsePaginationRole
	ToProtoResponsePaginationRoleDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, pbResponse []*response.RoleResponseDeleteAt) *pbrole.ApiResponsePaginationRoleDeleteAt
}

type CategoryProtoMapper interface {
	ToProtoResponseMonthlyTotalPrice(status string, message string, row []*response.CategoriesMonthlyTotalPriceResponse) *pbcategory.ApiResponseCategoryMonthlyTotalPrice
	ToProtoResponseYearlyTotalPrice(status string, message string, row []*response.CategoriesYearlyTotalPriceResponse) *pbcategory.ApiResponseCategoryYearlyTotalPrice
	ToProtoResponseCategoryMonthlyPrice(status string, message string, row []*response.CategoryMonthPriceResponse) *pbcategory.ApiResponseCategoryMonthPrice
	ToProtoResponseCategoryYearlyPrice(status string, message string, row []*response.CategoryYearPriceResponse) *pbcategory.ApiResponseCategoryYearPrice

	ToProtoResponsesCategory(status string, message string, pbResponse []*response.CategoryResponse) *pbcategory.ApiResponsesCategory
	ToProtoResponseCategoryDeleteAt(status string, message string, pbResponse *response.CategoryResponseDeleteAt) *pbcategory.ApiResponseCategoryDeleteAt

	ToProtoResponseCategoryAll(status string, message string) *pbcategory.ApiResponseCategoryAll
	ToProtoResponseCategory(status string, message string, pbResponse *response.CategoryResponse) *pbcategory.ApiResponseCategory
	ToProtoResponseCategoryDelete(status string, message string) *pbcategory.ApiResponseCategoryDelete
	ToProtoResponsePaginationCategoryDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, categories []*response.CategoryResponseDeleteAt) *pbcategory.ApiResponsePaginationCategoryDeleteAt
	ToProtoResponsePaginationCategory(pagination *pbcommon.PaginationMeta, status string, message string, categories []*response.CategoryResponse) *pbcategory.ApiResponsePaginationCategory
}

type CashierProtoMapper interface {
	ToProtoMonthlyTotalSales(status, message string, row []*response.CashierResponseMonthTotalSales) *pbcashier.ApiResponseCashierMonthlyTotalSales
	ToProtoYearlyTotalSales(status, message string, row []*response.CashierResponseYearTotalSales) *pbcashier.ApiResponseCashierYearlyTotalSales

	ToProtoResponseMonthlyTotalSales(status, message string, row []*response.CashierResponseMonthSales) *pbcashier.ApiResponseCashierMonthSales
	ToProtoResponseYearlyTotalSales(status, message string, row []*response.CashierResponseYearSales) *pbcashier.ApiResponseCashierYearSales

	ToProtoResponseCashier(status string, message string, pbResponse *response.CashierResponse) *pbcashier.ApiResponseCashier
	ToProtoResponseCashierDeleteAt(status string, message string, pbResponse *response.CashierResponseDeleteAt) *pbcashier.ApiResponseCashierDeleteAt
	ToProtoResponsesCashier(status string, message string, pbResponse []*response.CashierResponse) *pbcashier.ApiResponsesCashier
	ToProtoResponseCashierDelete(status string, message string) *pbcashier.ApiResponseCashierDelete
	ToProtoResponseCashierAll(status string, message string) *pbcashier.ApiResponseCashierAll
	ToProtoResponsePaginationCashierDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.CashierResponseDeleteAt) *pbcashier.ApiResponsePaginationCashierDeleteAt
	ToProtoResponsePaginationCashier(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.CashierResponse) *pbcashier.ApiResponsePaginationCashier
}

type MerchantProtoMapper interface {
	ToProtoResponseMerchant(status string, message string, pbResponse *response.MerchantResponse) *pbmerchant.ApiResponseMerchant
	ToProtoResponseMerchantDeleteAt(status string, message string, pbResponse *response.MerchantResponseDeleteAt) *pbmerchant.ApiResponseMerchantDeleteAt

	ToProtoResponsesMerchant(status string, message string, pbResponse []*response.MerchantResponse) *pbmerchant.ApiResponsesMerchant
	ToProtoResponseMerchantDelete(status string, message string) *pbmerchant.ApiResponseMerchantDelete
	ToProtoResponseMerchantAll(status string, message string) *pbmerchant.ApiResponseMerchantAll
	ToProtoResponsePaginationMerchantDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, merchants []*response.MerchantResponseDeleteAt) *pbmerchant.ApiResponsePaginationMerchantDeleteAt
	ToProtoResponsePaginationMerchant(pagination *pbcommon.PaginationMeta, status string, message string, merchants []*response.MerchantResponse) *pbmerchant.ApiResponsePaginationMerchant
}

type MerchantDocumentProtoMapper interface {
	ToProtoResponseMerchantDocument(status string, message string, doc *response.MerchantDocumentResponse) *pbmerchant_document.ApiResponseMerchantDocument
	ToProtoResponsesMerchantDocument(status string, message string, docs []*response.MerchantDocumentResponse) *pbmerchant_document.ApiResponsesMerchantDocument

	ToProtoResponsePaginationMerchantDocument(pagination *pbcommon.PaginationMeta, status string, message string, docs []*response.MerchantDocumentResponse) *pbmerchant_document.ApiResponsePaginationMerchantDocument
	ToProtoResponsePaginationMerchantDocumentDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, docs []*response.MerchantDocumentResponseDeleteAt) *pbmerchant_document.ApiResponsePaginationMerchantDocumentAt

	ToProtoResponseMerchantDocumentDelete(status string, message string) *pbmerchant_document.ApiResponseMerchantDocumentDelete

	ToProtoResponseMerchantDocumentAll(status string, message string) *pbmerchant_document.ApiResponseMerchantDocumentAll
}

type OrderItemProtoMapper interface {
	ToProtoResponseOrderItem(status string, message string, pbResponse *response.OrderItemResponse) *pb.ApiResponseOrderItem
	ToProtoResponsesOrderItem(status string, message string, pbResponse []*response.OrderItemResponse) *pb.ApiResponsesOrderItem
	ToProtoResponseOrderItemDelete(status string, message string) *pb.ApiResponseOrderItemDelete
	ToProtoResponseOrderItemAll(status string, message string) *pb.ApiResponseOrderItemAll
	ToProtoResponsePaginationOrderItemDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, orderItems []*response.OrderItemResponseDeleteAt) *pb.ApiResponsePaginationOrderItemDeleteAt
	ToProtoResponsePaginationOrderItem(pagination *pbcommon.PaginationMeta, status string, message string, orderItems []*response.OrderItemResponse) *pb.ApiResponsePaginationOrderItem
}

type OrderProtoMapper interface {
	ToProtoResponseMonthlyTotalRevenue(status string, message string, row []*response.OrderMonthlyTotalRevenueResponse) *pborder.ApiResponseOrderMonthlyTotalRevenue
	ToProtoResponseYearlyTotalRevenue(status string, message string, row []*response.OrderYearlyTotalRevenueResponse) *pborder.ApiResponseOrderYearlyTotalRevenue

	ToProtoResponseMonthlyRevenue(status string, message string, row []*response.OrderMonthlyResponse) *pborder.ApiResponseOrderMonthly
	ToProtoResponseYearlyRevenue(status string, message string, row []*response.OrderYearlyResponse) *pborder.ApiResponseOrderYearly

	ToProtoResponseOrder(status string, message string, pbResponse *response.OrderResponse) *pborder.ApiResponseOrder
	ToProtoResponseOrderDeleteAt(status string, message string, pbResponse *response.OrderResponseDeleteAt) *pborder.ApiResponseOrderDeleteAt
	ToProtoResponsesOrder(status string, message string, pbResponse []*response.OrderResponse) *pborder.ApiResponsesOrder
	ToProtoResponseOrderDelete(status string, message string) *pborder.ApiResponseOrderDelete
	ToProtoResponseOrderAll(status string, message string) *pborder.ApiResponseOrderAll
	ToProtoResponsePaginationOrderDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, orders []*response.OrderResponseDeleteAt) *pborder.ApiResponsePaginationOrderDeleteAt
	ToProtoResponsePaginationOrder(pagination *pbcommon.PaginationMeta, status string, message string, orders []*response.OrderResponse) *pborder.ApiResponsePaginationOrder
}

type ProductProtoMapper interface {
	ToProtoResponseProduct(status string, message string, pbResponse *response.ProductResponse) *pbproduct.ApiResponseProduct
	ToProtoResponseProductDeleteAt(status string, message string, pbResponse *response.ProductResponseDeleteAt) *pbproduct.ApiResponseProductDeleteAt

	ToProtoResponsesProduct(status string, message string, pbResponse []*response.ProductResponse) *pbproduct.ApiResponsesProduct
	ToProtoResponseProductDelete(status string, message string) *pbproduct.ApiResponseProductDelete
	ToProtoResponseProductAll(status string, message string) *pbproduct.ApiResponseProductAll
	ToProtoResponsePaginationProductDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, products []*response.ProductResponseDeleteAt) *pbproduct.ApiResponsePaginationProductDeleteAt
	ToProtoResponsePaginationProduct(pagination *pbcommon.PaginationMeta, status string, message string, products []*response.ProductResponse) *pbproduct.ApiResponsePaginationProduct
}

type TransactionProtoMapper interface {
	ToProtoResponseMonthAmountSuccess(status string, message string, row []*response.TransactionMonthlyAmountSuccessResponse) *pbtransaction.ApiResponseTransactionMonthAmountSuccess
	ToProtoResponseYearAmountSuccess(status string, message string, row []*response.TransactionYearlyAmountSuccessResponse) *pbtransaction.ApiResponseTransactionYearAmountSuccess
	ToProtoResponseMonthAmountFailed(status string, message string, row []*response.TransactionMonthlyAmountFailedResponse) *pbtransaction.ApiResponseTransactionMonthAmountFailed
	ToProtoResponseYearAmountFailed(status string, message string, row []*response.TransactionYearlyAmountFailedResponse) *pbtransaction.ApiResponseTransactionYearAmountFailed
	ToProtoResponseMonthMethod(status string, message string, row []*response.TransactionMonthlyMethodResponse) *pbtransaction.ApiResponseTransactionMonthPaymentMethod
	ToProtoResponseYearMethod(status string, message string, row []*response.TransactionYearlyMethodResponse) *pbtransaction.ApiResponseTransactionYearPaymentmethod

	ToProtoResponseTransaction(status string, message string, trans *response.TransactionResponse) *pbtransaction.ApiResponseTransaction
	ToProtoResponseTransactionDeleteAt(status string, message string, trans *response.TransactionResponseDeleteAt) *pbtransaction.ApiResponseTransactionDeleteAt
	ToProtoResponsesTransaction(status string, message string, transList []*response.TransactionResponse) *pbtransaction.ApiResponsesTransaction
	ToProtoResponseTransactionDelete(status string, message string) *pbtransaction.ApiResponseTransactionDelete
	ToProtoResponseTransactionAll(status string, message string) *pbtransaction.ApiResponseTransactionAll
	ToProtoResponsePaginationTransactionDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, transactions []*response.TransactionResponseDeleteAt) *pbtransaction.ApiResponsePaginationTransactionDeleteAt
	ToProtoResponsePaginationTransaction(pagination *pbcommon.PaginationMeta, status string, message string, transactions []*response.TransactionResponse) *pbtransaction.ApiResponsePaginationTransaction
}
