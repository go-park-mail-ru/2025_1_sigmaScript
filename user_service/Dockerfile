############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git

WORKDIR /service
COPY go.mod go.sum ./
COPY .env ./
RUN go mod download

COPY user_service/ ./user_service
COPY internal/errors/ ./internal/errors
# COPY internal/server/models/ ./internal/server/models
# COPY internal/server/mocks/ ./internal/server/mocks


WORKDIR /service/user_service

ENV CGO_ENABLED=0
RUN go build -o /service/user_service/cmd/main /service/user_service/cmd/main.go

############################
# STEP 2 build a small image
############################
FROM scratch

COPY --from=builder /service/user_service/cmd/main /bin/main
COPY --from=builder /service/.env /service/user_service/.env
COPY --from=builder /service/user_service/internal/config/config.yml /service/user_service/user_service/config/config.yml

WORKDIR /service/user_service
HEALTHCHECK --interval=10s --timeout=5s --start-period=10s --retries=5 CMD wget --spider http://localhost:8083/health || exit 1

ENTRYPOINT ["/bin/main"]
