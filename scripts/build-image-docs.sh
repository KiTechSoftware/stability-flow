#!/usr/bin/env bash
set -euo pipefail

CONTAINER_BIN="${CONTAINER_BIN:-docker}"
DOCS_IMAGE="${DOCS_IMAGE:-stability-flow-docs:local}"

echo "==> Building docs image: ${DOCS_IMAGE}"

if [[ "${CONTAINER_BIN}" == "docker" ]]; then
  "${CONTAINER_BIN}" build --load -f docker/docs/Dockerfile -t "${DOCS_IMAGE}" .
else
  "${CONTAINER_BIN}" build -f docker/docs/Dockerfile -t "${DOCS_IMAGE}" .
fi

echo "==> Docs image built successfully."
