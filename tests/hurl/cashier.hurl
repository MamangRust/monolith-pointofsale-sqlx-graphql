# Cashier E2E (GraphQL /query) — create merchant → create cashier → lifecycle + stats sales.
# Jalankan: hurl --variable baseUrl=http://localhost:5000 --variable uuid=<uuid> --test tests/hurl/cashier.hurl

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Cashier\", lastname: \"Owner\", email: \"cashier.owner.{{uuid}}@example.com\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"
[Captures]
ownerId: jsonpath "$.data.registerUser.data.id"

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { loginUser(input: {email: \"cashier.owner.{{uuid}}@example.com\", password: \"password123\"}) { status message data { access_token } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.loginUser.status" == "success"
[Captures]
accessToken: jsonpath "$.data.loginUser.data.access_token"

# ── Merchant parent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Toko Kasir {{uuid}}\", description: \"Toko untuk test cashier\", address: \"Jl. Sudirman No.2\", contact_email: \"kasir.{{uuid}}@example.com\", contact_phone: \"081298765432\", status: \"active\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchant.status" == "success"
[Captures]
merchantId: jsonpath "$.data.createMerchant.data.id"

# ── Create Cashier ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createCashier(input: {merchant_id: {{merchantId}}, user_id: {{ownerId}}, name: \"Kasir Utama {{uuid}}\"}) { status message data { id merchant_id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createCashier.status" == "success"
[Captures]
cashierId: jsonpath "$.data.createCashier.data.id"

# ── Find All / by ID / active / trashed / by merchant ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllCashier(input: {page: 1, page_size: 10}) { status message data { id name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllCashier.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByIdCashier(input: {id: {{cashierId}}}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByIdCashier.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByActiveCashier(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByActiveCashier.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByTrashedCashier(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByTrashedCashier.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByMerchantCashier(input: {merchant_id: {{merchantId}}, page: 1, page_size: 10}) { status message data { id merchant_id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByMerchantCashier.status" == "success"

# ── Stats: total sales (monthly/yearly + by merchant + by cashier) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalSales(input: {year: 2026, month: 6}) { status message data { month total_sales } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalSales.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalSales(input: {year: 2026}) { status message data { year total_sales } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalSales.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalSalesByMerchant(input: {year: 2026, month: 6, merchant_id: {{merchantId}}}) { status message data { month total_sales } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalSalesByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalSalesByMerchant(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { year total_sales } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalSalesByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalSalesById(input: {year: 2026, month: 6, cashier_id: {{cashierId}}}) { status message data { month total_sales } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalSalesById.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalSalesById(input: {year: 2026, cashier_id: {{cashierId}}}) { status message data { year total_sales } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalSalesById.status" == "success"

# ── Stats: sales ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthSales(input: {year: 2026}) { status message data { month cashier_id total_sales } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthSales.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearSales(input: {year: 2026}) { status message data { year cashier_id total_sales } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearSales.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthSalesByMerchant(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { month cashier_id total_sales } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthSalesByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearSalesByMerchant(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { year cashier_id total_sales } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearSalesByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthSalesById(input: {year: 2026, cashier_id: {{cashierId}}}) { status message data { month cashier_id total_sales } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthSalesById.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearSalesById(input: {year: 2026, cashier_id: {{cashierId}}}) { status message data { year cashier_id total_sales } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearSalesById.status" == "success"

# ── Update Cashier ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateCashier(input: {cashier_id: {{cashierId}}, name: \"Kasir Senior {{uuid}}\"}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateCashier.status" == "success"

# ── Trash → Restore ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedCashier(input: {id: {{cashierId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedCashier.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreCashier(input: {id: {{cashierId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreCashier.status" == "success"

# ── Delete permanent cashier 1 ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedCashier(input: {id: {{cashierId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedCashier.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteCashierPermanent(input: {id: {{cashierId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteCashierPermanent.status" == "success"

# ── Cashier kedua — coverage restoreAll + deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createCashier(input: {merchant_id: {{merchantId}}, user_id: {{ownerId}}, name: \"Kasir All {{uuid}}\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createCashier.status" == "success"
[Captures]
cashierId2: jsonpath "$.data.createCashier.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedCashier(input: {id: {{cashierId2}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedCashier.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllCashier { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllCashier.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedCashier(input: {id: {{cashierId2}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedCashier.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllCashierPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllCashierPermanent.status" == "success"
