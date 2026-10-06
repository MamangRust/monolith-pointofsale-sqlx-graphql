#!/usr/bin/env bash
# E2E semua endpoint via hurl.
#
# Alur:
#   1. Infra up (docker compose) — pg, redis, kafka, observability tools
#   2. Tunggu postgres + redis + kafka siap
#   3. Migrate + seeder
#   4. Build service lokal (just build → bin/)
#   5. Start service lokal (terhubung ke kafka compose; email consumer tidak dijalankan)
#   6. Jalankan hurl suite di tests/hurl/
#
# Prasyarat: docker + docker compose, hurl, just, go.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT"

BASE_URL="${BASE_URL:-http://localhost:5000}"
COMPOSE_FILE="deployments/local/docker-compose.infra.yml"
LOCAL_SERVICES=(auth role user category cashier merchant order_item order product transaction apigateway)
LOGS_DIR="$ROOT/logs"

mkdir -p "$LOGS_DIR"

# Bersihkan service lokal dari run sebelumnya (hindari port sudah terpakai)
pkill -f 'bin/(auth|role|user|category|cashier|merchant|order_item|order|product|transaction|apigateway)' 2>/dev/null || true
sleep 1

pids=()
cleanup() {
    echo ""
    echo "==> Membersihkan service lokal..."
    for pid in "${pids[@]:-}"; do
        kill "$pid" 2>/dev/null || true
    done
    pkill -f 'bin/(auth|role|user|category|cashier|merchant|order_item|order|product|transaction|apigateway)' 2>/dev/null || true
}
trap cleanup EXIT

echo "==> [1/6] Infra up (docker compose) — pg + redis + kafka + observability"
docker compose -f "$COMPOSE_FILE" up -d

echo "==> [2/6] Tunggu postgres, redis & kafka siap..."
for _ in $(seq 1 40); do
    pg_ok=false; redis_ok=false; kafka_ok=false
    if docker exec postgres pg_isready -U DRAGON -d POINTOFSALE >/dev/null 2>&1; then pg_ok=true; fi
    if docker exec redis_pointofsale redis-cli -a dragon_knight ping 2>/dev/null | grep -q PONG; then redis_ok=true; fi
    if docker exec my-kafka-pointofsale /opt/kafka/bin/kafka-broker-api-versions.sh --bootstrap-server localhost:9092 >/dev/null 2>&1; then kafka_ok=true; fi
    if $pg_ok && $redis_ok && $kafka_ok; then break; fi
    sleep 2
done
docker exec postgres pg_isready -U DRAGON -d POINTOFSALE >/dev/null 2>&1 || { echo "❌ postgres tidak siap"; exit 1; }
docker exec redis_pointofsale redis-cli -a dragon_knight ping 2>/dev/null | grep -q PONG || { echo "❌ redis tidak siap"; exit 1; }
docker exec my-kafka-pointofsale /opt/kafka/bin/kafka-broker-api-versions.sh --bootstrap-server localhost:9092 >/dev/null 2>&1 || { echo "❌ kafka tidak siap"; exit 1; }

echo "==> [3/6] Reset DB → Migrate → Seeder (idempotent)"
docker exec postgres psql -U DRAGON -d POINTOFSALE -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" >/dev/null
go run service/migrate/main.go up
go run service/seeder/main.go

echo "==> [4/6] Build service lokal"
mkdir -p "$ROOT/bin"
for s in "${LOCAL_SERVICES[@]}"; do
    if [ -f "service/$s/cmd/main.go" ]; then
        (cd "service/$s" && go build -o "$ROOT/bin/$s" ./cmd/main.go) || { echo "❌ build $s gagal"; exit 1; }
    else
        (cd "service/$s" && go build -o "$ROOT/bin/$s" ./main.go) || { echo "❌ build $s gagal"; exit 1; }
    fi
done

echo "==> [5/6] Start service lokal (terhubung kafka compose)"
# Export seluruh .env ke process env supaya os.Getenv("GRPC_*_ADDR") dll.
# terbaca (tanpa ini service memakai default hardcoded yang bisa salah port,
# mis. cashier → merchant default :50055 = port cashier sendiri).
set -a
# shellcheck disable=SC1091
source "$ROOT/.env"
set +a
export KAFKA_BROKERS=localhost:29092 APP_ENV=development
for s in "${LOCAL_SERVICES[@]}"; do
    echo "  → $s"
    nohup "./bin/$s" >> "$LOGS_DIR/$s.log" 2>&1 &
    pids+=("$!")
done

sleep 4
for i in "${!pids[@]}"; do
    if ! kill -0 "${pids[$i]}" 2>/dev/null; then
        echo "❌ ${LOCAL_SERVICES[$i]} mati saat startup — cek $LOGS_DIR/${LOCAL_SERVICES[$i]}.log"
        exit 1
    fi
done

echo "==> Tunggu apigateway :5000 + semua gRPC service siap..."
# /query dilayani langsung oleh apigateway (GraphQL), jadi menunggunya saja
# tidak cukup — service gRPC (auth, role, dst.) butuh waktu init DB pool lebih
# lama. Tunggu semua port gRPC listening dulu baru hurl (mencegah 503 spurious).
# Server dianggap siap bila /query membalas HTTP apapun (2xx/4xx/5xx) —
# bukan connection refused. Query tanpa token ditolak middleware (401),
# jadi cukup cek respons HTTP non-kosong.
probe_query() {
    code=$(curl -s -o /dev/null -w '%{http_code}' -X POST "$BASE_URL/query" \
        -H 'Content-Type: application/json' \
        -d '{"query":"query { __typename }"}' || true)
    [ -n "$code" ] && [ "$code" != "000" ]
}
GRPC_PORTS="50051 50052 50053 50054 50055 50056 50057 50058 50059 50060"
for _ in $(seq 1 40); do
    ready=1
    for p in $GRPC_PORTS; do
        if ! timeout 1 bash -c "</dev/tcp/127.0.0.1/$p" 2>/dev/null; then
            ready=0
            break
        fi
    done
    if [ "$ready" -eq 1 ] && probe_query; then
        break
    fi
    sleep 2
done
for p in $GRPC_PORTS; do
    timeout 1 bash -c "</dev/tcp/127.0.0.1/$p" 2>/dev/null || { echo "❌ service gRPC port $p tidak siap"; exit 1; }
done
probe_query || { echo "❌ apigateway tidak siap"; exit 1; }

echo "==> [6/6] Jalankan hurl suite (semua endpoint)"
cd "$ROOT/tests/hurl"
UUID_VAL="$(uuidgen)"
# trace_smoke.hurl dijalankan khusus oleh tests/smoke/trace_smoke.sh (butuh traceparent),
# jadi dikeluarkan dari suite E2E umum. Catatan: hurl tidak mendukung glob negasi
# ('--glob !x.hurl'), jadi pakai daftar file eksplisit (portabel lintas versi hurl).
HURL_FILES=$(ls *.hurl | grep -v '^trace_smoke\.hurl$')

# Jalankan per-file dengan jeda: apigateway memakai rate limiter 20 rps
# (burst 200), jadi melempar seluruh suite sekaligus akan kena 429.
PASS=0; FAIL=0; FAILED=()
for f in $HURL_FILES; do
    echo "  → $f"
    if hurl --variable baseUrl="$BASE_URL" --variable uuid="$UUID_VAL" --test "$f"; then
        PASS=$((PASS + 1))
        echo "  ✅ PASS $f"
    else
        FAIL=$((FAIL + 1))
        FAILED+=("$f")
        echo "  ❌ FAIL $f"
    fi
    sleep 2
done

echo
echo "== RESULT =="
echo "PASS: $PASS   FAIL: $FAIL   TOTAL: $((PASS + FAIL))"
if [ "$FAIL" -gt 0 ]; then
    echo
    echo "Failed suites:"
    for f in "${FAILED[@]}"; do echo "  - $f"; done
    exit 1
fi
echo "All E2E hurl suites passed."
