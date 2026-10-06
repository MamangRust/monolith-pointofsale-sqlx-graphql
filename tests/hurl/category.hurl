# Category E2E (GraphQL /query) — lifecycle + semua stats pricing.
# Jalankan: hurl --variable baseUrl=http://localhost:5000 --variable uuid=<uuid> --test tests/hurl/category.hurl

# Register admin (self-contained)
POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Category\", lastname: \"Admin\", email: \"category.admin.{{uuid}}@example.com\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"
[Captures]
ownerId: jsonpath "$.data.registerUser.data.id"

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { loginUser(input: {email: \"category.admin.{{uuid}}@example.com\", password: \"password123\"}) { status message data { access_token } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.loginUser.status" == "success"
[Captures]
accessToken: jsonpath "$.data.loginUser.data.access_token"

# ── Merchant (untuk stats by merchant) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Toko Kategori {{uuid}}\", description: \"Toko untuk stats kategori\", address: \"Jl. Kategori No.7\", contact_email: \"kategori.{{uuid}}@example.com\", contact_phone: \"081298765427\", status: \"active\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchant.status" == "success"
[Captures]
merchantId: jsonpath "$.data.createMerchant.data.id"

# ── Create Category ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createCategory(input: {name: \"Minuman {{uuid}}\", description: \"Kategori minuman segar\"}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createCategory.status" == "success"
[Captures]
categoryId: jsonpath "$.data.createCategory.data.id"

# ── Find All / by ID / active / trashed ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllCategory(input: {page: 1, page_size: 10}) { status message data { id name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllCategory.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByIdCategory(input: {id: {{categoryId}}}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByIdCategory.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByActiveCategory(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByActiveCategory.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByTrashedCategory(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByTrashedCategory.status" == "success"

# ── Stats: total pricing (monthly/yearly + by merchant/by id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalPrices(input: {year: 2026, month: 6}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalPrices.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalPrices(input: {year: 2026}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalPrices.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalPricesByMerchant(input: {year: 2026, month: 6, merchant_id: {{merchantId}}}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalPricesByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalPricesByMerchant(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalPricesByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalPricesById(input: {year: 2026, month: 6, category_id: {{categoryId}}}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalPricesById.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalPricesById(input: {year: 2026, category_id: {{categoryId}}}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalPricesById.status" == "success"

# ── Stats: pricing ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthPrice(input: {year: 2026}) { status message data { month category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthPrice.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearPrice(input: {year: 2026}) { status message data { year category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearPrice.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthPriceByMerchant(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { month category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthPriceByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearPriceByMerchant(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { year category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearPriceByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthPriceById(input: {year: 2026, category_id: {{categoryId}}}) { status message data { month category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthPriceById.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearPriceById(input: {year: 2026, category_id: {{categoryId}}}) { status message data { year category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearPriceById.status" == "success"

# ── Update Category ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateCategory(input: {category_id: {{categoryId}}, name: \"Minuman Segar {{uuid}}\", description: \"Kategori minuman segar updated\"}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateCategory.status" == "success"

# ── Trash → Restore ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedCategory(input: {id: {{categoryId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedCategory.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreCategory(input: {id: {{categoryId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreCategory.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedCategory(input: {id: {{categoryId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedCategory.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllCategory { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllCategory.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedCategory(input: {id: {{categoryId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedCategory.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllCategoryPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllCategoryPermanent.status" == "success"

# ── Category kedua — coverage deleteCategoryPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createCategory(input: {name: \"Minuman All {{uuid}}\", description: \"Kategori untuk test permanent/all\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createCategory.status" == "success"
[Captures]
permCategoryId: jsonpath "$.data.createCategory.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedCategory(input: {id: {{permCategoryId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedCategory.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteCategoryPermanent(input: {id: {{permCategoryId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteCategoryPermanent.status" == "success"
