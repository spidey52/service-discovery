### TypeScript Client Example

Defining interfaces and a client class to interact with the service discovery server.

```typescript
import axios, { AxiosInstance } from "axios";

export interface Metadata {
 environment: "dev" | "staging" | "prod";
 region: string;
 version: number;
 developer?: string;
 experimental?: boolean;
}

export interface Instance {
 serviceName: string;
 id: string;
 host: string;
 port: number;
 metadata: Metadata;
}

interface LookupFilter {
 service?: string;
 metadata?: Partial<Metadata>;
}

export class ServiceDiscoveryClient {
 private http: AxiosInstance;
 private heartbeatIntervalId?: NodeJS.Timeout;
 private heartbeatFailures = 0;

 constructor(private baseUrl: string) {
  this.http = axios.create({ baseURL: baseUrl, timeout: 5000 });
 }

 /** Register service instance */
 async register(instance: Instance): Promise<void> {
  this.validateInstance(instance);
  await this.http.post("/register", instance);
 }

 /** Send heartbeat */
 async heartbeat(serviceName: string, id: string): Promise<void> {
  await this.http.post("/heartbeat", { serviceName, id });
 }

 /** Start automatic heartbeat */
 startHeartbeat(serviceName: string, id: string, intervalMs = 10000) {
  if (this.heartbeatIntervalId) clearInterval(this.heartbeatIntervalId);

  this.heartbeatIntervalId = setInterval(async () => {
   try {
    await this.heartbeat(serviceName, id);
    this.heartbeatFailures = 0;
   } catch (err) {
    this.heartbeatFailures++;
    console.error("Heartbeat error:", err);

    if (this.heartbeatFailures >= 3) {
     console.warn("⚠ Stopping heartbeat due to repeated failures.");
     this.stopHeartbeat();
    }
   }
  }, intervalMs);
 }

 stopHeartbeat() {
  if (this.heartbeatIntervalId) clearInterval(this.heartbeatIntervalId);
 }

 /** Lookup */
 async lookup(filter: LookupFilter = {}): Promise<Instance[]> {
  try {
   const params: Record<string, string | number | boolean> = {};
   if (filter.service) params.service = filter.service;

   if (filter.metadata) {
    Object.entries(filter.metadata).forEach(([key, val]) => {
     if (val !== undefined) params[key] = val as any;
    });
   }

   const res = await this.http.get("/lookup", { params });
   return res.data;
  } catch (error) {
   console.error("Lookup error:", error);
   return [];
  }
 }

 /** Auto: register + start heartbeat */
 async autoRegister(service: Instance, heartbeatMs = 10_000) {
  console.log("📡 Registering with service discovery...");
  await this.register(service);

  console.log("❤️ Starting heartbeat...");
  this.startHeartbeat(service.serviceName, service.id, heartbeatMs);

  console.log(`🚀 Service Discovery active → ${service.serviceName} (${service.id})`);
 }

 /** Validate metadata */
 private validateInstance(instance: Instance) {
  if (!instance.metadata.environment || !instance.metadata.region || instance.metadata.version === undefined) {
   throw new Error("Missing required metadata fields: `environment`, `region`, `version`");
  }
 }
}
```

### Usage Example

```typescript
import { ServiceDiscoveryClient, Instance } from "./service-discovery";

const discovery = new ServiceDiscoveryClient("http://localhost:4000");

const myService: Instance = {
 serviceName: "order-service",
 id: "order-483",
 host: "127.0.0.1",
 port: 8080,
 metadata: {
  environment: "dev",
  region: "us-east",
  version: 2,
  developer: "john",
  experimental: false,
 },
};

async function start() {
 await discovery.autoRegister(myService, 5000);

 // Later — lookup example:
 const services = await discovery.lookup({
  service: "order-service",
  metadata: { environment: "dev" },
 });

 console.log("Lookup:", services);
}

start().catch(console.error);
```
