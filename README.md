# auth0-api-gateway

[![godoc for greencoda/confiq][godoc-badge]][godoc-url]
[![Go 1.22][goversion-badge]][goversion-url]
[![Build Status][actions-badge]][actions-url]
[![Go Coverage][gocoverage-badge]][gocoverage-url]
[![Go Report card][goreportcard-badge]][goreportcard-url]
[![Docker Hub][dockerhub-badge]][dockerhub-url]

`auth0-api-gateway` is a configurable reverse proxy API Gateway with Auth0 JWT authentication, built in Go. It provides a flexible way to route requests to multiple backend services while enforcing authentication and authorization policies through Auth0 scopes.

## Features

- **Auth0 Integration**: JWT token validation with Auth0 
- **Scope-based Authorization**: Fine-grained access control using Auth0 scopes
- **Reverse Proxy**: Route requests to multiple backend services
- **CORS Support**: Configurable CORS policies per route
- **Rate Limiting**: Built-in rate limiting capabilities
- **Request Logging**: Comprehensive request/response logging
- **Configuration-driven**: YAML-based configuration with live reload
- **Health Checks**: Built-in health monitoring
- **Docker Support**: Ready-to-use Docker container
- **Dependency Injection**: Clean architecture using Uber FX

## Install

```bash
go get -u github.com/greencoda/auth0-api-gateway
```

## Quick Start

### 1. Create Configuration

Create a `config.yaml` file based on the template:

```yaml
auth0:
  audience: "https://your-api.example.com"
  domain: "your-tenant.auth0.com"

server:
  address: ":8080"
  readTimeout: "15s"
  writeTimeout: "15s"
  idleTimeout: "15s"
  maxHeaderBytes: 1048576
  logCalls: true
  releaseStage: "development"
  logLevel: "info"

subrouters:
  - target_url: http://localhost:3001/api
    prefix: "/api/v1"
    strip_prefix: true
    name: "API Service"
    authorization_config:
      required_scopes:
        - "read:api"
        - "write:api"
    auth: true
    gzip: true
    cors:
      allow_credentials: true
      allowed_origins:
        - "https://yourdomain.com"
      allowed_headers:
        - "Authorization"
        - "Content-Type"
```

### 2. Run the Gateway

```bash
# Using Go directly
go run cmd/main.go

# Or using Make
make run

# Or using Docker Hub image (recommended)
docker run -p 8080:80 -v $(pwd)/config.yaml:/config.yaml greencoda/auth0-api-gateway:latest

# Or build locally
docker build -f docker/Dockerfile -t auth0-api-gateway .
docker run -p 8080:80 -v $(pwd)/config.yaml:/config.yaml auth0-api-gateway
```

### 3. Make Authenticated Requests

```bash
# Get an Auth0 token (example using curl)
TOKEN=$(curl -X POST "https://your-tenant.auth0.com/oauth/token" \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": "YOUR_CLIENT_ID",
    "client_secret": "YOUR_CLIENT_SECRET",
    "audience": "https://your-api.example.com",
    "grant_type": "client_credentials"
  }' | jq -r '.access_token')

# Make authenticated request through the gateway
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/api/v1/your-endpoint
```

## Configuration

The gateway is configured using a YAML file. Here's a comprehensive example:

### Auth0 Configuration

```yaml
auth0:
  audience: "https://your-api.example.com"      # Your API identifier in Auth0
  domain: "your-tenant.auth0.com"              # Your Auth0 domain
```

### Server Configuration

```yaml
server:
  address: ":8080"           # Server bind address
  readTimeout: "15s"         # Maximum duration for reading requests
  writeTimeout: "15s"        # Maximum duration for writing responses
  idleTimeout: "15s"         # Maximum duration for idle connections
  maxHeaderBytes: 1048576    # Maximum size of request headers
  logCalls: true            # Enable request/response logging
  releaseStage: "production" # Environment stage (local, development, staging, production)
  logLevel: "info"          # Log level (trace, debug, info, warn, error, fatal, panic)
```

### Subrouter Configuration

Each subrouter defines a route to a backend service:

```yaml
subrouters:
  - target_url: http://backend-service:3000    # Backend service URL
    prefix: "/api/users"                       # Route prefix
    strip_prefix: true                         # Remove prefix before forwarding
    name: "User Service"                       # Descriptive name
    authorization_config:
      required_scopes:                         # Required Auth0 scopes
        - "read:users"
        - "write:users"
    auth: true                                 # Enable authentication
    gzip: true                                 # Enable gzip compression
    rate_limit:                                # Optional rate limiting
      period: "1m"
      limit: 100
    cors:                                      # CORS configuration
      allow_credentials: true
      allowed_origins:
        - "https://yourdomain.com"
        - "https://admin.yourdomain.com"
      allowed_headers:
        - "Authorization"
        - "Content-Type"
        - "X-Requested-With"
      allowed_methods:
        - "GET"
        - "POST"
        - "PUT"
        - "DELETE"
      max_age: 86400
```

## Architecture

The gateway follows a clean architecture pattern with dependency injection:

```
cmd/
  main.go                 # Application entry point

internal/
  config/                 # Configuration structures
    auth0/               # Auth0 configuration
    server/              # Server configuration  
    subrouter/           # Subrouter configuration
    
  middleware/             # HTTP middleware components
    auth0/               # Auth0 JWT validation
    callLogger/          # Request/response logging
    cors/                # CORS handling
    rateLimit/           # Rate limiting
    
  server/                 # HTTP server and reverse proxy
    server.go            # Main server implementation
    reverseProxy.go      # Reverse proxy logic
    
  util/                   # Utility packages
    config/              # Configuration loading
    logging/             # Logging utilities
```

## Middleware

The gateway includes several built-in middleware components:

### Auth0 Middleware
- JWT token validation
- Scope-based authorization
- Automatic token refresh handling
- Comprehensive error responses

### CORS Middleware
- Configurable per-route CORS policies
- Support for preflight requests
- Credential handling

### Rate Limiting Middleware
- Token bucket algorithm
- Configurable limits per route
- Redis support for distributed rate limiting

### Call Logger Middleware
- Structured request/response logging
- Performance metrics
- Error tracking

## Development

### Prerequisites

- Go 1.24 or later
- Make (optional, for convenience commands)

### Setup

```bash
# Clone the repository
git clone https://github.com/greencoda/auth0-api-gateway.git
cd auth0-api-gateway

# Install dependencies
make deps

# Generate mocks (for testing)
make mock

# Run tests
make test

# Run with coverage
make test-cover

# Build binary
make build
```

### Testing

The project includes comprehensive tests with mocking:

```bash
# Run all tests
make test

# Run tests 100 times (for race condition detection)
make test-100

# Generate coverage report
make test-cover
```

## Docker

The Docker image is automatically built and published to Docker Hub on every release and push to master.

### Using Pre-built Image

```bash
# Pull the latest image from Docker Hub
docker pull greencoda/auth0-api-gateway:latest

# Run with local config file
docker run -p 8080:80 \
  -v $(pwd)/config.yaml:/config.yaml \
  greencoda/auth0-api-gateway:latest

# Run with environment-specific config
docker run -p 8080:80 \
  -v /path/to/production/config.yaml:/config.yaml \
  greencoda/auth0-api-gateway:latest

# Run specific version
docker run -p 8080:80 \
  -v $(pwd)/config.yaml:/config.yaml \
  greencoda/auth0-api-gateway:v1.0.0
```

### Building Locally

```bash
# Build the image locally
docker build -f docker/Dockerfile -t auth0-api-gateway .

# Run locally built image
docker run -p 8080:80 \
  -v $(pwd)/config.yaml:/config.yaml \
  auth0-api-gateway
```

### Available Tags

- `latest` - Latest stable release from master branch
- `v1.x.x` - Specific version tags (e.g., `v1.0.0`, `v1.1.0`)
- `master` - Latest build from master branch (development)

## Examples

### Basic API Gateway

Route all `/api/*` requests to a backend service with Auth0 authentication:

```yaml
subrouters:
  - target_url: http://api-backend:3000
    prefix: "/api"
    strip_prefix: true
    name: "Main API"
    authorization_config:
      required_scopes: ["api:access"]
    auth: true
```

### Multi-Service Gateway

Route different prefixes to different services:

```yaml
subrouters:
  - target_url: http://user-service:3001
    prefix: "/users"
    name: "User Service"
    authorization_config:
      required_scopes: ["read:users"]
    auth: true
    
  - target_url: http://order-service:3002
    prefix: "/orders"
    name: "Order Service"
    authorization_config:
      required_scopes: ["read:orders"]
    auth: true
    
  - target_url: http://public-api:3003
    prefix: "/public"
    name: "Public API"
    auth: false  # No authentication required
```

### Rate Limited Endpoints

```yaml
subrouters:
  - target_url: http://heavy-service:3000
    prefix: "/heavy"
    name: "Heavy Processing Service"
    rate_limit:
      period: "1m"
      limit: 10  # 10 requests per minute
    auth: true
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions:

- Create an [issue](https://github.com/greencoda/auth0-api-gateway/issues)

[godoc-badge]: https://pkg.go.dev/badge/github.com/greencoda/auth0-api-gateway
[godoc-url]: https://pkg.go.dev/github.com/greencoda/auth0-api-gateway
[actions-badge]: https://github.com/greencoda/auth0-api-gateway/actions/workflows/main.yml/badge.svg
[actions-url]: https://github.com/greencoda/auth0-api-gateway/actions/workflows/main.yml
[goversion-badge]: https://img.shields.io/badge/Go-1.24-%2300ADD8?logo=go
[goversion-url]: https://golang.org/doc/go1.24
[goreportcard-badge]: https://goreportcard.com/badge/github.com/greencoda/auth0-api-gateway
[goreportcard-url]: https://goreportcard.com/report/github.com/greencoda/auth0-api-gateway
[gocoverage-badge]: https://github.com/greencoda/auth0-api-gateway/wiki/coverage.svg
[gocoverage-url]: https://raw.githack.com/wiki/greencoda/auth0-api-gateway/coverage.html
[dockerhub-badge]: https://img.shields.io/docker/pulls/greencoda/auth0-api-gateway
[dockerhub-url]: https://hub.docker.com/r/greencoda/auth0-api-gateway