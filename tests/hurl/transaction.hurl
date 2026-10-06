# Transaction E2E (GraphQL /query) — full chain: merchant → cashier → category → product → order → create transaction.
# Catatan: amount dihitung ulang service (total + PPN 11%); request amount harus >= itu.
# Input transaction memakai camelCase (orderId, cashierId, paymentMethod, dst).
# Jalankan: hurl --variable baseUrl=http://localhost:5000 --variable uuid=<uuid> --test tests/hurl/transaction.hurl

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Transaction\", lastname: \"Owner\", email: \"tx.owner.{{uuid}}@example.com\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"
[Captures]
ownerId: jsonpath "$.data.registerUser.data.id"

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { loginUser(input: {email: \"tx.owner.{{uuid}}@example.com\", password: \"password123\"}) { status message data { access_token } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.loginUser.status" == "success"
[Captures]
accessToken: jsonpath "$.data.loginUser.data.access_token"

# ── Merchant + cashier + category + product parent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Toko Tx {{uuid}}\", description: \"Toko untuk test transaksi\", address: \"Jl. Tx No.5\", contact_email: \"tx.{{uuid}}@example.com\", contact_phone: \"081298765429\", status: \"active\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchant.status" == "success"
[Captures]
merchantId: jsonpath "$.data.createMerchant.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createCashier(input: {merchant_id: {{merchantId}}, user_id: {{ownerId}}, name: \"Kasir Tx {{uuid}}\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createCashier.status" == "success"
[Captures]
cashierId: jsonpath "$.data.createCashier.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createCategory(input: {name: \"Minuman Tx {{uuid}}\", description: \"Kategori minuman\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createCategory.status" == "success"
[Captures]
categoryId: jsonpath "$.data.createCategory.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: CreateProductInput!){createProduct(input: $input){status message data{id}}}","variables":{"input":{"merchantId":{{merchantId}},"categoryId":{{categoryId}},"name":"Es Teh {{uuid}}","description":"Es teh manis","price":10000,"countInStock":100,"brand":"Tehku","weight":200,"image":null}}}
map: {"0":["variables.input.image"]}
0: file,fixtures/product.jpg; image/jpeg
HTTP 200
[Asserts]
jsonpath "$.data.createProduct.status" == "success"
[Captures]
productId: jsonpath "$.data.createProduct.data.id"

# ── Create Order: total 20000 + PPN 2200 = 22200 → amount 30000 cukup ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createOrder(input: {merchant_id: {{merchantId}}, cashier_id: {{cashierId}}, items: [{product_id: {{productId}}, quantity: 2}]}) { status message data { id total_price } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createOrder.status" == "success"
[Captures]
orderId: jsonpath "$.data.createOrder.data.id"

# ── Create Transaction ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createTransaction(input: {orderId: {{orderId}}, cashierId: {{cashierId}}, paymentMethod: \"cash\", amount: 30000, paymentStatus: \"success\"}) { status message data { id orderId paymentMethod amount paymentStatus } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createTransaction.status" == "success"
[Captures]
transactionId: jsonpath "$.data.createTransaction.data.id"

# ── Update Transaction — payment_status hasil create = \"success\" → boleh update ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateTransaction(input: {transactionId: {{transactionId}}, orderId: {{orderId}}, cashierId: {{cashierId}}, paymentMethod: \"cash\", amount: 30000, paymentStatus: \"success\"}) { status message data { id amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateTransaction.status" == "success"

# ── Find All / by ID / by merchant / active / trashed ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllTransaction(input: {page: 1, pageSize: 10}) { status message data { id orderId amount paymentStatus } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllTransaction.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByIdTransaction(input: {id: {{transactionId}}}) { status message data { id orderId amount paymentStatus } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByIdTransaction.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByMerchantTransaction(input: {merchantId: {{merchantId}}, page: 1, pageSize: 10}) { status message data { id orderId amount } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByMerchantTransaction.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByActiveTransaction(input: {page: 1, pageSize: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByActiveTransaction.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByTrashedTransaction(input: {page: 1, pageSize: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByTrashedTransaction.status" == "success"

# ── Stats: status ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthStatusSuccess(input: {year: 2026, month: 6}) { status message data { month totalSuccess totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthStatusSuccess.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearStatusSuccess(input: {year: 2026}) { status message data { year totalSuccess totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearStatusSuccess.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthStatusFailed(input: {year: 2026, month: 6}) { status message data { month totalFailed totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthStatusFailed.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearStatusFailed(input: {year: 2026}) { status message data { year totalFailed totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearStatusFailed.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthStatusSuccessByMerchant(input: {year: 2026, month: 6, merchantId: {{merchantId}}}) { status message data { month totalSuccess totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthStatusSuccessByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearStatusSuccessByMerchant(input: {year: 2026, merchantId: {{merchantId}}}) { status message data { year totalSuccess totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearStatusSuccessByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthStatusFailedByMerchant(input: {year: 2026, month: 6, merchantId: {{merchantId}}}) { status message data { month totalFailed totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthStatusFailedByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearStatusFailedByMerchant(input: {year: 2026, merchantId: {{merchantId}}}) { status message data { year totalFailed totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearStatusFailedByMerchant.status" == "success"

# ── Stats: method ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthMethodSuccess(input: {year: 2026, month: 6}) { status message data { month paymentMethod totalTransactions totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthMethodSuccess.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearMethodSuccess(input: {year: 2026}) { status message data { year paymentMethod totalTransactions totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearMethodSuccess.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthMethodFailed(input: {year: 2026, month: 6}) { status message data { month paymentMethod totalTransactions totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthMethodFailed.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearMethodFailed(input: {year: 2026}) { status message data { year paymentMethod totalTransactions totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearMethodFailed.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthMethodByMerchantSuccess(input: {year: 2026, month: 6, merchantId: {{merchantId}}}) { status message data { month paymentMethod totalTransactions totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthMethodByMerchantSuccess.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearMethodByMerchantSuccess(input: {year: 2026, merchantId: {{merchantId}}}) { status message data { year paymentMethod totalTransactions totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearMethodByMerchantSuccess.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthMethodByMerchantFailed(input: {year: 2026, month: 6, merchantId: {{merchantId}}}) { status message data { month paymentMethod totalTransactions totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthMethodByMerchantFailed.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearMethodByMerchantFailed(input: {year: 2026, merchantId: {{merchantId}}}) { status message data { year paymentMethod totalTransactions totalAmount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearMethodByMerchantFailed.status" == "success"

# ── Trash → Restore ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedTransaction(input: {id: {{transactionId}}}) { status message data { id deletedAt } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedTransaction.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreTransaction(input: {id: {{transactionId}}}) { status message data { id deletedAt } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreTransaction.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedTransaction(input: {id: {{transactionId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedTransaction.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllTransaction { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllTransaction.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedTransaction(input: {id: {{transactionId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedTransaction.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllTransactionPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllTransactionPermanent.status" == "success"

# ── Transaction kedua — coverage deleteTransactionPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createTransaction(input: {orderId: {{orderId}}, cashierId: {{cashierId}}, paymentMethod: \"cash\", amount: 30000, paymentStatus: \"success\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createTransaction.status" == "success"
[Captures]
transactionId2: jsonpath "$.data.createTransaction.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedTransaction(input: {id: {{transactionId2}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedTransaction.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteTransactionPermanent(input: {id: {{transactionId2}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteTransactionPermanent.status" == "success"
