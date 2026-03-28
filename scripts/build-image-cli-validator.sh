#!/usr/bin/env bash
set -euo pipefail

CONTAINER_BIN="${CONTAINER_BIN:-docker}"
VALIDATOR_IMAGE="${VALIDATOR_IMAGE:-stability-flow-validator:local}"
PLATFORM="${PLATFORM:-linux/amd64}"

echo "==> Building validator image: ${VALIDATOR_IMAGE}"
echo "==> Platform: ${PLATFORM}"

if [[ "${CONTAINER_BIN}" == "docker" ]]; then
  "${CONTAINER_BIN}" buildx build \
    --platform "${PLATFORM}" \
    --load \
    -f docker/validator/Dockerfile \
    -t "${VALIDATOR_IMAGE}" \
    .
else
  "${CONTAINER_BIN}" build \
    -f docker/validator/Dockerfile \
    -t "${VALIDATOR_IMAGE}" \
    .
fi

echo "==> Validator image built successfully."