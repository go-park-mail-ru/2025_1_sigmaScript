COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html
COVERIGNORE_FILE := cover_ignore.txt
APP_NAME := filmlook.a

UNAME := $(shell uname -s)
ifeq ($(UNAME), Linux)
    OPEN_CMD = xdg-open
else ifeq ($(UNAME), Darwin)
    OPEN_CMD = open
else
    OPEN_CMD = start
endif

TEST_PACKAGES := $(shell go list ./... | grep -v -Ff $(COVERIGNORE_FILE))

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

build:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(APP_NAME)$(if $(filter windows,$(GOOS)),.exe,) ./cmd/main.go

build-windows:
	@make build GOOS=windows GOARCH=amd64

build-linux:
	@make build GOOS=linux GOARCH=amd64

build-darwin:
	@make build GOOS=darwin GOARCH=arm64

run: build
	@./$(APP_NAME)

test:
	@go test -coverprofile=$(COVERAGE_FILE) -covermode=atomic $(TEST_PACKAGES)

html: test
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	$(OPEN_CMD) $(COVERAGE_HTML)

coverage: test
	@go tool cover -func=$(COVERAGE_FILE) | grep total:

clean:
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)

.PHONY: build cross-build run test html coverage clean
