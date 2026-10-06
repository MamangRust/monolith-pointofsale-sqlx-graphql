# Order E2E (GraphQL /query) — full chain: merchant → cashier → category → product → create order → lifecycle.
# Catatan: total dihitung service (2 × qty 2 × 30000 = 120000).
# Jalankan: hurl --variable baseUrl=http://localhost:5000 --variable uuid=<uuid> --test tests/hurl/order.hurl

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Order\", lastname: \"Owner\", email: \"order.owner.{{uuid}}@example.com\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"
[Captures]
ownerId: jsonpath "$.data.registerUser.data.id"

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { loginUser(input: {email: \"order.owner.{{uuid}}@example.com\", password: \"password123\"}) { status message data { access_token } } }"
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
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Toko Order {{uuid}}\", description: \"Toko untuk test order\", address: \"Jl. Order No.4\", contact_email: \"order.{{uuid}}@example.com\", contact_phone: \"081298765430\", status: \"active\"}) { status message data { id } } }"
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
    "query": "mutation { createCashier(input: {merchant_id: {{merchantId}}, user_id: {{ownerId}}, name: \"Kasir Order {{uuid}}\"}) { status message data { id } } }"
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
    "query": "mutation { createCategory(input: {name: \"Makanan {{uuid}}\", description: \"Kategori makanan\"}) { status message data { id } } }"
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
operations: {"query":"mutation($input: CreateProductInput!){createProduct(input: $input){status message data{id}}}","variables":{"input":{"merchantId":{{merchantId}},"categoryId":{{categoryId}},"name":"Nasi Goreng {{uuid}}","description":"Nasi goreng spesial","price":30000,"countInStock":50,"brand":"Dapur","weight":300,"image":null}}}
map: {"0":["variables.input.image"]}
0: file,fixtures/product.jpg; image/jpeg
HTTP 200
[Asserts]
jsonpath "$.data.createProduct.status" == "success"
[Captures]
productId: jsonpath "$.data.createProduct.data.id"

# ── Create Order — total dihitung service: 2 × qty 2 × 30000 = 120000 ──
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

# ── Update Order — butuh order_item_id (dari findOrderItemByOrder) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findOrderItemByOrder(input: {id: {{orderId}}}) { status message data { id product_id quantity } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findOrderItemByOrder.status" == "success"
[Captures]
orderItemId: jsonpath "$.data.findOrderItemByOrder.data[0].id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateOrder(input: {order_id: {{orderId}}, items: [{order_item_id: {{orderItemId}}, product_id: {{productId}}, quantity: 3}]}) { status message data { id total_price } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateOrder.status" == "success"

# ── Find All / by ID / active / trashed / by merchant ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllOrder(input: {page: 1, page_size: 10}) { status message data { id total_price } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllOrder.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByIdOrder(input: {id: {{orderId}}}) { status message data { id total_price } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByIdOrder.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByMerchantOrder(input: {merchant_id: {{merchantId}}, page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByMerchantOrder.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByActiveOrder(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByActiveOrder.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByTrashedOrder(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByTrashedOrder.status" == "success"

# ── Stats: total revenue ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalRevenue(input: {year: 2026, month: 6}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalRevenue.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalRevenue(input: {year: 2026}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalRevenue.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalRevenueByMerchant(input: {year: 2026, month: 6, merchant_id: {{merchantId}}}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalRevenueByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalRevenueByMerchant(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalRevenueByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalRevenueById(input: {year: 2026, month: 6, order_id: {{orderId}}}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalRevenueById.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalRevenueById(input: {year: 2026, order_id: {{orderId}}}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalRevenueById.status" == "success"

# ── Stats: revenue ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyRevenue(input: {year: 2026}) { status message data { month order_count total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyRevenue.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyRevenue(input: {year: 2026}) { status message data { year order_count total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyRevenue.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyRevenueByMerchant(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { month order_count total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyRevenueByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyRevenueByMerchant(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { year order_count total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyRevenueByMerchant.status" == "success"

# ── Trash → Restore ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedOrder(input: {id: {{orderId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedOrder.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreOrder(input: {id: {{orderId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreOrder.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedOrder(input: {id: {{orderId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedOrder.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllOrder { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllOrder.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedOrder(input: {id: {{orderId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedOrder.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllOrderPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllOrderPermanent.status" == "success"

# ── Order kedua — coverage deleteOrderPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createOrder(input: {merchant_id: {{merchantId}}, cashier_id: {{cashierId}}, items: [{product_id: {{productId}}, quantity: 1}]}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createOrder.status" == "success"
[Captures]
orderId2: jsonpath "$.data.createOrder.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedOrder(input: {id: {{orderId2}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedOrder.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteOrderPermanent(input: {id: {{orderId2}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteOrderPermanent.status" == "success"
