# Product E2E (GraphQL /query) — create merchant+category → create product (upload image) → lifecycle.
# image: Upload! required pada create → request multipart (GraphQL multipart spec).
# Jalankan: hurl --variable baseUrl=http://localhost:5000 --variable uuid=<uuid> --test tests/hurl/product.hurl

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Product\", lastname: \"Owner\", email: \"product.owner.{{uuid}}@example.com\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"
[Captures]
ownerId: jsonpath "$.data.registerUser.data.id"

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { loginUser(input: {email: \"product.owner.{{uuid}}@example.com\", password: \"password123\"}) { status message data { access_token } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.loginUser.status" == "success"
[Captures]
accessToken: jsonpath "$.data.loginUser.data.access_token"

# ── Merchant + category parent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Toko Produk {{uuid}}\", description: \"Toko untuk test produk\", address: \"Jl. Produk No.3\", contact_email: \"produk.{{uuid}}@example.com\", contact_phone: \"081298765431\", status: \"active\"}) { status message data { id } } }"
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
    "query": "mutation { createCategory(input: {name: \"Kopi-{{uuid}}\", description: \"Kategori kopi\"}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createCategory.status" == "success"
[Captures]
categoryId: jsonpath "$.data.createCategory.data.id"
categoryName: jsonpath "$.data.createCategory.data.name"

# ── Create Product (multipart + image upload) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: CreateProductInput!){createProduct(input: $input){status message data{id name price}}}","variables":{"input":{"merchantId":{{merchantId}},"categoryId":{{categoryId}},"name":"Kopi Susu {{uuid}}","description":"Es kopi susu gula aren","price":25000,"countInStock":100,"brand":"Nusantara","weight":250,"image":null}}}
map: {"0":["variables.input.image"]}
0: file,fixtures/product.jpg; image/jpeg
HTTP 200
[Asserts]
jsonpath "$.data.createProduct.status" == "success"
[Captures]
productId: jsonpath "$.data.createProduct.data.id"

# ── Find All Products ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllProduct(input: {page: 1, pageSize: 10}) { status message data { id name price } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllProduct.status" == "success"

# ── Find Product by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByIdProduct(input: {id: {{productId}}}) { status message data { id name price } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByIdProduct.status" == "success"

# ── Find by Merchant / Category ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByMerchantProduct(input: {merchantId: {{merchantId}}, page: 1, pageSize: 10, minPrice: 0, maxPrice: 0}) { status message data { id name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByMerchantProduct.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByCategoryProduct(input: {categoryName: \"{{categoryName}}\", page: 1, pageSize: 10, minPrice: 0, maxPrice: 0}) { status message data { id name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByCategoryProduct.status" == "success"

# ── Find Active Products ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByActiveProduct(input: {page: 1, pageSize: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByActiveProduct.status" == "success"

# ── Find Trashed Products ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByTrashedProduct(input: {page: 1, pageSize: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByTrashedProduct.status" == "success"

# ── Update Product (image opsional — tidak dikirim) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateProduct(input: {productId: {{productId}}, merchantId: {{merchantId}}, categoryId: {{categoryId}}, name: \"Kopi Susu Premium {{uuid}}\", description: \"Es kopi susu gula aren premium\", price: 30000, countInStock: 90, brand: \"Nusantara\", weight: 250}) { status message data { id name price } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateProduct.status" == "success"

# ── Trash → Restore ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedProduct(input: {id: {{productId}}}) { status message data { id deletedAt } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedProduct.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreProduct(input: {id: {{productId}}}) { status message data { id deletedAt } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreProduct.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedProduct(input: {id: {{productId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedProduct.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllProduct { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllProduct.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedProduct(input: {id: {{productId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedProduct.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllProductPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllProductPermanent.status" == "success"

# ── Product kedua — coverage deleteProductPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: CreateProductInput!){createProduct(input: $input){status message data{id}}}","variables":{"input":{"merchantId":{{merchantId}},"categoryId":{{categoryId}},"name":"Kopi Perm {{uuid}}","description":"Produk untuk test permanent","price":15000,"countInStock":20,"brand":"Nusantara","weight":100,"image":null}}}
map: {"0":["variables.input.image"]}
0: file,fixtures/product.jpg; image/jpeg
HTTP 200
[Asserts]
jsonpath "$.data.createProduct.status" == "success"
[Captures]
permProductId: jsonpath "$.data.createProduct.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedProduct(input: {id: {{permProductId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedProduct.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteProductPermanent(input: {id: {{permProductId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteProductPermanent.status" == "success"
