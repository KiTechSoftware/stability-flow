set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

default:
    @just --list

help:
    @just --list

build-image-docs:
    ./scripts/build-image-docs.sh

build-image-validator:
    ./scripts/build-image-cli-validator.sh

run-image-validator *args:
    ./scripts/run-image-cli-validator.sh {{args}}

run-image-docs:
    ./scripts/run-image-docs.sh

test-image-docs:
    ./scripts/test-image-docs.sh

test-image-validator:
    ./scripts/test-image-cli-validator.sh

test-images:
    ./scripts/test-images.sh

test-flow:
    ./scripts/test-stability-flow.sh
