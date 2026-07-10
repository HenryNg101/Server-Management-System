# --- Build stage ---
FROM golang:1.26.3-alpine3.23 AS builder

WORKDIR /app

# Copy ONLY auth-service module files first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of auth-service source code
COPY . .

# Build the binary
# RUN --mount=type=cache,target=/root/.cache/go-build \
#     --mount=type=cache,target=/go/pkg/mod \
#     go build -p=1 -o app ./cmd

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o app ./cmd

# RUN go build -p=1 -o app ./cmd

# --- Runtime stage ---
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8083

CMD ["./app"]