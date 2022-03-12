PACKAGE = github.com/Asmodai/gotools
VERSION ?= 0.0.0

CMD_DIR  = $(PACKAGE)/cmd
ROOT_DIR = $(PWD)

LOGREAD = logread

.phony: configs clean

all: deps build

deps:
	@echo Getting dependencies
	@go mod vendor

tidy:
	@echo Tidying dependencies
	@go mod tidy

listdeps:
	@echo Listing dependencies
	@go list -m all

build: logread

logread: deps
	$(eval GIT_VERSION = $(shell scripts/tag2semver.sh $(VERSION) 2>/dev/null))
	@echo "Building $(LOGREAD) with '$(GIT_VERSION)'"
	@go build                                               \
		-tags=go_json                                   \
		-o $(ROOT_DIR)/bin/$(LOGREAD)                   \
		-ldflags "-s -w $(GIT_VERSION)"                 \
		$(CMD_DIR)/$(LOGREAD)

test: deps
	@echo Running tests
	@go test --tags=testing $$(go list ./...) -coverprofile=tests.out
	@go tool cover -html=tests.out -o coverage.html

.PHONY: clean
clean:
	@echo Cleaning $(APP)
	@-rm $(ROOT_DIR)/bin/* 2>/dev/null
	@-rm tests.out 2>/dev/null
	@-rm coverage.html 2>/dev/null
