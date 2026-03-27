#!/usr/bin/env bash
set -euo pipefail

CONTAINER_BIN="${CONTAINER_BIN:-docker}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VALIDATOR_IMAGE="${VALIDATOR_IMAGE:-stability-flow-validator:local}"

"${ROOT_DIR}/scripts/build-image-cli-validator.sh"

if [[ $# -eq 0 ]]; then
  echo "Usage: $0 <validator args...>" >&2
  echo "Example: $0 validate-branch-name --branch feat/add-authentication" >&2
  exit 2
fi

exec "${CONTAINER_BIN}" run --rm "${VALIDATOR_IMAGE}" "$@"
