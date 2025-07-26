### Stage 1: Build
FROM golang:1.24.5-alpine AS builder

# Set working directory
ADD . /opt/auth0-api-gateway
WORKDIR /opt/auth0-api-gateway

RUN apk add --no-cache --virtual webstuff ca-certificates \
    curl \
    git \
    && go mod download -x \
    && go mod vendor \
    && apk del webstuff

# Build a static binary
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -a -ldflags '-s -w' -o /usr/local/bin/auth0-api-gateway cmd/main.go

### Stage 2: Final image
FROM alpine:3.22.1

# Copy static binary from build stage
COPY --from=builder /usr/local/bin/auth0-api-gateway /usr/local/bin/auth0-api-gateway

EXPOSE 80

CMD ["/usr/local/bin/auth0-api-gateway"]