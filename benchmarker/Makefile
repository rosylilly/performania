GO?=go

PREFIX?=

GO_ROOT_PACKAGE:=$(shell $(GO) list -m)
GO_ALL_PACKAGES:=$(shell $(GO) list ./... | grep -v /vendor/)
GO_CMD_PACKAGES:=$(shell $(GO) list ./cmd/... | grep -v /vendor/)

GO_FILES:=$(shell find . -name "*.go" -type f -not -path "./vendor/*")

GO_CMD_ARTIFACTS:=$(addprefix $(PREFIX)bin/,$(notdir $(GO_CMD_PACKAGES)))

.PHONY: build
build: ## Build all binaries
build: ${GO_CMD_ARTIFACTS}

define build
$(1): ${GO_FILES}
	$(GO) build -o $(1) $(2)
endef

$(foreach cmd,$(GO_CMD_PACKAGES),$(eval $(call build,$(PREFIX)bin/$(notdir $(cmd)),$(cmd))))
