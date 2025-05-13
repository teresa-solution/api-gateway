# Teresa API Gateway
Centralized entry point for handling all client HTTP requests within the VMware infrastructure. The API Gateway routes incoming traffic to appropriate backend services using gRPC, enforces security policies, and supports observability via structured logs, traces, and metrics. Integrated with Zerolog, Jaeger, and Prometheus for monitoring and troubleshooting.
