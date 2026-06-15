# --- Build stage ---
FROM golang:1.26.3-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build argument to choose which cmd to build
ARG CMD_PATH

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o app ${CMD_PATH}

# --- Runtime stage ---
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 3000

CMD ["./app"]