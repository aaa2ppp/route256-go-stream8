SHELL := bash -euo pipefail
.DELETE_ON_ERROR:


GOVER   := $(shell go env GOVERSION)
MODULE  := $(shell sed -n '1s/module //p' go.mod)
PACKAGE := $(MODULE)/cmd/app

# Для виндовс содержить расширение .exe или пустую строку для других систем
EXE := $(if $(filter windows,$(shell go env GOOS)),.exe)

SMARTIMPORTS := bin/smartimports_$(GOVER)$(EXE)
LINTVER := v1.64.4
LINTBIN := bin/lint_$(GOVER)_$(LINTVER)$(EXE)

all: format build test lint

build: generate sqlc-generate tidy
	GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build -o bin/app $(PACKAGE)

sqlc-generate:
	sqlc generate

test:
	go test ./...

run:
	go run $(PACKAGE)

tidy:
	go mod tidy

lint: install-lint
	$(LINTBIN) run

precommit: format build test lint
	@echo "OK"

format: install-smartimports
	$(SMARTIMPORTS) -exclude internal/mocks

install-lint: $(LINTBIN)

install-smartimports: $(SMARTIMPORTS)

# Инсталяция бинарных зависимостей

# Интсталируем в локальный bin
GO_INSTALL := GOBIN=$(PWD)/bin go install

# Зависимости которые нужно устновить
BIN_DEPS := \
	github.com/golangci/golangci-lint/cmd/golangci-lint@$(LINTVER);$(LINTBIN) \
	github.com/pav5000/smartimports/cmd/smartimports@latest;$(SMARTIMPORTS)

define GEN_BIN_RULE
$$(lastword $$(subst ;, ,$(1))):
	$(GO_INSTALL) $$(firstword $$(subst ;, ,$(1)))
	$$(eval BIN := bin/$$(notdir $$(firstword $$(subst @, ,$(1))))$(EXE))
	@test -f $$(BIN) || { echo "$$(BIN) not found"; exit 1; }
	mv $$(BIN) $$@

BIN_TARGETS += $$(lastword $$(subst ;, ,$(1)))
endef

$(foreach dep,$(BIN_DEPS),$(eval $(call GEN_BIN_RULE,$(dep))))

.bin-targets:
	@echo $(BIN_TARGETS)

.bin-deps: $(BIN_TARGETS)
	touch $@

.PHONY: clean-bin-deps
clean-bin-deps:
	-rm -fr $(BIN_TARGETS) .bin-deps


# Инсталяция плагинов protoc

# Список плагинов protoc, их путей и версий
PROTOC_DEPS := \
	google.golang.org/protobuf/cmd/protoc-gen-go@latest \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest \
	github.com/envoyproxy/protoc-gen-validate@latest


define GEN_PROTOC_RULE
bin/$$(notdir $$(firstword $$(subst @, ,$(1))))$(EXE):
	$(GO_INSTALL) $(1)
	$$(eval BIN := bin/$$(notdir $$(firstword $$(subst @, ,$(1))))$(EXE))
	@test -f $$(BIN) || { echo "$$(BIN) not found"; exit 1; }

PROTOC_TARGETS += bin/$$(notdir $$(firstword $$(subst @, ,$(1))))$(EXE)
endef

$(foreach dep,$(PROTOC_DEPS),$(eval $(call GEN_PROTOC_RULE,$(dep))))

.protoc-targets:
	@echo $(PROTOC_TARGETS)

.protoc-deps: $(PROTOC_TARGETS)
	touch $@

.PHONY: clean-protoc-deps
clean-protoc-deps:
	-rm -fr $(PROTOC_TARGETS) .protoc-deps


# Вендоринг внешних proto файлов

VENDOR_DIR = vendor-proto
GIT_CLONE = git clone -q --single-branch -n --filter=tree:0
GIT_CHECKOUT = git sparse-checkout set --no-cone

# Список зависимостей для вендоринга
VENDOR_DEPS = \
    google/api;master;https://github.com/googleapis/googleapis;google/api \
    protoc-gen-openapiv2/options;main;https://github.com/grpc-ecosystem/grpc-gateway;protoc-gen-openapiv2/options \
    google/protobuf;main;https://github.com/protocolbuffers/protobuf;src/google/protobuf \
    validate;main;https://github.com/bufbuild/protoc-gen-validate;validate

define GEN_VENDOR_RULE
$(VENDOR_DIR)/$$(firstword $$(subst ;, ,$(1))):
	$$(eval TEMP_DIR := $$(shell mkdir -p $(VENDOR_DIR) && mktemp -d $(VENDOR_DIR)/tmp.XXXXXXXX))
	$$(eval PARAMS   := $$(wordlist 2, 999, $(subst ;, ,$(1))))
	$$(eval BRANCH   := $$(word 1,$$(PARAMS)))
	$$(eval REPO_URL := $$(word 2,$$(PARAMS)))
	$$(eval SRC_PATH := $$(word 3,$$(PARAMS)))

	$(GIT_CLONE) -b $$(BRANCH) --depth=1 $$(REPO_URL) $$(TEMP_DIR)
	cd $$(TEMP_DIR) && $(GIT_CHECKOUT) $$(SRC_PATH) && git checkout -q
	mkdir -p $$(@D) && mv $$(TEMP_DIR)/$$(SRC_PATH) $$@
	rm -fr $$(TEMP_DIR)

VENDOR_TARGETS += $(VENDOR_DIR)/$$(firstword $$(subst ;, ,$(1)))
endef

$(foreach dep,$(VENDOR_DEPS),$(eval $(call GEN_VENDOR_RULE,$(dep))))

.vendor-targets:
	@echo $(VENDOR_TARGETS)

.vendor-deps: $(VENDOR_TARGETS)
	touch .vendor-deps

.PHONY: clean-vendor-deps
clean-vendor-deps:
	-rm -fr $(VENDOR_TARGETS) .vendor-deps


# Генерация сервисов по их ptoto описаниям

# Список proto описаний сервисов
SERVICE_PROTOS := \
	api/order/v1/order.proto \
	api/stock/v1/stock.proto

define GEN_SERVICE_RULE
pkg/$$(dir $(1)): $(1) .protoc-deps .vendor-deps
	mkdir -p $$@
	-rm -fr $$@/*
	protoc -I $$(dir $$<) -I vendor-proto \
		--plugin=bin/protoc-gen-go$$(EXE) --go_out $$@ --go_opt paths=source_relative \
		--plugin=bin/protoc-gen-go-grpc$$(EXE) --go-grpc_out $$@ --go-grpc_opt paths=source_relative \
		$$<

SERVICE_TARGETS += pkg/$$(dir $(1))
endef

$(foreach service,$(SERVICE_PROTOS),$(eval $(call GEN_SERVICE_RULE,$(service))))

.service-targets:
	@echo $(SERVICE_TARGETS)

.PHONY: generate
generate: .generate

.PHONY: clean-generate
.generate: $(SERVICE_TARGETS) .protoc-deps .vendor-deps
	touch $@
clean-generate:
	-rm -fr $(SERVICE_TARGETS) .generate


# clean

.PHONY: clean
clean: clean-bin-deps clean-protoc-deps clean-vendor-deps clean-generate
	-rm -fr bin
