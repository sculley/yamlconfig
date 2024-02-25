OS = $(shell uname | tr A-Z a-z)
export PATH := $(abspath bin/):${PATH}

# Build variables
BUILD_DIR ?= build
export CGO_ENABLED ?= 0
export GOOS = $(shell go env GOOS)
ifeq (${VERBOSE}, 1)
ifeq ($(filter -v,${GOARGS}),)
	GOARGS += -v
endif
TEST_FORMAT = short-verbose
endif

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt:
	golangci-lint run --fix

.PHONE: test
test:
	go test ./... -v --coverprofile=./coverage.html

.PHONE: coverage
coverage:
	go tool cover -html=./coverage.html