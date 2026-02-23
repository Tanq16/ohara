.PHONY: help assets verify-assets clean build build-for build-all run docker-build docker-push version

# =============================================================================
# Variables
# =============================================================================
APP_NAME := ohara
DOCKER_USER := tanq16

# Build variables (set by CI or use defaults)
VERSION ?= dev-build
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Asset versions
HLJS_VERSION := 11.11.1
CHARTJS_VERSION := 4.5.1
MARKED_VERSION := 17.0.3
MARKED_HL_VERSION := 2.2.3
MERMAID_VERSION := 11.12.3

# Directories
STATIC_DIR := internal/server/static
JS_DIR := $(STATIC_DIR)/js
CSS_DIR := $(STATIC_DIR)/css
FONTS_DIR := $(STATIC_DIR)/fonts

# Console colors
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m

# =============================================================================
# Help
# =============================================================================
help: ## Show this help
	@echo "$(CYAN)Available targets:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

# =============================================================================
# Assets
# =============================================================================
assets: ## Download static assets
	@echo "$(CYAN)Downloading assets...$(NC)"
	@mkdir -p $(JS_DIR) $(CSS_DIR) $(FONTS_DIR)
	@curl -sL "https://cdn.tailwindcss.com" -o "$(JS_DIR)/tailwindcss.js"
	@curl -sL "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/$(HLJS_VERSION)/highlight.min.js" -o "$(JS_DIR)/highlight.min.js"
	@curl -sL "https://cdn.jsdelivr.net/npm/chart.js@$(CHARTJS_VERSION)/dist/chart.umd.min.js" -o "$(JS_DIR)/chart.umd.min.js"
	@curl -sL "https://cdn.jsdelivr.net/npm/marked@$(MARKED_VERSION)/lib/marked.umd.js" -o "$(JS_DIR)/marked.umd.js"
	@curl -sL "https://cdn.jsdelivr.net/npm/marked-highlight@$(MARKED_HL_VERSION)/lib/index.umd.js" -o "$(JS_DIR)/marked-highlight.umd.js"
	@curl -sL "https://cdn.jsdelivr.net/npm/mermaid@$(MERMAID_VERSION)/dist/mermaid.min.js" -o "$(JS_DIR)/mermaid.min.js"
	@curl -sL "https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" -H "User-Agent: Mozilla/5.0" -o "$(CSS_DIR)/inter.css"
	@grep -o "https://fonts.gstatic.com/[^)']*" "$(CSS_DIR)/inter.css" | sort -u | while read url; do \
		filename=$$(basename "$$url" | sed 's/?.*//'); \
		curl -sL "$$url" -o "$(FONTS_DIR)/$$filename"; \
	done
	@sed -i.bak -E 's|https://fonts.gstatic.com/s/inter/[^/]+/||g' "$(CSS_DIR)/inter.css" && rm -f "$(CSS_DIR)/inter.css.bak"
	@sed -i.bak 's|src: url(|src: url(/static/fonts/|g' "$(CSS_DIR)/inter.css" && rm -f "$(CSS_DIR)/inter.css.bak"
	@curl -sL "https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500&display=swap" -H "User-Agent: Mozilla/5.0" -o "$(CSS_DIR)/jetbrains-mono.css"
	@grep -o "https://fonts.gstatic.com/[^)']*" "$(CSS_DIR)/jetbrains-mono.css" | sort -u | while read url; do \
		filename=$$(basename "$$url" | sed 's/?.*//'); \
		curl -sL "$$url" -o "$(FONTS_DIR)/$$filename"; \
	done
	@sed -i.bak -E 's|https://fonts.gstatic.com/s/jetbrainsmono/[^/]+/||g' "$(CSS_DIR)/jetbrains-mono.css" && rm -f "$(CSS_DIR)/jetbrains-mono.css.bak"
	@sed -i.bak 's|src: url(|src: url(/static/fonts/|g' "$(CSS_DIR)/jetbrains-mono.css" && rm -f "$(CSS_DIR)/jetbrains-mono.css.bak"
	@echo "$(GREEN)Assets downloaded$(NC)"

verify-assets: ## Verify required assets exist
	@test -f $(JS_DIR)/tailwindcss.js || (echo "$(YELLOW)tailwindcss.js missing. Run 'make assets'$(NC)" && exit 1)
	@test -f $(JS_DIR)/highlight.min.js || (echo "$(YELLOW)highlight.min.js missing. Run 'make assets'$(NC)" && exit 1)
	@test -f $(JS_DIR)/chart.umd.min.js || (echo "$(YELLOW)chart.umd.min.js missing. Run 'make assets'$(NC)" && exit 1)
	@test -f $(JS_DIR)/marked.umd.js || (echo "$(YELLOW)marked.umd.js missing. Run 'make assets'$(NC)" && exit 1)
	@test -f $(JS_DIR)/marked-highlight.umd.js || (echo "$(YELLOW)marked-highlight.umd.js missing. Run 'make assets'$(NC)" && exit 1)
	@test -f $(JS_DIR)/mermaid.min.js || (echo "$(YELLOW)mermaid.min.js missing. Run 'make assets'$(NC)" && exit 1)
	@test -f $(CSS_DIR)/inter.css || (echo "$(YELLOW)inter.css missing. Run 'make assets'$(NC)" && exit 1)
	@test -f $(CSS_DIR)/jetbrains-mono.css || (echo "$(YELLOW)jetbrains-mono.css missing. Run 'make assets'$(NC)" && exit 1)
	@echo "$(GREEN)Assets verified$(NC)"

clean: ## Remove built artifacts and downloaded assets
	@rm -f $(APP_NAME) $(APP_NAME)-*
	@rm -f $(JS_DIR)/tailwindcss.js $(JS_DIR)/highlight.min.js $(JS_DIR)/chart.umd.min.js $(JS_DIR)/marked.umd.js $(JS_DIR)/marked-highlight.umd.js $(JS_DIR)/mermaid.min.js
	@rm -rf $(CSS_DIR) $(FONTS_DIR)
	@echo "$(GREEN)Cleaned$(NC)"

# =============================================================================
# Build
# =============================================================================
build: assets verify-assets ## Build binary for current platform
	@go build -ldflags="-s -w -X 'github.com/tanq16/ohara/cmd.AppVersion=$(VERSION)'" -o $(APP_NAME) .
	@echo "$(GREEN)Built: ./$(APP_NAME)$(NC)"

build-for: verify-assets ## Build binary for specified GOOS/GOARCH
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-s -w -X 'github.com/tanq16/ohara/cmd.AppVersion=$(VERSION)'" -o $(APP_NAME)-$(GOOS)-$(GOARCH) .
	@echo "$(GREEN)Built: ./$(APP_NAME)-$(GOOS)-$(GOARCH)$(NC)"

build-all: assets verify-assets ## Build all platform binaries
	@$(MAKE) build-for GOOS=linux GOARCH=amd64
	@$(MAKE) build-for GOOS=linux GOARCH=arm64
	@$(MAKE) build-for GOOS=darwin GOARCH=amd64
	@$(MAKE) build-for GOOS=darwin GOARCH=arm64

run: ## Run the server locally
	@go run . --data-dir ./data --port 8080

# =============================================================================
# Docker
# =============================================================================
docker-build: ## Build Docker image
	@docker build -t $(DOCKER_USER)/$(APP_NAME):$(VERSION) .
	@docker tag $(DOCKER_USER)/$(APP_NAME):$(VERSION) $(DOCKER_USER)/$(APP_NAME):latest
	@echo "$(GREEN)Docker image built$(NC)"

docker-push: docker-build ## Push Docker image to Docker Hub
	@docker push $(DOCKER_USER)/$(APP_NAME):$(VERSION)
	@docker push $(DOCKER_USER)/$(APP_NAME):latest
	@echo "$(GREEN)Docker image pushed$(NC)"

# =============================================================================
# Version
# =============================================================================
version: ## Calculate next version from commit message
	@LATEST_TAG=$$(git tag --sort=-v:refname | head -n1 || echo "0.0.0"); \
	LATEST_TAG=$${LATEST_TAG#v}; \
	MAJOR=$$(echo "$$LATEST_TAG" | cut -d. -f1); \
	MINOR=$$(echo "$$LATEST_TAG" | cut -d. -f2); \
	PATCH=$$(echo "$$LATEST_TAG" | cut -d. -f3); \
	MAJOR=$${MAJOR:-0}; MINOR=$${MINOR:-0}; PATCH=$${PATCH:-0}; \
	COMMIT_MSG="$$(git log -1 --pretty=%B)"; \
	if echo "$$COMMIT_MSG" | grep -q "\[major-release\]"; then \
		MAJOR=$$((MAJOR + 1)); MINOR=0; PATCH=0; \
	elif echo "$$COMMIT_MSG" | grep -q "\[minor-release\]"; then \
		MINOR=$$((MINOR + 1)); PATCH=0; \
	else \
		PATCH=$$((PATCH + 1)); \
	fi; \
	echo "v$${MAJOR}.$${MINOR}.$${PATCH}"
