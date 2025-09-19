FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/main.go \
 && go build -o kafka-emitter ./scripts/kafka-emitter

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/internal/web ./internal/web
COPY --from=builder /app/kafka-emitter ./kafka-emitter

CMD ["./server"]
