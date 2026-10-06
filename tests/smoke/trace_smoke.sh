#!/usr/bin/env bash
#
# trace_smoke.sh — Smoke test: kontinuitas distributed trace
#
# Memverifikasi bahwa satu trace_id bertahan melewati seluruh lifecycle request:
#
#   HTTP client ──▶ API Gateway ──▶ gRPC transaction-service
#        │                └────────────▶ Kafka (email-service-topic-transaction-create)
#        │                                        └──▶ email-service consumer
#
# Alur verifikasi (via Jaeger query API):
#   1. Kirim flow create-transaction dengan header `traceparent` (trace-id tetap).
#   2. Ambil trace tsb dari Jaeger — semua span harus berada di trace-id yang sama.
#   3. Pastikan span kunci ada:
#        - post-transaction-create                     (apigateway)
#        - CreateTransaction                           (transaction-service)
#        - consume:email-service-topic-transaction-create (email-service)
#   4. Pastikan rantai ancestor email span menembus CreateTransaction dan
#      post-transaction-create (CHILD_OF) — membuktikan hop gRPC & Kafka
#      meneruskan konteks trace.
#
# Prasyarat stack lokal yang berjalan: gateway :5000, Kafka :29092, Jaeger :16686,
# OTel collector :4317. email-service di-start otomatis bila belum berjalan.
# Wajib: hurl, jq, openssl.
#
# Usage: ./tests/smoke/trace_smoke.sh     (atau: just smoke-trace)
#
# Env override: BASE_URL, JAEGER_URL, KAFKA_PORT, ENV_FILE, EMAIL_BIN, EMAIL_LOG,
# POLL_ATTEMPTS, POLL_INTERVAL

set -uo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
# dotenv.Viper() membaca .env dari CWD → jalankan semua dari project root.
cd "$ROOT"

BASE_URL="${BASE_URL:-http://localhost:5000}"
JAEGER_URL="${JAEGER_URL:-http://localhost:16686}"
KAFKA_PORT="${KAFKA_PORT:-29092}"
ENV_FILE="${ENV_FILE:-$ROOT/.env.local}"
EMAIL_BIN="${EMAIL_BIN:-/tmp/pos-bin/email}"
EMAIL_LOG="${EMAIL_LOG:-/tmp/pos-logs/email.log}"
POLL_ATTEMPTS="${POLL_ATTEMPTS:-45}"
POLL_INTERVAL="${POLL_INTERVAL:-1}"
SPAN_POLL_ATTEMPTS="${SPAN_POLL_ATTEMPTS:-20}"
CONSUMER_READY_WAIT="${CONSUMER_READY_WAIT:-30}"

GREEN=$'\033[0;32m'; RED=$'\033[0;31m'; YELLOW=$'\033[1;33m'; NC=$'\033[0m'
pass() { printf "%s✅ %s%s\n" "$GREEN" "$1" "$NC"; }
warn() { printf "%s⚠️  %s%s\n" "$YELLOW" "$1" "$NC"; }
fail() { printf "%s❌ %s%s\n" "$RED" "$1" "$NC"; }
die()  { fail "$1"; exit 1; }

port_open() { timeout 3 bash -c "cat < /dev/null > /dev/tcp/$1/$2" 2>/dev/null; }

preflight() {
  local ok=1 code
  code="$(curl -s -o /dev/null -w '%{http_code}' --max-time 3 -X POST "$BASE_URL/api/auth/login")"
  [ "$code" != "000" ] || { fail "API gateway tidak reachable ($BASE_URL, code=$code)"; ok=0; }
  port_open localhost "$KAFKA_PORT" || { fail "Kafka tidak reachable (localhost:$KAFKA_PORT)"; ok=0; }
  curl -sf --max-time 3 "$JAEGER_URL/api/services" >/dev/null || { fail "Jaeger tidak reachable ($JAEGER_URL)"; ok=0; }
  command -v hurl >/dev/null 2>&1 || { fail "hurl tidak terinstall"; ok=0; }
  command -v jq  >/dev/null 2>&1 || { fail "jq tidak terinstall"; ok=0; }
  command -v openssl >/dev/null 2>&1 || { fail "openssl tidak terinstall"; ok=0; }
  [ "$ok" = 1 ] || die "Preflight gagal — pastikan stack lokal berjalan."
  pass "Preflight OK (gateway, kafka, jaeger)"
}

# Memastikan email-service berjalan dengan binary yang punya marker
# 'Consumer session ready' (Setup handler). Binary lama / proses tanpa marker
# akan di-restart agar smoke test deterministik.
start_email_service() {
  warn "membuild & memulai email-service..."
  (cd "$ROOT/service/email" && go build -o "$EMAIL_BIN" ./cmd/) || die "Gagal build email-service"
  mkdir -p "$(dirname "$EMAIL_LOG")"
  # .env.local tidak mendefinisikan METRIC_EMAIL_ADDR → beri default agar
  # server /metrics tidak crash (ListenAndServe pada port kosong).
  set -a; . "$ENV_FILE"; set +a
  export METRIC_EMAIL_ADDR="${METRIC_EMAIL_ADDR:-8090}"
  # Log ditruncate agar marker 'Consumer session ready' hanya berasal dari run ini.
  (cd "$ROOT" && setsid -f "$EMAIL_BIN" > "$EMAIL_LOG" 2>&1 < /dev/null)
  sleep 3
  local pid
  # Anchor regex: cmdline proses email persis = path binary, hindari self-match
  # bila wrapper/CI punya path ini di command line-nya.
  pid="$(pgrep -f "^${EMAIL_BIN}$" | head -1)"
  [ -n "$pid" ] || die "email-service gagal start — lihat $EMAIL_LOG"
  pass "email-service dimulai (pid $pid, log $EMAIL_LOG)"
}

ensure_email_service() {
  local pid
  pid="$(pgrep -f "^${EMAIL_BIN}$" | head -1)"
  if [ -n "$pid" ]; then
    if grep -q 'Consumer session ready' "$EMAIL_LOG" 2>/dev/null; then
      pass "email-service sudah berjalan & konsumen siap (pid $pid)"
      return 0
    fi
    warn "email-service berjalan tapi konsumen belum siap (binary lama?) — restart..."
    kill "$pid" 2>/dev/null; sleep 1
  fi
  start_email_service
}

# Consumer group sarama memakai OffsetNewest: grup baru yang join setelah event
# dipublish akan melewatkannya. Tunggu marker readiness (dari Setup handler)
# sebelum menjalankan flow agar konsumen dijamin sudah join.
wait_consumer_ready() {
  local i attempt pid
  for attempt in 1 2; do
    echo "== Menunggu email consumer siap (percobaan $attempt, max ${CONSUMER_READY_WAIT}s) =="
    for i in $(seq 1 "$CONSUMER_READY_WAIT"); do
      if grep -q 'Consumer session ready' "$EMAIL_LOG" 2>/dev/null; then
        pass "email consumer siap (setelah ${i}s)"
        return 0
      fi
      sleep 1
    done
    if [ "$attempt" = 1 ]; then
      warn "Marker consumer belum muncul — restart email-service sekali..."
      pid="$(pgrep -f "^${EMAIL_BIN}$" | head -1)"
      [ -n "$pid" ] && kill "$pid" 2>/dev/null
      sleep 1
      start_email_service
    fi
  done
  warn "Marker 'Consumer session ready' tetap tidak muncul — lanjut (span email jadi verdict akhir)"
}

run_flow() {
  local traceparent="$1" uuid
  uuid="$(cat /proc/sys/kernel/random/uuid)"
  echo "== Menjalankan transaction flow (uuid=$uuid) =="
  if hurl --variable baseUrl="$BASE_URL" --variable uuid="$uuid" \
          --variable traceparent="$traceparent" \
          --test "$ROOT/tests/hurl/trace_smoke.hurl" >/tmp/trace_smoke_hurl.out 2>&1; then
    pass "Transaction flow e2e sukses"
  else
    warn "Transaction flow hurl gagal — tetap cek trace di Jaeger:"
    tail -25 /tmp/trace_smoke_hurl.out | sed 's/^/    /'
  fi
}

fetch_trace() { curl -sf --max-time 3 "$JAEGER_URL/api/traces/$1" 2>/dev/null; }

wait_for_trace() {
  local trace_id="$1" j
  for _ in $(seq 1 "$POLL_ATTEMPTS"); do
    j="$(fetch_trace "$trace_id")"
    if [ -n "$j" ] && ! echo "$j" | jq -e '.data | length == 0' >/dev/null 2>&1; then
      echo "$j"
      return 0
    fi
    sleep "$POLL_INTERVAL"
  done
  return 1
}

# wait_for_target_spans <trace-json>: re-fetch trace sampai ketiga span target
# hadir (OTel batch exporter membuat span tiba tidak serentak).
wait_for_target_spans() {
  local j="$1" i
  for i in $(seq 1 "$SPAN_POLL_ATTEMPTS"); do
    if span_present "$j" "post-transaction-create" "apigateway" \
       && span_present "$j" "CreateTransaction" "transaction-service" \
       && span_present "$j" "consume:email-service-topic-transaction-create" "email-service"; then
      echo "$j"
      return 0
    fi
    sleep 1
    j="$(fetch_trace "$TRACE_ID")"
    [ -n "$j" ] || j="$1"
  done
  echo "$j"
  return 1
}

# span_present <trace-json> <operation> <service>
span_present() {
  echo "$1" | jq -e --arg op "$2" --arg svc "$3" \
    '.data[0] as $t | $t.spans[] | select(.operationName == $op and $t.processes[.processID].serviceName == $svc)' >/dev/null 2>&1
}

# span_ops_chain <trace-json> <leaf-operation>: daftar operationName dari leaf
# span naik ke root (satu rantai CHILD_OF). Output: operationName per baris,
# dari daun ke akar.
span_ops_chain() {
  local json="$1" leaf_op="$2"
  local op_of parent_of leaf_sid sid opn chain=""
  op_of="$(echo "$json" | jq -r '.data[0].spans[] | "\(.spanID) \(.operationName)"')"
  parent_of="$(echo "$json" | jq -r '.data[0].spans[] | "\(.spanID) \(.references[]? | select(.refType == "CHILD_OF") | .spanID)"')"
  leaf_sid="$(echo "$op_of" | awk -v op="$leaf_op" '$2 == op {print $1; exit}')"
  [ -n "$leaf_sid" ] || return 1
  sid="$leaf_sid"
  for _ in $(seq 1 15); do
    opn="$(echo "$op_of" | awk -v id="$sid" '$1 == id {print $2; exit}')"
    [ -n "$opn" ] || break
    chain="${chain}${opn}\n"
    sid="$(echo "$parent_of" | awk -v id="$sid" '$1 == id {print $2; exit}')"
    [ -n "$sid" ] || break
  done
  printf "%b" "$chain"
}

main() {
  preflight
  ensure_email_service
  wait_consumer_ready

  TRACE_ID="$(openssl rand -hex 16)"
  SPAN_ID="$(openssl rand -hex 8)"
  TRACEPARENT="00-${TRACE_ID}-${SPAN_ID}-01"
  echo "== Trace ID: $TRACE_ID (traceparent: $TRACEPARENT) =="

  run_flow "$TRACEPARENT"

  echo "== Menunggu trace di Jaeger (max ${POLL_ATTEMPTS}s) =="
  TRACE_JSON="$(wait_for_trace "$TRACE_ID")" \
    || die "Trace $TRACE_ID tidak muncul di Jaeger dalam ${POLL_ATTEMPTS}s"
  pass "Trace $TRACE_ID ditemukan di Jaeger"

  TRACE_JSON="$(wait_for_target_spans "$TRACE_JSON")" \
    || warn "Span target belum lengkap setelah poll tambahan — detail di bawah"

  status=0

  span_present "$TRACE_JSON" "post-transaction-create" "apigateway" \
    && pass "Span post-transaction-create (apigateway)" \
    || { fail "Span post-transaction-create (apigateway) TIDAK ditemukan"; status=1; }

  span_present "$TRACE_JSON" "CreateTransaction" "transaction-service" \
    && pass "Span CreateTransaction (transaction-service)" \
    || { fail "Span CreateTransaction (transaction-service) TIDAK ditemukan"; status=1; }

  span_present "$TRACE_JSON" "consume:email-service-topic-transaction-create" "email-service" \
    && pass "Span consume:email-service-topic-transaction-create (email-service)" \
    || { fail "Span consume:... (email-service) TIDAK ditemukan"; status=1; }

  # Rantai ancestor email span harus menembus CreateTransaction & span gateway
  CHAIN="$(span_ops_chain "$TRACE_JSON" "consume:email-service-topic-transaction-create")"
  if echo "$CHAIN" | grep -q 'CreateTransaction' && echo "$CHAIN" | grep -q 'post-transaction-create'; then
    pass "Rantai CHILD_OF: email ← … ← CreateTransaction ← … ← post-transaction-create (hop gRPC + Kafka kontinu)"
  else
    fail "Rantai CHILD_OF terputus — span email tidak menurun dari CreateTransaction/gateway"
    status=1
  fi

  echo "== Spans dalam trace $TRACE_ID =="
  echo "$TRACE_JSON" | jq -r '.data[0] as $t | $t.spans[] | "  " + $t.processes[.processID].serviceName + "  " + .operationName'

  if [ "$status" = 0 ]; then
    pass "trace_id $TRACE_ID kontinu: HTTP → gRPC → Kafka → email"
    exit 0
  fi
  die "Kontinuitas trace GAGAL (trace_id $TRACE_ID)"
}

main "$@"
