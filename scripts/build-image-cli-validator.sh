#!/usr/bin/env bash
set -euo pipefail

CONTAINER_BIN="${CONTAINER_BIN:-docker}"
VALIDATOR_IMAGE="${VALIDATOR_IMAGE:-stability-flow-validator:local}"

echo "==> Building validator image: ${VALIDATOR_IMAGE}"

if [[ "${CONTAINER_BIN}" == "docker" ]]; then
  "${CONTAINER_BIN}" build --load -f docker/validator/Dockerfile -t "${VALIDATOR_IMAGE}" .
else
  "${CONTAINER_BIN}" build -f docker/validator/Dockerfile -t "${VALIDATOR_IMAGE}" .
fi

echo "==> Validator image built successfully."
