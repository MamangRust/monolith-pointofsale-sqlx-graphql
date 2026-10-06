# Order Item E2E (GraphQL /query) — query only (tidak ada mutation di order_item service).
# Jalankan: hurl --variable baseUrl=http://localhost:5000 --variable uuid=<uuid> --test tests/hurl/order_item.hurl

# Register + login (self-contained)
POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"OrderItem\", lastname: \"Viewer\", email: \"orderitem.{{uuid}}@example.com\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { loginUser(input: {email: \"orderitem.{{uuid}}@example.com\", password: \"password123\"}) { status message data { access_token } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.loginUser.status" == "success"
[Captures]
accessToken: jsonpath "$.data.loginUser.data.access_token"

# ── Find all order items ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllOrderItem(input: {page: 1, page_size: 10}) { status message data { id order_id product_id quantity } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllOrderItem.status" == "success"

# ── Find active / trashed ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByActiveOrderItem(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByActiveOrderItem.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByTrashedOrderItem(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByTrashedOrderItem.status" == "success"

# ── Find by order id (1 = hasil seeder, kosong bila tidak ada — HTTP 200) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findOrderItemByOrder(input: {id: 1}) { status message data { id order_id product_id quantity } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findOrderItemByOrder.status" == "success"
