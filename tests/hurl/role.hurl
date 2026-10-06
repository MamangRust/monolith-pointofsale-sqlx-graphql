# Role E2E (GraphQL /query) — lifecycle + active/trashed/find-by-user + restore-all/permanent-all.
# Register auto-assign ROLE_ADMIN sehingga role query bisa diakses user baru.
# Jalankan: hurl --variable baseUrl=http://localhost:5000 --variable uuid=<uuid> --test tests/hurl/role.hurl

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Role\", lastname: \"Admin\", email: \"role.admin.{{uuid}}@example.com\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"
[Captures]
adminId: jsonpath "$.data.registerUser.data.id"

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { loginUser(input: {email: \"role.admin.{{uuid}}@example.com\", password: \"password123\"}) { status message data { access_token } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.loginUser.status" == "success"
[Captures]
accessToken: jsonpath "$.data.loginUser.data.access_token"

# ── Create Role ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createRole(input: {name: \"cashier-{{uuid}}\"}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createRole.status" == "success"
[Captures]
roleId: jsonpath "$.data.createRole.data.id"

# ── Find All Roles ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllRole(input: {page: 1, page_size: 10}) { status message data { id name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllRole.status" == "success"

# ── Find Role by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByIdRole(input: {role_id: {{roleId}}}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByIdRole.status" == "success"

# ── Find Active Roles ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByActiveRole(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByActiveRole.status" == "success"

# ── Find Trashed Roles ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByTrashedRole(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByTrashedRole.status" == "success"

# ── Find Roles by User ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByUserIdRole(input: {user_id: {{adminId}}}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByUserIdRole.status" == "success"

# ── Update Role ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateRole(input: {id: {{roleId}}, name: \"upd-{{uuid}}\"}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateRole.status" == "success"

# ── Trash → Restore ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedRole(input: {role_id: {{roleId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedRole.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreRole(input: {role_id: {{roleId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreRole.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedRole(input: {role_id: {{roleId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedRole.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllRole { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllRole.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedRole(input: {role_id: {{roleId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedRole.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllRolePermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllRolePermanent.status" == "success"

# ── Role kedua — coverage deleteRolePermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createRole(input: {name: \"perm-{{uuid}}\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createRole.status" == "success"
[Captures]
permRoleId: jsonpath "$.data.createRole.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedRole(input: {role_id: {{permRoleId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedRole.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteRolePermanent(input: {role_id: {{permRoleId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteRolePermanent.status" == "success"
