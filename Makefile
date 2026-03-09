BIN_DIR := ./bin
DEFAULT_SEEDS_DIR := ./seeds

.PHONY: migrator migrate-up migrate-down seeder seed-up seed-up-default seed-down seed-down-default

migrator:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/migrator ./cmd/migrator

migrate-up: migrator
	$(BIN_DIR)/migrator --up

migrate-down: migrator
	$(BIN_DIR)/migrator --down

seeder:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/seeder ./cmd/seeder

seed-up: seeder
	@if [ -z "$(FILE)" ]; then \
		echo "usage: make seed-up FILE=<filepath>"; \
		exit 1; \
	fi
	$(BIN_DIR)/seeder --up --file "$(FILE)"

seed-up-default: seeder
	./scripts/seed_default.sh up "$(BIN_DIR)/seeder" "$(DEFAULT_SEEDS_DIR)"

seed-down: seeder
	@if [ -z "$(FILE)" ]; then \
		echo "usage: make seed-down FILE=<filepath>"; \
		exit 1; \
	fi
	$(BIN_DIR)/seeder --down --file "$(FILE)"

seed-down-default: seeder
	./scripts/seed_default.sh down "$(BIN_DIR)/seeder" "$(DEFAULT_SEEDS_DIR)"


GOLANGCI_LINT_VERSION ?= 2.11.2
GOLANGCI_LINT_BIN ?= $(shell go env GOPATH)/bin/golangci-lint
GOLANGCI_LINT_CONFIG ?= .golangci.yml
GO_CMD = go
GO_MOD = $(GO_CMD) mod

# Цвета
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
BLUE := \033[0;34m
NC := \033[0m

.PHONY: go-tools, lint, lint-fix, version

go-tools:
	@printf "$(BLUE)╔════════════════════════════════════╗$(NC)\n"
	@printf "$(BLUE)║    Установка инструментов Go       ║$(NC)\n"
	@printf "$(BLUE)╚════════════════════════════════════╝$(NC)\n"
	
	@printf "$(YELLOW)📦 Проверка golangci-lint...$(NC)\n"
	@if ! command -v $(GOLANGCI_LINT_BIN) >/dev/null 2>&1; then \
		printf "   Установка golangci-lint $(GOLANGCI_LINT_VERSION)...\n"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
		printf "$(GREEN)   ✓ golangci-lint установлен$(NC)\n"; \
	else \
		printf "$(GREEN)   ✓ golangci-lint уже установлен$(NC)\n"; \
	fi
	
	@printf "$(YELLOW)📥 Управление зависимостями...$(NC)\n"
	@$(GO_MOD) download && \
	 printf "$(GREEN)   ✓ Зависимости скачаны$(NC)\n" || \
	 printf "$(RED)   ✗ Ошибка при скачивании зависимостей$(NC)\n"
	
	@$(GO_MOD) verify && \
	 printf "$(GREEN)   ✓ Зависимости проверены$(NC)\n" || \
	 printf "$(RED)   ✗ Ошибка при проверке зависимостей$(NC)\n"
	
	@printf "$(GREEN)✅ Инструменты готовы к работе$(NC)\n"

lint: go-tools
	@printf "$(BLUE)╔════════════════════════════════════╗$(NC)\n"
	@printf "$(BLUE)║         Запуск линтеров            ║$(NC)\n"
	@printf "$(BLUE)╚════════════════════════════════════╝$(NC)\n"
	
	@if [ -f "$(GOLANGCI_LINT_CONFIG)" ]; then \
		printf "$(YELLOW)🔍 Использую конфиг: $(GOLANGCI_LINT_CONFIG)$(NC)\n"; \
		$(GOLANGCI_LINT_BIN) run --config=$(GOLANGCI_LINT_CONFIG) --timeout=5m ./...; \
	else \
		printf "$(YELLOW)🔍 Конфиг не найден, использую стандартные настройки$(NC)\n"; \
		$(GOLANGCI_LINT_BIN) run --timeout=5m ./...; \
	fi; \
	if [ $$? -eq 0 ]; then \
		printf "$(GREEN)✅ Линтинг прошел успешно!$(NC)\n"; \
	else \
		printf "$(RED)❌ Найдены проблемы в коде$(NC)\n"; \
		printf "$(YELLOW)   Попробуйте 'make lint-fix' для автоисправления$(NC)\n"; \
	fi

lint-fix: go-tools
	@printf "$(BLUE)╔════════════════════════════════════╗$(NC)\n"
	@printf "$(BLUE)║     Автоисправление линтеров       ║$(NC)\n"
	@printf "$(BLUE)╚════════════════════════════════════╝$(NC)\n"
	
	@printf "$(YELLOW)🔧 Запуск автоисправления...$(NC)\n"
	@$(GOLANGCI_LINT_BIN) run --fix --timeout=5m ./... || true
	
	@if [ $$? -eq 0 ]; then \
		printf "$(GREEN)✅ Автоисправление завершено$(NC)\n"; \
	else \
		printf "$(RED)❌ Ошибка при автоисправлении$(NC)\n"; \
	fi

# Версия
version:
	@printf "$(YELLOW)Инструменты:$(NC)\n"
	@printf "  Go: $$(go version)\n"
	@printf "  golangci-lint: $$($(GOLANGCI_LINT_BIN) version 2>/dev/null || echo 'не установлен')\n"