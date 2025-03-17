GOVER   := $(shell go env GOVERSION)
MODULE  := $(shell sed -n '1s/module //p' go.mod)
PACKAGE := $(MODULE)/cmd/app

# Для виндовс содержить расширение .exe или пустую строку для других систем
EXE := $(if $(filter windows,$(shell go env GOOS)),.exe)

SMARTIMPORTS := bin/smartimports_$(GOVER)$(EXE)
LINTVER := v1.64.4
LINTBIN := bin/lint_$(GOVER)_$(LINTVER)$(EXE)

all: format build test lint

build:
	GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build -o bin/app $(PACKAGE)

test:
	go test ./...

run:
	go run $(PACKAGE)

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
	mv bin/$$(notdir $$(firstword $$(subst @, ,$(1))))$(EXE) $$@

BIN_TARGETS += $$(lastword $$(subst ;, ,$(1)))
endef

$(foreach dep,$(BIN_DEPS),$(eval $(call GEN_BIN_RULE,$(dep))))

.bin-targets:
	@echo $(BIN_TARGETS)

.bin-deps: $(BIN_TARGETS)
	@touch $@

.PHONY: clean-bin-deps
clean-bin-deps:
	-rm -fr $(BIN_TARGETS) .bin-deps

generate:
	:

.PHONY: clean
clean: clean-bin-deps
	-rm -fr bin
