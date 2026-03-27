SHELL := /usr/bin/env bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help

.PHONY: help \
	build-image-docs \
	build-image-validator \
	run-image-validator \
	run-image-docs \
	test-image-docs \
	test-image-validator \
	test-images \
	test-flow

help:
	@printf '%s\n' \
		'Available targets:' \
		'  build-image-docs' \
		'  build-image-validator' \
		'  run-image-validator ARGS="..."' \
		'  run-image-docs' \
		'  test-image-docs' \
		'  test-image-validator' \
		'  test-images' \
		'  test-flow'

build-image-docs:
	./scripts/build-image-docs.sh

build-image-validator:
	./scripts/build-image-cli-validator.sh

run-image-validator:
	./scripts/run-image-cli-validator.sh $(ARGS)

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
