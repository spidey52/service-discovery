# Service Discovery TypeScript SDK

A TypeScript client SDK for the Service Discovery system. This SDK provides a convenient way to interact with the Service Discovery server from Node.js and browser applications.

## Installation

```bash
npm install @spidey52/service-discovery-sdk
```

## Usage

### Basic Setup

```typescript
import { ServiceDiscoveryClient, Instance } from "@spidey52/service-discovery-sdk";

const client = new ServiceDiscoveryClient({
 baseUrl: "http://localhost:4000",
 timeout: 5000, // optional
});
```

### Registering a Service

```typescript
const myService: Instance = {
 serviceName: "order-service",
 id: "order-001",
 host: "127.0.0.1",
 port: 8080,
 mode: "dev",
 metadata: {
  environment: "dev",
  region: "us-east",
  version: 1,
  developer: "john.doe",
  experimental: false,
 },
};

await client.register(myService);
```

### Automatic Registration with Heartbeat

```typescript
// Register and start automatic heartbeat (every 10 seconds)
await client.autoRegister(myService, 10000);
```

### Manual Heartbeat

```typescript
// Send heartbeat manually
await client.heartbeat("order-service", "order-001");

// Or start/stop automatic heartbeat
client.startHeartbeat("order-service", "order-001", 15000);
client.stopHeartbeat();
```

### Service Lookup

```typescript
// Find all instances of a service
const services = await client.lookup({
 service: "order-service",
});

// Find services with metadata filters
const devServices = await client.lookup({
 service: "order-service",
 metadata: {
  environment: "dev",
  region: "us-east",
 },
});
```

### Advanced Configuration

```typescript
const client = new ServiceDiscoveryClient({
 baseUrl: "https://discovery.example.com",
 timeout: 10000,
 maxHeartbeatFailures: 5, // Stop heartbeat after 5 failures
});
```

### Heartbeat Status

```typescript
const status = client.getHeartbeatStatus();
console.log(`Heartbeat running: ${status.isRunning}`);
console.log(`Failure count: ${status.failureCount}`);
```

## API Reference

### ServiceDiscoveryClient

#### Constructor

```typescript
new ServiceDiscoveryClient(config: ServiceDiscoveryConfig)
```

#### Methods

- `register(instance: Instance): Promise<void>` - Register a service instance
- `heartbeat(serviceName: string, id: string): Promise<void>` - Send heartbeat
- `startHeartbeat(serviceName: string, id: string, intervalMs?: number): void` - Start automatic heartbeat
- `stopHeartbeat(): void` - Stop automatic heartbeat
- `lookup(filter?: LookupFilter): Promise<Instance[]>` - Lookup services
- `autoRegister(service: Instance, heartbeatMs?: number): Promise<void>` - Register and start heartbeat
- `getHeartbeatStatus(): { isRunning: boolean; failureCount: number }` - Get heartbeat status

### Types

#### Instance

```typescript
interface Instance {
 serviceName: string;
 id: string;
 host: string;
 port: number;
 mode: "dev" | "staging" | "prod";
 metadata: Metadata;
 health?: string;
 lastHeartbeat?: Date;
}
```

#### Metadata

```typescript
interface Metadata {
 environment: "dev" | "staging" | "prod";
 region: string;
 version: number;
 developer?: string;
 experimental?: boolean;
}
```

#### LookupFilter

```typescript
interface LookupFilter {
 service?: string;
 metadata?: Partial<Metadata>;
}
```

#### ServiceDiscoveryConfig

```typescript
interface ServiceDiscoveryConfig {
 baseUrl: string;
 timeout?: number;
 maxHeartbeatFailures?: number;
}
```

## Error Handling

The SDK throws descriptive errors for different failure scenarios:

- **Registration errors**: Invalid data, server errors
- **Heartbeat errors**: Network issues, service not found
- **Lookup errors**: Network issues (returns empty array on failure)

```typescript
try {
 await client.register(myService);
} catch (error) {
 console.error("Registration failed:", error.message);
}
```

## Development

### Building

```bash
npm run build
```

### Testing

```bash
npm test
```

### Linting

```bash
npm run lint
```

## License

MIT
