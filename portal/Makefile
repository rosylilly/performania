GO?=go

PREFIX?=

GO_ROOT_PACKAGE:=$(shell $(GO) list -m)
GO_ALL_PACKAGES:=$(shell $(GO) list ./... | grep -v /vendor/)
GO_CMD_PACKAGES:=$(shell $(GO) list ./cmd/... | grep -v /vendor/)

GO_FILES:=$(shell find . -name "*.go" -type f -not -path "./vendor/*")

GO_CMD_ARTIFACTS:=$(addprefix $(PREFIX)bin/,$(notdir $(GO_CMD_PACKAGES)))

ifndef BUILD
	BUILD:=debug
endif
ifeq (${BUILD},release)
	ifndef STATIC
		STATIC:=true
	endif
endif

# Go build flags
# symbol table and dwarf
GO_LDFLAGS_SYMBOL:=
ifeq (${BUILD},release)
	GO_LDFLAGS_SYMBOL:=-w -s
endif
# static ldflag
GO_LDFLAGS_STATIC:=
ifeq (${STATIC},true)
	GO_LDFLAGS_STATIC:=-extldflags '-static'
endif
# build ldflags
GO_LDFLAGS:=${GO_LDFLAGS_SYMBOL} ${GO_LDFLAGS_STATIC}
# build tags
GO_BUILD_TAGS:=debug
ifeq (${BUILD},release)
	GO_BUILD_TAGS:=release
endif
# race detector
GO_BUILD_RACE:=-race
ifeq (${BUILD},release)
	GO_BUILD_RACE:=
endif
# static build flag
GO_BUILD_STATIC:=
ifeq (${STATIC},true)
	GO_BUILD_STATIC:=-trimpath
	GO_BUILD_TAGS:=${GO_BUILD_TAGS},netgo
endif
# go build arguments
GO_BUILD_ARGS:=-v $(strip -tags=${GO_BUILD_TAGS} ${GO_BUILD_RACE} ${GO_BUILD_STATIC} -ldflags "$(strip ${GO_LDFLAGS})")

.PHONY: build
build: ## Build all binaries
build: ${GO_CMD_ARTIFACTS}

define build
$(1): ${GO_FILES}
	$(GO) build ${GO_BUILD_ARGS} -o $(1) $(2)
endef

$(foreach cmd,$(GO_CMD_PACKAGES),$(eval $(call build,$(PREFIX)bin/$(notdir $(cmd)),$(cmd))))
