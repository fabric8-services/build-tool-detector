PROJECT_NAME=build-tool-detector
PACKAGE_NAME:=github.com/fabric8-services/$(PROJECT_NAME)
CUR_DIR=$(shell pwd)
TMP_PATH=$(CUR_DIR)/tmp
INSTALL_PREFIX=$(CUR_DIR)/bin
VENDOR_DIR=vendor
SOURCE_DIR ?= .
SOURCES := $(shell find $(SOURCE_DIR) -path $(SOURCE_DIR)/vendor -prune -o -name '*.go' -print)
DESIGN_DIR=design
DESIGNS := $(shell find $(SOURCE_DIR)/$(DESIGN_DIR) -path $(SOURCE_DIR)/vendor -prune -o -name '*.go' -print)
GOBIN=$(INSTALL_PREFIX)

# This pattern excludes the listed folders while running tests
TEST_PKGS_EXCLUDE_PATTERN = "vendor|app|tool\/build-tool-detector-cli|design|client|test"

include ./.make/docker.mk
ifeq ($(OS),Windows_NT)
include ./.make/Makefile.win
else
include ./.make/Makefile.lnx
endif

# This is a fix for a non-existing user in passwd file when running in a docker
# container and trying to clone repos of dependencies
GIT_COMMITTER_NAME ?= "user"
GIT_COMMITTER_EMAIL ?= "user@example.com"
export GIT_COMMITTER_NAME
export GIT_COMMITTER_EMAIL
export GOBIN

COMMIT=$(shell git rev-parse HEAD 2>/dev/null)
GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(GITUNTRACKEDCHANGES),)
COMMIT := $(COMMIT)-dirty
endif
BUILD_TIME=`date -u '+%Y-%m-%dT%H:%M:%SZ'`

.DEFAULT_GOAL := help

# Call this function with $(call log-info,"Your message")
define log-info =
@echo "INFO: $(1)"
endef

# -------------------------------------------------------------------
# help!
# -------------------------------------------------------------------

.PHONY: help
help: ## Prints this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

# -------------------------------------------------------------------
# required tools
# -------------------------------------------------------------------

# Find all required tools:
GIT_BIN := $(shell command -v $(GIT_BIN_NAME) 3> /dev/null)

GO_BIN := $(shell command -v $(GO_BIN_NAME) 2> /dev/null)

$(INSTALL_PREFIX):
	mkdir -p $(INSTALL_PREFIX)
$(TMP_PATH):
	mkdir -p $(TMP_PATH)

.PHONY: test-deps
test-deps: $(GINKGO_BIN)

# install ginkgo cli
$(GINKGO_BIN):
	$(GO_BIN) install github.com/onsi/ginkgo/ginkgo
	@chmod +x $(GINKGO_BIN)

# -------------------------------------------------------------------
# support for generating goa code
# -------------------------------------------------------------------
$(GOAGEN_BIN):
	@echo "Building goagen tool"
	$(GO_BIN) install github.com/goadesign/goa/goagen

# -------------------------------------------------------------------
# clean
# -------------------------------------------------------------------

# For the global "clean" target all targets in this variable will be executed
CLEAN_TARGETS =

CLEAN_TARGETS += clean-artifacts
.PHONY: clean-artifacts
## Removes the ./bin directory.
clean-artifacts:
	-rm -rf $(INSTALL_PREFIX)

CLEAN_TARGETS += clean-object-files
.PHONY: clean-object-files
## Runs go clean to remove any executables or other object files.
clean-object-files:
	go clean ./...

CLEAN_TARGETS += clean-generated
.PHONY: clean-generated
## Removes all generated code.
clean-generated:
	-rm -rf ./app
	-rm -rf ./swagger/

CLEAN_TARGETS += clean-vendor
.PHONY: clean-vendor
## Removes the ./vendor directory.
clean-vendor:
	-rm -rf $(VENDOR_DIR)

CLEAN_TARGETS += clean-tmp
.PHONY: clean-tmp
## Removes the ./vendor directory.
clean-tmp:
	-rm -rf $(TMP_DIR)

# Keep this "clean" target here after all `clean-*` sub tasks
.PHONY: clean
clean: $(CLEAN_TARGETS) ## Runs all clean-* targets.

# -------------------------------------------------------------------
# build the binary executable (to ship in prod)
# -------------------------------------------------------------------
LDFLAGS=-ldflags "-X ${PACKAGE_NAME}/app.Commit=${COMMIT} -X ${PACKAGE_NAME}/app.BuildTime=${BUILD_TIME}"

$(SERVER_BIN): ## Build the server
	@echo "building $(SERVER_BIN)..."
	$(GO_BIN) build -v $(LDFLAGS) -o $(SERVER_BIN)

.PHONY: build
build: $(SERVER_BIN) ## Build the server

.PHONY: generate
generate: prebuild-check $(DESIGNS) $(GOAGEN_BIN) ## Generate GOA sources. Only necessary after clean of if changed `design` folder.
	@echo "Generating goa artifacts"
	$(GOAGEN_BIN) app -d ${PACKAGE_NAME}/${DESIGN_DIR}
	$(GOAGEN_BIN) controller -d ${PACKAGE_NAME}/${DESIGN_DIR} -o controllers/ --pkg controllers --app-pkg ${PACKAGE_NAME}/app
	$(GOAGEN_BIN) gen -d ${PACKAGE_NAME}/${DESIGN_DIR} --pkg-path=github.com/fabric8-services/fabric8-common/goasupport/status --out app
	$(GOAGEN_BIN) swagger -d ${PACKAGE_NAME}/${DESIGN_DIR}
	
.PHONY: test 
test: test-deps generate  ## Executes all tests
	$(eval TEST_PACKAGES:=$(shell go list ./... | grep -v -E $(TEST_PKGS_EXCLUDE_PATTERN)))
	go test -vet off $(TEST_PACKAGES) -v

.PHONY: format ## Removes unneeded imports and formats source code
format: 
	@goimports -l -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

GOFORMAT_FILES := $(shell find  . -name '*.go' | grep -vEf .gofmt_exclude)

.PHONY: show-info
show-info:
	$(call log-info,"$(shell go version)")
	$(call log-info,"$(shell go env)")

.PHONY: prebuild-check
prebuild-check: $(TMP_PATH) $(INSTALL_PREFIX) show-info
# Check that all tools where found
ifndef GIT_BIN
	$(error The "$(GIT_BIN_NAME)" executable could not be found in your PATH)
endif

.PHONY: check-go-format
## Exists with an error if there are files whose formatting differs from gofmt's
check-go-format: prebuild-check
	@gofmt -s -l ${GOFORMAT_FILES} 2>&1 \
		| tee /tmp/gofmt-errors \
		| read \
	&& echo "ERROR: These files differ from gofmt's style (run 'make format-go-code' to fix this):" \
	&& cat /tmp/gofmt-errors \
	&& exit 1 \
	|| true

.PHONY: format-go-code
## Formats any go file that differs from gofmt's style
format-go-code: prebuild-check
	@gofmt -s -l -w ${GOFORMAT_FILES}


.PHONY: check
check: ## Concurrently runs a whole bunch of static analysis tools
	@gometalinter --enable=misspell --enable=gosimple --enable-gc --vendor --skip=app --skip=client --skip=tool --exclude ^app/test/ --deadline 300s ./...

.PHONY: run
run: build ## runs the service locally
	$(SERVER_BIN)
