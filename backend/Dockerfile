FROM golang:1.26.2-alpine AS builder

WORKDIR /build

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.27.0 && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN sqlc generate
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o server ./cmd/server

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/server .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /build/db/migrations ./migrations
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]