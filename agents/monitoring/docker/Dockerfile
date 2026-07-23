# Dockerfile.agent

FROM golang:1.26.3-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Cache the build files and downloaded modules
# Limit the number of processors to be used by compiler, otherwise it would occupy the entire Docker
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -p=2 -o agent ./

# ---- runtime ----
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/agent .

RUN adduser -D appuser

# Create + fix permissions BEFORE switching user
RUN mkdir -p /app/data && chown -R appuser:appuser /app

USER appuser

CMD ["./agent"]