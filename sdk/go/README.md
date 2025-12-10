# Service Discovery Go SDK

A Go client SDK for the Service Discovery system. This SDK provides a convenient way to interact with the Service Discovery server from Go applications.

## Installation

```bash
go get github.com/spidey52/service-discovery-sdk
```

## Usage

### Basic Setup

```go
package main

import (
    "log"
    servicediscovery "github.com/spidey52/service-discovery-sdk"
)

func main() {
    client, err := servicediscovery.NewClient(servicediscovery.DefaultConfig("http://localhost:4000"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
}
```

### Registering a Service

```go
instance := servicediscovery.Instance{
    ServiceName: "order-service",
    ID:          "order-001",
    Host:        "127.0.0.1",
    Port:        8080,
    Mode:        servicediscovery.EnvDev,
    Metadata: servicediscovery.Metadata{
        Environment: servicediscovery.EnvDev,
        Region:      "us-east",
        Version:     1,
        Developer:   "john.doe",
        Experimental: false,
    },
}

err := client.Register(context.Background(), instance)
if err != nil {
    log.Fatal(err)
}
```

### Automatic Registration with Heartbeat

```go
// Register and start automatic heartbeat (every 10 seconds)
err := client.AutoRegister(context.Background(), instance, 10*time.Second)
if err != nil {
    log.Fatal(err)
}
```

### Manual Heartbeat

```go
// Send heartbeat manually
err := client.Heartbeat(context.Background(), "order-service", "order-001")
if err != nil {
    log.Printf("Heartbeat failed: %v", err)
}

// Or start/stop automatic heartbeat
client.StartHeartbeat("order-service", "order-001", 15*time.Second)
client.StopHeartbeat()
```

### Service Lookup

```go
// Find all instances of a service
services, err := client.Lookup(context.Background(), servicediscovery.LookupFilter{
    Service: "order-service",
})
if err != nil {
    log.Fatal(err)
}

// Find services with metadata filters
devServices, err := client.Lookup(context.Background(), servicediscovery.LookupFilter{
    Service: "order-service",
    Metadata: map[string]interface{}{
        "environment": "dev",
        "region":      "us-east",
    },
})
if err != nil {
    log.Fatal(err)
}
```

### Advanced Configuration

```go
config := &servicediscovery.Config{
    BaseURL:             "https://discovery.example.com",
    Timeout:             10 * time.Second,
    MaxHeartbeatFailures: 5, // Stop heartbeat after 5 failures
}

client, err := servicediscovery.NewClient(config)
```

### Heartbeat Status

```go
isRunning, failureCount := client.GetHeartbeatStatus()
fmt.Printf("Heartbeat running: %v\n", isRunning)
fmt.Printf("Failure count: %d\n", failureCount)
```

## API Reference

### Types

#### Instance

```go
type Instance struct {
    ServiceName   string      `json:"serviceName"`
    ID            string      `json:"id"`
    Host          string      `json:"host"`
    Port          int         `json:"port"`
    Mode          Environment `json:"mode"`
    Metadata      Metadata    `json:"metadata"`
    Health        string      `json:"health,omitempty"`
    LastHeartbeat time.Time   `json:"lastHeartbeat,omitempty"`
}
```

#### Metadata

```go
type Metadata struct {
    Environment  Environment `json:"environment"`
    Region       string      `json:"region"`
    Version      int         `json:"version"`
    Developer    string      `json:"developer,omitempty"`
    Experimental bool        `json:"experimental,omitempty"`
}
```

#### LookupFilter

```go
type LookupFilter struct {
    Service  string
    Metadata map[string]interface{}
}
```

#### Config

```go
type Config struct {
    BaseURL             string
    Timeout             time.Duration
    MaxHeartbeatFailures int
}
```

### Client Methods

- `NewClient(config *Config) (*Client, error)` - Create a new client
- `Register(ctx context.Context, instance Instance) error` - Register a service instance
- `Heartbeat(ctx context.Context, serviceName, id string) error` - Send heartbeat
- `StartHeartbeat(serviceName, id string, interval time.Duration)` - Start automatic heartbeat
- `StopHeartbeat()` - Stop automatic heartbeat
- `Lookup(ctx context.Context, filter LookupFilter) ([]Instance, error)` - Lookup services
- `AutoRegister(ctx context.Context, instance Instance, heartbeatInterval time.Duration) error` - Register and start heartbeat
- `GetHeartbeatStatus() (isRunning bool, failureCount int)` - Get heartbeat status
- `Close()` - Gracefully shut down the client

## Error Handling

The SDK returns descriptive errors for different failure scenarios:

- **Registration errors**: Invalid data, network errors, server errors
- **Heartbeat errors**: Network issues, service not found
- **Lookup errors**: Network issues

```go
err := client.Register(context.Background(), instance)
if err != nil {
    log.Printf("Registration failed: %v", err)
    return
}
```

## Testing

```bash
go test ./...
```

## License

MIT
