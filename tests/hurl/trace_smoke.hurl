# Trace smoke E2E (GraphQL /query) — full chain: merchant → cashier → category → product → order → create transaction.
# Jalankan (dari tests/smoke/trace_smoke.sh, bukan manual):
#   hurl --variable baseUrl=http://localhost:5000 --variable uuid=$UUID \
#        --variable traceparent="00-<trace-id>-<span-id>-01" \
#        --test tests/hurl/trace_smoke.hurl
#
# Khusus request create-transaction membawa header `traceparent` dengan
# trace-id tetap. Trace yang dihasilkan harus kontinu di Jaeger:
#   apigateway:createTransaction → transaction-service:CreateTransaction
#   → Kafka email-service-topic-transaction-create → email-service:consume:...
# Catatan: amount dihitung ulang service (total + PPN 11%); request amount harus >= itu.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Trace\", lastname: \"Smoke\", email: \"trace.smoke.{{uuid}}@example.com\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"
[Captures]
ownerId: jsonpath "$.data.registerUser.data.id"

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { loginUser(input: {email: \"trace.smoke.{{uuid}}@example.com\", password: \"password123\"}) { status message data { access_token } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.loginUser.status" == "success"
[Captures]
accessToken: jsonpath "$.data.loginUser.data.access_token"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Toko Trace {{uuid}}\", description: \"Toko untuk smoke test trace\", address: \"Jl. Trace No.1\", contact_email: \"trace.{{uuid}}@example.com\", contact_phone: \"081298765429\", status: \"active\"}) { status message data { id } } }"
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
    "query": "mutation { createCashier(input: {merchant_id: {{merchantId}}, user_id: {{ownerId}}, name: \"Kasir Trace {{uuid}}\"}) { status message data { id } } }"
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
    "query": "mutation { createCategory(input: {name: \"Minuman Trace {{uuid}}\", description: \"Kategori minuman\"}) { status message data { id } } }"
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
operations: {"query":"mutation($input: CreateProductInput!){createProduct(input: $input){status message data{id}}}","variables":{"input":{"merchantId":{{merchantId}},"categoryId":{{categoryId}},"name":"Es Teh Trace {{uuid}}","description":"Es teh manis","price":10000,"countInStock":100,"brand":"Tehku","weight":200,"image":null}}}
map: {"0":["variables.input.image"]}
0: file,fixtures/product.jpg; image/jpeg
HTTP 200
[Asserts]
jsonpath "$.data.createProduct.status" == "success"
[Captures]
productId: jsonpath "$.data.createProduct.data.id"

# Order: total 20000 + PPN 2200 = 22200 → amount 30000 cukup
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createOrder(input: {merchant_id: {{merchantId}}, cashier_id: {{cashierId}}, items: [{product_id: {{productId}}, quantity: 2}]}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createOrder.status" == "success"
[Captures]
orderId: jsonpath "$.data.createOrder.data.id"

# Request kunci: membawa traceparent agar trace-id dapat diverifikasi di Jaeger
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
traceparent: {{traceparent}}
Content-Type: application/json
{
    "query": "mutation { createTransaction(input: {orderId: {{orderId}}, cashierId: {{cashierId}}, paymentMethod: \"cash\", amount: 30000, paymentStatus: \"success\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createTransaction.status" == "success"
[Captures]
transactionId: jsonpath "$.data.createTransaction.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllTransaction(input: {page: 1, pageSize: 10}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllTransaction.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByIdTransaction(input: {id: {{transactionId}}}) { status message data { id orderId amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByIdTransaction.status" == "success"

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
    "query": "mutation { trashedTransaction(input: {id: {{transactionId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedTransaction.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreTransaction(input: {id: {{transactionId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreTransaction.status" == "success"

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
    "query": "mutation { deleteTransactionPermanent(input: {id: {{transactionId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteTransactionPermanent.status" == "success"
