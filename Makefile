COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html
COVERIGNORE_FILE := cover_ignore.txt

UNAME := $(shell uname -s)
ifeq ($(UNAME), Linux)
    OPEN_CMD = xdg-open
else ifeq ($(UNAME), Darwin)
    OPEN_CMD = open
else
    OPEN_CMD = start
endif

TEST_PACKAGES := $(shell go list ./... | grep -v -Ff $(COVERIGNORE_FILE))

test:
	@go test -coverprofile=$(COVERAGE_FILE) -covermode=atomic $(TEST_PACKAGES)

html: test
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	$(OPEN_CMD) $(COVERAGE_HTML)

coverage: test
	@go tool cover -func=$(COVERAGE_FILE) | grep total:

clean:
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)

.PHONY: test html coverage clean
