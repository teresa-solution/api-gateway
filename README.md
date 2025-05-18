# Teresa API Gateway 🚀

[![Go Report Card](https://goreportcard.com/badge/github.com/teresa-solution/api-gateway)](https://goreportcard.com/report/github.com/teresa-solution/api-gateway)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/teresa-solution/api-gateway)](https://golang.org/)

<div align="center">
  <img src="https://via.placeholder.com/800x400?text=Teresa+API+Gateway" alt="Teresa API Gateway Logo" width="600" />
</div>

> A high-performance, secure, multi-tenant API Gateway service that provides authentication, rate limiting, metrics, and proxying capabilities for backend gRPC services within the Teresa Solution ecosystem.

## 🌟 Features

- **🔐 Secure HTTPS Endpoints**: TLS encryption for all API traffic
- **🏢 Multi-tenancy**: Support for tenant isolation via subdomain headers with proper context propagation
- **🔑 Authentication Middleware**: Request validation before forwarding to backend services
- **⚡ Rate Limiting**: Protects backend services from excessive traffic
- **📊 Comprehensive Metrics**: Prometheus integration for monitoring all gateway activities
- **🔄 gRPC-Gateway Integration**: Seamless translation between RESTful HTTP APIs and gRPC services
- **✅ Graceful Shutdown**: Ensures no requests are dropped during deployment
- **💻 Developer-Friendly**: Easy to extend with new services and features

## 🧩 Teresa Ecosystem Integration

The API Gateway is the entry point for all Teresa Solution services:

* **[Tenant Management Service](https://github.com/teresa-solution/tenant-management-service)**: API Gateway routes tenant provisioning and management requests to this service
* **[Connection Pool Manager](https://github.com/teresa-solution/connection-pool-manager)**: Database connection pooling for multi-tenant applications

<div align="center">
  <kbd>
    <img src="https://via.placeholder.com/800x400?text=Teresa+Architecture+Diagram" alt="Teresa API Gateway Architecture" width="600" />
  </kbd>
</div>

```
┌─────────────┐       HTTPS       ┌───────────────────────────────────┐         gRPC        ┌─────────────────────┐
│             │                   │           API Gateway             │                     │  Backend Services   │
│   Clients   │◄─────────────────►│  ┌─────────┐ ┌────────┐ ┌──────┐  │◄───────────────────►│  - Tenant Mgmt     │
│ (Web/Mobile)│      (8080)       │  │   TLS   │ │  Auth  │ │ Rate │  │       (50051)       │  - Connection Pool │
└─────────────┘                   │  │Terminator│ │Middleware│ │Limit │  │                     └─────────────────────┘
                                 │  └─────────┘ └────────┘ └──────┘  │
                                 │         │          │         │      │
                                 │         ▼          ▼         ▼      │
                                 │       ┌─────────────────────────┐   │
                                 │       │   Multi-tenant Router   │   │
                                 │       └─────────────────────────┘   │
                                 │                   │                 │
                                 │                   ▼                 │
                                 │       ┌─────────────────────────┐   │
                                 │       │     Metrics Collector   │───┼──────► Prometheus
                                 │       └─────────────────────────┘   │
                                 └───────────────────────────────────┘
```

## 📋 Prerequisites

- **Go 1.18+**
- **TLS certificates** in the `certs/` directory
- **Running gRPC backend services**:
  - Tenant Management Service (port 50051)
  - Connection Pool Manager (port 50052)

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/teresa-solution/api-gateway.git
cd api-gateway

# Install dependencies
go mod download

# Build the application
go build -o api-gateway
```

### Configuration

The API Gateway accepts the following command-line flags:

| Flag | Description | Default |
|------|-------------|---------|
| `--http-port` | HTTP server port | `8080` |
| `--tenant-svc-addr` | Tenant service address | `127.0.0.1:50051` |
| `--pool-mgr-addr` | Connection pool manager address | `127.0.0.1:50052` |

### TLS Certificates

The gateway requires TLS certificates for HTTPS. Place them in the following locations:

- `certs/cert.pem`: Certificate file
- `certs/key.pem`: Private key file

Generate self-signed certificates for development:

```bash
mkdir -p certs
openssl req -x509 -newkey rsa:4096 -keyout certs/key.pem -out certs/cert.pem -days 365 -nodes
```

### Running the API Gateway

```bash
# Start with default configuration
./api-gateway

# Start with custom configuration
./api-gateway --http-port=9090 --tenant-svc-addr=tenant-svc:50051 --pool-mgr-addr=pool-mgr:50052
```

## 🔍 Key Components

### Multi-tenancy

The API Gateway supports multi-tenancy through the `X-Tenant-Subdomain` header. This header value is forwarded to backend services as metadata, allowing tenant-specific processing.

Example usage:
```bash
curl -H "X-Tenant-Subdomain: customer1" https://api.example.com/v1/resources
```

### Authentication

The `AuthMiddleware` validates incoming requests before allowing them to proceed to backend services. The middleware:

- Verifies JWT tokens
- Checks permissions
- Enforces tenant isolation

### Rate Limiting

Protects your backend services from being overwhelmed by limiting the number of requests per client:

- Per-client IP rate limiting
- Per-tenant rate limiting
- Customizable limit thresholds

### Metrics

Prometheus metrics are exposed at `/metrics`. The following metrics are available:

| Metric | Description | Labels |
|--------|-------------|--------|
| `http_requests_total` | Total number of HTTP requests | method, path, status, tenant |
| `http_request_duration_seconds` | Duration of HTTP requests | method, path, status, tenant |
| `http_response_size_bytes` | Size of HTTP responses | method, path, status, tenant |
| `rate_limit_exceeded_total` | Rate limit violations | client_ip, tenant |
| `auth_failures_total` | Authentication failures | reason, path, tenant |
| `api_gateway_uptime_seconds` | API gateway uptime | - |
| `active_connections` | Active connections | - |
| `backend_requests_total` | Backend service requests | service, method, status, tenant |
| `backend_request_duration_seconds` | Backend request duration | service, method, status, tenant |
| `backend_failures_total` | Backend request failures | service, reason, tenant |

## 📁 Project Structure

```
api-gateway/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── handler/                # Service handlers
│   │   ├── tenant.go           # Tenant service registration
│   │   └── pool.go             # Connection pool manager registration
│   ├── middleware/             # HTTP middleware components
│   │   ├── auth.go             # Authentication middleware
│   │   ├── metrics.go          # Metrics collection middleware
│   │   └── ratelimit.go        # Rate limiting middleware
│   └── monitoring/             # Metrics and monitoring
│       └── metrics.go          # Prometheus metrics definition
├── proto/                      # Protocol buffer definitions
│   ├── tenant.proto            # Tenant service definition
│   ├── pool.proto              # Connection pool manager definition
│   └── gen/                    # Generated gRPC code
├── certs/                      # TLS certificates
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksum
└── README.md                   # This file
```

## 🔧 Development

### Adding New Services

To add a new service to the API Gateway:

1. Define your service using Protocol Buffers in the `proto/` directory
2. Generate the gRPC gateway code:
   ```bash
   protoc -I=./proto \
       --go_out=./proto/gen --go_opt=paths=source_relative \
       --go-grpc_out=./proto/gen --go-grpc_opt=paths=source_relative \
       --grpc-gateway_out=./proto/gen --grpc-gateway_opt=paths=source_relative \
       ./proto/your_service.proto
   ```
3. Register the new service handler in `handler/register.go`:
   ```go
   if err := yourservicepb.RegisterYourServiceHandlerFromEndpoint(
       ctx, mux, grpcAddr, opts,
   ); err != nil {
       return err
   }
   ```

## 📖 Documentation

Full API documentation is available at:

- [API Gateway Documentation](https://github.com/teresa-solution/api-gateway/wiki)
- [API Specification](https://github.com/teresa-solution/api-gateway/wiki/API-Specification)

## 💡 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Teresa Solution Team** - *Initial work* - [Teresa Solution](https://github.com/teresa-solution)

## 🙏 Acknowledgments

- [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway) - For the RESTful HTTP to gRPC translation
- [Prometheus](https://prometheus.io/) - For metrics and monitoring
- [Zerolog](https://github.com/rs/zerolog) - For structured logging
