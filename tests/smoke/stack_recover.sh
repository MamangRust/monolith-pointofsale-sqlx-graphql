#!/usr/bin/env bash
# Recovery stack POS setelah HOST REBOOT (2026-08-08, §14.4 SUPER_PLANNING.md)
# Akibat reboot:
#   - binary app /tmp/pos-bin/* hilang (tmpfs ter-clear)
#   - postgres/kafka/jaeger exited (255), restart policy = no
#   - redis_* + otel-collector auto-start (restart policy)
#   - volume postgres/kafka/redis INTACT -> data aman
# Jalankan: ./tests/smoke/stack_recover.sh
set -u
cd "$(dirname "$0")/../.." || exit 1

echo "=== 1. START INFRA (postgres kafka jaeger) ==="
docker compose -f deployments/local/docker-compose.yml up -d postgres kafka jaeger 2>&1 | tail -3

echo "=== WAIT POSTGRES ==="
ok=0
for i in $(seq 1 60); do
  if docker exec postgres pg_isready -U DRAGON -d POINTOFSALE >/dev/null 2>&1; then
    echo "PG_READY after ${i}s"; ok=1; break
  fi
  sleep 1
done
[ "$ok" = 1 ] || { echo "PG_NOT_READY"; exit 1; }

echo "=== WAIT KAFKA (KRaft cold start) ==="
ok=0
for i in $(seq 1 360); do
  st=$(docker inspect my-kafka-pointofsale --format '{{.State.Health.Status}}' 2>/dev/null)
  if [ "$st" = "healthy" ]; then
    echo "KAFKA_HEALTHY after ${i}s"; ok=1; break
  fi
  sleep 1
done
[ "$ok" = 1 ] || { echo "KAFKA_NOT_HEALTHY"; exit 1; }

echo "=== 2. BUILD BINARIES ==="
mkdir -p /tmp/pos-bin /tmp/pos-logs
fail=0
for s in apigateway auth user role category merchant cashier product order order_item transaction email; do
  if (cd "service/$s" && go build -o "/tmp/pos-bin/$s" ./cmd/ 2>/tmp/build_err.txt); then
    echo "$s BUILD_OK"
  else
    echo "$s BUILD_FAIL"; head -3 /tmp/build_err.txt; fail=1
  fi
done
[ "$fail" = 0 ] || exit 1

echo "=== 3. START SERVICES ==="
set -a; . ./.env.local; set +a
for s in apigateway auth user role category merchant cashier product order order_item transaction email; do
  setsid -f "/tmp/pos-bin/$s" >> "/tmp/pos-logs/$s.log" 2>&1 < /dev/null
done
sleep 8

echo "=== 4. VERIFY PROCS ==="
up=0; down=0
for s in apigateway auth user role category merchant cashier product order order_item transaction email; do
  if pgrep -f "^/tmp/pos-bin/$s$" >/dev/null; then echo "$s UP"; up=$((up+1)); else echo "$s DOWN"; down=$((down+1)); fi
done
echo "UP=$up DOWN=$down"

echo "=== 5. QUICK E2E (auth) ==="
UUID=$(cat /proc/sys/kernel/random/uuid)
hurl --variable baseUrl=http://localhost:5000 --variable uuid="$UUID" --test tests/hurl/auth.hurl 2>&1 | grep -E '^Success|^Failure' | head -1

echo "=== DONE ==="
