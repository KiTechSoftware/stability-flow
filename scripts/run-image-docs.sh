#!/usr/bin/env bash
set -euo pipefail

CONTAINER_BIN="${CONTAINER_BIN:-docker}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

DOCS_IMAGE="${DOCS_IMAGE:-stability-flow-docs:local}"
DOCS_PORT="${DOCS_PORT:-8080}"

"${ROOT_DIR}/scripts/build-image-docs.sh"

echo
echo "==> Running docs image on http://localhost:${DOCS_PORT}"
exec "${CONTAINER_BIN}" run --rm -p "${DOCS_PORT}:80" "${DOCS_IMAGE}"
