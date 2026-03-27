#!/usr/bin/env bash
set -euo pipefail

CONTAINER_BIN="${CONTAINER_BIN:-docker}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

DOCS_IMAGE="${DOCS_IMAGE:-stability-flow-docs:local}"
VALIDATOR_IMAGE="${VALIDATOR_IMAGE:-stability-flow-validator:local}"
DOCS_PORT="${DOCS_PORT:-8080}"

cleanup() {
  if [[ -n "${DOCS_CONTAINER_ID:-}" ]]; then
    "${CONTAINER_BIN}" rm -f "${DOCS_CONTAINER_ID}" >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT

echo "==> Building images using ${CONTAINER_BIN}..."
"${ROOT_DIR}/scripts/build-image-docs.sh"
"${ROOT_DIR}/scripts/build-image-cli-validator.sh"

echo "==> Testing validator image..."
"${CONTAINER_BIN}" run --rm "${VALIDATOR_IMAGE}" validate-branch-name --branch feat/add-authentication
"${CONTAINER_BIN}" run --rm "${VALIDATOR_IMAGE}" validate-origin --branch hotfix/1.2.4 --base main
"${CONTAINER_BIN}" run --rm "${VALIDATOR_IMAGE}" validate-merge --source release/1.2.3 --target main --format json
"${CONTAINER_BIN}" run --rm "${VALIDATOR_IMAGE}" validate-commit --mode squash --message "feat: complete validator v1"

echo "==> Testing docs image..."
DOCS_CONTAINER_ID="$("${CONTAINER_BIN}" run -d -p "${DOCS_PORT}:80" "${DOCS_IMAGE}")"

for _ in {1..20}; do
  if curl -fsSL "http://localhost:${DOCS_PORT}" >/dev/null; then
    echo "==> Docs image responded successfully."
    echo "==> "${CONTAINER_BIN}" images built and tested successfully."
    exit 0
  fi
  sleep 1
done

echo "[FAIL] Docs image did not respond on port ${DOCS_PORT}" >&2
exit 1
