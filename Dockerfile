FROM golang:1.25-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/go-server ./cmd/api

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && \
    adduser -S -G app app && \
    mkdir -p /app/logs /app/uploads /app/tmp && \
    chown -R app:app /app

COPY --from=builder /out/go-server /app/go-server
COPY --chown=app:app config.yaml /app/config.yaml
COPY --chown=app:app configs /app/configs
COPY --chown=app:app migrations /app/migrations

USER app

EXPOSE 8080

CMD ["./go-server"]
