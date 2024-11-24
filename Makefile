
LDFLAGS=

BUILD_ENVPARMS:=CGO_ENABLED=0

LOCAL_BIN:=$(CURDIR)/bin

##################### GOX #####################
GOX_BIN:=$(LOCAL_BIN)/gox

# local gox
ifeq ($(wildcard $(GOX_BIN)),)
GOX_BIN:=
endif

# Check global bin version
ifneq (, $(shell which gox))
GOX_BIN:=$(shell which gox)
endif

##################### GOLANG-CI RELATED CHECKS #####################
# Check global GOLANGCI-LINT
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint ## local linter binary path
GOLANGCI_TAG:=1.61.0 ## linter version to use
GOLANGCI_LINTER_IMAGE:="golangci/golangci-lint" ## pipeline linter image to use in ci-lint target

# Check local bin version
ifneq ($(wildcard $(GOLANGCI_BIN)),)
GOLANGCI_BIN_VERSION:=$(shell $(GOLANGCI_BIN) --version)
ifneq ($(GOLANGCI_BIN_VERSION),)
GOLANGCI_BIN_VERSION_SHORT:=$(shell echo "$(GOLANGCI_BIN_VERSION)" | sed -E 's/.* version (.*) built from .* on .*/\1/g')
else
GOLANGCI_BIN_VERSION_SHORT:=0
endif
ifneq "$(GOLANGCI_TAG)" "$(word 1, $(sort $(GOLANGCI_TAG) $(GOLANGCI_BIN_VERSION_SHORT)))"
GOLANGCI_BIN:=
endif
endif

# Check global bin version
ifneq (, $(shell which golangci-lint))
GOLANGCI_VERSION:=$(shell golangci-lint --version 2> /dev/null )
ifneq ($(GOLANGCI_VERSION),)
GOLANGCI_VERSION_SHORT:=$(shell echo "$(GOLANGCI_VERSION)"|sed -E 's/.* version (.*) built from .* on .*/\1/g')
else
GOLANGCI_VERSION_SHORT:=0
endif
ifeq "$(GOLANGCI_TAG)" "$(word 1, $(sort $(GOLANGCI_TAG) $(GOLANGCI_VERSION_SHORT)))"
GOLANGCI_BIN:=$(shell which golangci-lint)
endif
endif

.PHONY: install-lint
install-lint: ## install golangci-lint binary
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	$(info Downloading golangci-lint v$(GOLANGCI_TAG))
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_TAG)
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
endif

.PHONY: .lint-full
.lint-full: install-lint
	$(GOLANGCI_BIN) run --config=.golangci.yml ./...

.PHONY: lint-full
lint-full: .lint-full

.PHONY: .bin-deps
.bin-deps:
	mkdir -p bin
	$(info Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install github.com/mitchellh/gox@v1.0.1  && \
	GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/goimports@v0.1.9 && \

.PHONY: .deps
.deps:
	$(info Install dependencies...)
	go mod download

.PHONY: update-deps
update-deps: .deps .bin-deps

.PHONY: all
all: test build ## default scratch target: test and build

.PHONY: .test
.test:
	$(info Running tests...)
	go test ./...

.PHONY: test
test: .test ## run unit tests

# CMD_LIST список таргетов (через пробел) которые надо собрать
# можно переопределить в Makefile, по дефолту все из ./cmd кроме основного пакета
# пример переопределения CMD_LIST:= ./cmd/example ./cmd/app ./cmd/cron
ifndef CMD_LIST
CMD_LIST:=$(shell ls ./cmd | sed -e 's/^/.\/cmd\//')
endif
# определение текущий ос
ifndef HOSTOS
HOSTOS:=$(shell go env GOHOSTOS)
endif
# определение текущий архитектуры
ifndef HOSTARCH
HOSTARCH:=$(shell go env GOHOSTARCH)
endif

ifndef BIN_DIR
BIN_DIR=./bin
endif

# если нужно собрать только основной сервис, можно указать в Makefile SINGLE_BUILD=1
DISABLE_CMD_LIST_BUILD?=0

.PHONY: .build
.build:
# сначала собирается основной сервис, скачиваются нужные пакеты и все кладется в кеш для дальнейшего использования
	$(info Building...)
	@if [ -n "$(CMD_LIST)" ] && [ "$(DISABLE_CMD_LIST_BUILD)" != 1 ]; then\
		$(BUILD_ENVPARMS) $(GOX_BIN) -output="$(BIN_DIR)" -osarch="$(HOSTOS)/$(HOSTARCH)" -ldflags "$(LDFLAGS)" $(CMD_LIST);\
	fi

.PHONY: build
build: .build ## build project
