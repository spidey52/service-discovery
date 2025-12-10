# Service Discovery

A lightweight, high-performance service discovery system built with Go, MongoDB, and Gin. It provides automatic service registration, health monitoring via heartbeats, and flexible service lookup with metadata filtering.

## Features

- **Service Registration**: Register service instances with metadata
- **Heartbeat Monitoring**: Automatic cleanup of dead services based on heartbeat TTL
- **Flexible Lookup**: Query services by name, environment, region, version, and custom metadata
- **RESTful API**: Simple HTTP endpoints for all operations
- **MongoDB Storage**: Persistent storage with efficient indexing
- **Graceful Shutdown**: Proper cleanup on termination
- **SDKs**: Official client libraries for TypeScript and Go applications

## Prerequisites

- Go 1.19 or later
- MongoDB 4.0 or later
- (Optional) Node.js 16+ for TypeScript SDK
- (Optional) Go 1.19+ for Go SDK

## Installation

1. Clone the repository:

```bash
git clone https://github.com/spidey52/service-discovery.git
cd service-discovery
```

2. Install dependencies:

```bash
go mod tidy
```

3. Ensure MongoDB is running locally on port 27017, or update the connection string in `main.go`.

## Usage

### Running the Server

```bash
go run main.go
```

The server will start on port 4000 and display:

```
Service discovery running on :4000
```

### Configuration

The server uses the following default configuration (configurable in `main.go`):

- **MongoDB URI**: `mongodb://localhost:27017`
- **Database**: `service_registry`
- **Collection**: `registry`
- **Heartbeat TTL**: 30 seconds
- **Cleanup Interval**: 10 seconds
- **Port**: 4000

## API Endpoints

### Register Service

Register a new service instance.

```http
POST /register
Content-Type: application/json

{
  "serviceName": "order-service",
  "id": "order-483",
  "host": "127.0.0.1",
  "port": 8080,
  "mode": "dev",
  "metadata": {
    "environment": "dev",
    "region": "us-east",
    "version": 2,
    "developer": "john",
    "experimental": false
  }
}
```

**Response:**

```json
{
 "serviceName": "order-service",
 "id": "order-483",
 "host": "127.0.0.1",
 "port": 8080,
 "mode": "dev",
 "metadata": {
  "environment": "dev",
  "region": "us-east",
  "version": 2,
  "developer": "john",
  "experimental": false
 },
 "health": "",
 "lastHeartbeat": "2025-12-10T10:30:00Z"
}
```

### Send Heartbeat

Update the last heartbeat timestamp for a service instance.

```http
POST /heartbeat
Content-Type: application/json

{
  "serviceName": "order-service",
  "id": "order-483"
}
```

**Response:**

```json
{
 "message": "heartbeat ok"
}
```

### Lookup Services

Find service instances with optional filtering.

```http
GET /lookup?service=order-service&environment=dev&region=us-east&version=2
```

**Query Parameters:**

- `service`: Service name (required)
- `mode`: Environment mode filter
- Additional metadata filters (environment, region, version, developer, experimental)

**Response:**

```json
[
 {
  "serviceName": "order-service",
  "id": "order-483",
  "host": "127.0.0.1",
  "port": 8080,
  "mode": "dev",
  "metadata": {
   "environment": "dev",
   "region": "us-east",
   "version": 2,
   "developer": "john",
   "experimental": false
  },
  "health": "",
  "lastHeartbeat": "2025-12-10T10:30:00Z"
 }
]
```

## SDKs

Official client libraries are available for both TypeScript and Go applications.

### TypeScript SDK

A complete TypeScript SDK is available in the `sdk/typescript/` directory. Install it with:

```bash
npm install @spidey52/service-discovery-sdk
```

See `sdk/typescript/README.md` for detailed usage instructions.

### Go SDK

A complete Go SDK is available in the `sdk/go/` directory. Install it with:

```bash
go get github.com/spidey52/service-discovery-sdk
```

See `sdk/go/README.md` for detailed usage instructions.

## Architecture

The system consists of three main components:

1. **HTTP Handlers** (`handlers/http.go`): REST API endpoints
2. **Repository Layer** (`repository/mongo_repo.go`): MongoDB operations
3. **Models** (`models/instance.go`): Data structures and validation

### Service Instance Model

```go
type Instance struct {
    ServiceName   string    `json:"serviceName"`
    ID            string    `json:"id"`
    Host          string    `json:"host"`
    Port          int       `json:"port"`
    Mode          string    `json:"mode"` // dev, staging, prod
    Metadata      Metadata  `json:"metadata"`
    Health        string    `json:"health"`
    LastHeartbeat time.Time `json:"lastHeartbeat"`
}
```

### Metadata Model

```go
type Metadata struct {
    Environment  string `json:"environment"` // dev, staging, prod
    Region       string `json:"region"`
    Version      int    `json:"version"`
    Developer    string `json:"developer,omitempty"`
    Experimental bool   `json:"experimental,omitempty"`
}
```

## Health Monitoring

- Services must send periodic heartbeats to stay alive
- Dead services are automatically cleaned up after TTL expires
- Default heartbeat TTL is 30 seconds
- Cleanup runs every 10 seconds

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o service-discovery main.go
```

### Docker Support

To run with Docker, ensure MongoDB is available and update the connection string accordingly.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For questions or issues, please open an issue on GitHub.
