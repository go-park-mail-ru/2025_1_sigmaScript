FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git

WORKDIR /service
COPY go.mod go.sum ./
RUN go mod download

COPY . ./main_server
COPY prometheus/prometheus.yml ./prometheus/prometheus.yml
WORKDIR /service/main_server

ENV CGO_ENABLED=0
RUN go build -o /service/main_server/cmd/main cmd/main.go

FROM scratch

COPY --from=builder /service/main_server/cmd/main /bin/main
COPY --from=builder /service/main_server/.env /service/main_server/.env
COPY --from=builder /service/main_server/internal/config/config.yml /service/main_server/internal/config/config.yml

COPY --from=builder /service/main_server/.env /bin/.env
COPY --from=builder /service/main_server/.env /service/main_server/.env
COPY --from=builder /service/main_server/internal/config/config.yml /bin/config/config.yml


WORKDIR /service/main_server

ENTRYPOINT ["/bin/main"]
