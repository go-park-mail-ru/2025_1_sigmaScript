COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html
COVERIGNORE_FILE := cover_ignore.txt
APP_NAME := filmlook.a
AUTH_SERVICE_NAME := filmlook_auth_service.a
USER_SERVICE_NAME := filmlook_user_service.a
MOVIE_SERVICE_NAME := filmlook_movie_service.a

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

build-bin:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(APP_NAME)$(if $(filter windows,$(GOOS)),.exe,) ./cmd/main.go

build-windows:
	@make build-bin GOOS=windows GOARCH=amd64

build-linux:
	@make build-bin GOOS=linux GOARCH=amd64

build-darwin:
	@make build-bin GOOS=darwin GOARCH=arm64

build:
	@./start_app.sh --build

run:
	@./start_app.sh

stop:
	@./stop_app.sh

remove:
	@./stop_app.sh --remove

test:
	@go test -coverprofile=$(COVERAGE_FILE) -covermode=atomic $(TEST_PACKAGES)

html: test
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	$(OPEN_CMD) $(COVERAGE_HTML)

coverage: test
	@go tool cover -func=$(COVERAGE_FILE) | grep total:

clean:
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)


build-auth-service:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(AUTH_SERVICE_NAME)$(if $(filter windows,$(GOOS)),.exe,) ./auth_service/cmd/main.go

build-user-service:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(USER_SERVICE_NAME)$(if $(filter windows,$(GOOS)),.exe,) ./user_service/cmd/main.go

run-auth-service: build-auth-service
	@./$(AUTH_SERVICE_NAME)

run-user-service: build-user-service
	@./$(USER_SERVICE_NAME)

build-movie-service:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(MOVIE_SERVICE_NAME)$(if $(filter windows,$(GOOS)),.exe,) ./movie_service/cmd/main.go

run-movie-service: build-movie-service
	@./$(MOVIE_SERVICE_NAME)

run-auth-service: build-auth-service
	@./$(AUTH_SERVICE_NAME)

.PHONY: build build-bin cross-build run stop test html coverage clean build-auth-service run-auth-service build-user-service run-user-service
