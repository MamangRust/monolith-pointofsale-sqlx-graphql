# User E2E (GraphQL /query) — lifecycle + active/trashed + restore-all/permanent-all.
# Jalankan: hurl --variable baseUrl=http://localhost:5000 --variable uuid=<uuid> --test tests/hurl/user.hurl

# Register admin (self-contained)
POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"User\", lastname: \"Admin\", email: \"user.admin.{{uuid}}@example.com\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"
[Captures]
adminId: jsonpath "$.data.registerUser.data.id"

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { loginUser(input: {email: \"user.admin.{{uuid}}@example.com\", password: \"password123\"}) { status message data { access_token } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.loginUser.status" == "success"
[Captures]
accessToken: jsonpath "$.data.loginUser.data.access_token"

# ── Create User ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createUser(input: {firstname: \"Jane\", lastname: \"Doe\", email: \"jane.{{uuid}}@example.com\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id email } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createUser.status" == "success"
[Captures]
newUserId: jsonpath "$.data.createUser.data.id"

# ── Find All / Active / Trashed ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllUsers(input: {page: 1, page_size: 10}) { status message data { id email } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllUsers.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByActiveUsers(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByActiveUsers.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByTrashedUsers(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByTrashedUsers.status" == "success"

# ── Find User by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByIdUser(input: {id: {{newUserId}}}) { status message data { id email } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByIdUser.status" == "success"

# ── Update User ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateUser(input: {id: {{newUserId}}, firstname: \"Jane\", lastname: \"Smith\", email: \"jane.smith.{{uuid}}@example.com\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id firstname } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateUser.status" == "success"

# ── Trash → Restore ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedUser(input: {id: {{newUserId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedUser.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreUser(input: {id: {{newUserId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreUser.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedUser(input: {id: {{newUserId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedUser.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllUser { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllUser.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedUser(input: {id: {{newUserId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedUser.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllUserPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllUserPermanent.status" == "success"

# ── User kedua — coverage deleteUserPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createUser(input: {firstname: \"Perm\", lastname: \"User\", email: \"perm.{{uuid}}@example.com\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createUser.status" == "success"
[Captures]
permUserId: jsonpath "$.data.createUser.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedUser(input: {id: {{permUserId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedUser.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteUserPermanent(input: {id: {{permUserId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteUserPermanent.status" == "success"
