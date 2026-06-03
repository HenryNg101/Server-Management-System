# --- Build stage ---
FROM golang:1.26.3-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build argument to choose which cmd to build
ARG CMD_PATH

# RUN go build -o app ${CMD_PATH}
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o app ${CMD_PATH}

# --- Runtime stage ---
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]

# FROM golang:1.26.3-alpine3.23 AS builder

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN --mount=type=cache,target=/go/pkg/mod \
#     go mod download

# COPY . .

# RUN --mount=type=cache,target=/root/.cache/go-build \
#     --mount=type=cache,target=/go/pkg/mod \
#     go build -o /bin/api ./cmd/api/ && \
#     go build -o /bin/servers_checker ./cmd/worker/servers_checker/ && \
#     go build -o /bin/emailer ./cmd/worker/emailer/

# FROM alpine:latest

# WORKDIR /app

# COPY --from=builder /bin/* /app/

# CMD ["sh"]