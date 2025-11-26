ARG GO_VERSION=1

FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app

# Preload modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Build code
COPY . .
RUN go build -v -o run-app .


# Final runtime image
FROM debian:bookworm

# REQUIRED — without this MikroTik TLS fails!
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

# Copy binary
COPY --from=builder /usr/src/app/run-app /usr/local/bin/run-app

# Expose port (optional but Fly likes it)
EXPOSE 8080

CMD ["run-app"]
