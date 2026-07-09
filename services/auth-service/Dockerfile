# # --- Build stage ---
# FROM golang:1.26.3-alpine3.23 AS builder

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# # Build argument to choose which cmd to build
# ARG CMD_PATH

# # Cache the build files and downloaded modules
# # Limit the number of processors to be used by compiler, otherwise it would occupy the entire Docker
# RUN --mount=type=cache,target=/root/.cache/go-build \
#     --mount=type=cache,target=/go/pkg/mod \
#     go build -p=1 -o app ${CMD_PATH}

# # --- Runtime stage ---
# FROM alpine:latest

# WORKDIR /app

# COPY --from=builder /app/app .

# EXPOSE 8080

# CMD ["./app"]

# --- Build stage ---
FROM golang:1.26.3-alpine3.23 AS builder

WORKDIR /app

# Copy ONLY auth-service module files first (for caching)
COPY services/auth-service/go.mod services/auth-service/go.sum ./
RUN go mod download

# Copy the rest of auth-service source code
COPY services/auth-service/ .

# Build the binary
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -p=1 -o app ./cmd

# --- Runtime stage ---
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8081

CMD ["./app"]