# Hurl E2E — API Gateway (localhost:5000, GraphQL `/query`)

Suite e2e untuk `service/apigateway` (GraphQL :5000, endpoint `/query`).
Semua file memakai GraphQL mutation/query dengan pola response
`{ status, message, data }` ala API gateway ini. Register auto-verify di
environment non-production (`APP_ENV != production/kubernetes`), sehingga
register → login langsung berhasil.

## Prasyarat

- Stack berjalan: `just up` (atau service + gateway lokal).
- Seeder dijalankan agar data referensi (merchant, dll.) ada: `just seeder`.
- Email test memakai `{{uuid}}` sehingga aman dijalankan berulang (data
  unik per run).

## Menjalankan

Semua file memakai variabel `baseUrl` dan `uuid`:

```bash
hurl --variable baseUrl=http://localhost:5000 --variable uuid=$(uuidgen) --test tests/hurl/
```

Satu file saja:

```bash
hurl --variable baseUrl=http://localhost:5000 --variable uuid=$(uuidgen) --test tests/hurl/auth.hurl
```

## Isi suite

| File | Cakupan |
|---|---|
| `auth.hurl` | register → login → me → 401 → refresh-token → forgot-password → reset-password (invalid) → verify-code (invalid) |
| `user.hurl` | lifecycle user + active/trashed + restore-all/permanent-all + 404 |
| `role.hurl` | lifecycle role + active/trashed/find-by-user + restore-all/permanent-all |
| `category.hurl` | lifecycle category + semua stats (total-pricing & pricing, merchant & by-id) |
| `cashier.hurl` | create merchant → create cashier → lifecycle + semua stats sales |
| `merchant.hurl` | lifecycle merchant + active/trashed + restore-all/permanent-all |
| `merchant_document.hurl` | lifecycle document + update-status + regression DocumentID vs MerchantID |
| `product.hurl` | merchant + category + create product (upload multipart) → lifecycle + filter merchant/category + restore-all/permanent-all |
| `order.hurl` | chain penuh → create order → lifecycle + semua stats revenue |
| `transaction.hurl` | chain penuh → create transaction → lifecycle + semua stats status/method |
| `order_item.hurl` | query order item (query only) |
| `trace_smoke.hurl` | chain penuh dengan header `traceparent` (dijalankan khusus oleh `tests/smoke/trace_smoke.sh`) |

Semua file **self-contained** (register user sendiri di awal), sehingga bisa
dijalankan independen.

## E2E otomatis (docker compose infra + service lokal)

`just e2e-hurl` menjalankan `tests/hurl/run_e2e.sh`:

1. Infra up via `docker compose -f deployments/local/docker-compose.infra.yml`
   (postgres, redis, kafka, dan observability tools — service Go TIDAK
   di-container, dijalankan lokal).
2. Tunggu postgres/redis/kafka → reset DB → `go run service/migrate` → seeder.
3. Build service lokal → start service lokal (terhubung kafka compose).
4. Jalankan setiap `*.hurl` satu per satu (dengan jeda antar file agar tidak
   kena rate limiter gateway).

Manual:

```bash
just e2e-hurl                    # semuanya otomatis
# atau step-by-step:
just infra-up
just migrate && just seeder
just build && just services-up
hurl --variable baseUrl=http://localhost:5000 --variable uuid=$(uuidgen) --test tests/hurl/
just services-down && just infra-down
```

## Catatan

- Order/transaction memakai chain: register → merchant → cashier →
  category → product → order → transaction. Setiap `[Captures]` mengambil
  id dari `$.data.<operation>.data.id`; bila response mapper berubah,
  sesuaikan jsonpath.
- Product `image: Upload!` wajib multipart (GraphQL multipart spec) — file
  `product.hurl` memakai fixture `fixtures/product.jpg`.
- `transaction` menghitung ulang `amount` server-side (total + PPN 11%);
  request amount harus >= nilai tersebut (file sudah memakai nilai cukup).
- Status negatif: GraphQL error dibawa lewat `$.errors` (HTTP tetap 200)
  kecuali auth middleware yang mengembalikan `401`.
