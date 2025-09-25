# AI Agent Instructions for grpc-api-gateway

## Project Overview
This is a high-performance gRPC-based API Gateway implemented in Go that includes a load balancer with dynamic server health checks.

## Key Components

### Load Balancer (`loadbalancer/loadbalancer.go`)
- Implements round-robin load balancing with health checks
- Uses Go's `httputil.ReverseProxy` for request forwarding
- Health check mechanism through `testServer()` function
- Default configuration: Listens on port `:8090` with endpoint `/loadBalancer`
- Servers are expected at `http://localhost:808{0-9}`

### Server Management (`servers/servers.go`)
- Dynamic server creation and management
- Concurrent server operation using goroutines and `sync.WaitGroup`
- Server ports allocated in range `808{0-9}` (max 10 servers)
- Each server provides:
  - Root endpoint `/` returning server port info
  - Shutdown endpoint `/shutdown` for graceful termination

## Common Operations

### Starting the Gateway
The gateway is initialized in `main.go` with a specified number of backend servers:
```go
loadbalancer.MakeLoadBalancer(5) // Creates load balancer with 5 backend servers
```

### Server Health Checks
Load balancer performs health checks before forwarding requests:
- Checks server availability via HTTP GET
- Expects HTTP 200 OK response
- Automatically rotates to next server on failure

## Code Patterns

### Error Handling
- Server errors are handled through HTTP status codes
- Critical failures (e.g., exceeding max servers) use `log.Fatal`
- Network errors in health checks return `false` for graceful degradation

### Concurrency
- Server operations use goroutines for concurrent handling
- Synchronization through `sync.WaitGroup`
- Context-based shutdown for graceful termination

### URL Management
- URLs are managed through Go's `url.URL` type
- Dynamic endpoint creation in `createEndpoint()`
- Round-robin rotation through `Endpoints.Shuffle()`

## Integration Points
- Load balancer entry point: `http://localhost:8090/loadBalancer`
- Backend servers: `http://localhost:808{0-9}`
- Health check: HTTP GET to backend server root path

## Development Workflow
1. Implement new features in appropriate component:
   - Load balancing logic in `loadbalancer.go`
   - Server management in `servers.go`
2. Ensure health check compatibility for new endpoints
3. Update main.go for any configuration changes