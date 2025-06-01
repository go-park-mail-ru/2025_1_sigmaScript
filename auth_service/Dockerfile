############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git

WORKDIR /service
COPY go.mod ./
COPY go.sum ./
COPY .env ./
RUN go mod download

COPY pkg/ ./pkg
COPY auth_service/ ./auth_service


WORKDIR /service/auth_service

ENV CGO_ENABLED=0
RUN go build -o /service/auth_service/cmd/main /service/auth_service/cmd/main.go

############################
# STEP 2 build a small image
############################
FROM scratch

COPY --from=builder /service/auth_service/cmd/main /bin/main
COPY --from=builder /service/.env /service/auth_service/.env
COPY --from=builder /service/auth_service/internal/config/config.yml /service/auth_service/internal/config/config.yml

WORKDIR /service/auth_service
HEALTHCHECK --interval=10s --timeout=5s --start-period=10s --retries=5 CMD wget --spider http://localhost:8081/health || exit 1

ENTRYPOINT ["/bin/main"]
