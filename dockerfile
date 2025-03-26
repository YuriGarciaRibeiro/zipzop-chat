FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o zipzop ./cmd/zipzop/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/zipzop .
COPY --from=builder /app/.env .
COPY --from=builder /app/public ./public

EXPOSE ${PORT}
CMD ["./zipzop"]