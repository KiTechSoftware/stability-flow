#!/usr/bin/env bash
set -euo pipefail

CONTAINER_BIN="${CONTAINER_BIN:-docker}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VALIDATOR_IMAGE="${VALIDATOR_IMAGE:-stability-flow-validator:local}"

"${ROOT_DIR}/scripts/build-image-cli-validator.sh"

echo "==> Testing validator image..."
"${CONTAINER_BIN}" run --rm "${VALIDATOR_IMAGE}" validate-branch-name --branch feat/add-authentication
"${CONTAINER_BIN}" run --rm "${VALIDATOR_IMAGE}" validate-origin --branch hotfix/1.2.4 --base main
"${CONTAINER_BIN}" run --rm "${VALIDATOR_IMAGE}" validate-merge --source release/1.2.3 --target main --format json
"${CONTAINER_BIN}" run --rm "${VALIDATOR_IMAGE}" validate-commit --mode squash --message "feat: complete validator v1"

echo "==> Testing validator negative case..."
if "${CONTAINER_BIN}" run --rm "${VALIDATOR_IMAGE}" validate-merge --source feat/add-authentication --target main; then
  echo "[FAIL] Expected invalid merge to fail" >&2
  exit 1
fi
