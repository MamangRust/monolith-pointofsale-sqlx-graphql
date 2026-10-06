#!/usr/bin/env bash
#
# stack_lifecycle.sh — Fase 6: Local deployment & data lifecycle
#
# Memvalidasi Docker Compose infra lokal (postgres/kafka/redis×11/otel/jaeger):
#   1. start dari clean volume (down -v → up → migrate → probe data)
#   2. restart  → data probe tetap ada
#   3. stop/start → data probe tetap ada
#   4. down → up (tanpa -v) → data probe tetap ada
#
# DANGER: `docker compose down -v` menghapus volume infra (data dev hilang).
# Hanya jalankan bila data boleh dibuang (env test/dev).
#
# Guard: butuh env STACK_LIFECYCLE_CONFIRM=1 (atau arg --confirm) untuk
# menjalankan `down -v`. Hindari aksi destruktif tak sengaja.
#
# Exit code 0 = PASS, 1 = FAIL.

set -uo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
COMPOSE_FILE="$ROOT/deployments/local/docker-compose.yml"
COMPOSE="docker compose -f $COMPOSE_FILE"

# Hanya service infra (service app berjalan sebagai binary lokal).
INFRA="postgres kafka redis_apigateway redis_auth redis_role redis_user \
redis_category redis_cashier redis_merchant redis_orderitem redis_order \
redis_product redis_transaction otel-collector jaeger"

PGUSER=DRAGON
PGDB=POINTOFSALE
REDIS_PASS=dragon_knight

CONFIRMED=0
[ "${STACK_LIFECYCLE_CONFIRM:-}" = 1 ] && CONFIRMED=1
[ "${1:-}" = "--confirm" ] && CONFIRMED=1
if [ "$CONFIRMED" = 0 ]; then
    echo "⚠  Guard: aksi ini menjalankan 'docker compose down -v' (menghapus volume infra)."
    echo "    Jalankan ulang dengan: STACK_LIFECYCLE_CONFIRM=1 $0"
    echo "    (atau: $0 --confirm)"
    exit 1
fi

PASS=0
FAIL=0
STEP_NUM=0

step() {
    STEP_NUM=$((STEP_NUM + 1))
    echo ""
    echo "== [$STEP_NUM] $1 =="
}

ok()   { PASS=$((PASS + 1)); echo "  ✅ $1"; }
bad()  { FAIL=$((FAIL + 1)); echo "  ❌ $1"; }

preflight() {
    step "Preflight"
    command -v docker >/dev/null || { bad "docker tidak ada"; return 1; }
    docker info >/dev/null 2>&1 || { bad "docker daemon tidak jalan"; return 1; }
    $COMPOSE config --quiet && ok "compose config valid" || { bad "compose config invalid"; return 1; }
    return 0
}

wait_healthy() {
    # $1 = container name; $2 = label; $3 = timeout (s, default 120).
    local name="$1" label="$2" timeout="${3:-120}" i
    local tries=$((timeout / 2))
    for i in $(seq 1 "$tries"); do
        local st
        st=$($COMPOSE ps -q "$name" 2>/dev/null | head -1)
        if [ -n "$st" ]; then
            local health
            health=$(docker inspect --format='{{.State.Health.Status}}' "$st" 2>/dev/null)
            if [ "$health" = "healthy" ]; then
                return 0
            fi
        fi
        sleep 2
    done
    echo "  ⏳ $label belum healthy setelah ${timeout}s"
    return 1
}

wait_tcp() {
    # $1 = port; $2 = label. Poll koneksi TCP.
    local port="$1" label="$2" i
    for i in $(seq 1 60); do
        if (exec 3<>"/dev/tcp/localhost/$port") 2>/dev/null; then
            exec 3>&- 3<&- 2>/dev/null
            return 0
        fi
        sleep 2
    done
    echo "  ⏳ $label belum listening setelah 120s"
    return 1
}

compose_quiet() {
    # Jalankan compose, tampilkan stderr hanya saat gagal (agar error terlihat).
    local out
    out=$("$@" 2>&1) || { echo "$out" | tail -5 | sed 's/^/    /'; return 1; }
    return 0
}

infra_up() {
    compose_quiet $COMPOSE up -d $INFRA
}

infra_down_v() {
    # Hapus volume infra sekaligus (clean volume).
    compose_quiet $COMPOSE down -v $INFRA
}

infra_down() {
    compose_quiet $COMPOSE down $INFRA
}

infra_restart() {
    compose_quiet $COMPOSE restart $INFRA
}

infra_stop() {
    compose_quiet $COMPOSE stop $INFRA
}

infra_start() {
    compose_quiet $COMPOSE start $INFRA
}

wait_infra_ready() {
    local okall=1
    wait_healthy postgres "postgres" || okall=0
    # Kafka KRaft cold start dari clean volume bisa >2 menit (format log dirs + boot).
    wait_healthy kafka "kafka" 360 || okall=0
    for p in 6379 6380 6381 6382 6383 6384 6385 6386 6387 6388 6389; do
        wait_tcp "$p" "redis:$p" || okall=0
    done
    wait_tcp 4317 "otel-collector" || okall=0
    wait_tcp 16686 "jaeger" || okall=0
    [ "$okall" = 1 ]
}

pg() {
    # -q quiet: RETURNING hanya menghasilkan baris data (tanpa command tag).
    docker exec postgres psql -q -U "$PGUSER" -d "$PGDB" -tAc "$1" 2>&1
}

# ---- probe data -------------------------------------------------------------

PROBE_ROW_ID=""
PROBE_REDIS=""

seed_probe() {
    # PostgreSQL: tabel probe + 1 row; Redis: key probe.
    pg "CREATE TABLE IF NOT EXISTS lifecycle_probe (id serial PRIMARY KEY, val text NOT NULL, created_at timestamptz DEFAULT now())" >/dev/null
    PROBE_ROW_ID=$(pg "INSERT INTO lifecycle_probe (val) VALUES ('fase6') RETURNING id" | tr -d ' \r')
    docker exec redis_apigateway redis-cli -a "$REDIS_PASS" SET lifecycle_probe fase6 >/dev/null 2>&1
    PROBE_REDIS=$(docker exec redis_apigateway redis-cli -a "$REDIS_PASS" GET lifecycle_probe 2>/dev/null | tr -d '\r')
    [ -n "$PROBE_ROW_ID" ] && [ "$PROBE_REDIS" = "fase6" ]
}

seed_roles() {
    # Seeder penuh butuh ~5 menit (10 entity × 30s delay); untuk lifecycle
    # cukup seed role yang dibutuhkan e2e: ROLE_ADMIN (auth register) dan
    # 'Admin Access 1' (user create, defaultRoleName di userCommandService).
    local n
    n=$(pg "SELECT count(*) FROM roles WHERE role_name IN ('ROLE_ADMIN','Admin Access 1')" | tr -d ' \r')
    if [ "$n" = 0 ]; then
        pg "INSERT INTO roles (role_name) VALUES ('ROLE_ADMIN'),('Admin Access 1')" >/dev/null
    fi
    [ "$(pg "SELECT count(*) FROM roles WHERE role_name='ROLE_ADMIN'" | tr -d ' \r')" = 1 ] && \
        [ "$(pg "SELECT count(*) FROM roles WHERE role_name='Admin Access 1'" | tr -d ' \r')" = 1 ]
}

verify_roles() {
    local n
    n=$(pg "SELECT count(*) FROM roles WHERE role_name IN ('ROLE_ADMIN','Admin Access 1')" | tr -d ' \r')
    if [ "$n" = 2 ]; then
        ok "role ROLE_ADMIN + Admin Access 1 tetap ada"
        return 0
    fi
    bad "role hilang! (ada $n dari 2)"
    return 1
}

verify_probe() {
    local row
    row=$(pg "SELECT val FROM lifecycle_probe WHERE id = $PROBE_ROW_ID" | tr -d ' \r')
    local rp
    rp=$(docker exec redis_apigateway redis-cli -a "$REDIS_PASS" GET lifecycle_probe 2>/dev/null | tr -d '\r')
    if [ "$row" = "fase6" ] && [ "$rp" = "fase6" ]; then
        ok "data tetap ada (pg row $PROBE_ROW_ID + redis key)"
        return 0
    fi
    bad "data hilang! pg=[$row] redis=[$rp]"
    return 1
}

run_migrate() {
    # Migrate dari root (pakai .env → DB_HOST=localhost). Butuh go.
    (cd "$ROOT" && go run service/migrate/main.go up >/dev/null 2>&1)
}

# ---- main -------------------------------------------------------------------

echo "=== Fase 6: stack lifecycle (infra compose) ==="
echo "Compose file : $COMPOSE_FILE"
echo "⚠  Script ini menjalankan 'down -v' (menghapus volume infra)."

preflight || { echo "TOTAL: PASS=$PASS FAIL=$FAIL"; exit 1; }

step "Clean volume start"
infra_down_v
infra_up || { bad "compose up gagal"; exit 1; }
wait_infra_ready || { bad "infra tidak siap setelah clean start"; exit 1; }
ok "infra start dari clean volume & healthy"

step "Migrate"
if run_migrate; then
    ok "migrate up sukses"
else
    bad "migrate gagal — cek log"
    exit 1
fi

step "Seed probe data (pg + redis) + role dasar"
if seed_probe; then
    ok "probe: pg row id=$PROBE_ROW_ID, redis=$PROBE_REDIS"
else
    bad "probe gagal di-seed"
    exit 1
fi
if seed_roles; then
    ok "role ROLE_ADMIN/ROLE_USER ter-seed"
else
    bad "seed role gagal"
    exit 1
fi

step "Restart (tanpa kehilangan data)"
infra_restart
wait_infra_ready || { bad "infra tidak siap setelah restart"; exit 1; }
verify_probe
verify_roles

step "Stop → Start"
infra_stop
sleep 3
infra_start
wait_infra_ready || { bad "infra tidak siap setelah stop/start"; exit 1; }
verify_probe
verify_roles

step "Down → Up (tanpa -v)"
infra_down
infra_up
wait_infra_ready || { bad "infra tidak siap setelah down/up"; exit 1; }
verify_probe
verify_roles

echo ""
echo "=== HASIL ==="
echo "PASS=$PASS FAIL=$FAIL"
[ "$FAIL" -eq 0 ] && echo "✅ LIFECYCLE TEST PASS" || echo "❌ LIFECYCLE TEST FAIL"
exit $([ "$FAIL" -eq 0 ] && echo 0 || echo 1)
