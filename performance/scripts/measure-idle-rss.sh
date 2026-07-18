#!/usr/bin/env bash
# Measure idle RSS (MB) of gateforge-iam-server after warm connect.
# Usage: PERF_PID=<pid> ./scripts/measure-idle-rss.sh
set -euo pipefail

PERF_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BASE_URL="${PERF_BASE_URL:-http://127.0.0.1:3000}"
SAMPLES="${PERF_RSS_SAMPLES:-10}"
IDLE_SECS="${PERF_RSS_IDLE_SECS:-30}"
OUT="${PERF_RSS_OUT:-$PERF_DIR/.data/idle-rss.txt}"

mkdir -p "$(dirname "$OUT")"

pid="${PERF_PID:-}"
if [[ -z "${pid}" ]]; then
  pid="$(ps aux | awk '/bin\/gateforge-iam-server/ && !/awk|zsh|measure-idle/ {print $2; exit}')"
fi
if [[ -z "${pid}" ]]; then
  echo "no gateforge-iam-server process; pass PERF_PID=<pid>" >&2
  exit 1
fi

echo "warming ${IDLE_SECS}s (pid=${pid})..."
curl -fsS "${BASE_URL}/api/v1/" >/dev/null || true
sleep "${IDLE_SECS}"

tmp="$(mktemp)"
for _ in $(seq 1 "${SAMPLES}"); do
  ps -o rss= -p "${pid}" | tr -d ' ' >>"${tmp}"
  sleep 1
done

median_kb="$(sort -n "${tmp}" | awk -v n="${SAMPLES}" '{
  a[NR]=$1
}
END {
  if (n % 2) print a[int(n/2)+1]
  else print int((a[n/2]+a[n/2+1])/2)
}')"
rm -f "${tmp}"
median_mb="$(awk -v kb="${median_kb}" 'BEGIN { printf "%.1f", kb/1024 }')"

{
  echo "pid=${pid}"
  echo "median_kb=${median_kb}"
  echo "median_mb=${median_mb}"
  echo "measured_at=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
} | tee "${OUT}"

echo "Idle RSS median: ${median_mb} MB (app process only)"
