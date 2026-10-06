# Merchant E2E (GraphQL /query) — lifecycle: create → update → trash → restore → delete permanent.
# Jalankan: hurl --variable baseUrl=http://localhost:5000 --variable uuid=<uuid> --test tests/hurl/merchant.hurl

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Merchant\", lastname: \"Owner\", email: \"owner.{{uuid}}@example.com\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"
[Captures]
ownerId: jsonpath "$.data.registerUser.data.id"

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { loginUser(input: {email: \"owner.{{uuid}}@example.com\", password: \"password123\"}) { status message data { access_token } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.loginUser.status" == "success"
[Captures]
accessToken: jsonpath "$.data.loginUser.data.access_token"

# ── Create Merchant ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Toko Kopi {{uuid}}\", description: \"Kopi nusantara\", address: \"Jl. Merdeka No.1\", contact_email: \"toko.{{uuid}}@example.com\", contact_phone: \"081234567890\", status: \"active\"}) { status message data { id name status } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchant.status" == "success"
[Captures]
merchantId: jsonpath "$.data.createMerchant.data.id"

# ── Find All Merchants ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllMerchant(input: {page: 1, page_size: 10}) { status message data { id name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllMerchant.status" == "success"

# ── Find Merchant by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByIdMerchant(input: {id: {{merchantId}}}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByIdMerchant.status" == "success"

# ── Update Merchant ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateMerchant(input: {merchant_id: {{merchantId}}, user_id: {{ownerId}}, name: \"Toko Kopi Premium {{uuid}}\", description: \"Kopi nusantara premium\", address: \"Jl. Kopi No. 1\", contact_email: \"kopi.{{uuid}}@example.com\", contact_phone: \"081298765432\", status: \"active\"}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateMerchant.status" == "success"

# ── Find Active Merchants ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByActiveMerchant(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByActiveMerchant.status" == "success"

# ── Find Trashed Merchants ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByTrashedMerchant(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByTrashedMerchant.status" == "success"

# ── Trash → Restore ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedMerchant(input: {id: {{merchantId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreMerchant(input: {id: {{merchantId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreMerchant.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedMerchant(input: {id: {{merchantId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllMerchant { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedMerchant(input: {id: {{merchantId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllMerchantPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllMerchantPermanent.status" == "success"

# ── Merchant kedua — coverage deleteMerchantPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Toko Perm {{uuid}}\", description: \"Toko untuk test permanent\", address: \"Jl. Perm No.9\", contact_email: \"perm.{{uuid}}@example.com\", contact_phone: \"081298765431\", status: \"active\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchant.status" == "success"
[Captures]
permMerchantId: jsonpath "$.data.createMerchant.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedMerchant(input: {id: {{permMerchantId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteMerchantPermanent(input: {id: {{permMerchantId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteMerchantPermanent.status" == "success"
