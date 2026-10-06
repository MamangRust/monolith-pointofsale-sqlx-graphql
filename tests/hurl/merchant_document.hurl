# Merchant Document E2E (GraphQL /query) — lifecycle + regression DocumentID vs MerchantID.
# Jalankan: hurl --variable baseUrl=http://localhost:5000 --variable uuid=<uuid> --test tests/hurl/merchant_document.hurl

# Register owner
POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Doc\", lastname: \"Owner\", email: \"doc.owner.{{uuid}}@example.com\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"
[Captures]
ownerId: jsonpath "$.data.registerUser.data.id"

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { loginUser(input: {email: \"doc.owner.{{uuid}}@example.com\", password: \"password123\"}) { status message data { access_token } } }"
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
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Toko Dokumen {{uuid}}\", description: \"Toko untuk test dokumen\", address: \"Jl. Dokumen No.6\", contact_email: \"dokumen.{{uuid}}@example.com\", contact_phone: \"081298765428\", status: \"active\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchant.status" == "success"
[Captures]
merchantId: jsonpath "$.data.createMerchant.data.id"

# ── Document 1 ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchantDocument(input: {merchant_id: {{merchantId}}, document_type: \"license\", document_url: \"https://example.com/license-{{uuid}}.pdf\"}) { status message data { document_id merchant_id document_type } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchantDocument.status" == "success"
[Captures]
documentId: jsonpath "$.data.createMerchantDocument.data.document_id"

# ── Document 2 — ID-nya beda dari merchant_id; update harus menarget DocumentID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchantDocument(input: {merchant_id: {{merchantId}}, document_type: \"id_card\", document_url: \"https://example.com/id-{{uuid}}.pdf\"}) { status message data { document_id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchantDocument.status" == "success"
[Captures]
documentId2: jsonpath "$.data.createMerchantDocument.data.document_id"

# ── Find All Documents ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllMerchantDocuments(input: {page: 1, page_size: 10}) { status message data { document_id merchant_id document_type } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllMerchantDocuments.status" == "success"

# ── Find Document by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMerchantDocumentById(input: {document_id: {{documentId}}}) { status message data { document_id document_type } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMerchantDocumentById.status" == "success"

# ── Find Active Documents ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllActiveMerchantDocuments(input: {page: 1, page_size: 10}) { status message data { document_id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllActiveMerchantDocuments.status" == "success"

# ── Update Document 2 (status/note ikut dikirim) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateMerchantDocument(input: {document_id: {{documentId2}}, merchant_id: {{merchantId}}, document_type: \"id_card_updated\", document_url: \"https://example.com/id-updated-{{uuid}}.pdf\", status: \"pending\", note: \"reupload\"}) { status message data { document_id document_type status } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateMerchantDocument.status" == "success"

# ── Update Status Document 2 ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateMerchantDocumentStatus(input: {document_id: {{documentId2}}, merchant_id: {{merchantId}}, status: \"verified\", note: \"Approved\"}) { status message data { document_id status } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateMerchantDocumentStatus.status" == "success"

# ── Trash → list trashed → restore (dokumen 1) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedMerchantDocument(input: {document_id: {{documentId}}}) { status message data { document_id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedMerchantDocument.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllTrashedMerchantDocuments(input: {page: 1, page_size: 10}) { status message data { document_id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllTrashedMerchantDocuments.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreMerchantDocument(input: {document_id: {{documentId}}}) { status message data { document_id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreMerchantDocument.status" == "success"

# ── Delete permanent dokumen 2 ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedMerchantDocument(input: {document_id: {{documentId2}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedMerchantDocument.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteMerchantDocumentPermanent(input: {document_id: {{documentId2}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteMerchantDocumentPermanent.status" == "success"

# ── Coverage restoreAll + deleteAllPermanent — dokumen 1 di-trash dulu ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedMerchantDocument(input: {document_id: {{documentId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedMerchantDocument.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllMerchantDocuments { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllMerchantDocuments.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedMerchantDocument(input: {document_id: {{documentId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedMerchantDocument.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllMerchantDocumentsPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllMerchantDocumentsPermanent.status" == "success"
