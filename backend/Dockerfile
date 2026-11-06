# Build stage
FROM golang:1.27.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/stockflow \
    ./cmd/app

# Runtime stage
FROM alpine:3.24

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/stockflow ./stockflow

EXPOSE 3000

CMD ["./stockflow"]