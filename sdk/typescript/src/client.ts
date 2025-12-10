import axios, { AxiosError, AxiosInstance } from "axios";
import { HeartbeatRequest, Instance, LookupFilter, Metadata, ServiceDiscoveryConfig } from "./types";

/**
 * Service Discovery Client
 *
 * A TypeScript client for interacting with the Service Discovery server.
 * Provides methods for service registration, heartbeat monitoring, and service lookup.
 */
export class ServiceDiscoveryClient {
	private http: AxiosInstance;
	private heartbeatIntervalId?: NodeJS.Timeout;
	private heartbeatFailures = 0;
	private config: Required<ServiceDiscoveryConfig>;

	/**
	 * Create a new Service Discovery client
	 *
	 * @param config - Client configuration
	 */
	constructor(config: ServiceDiscoveryConfig) {
		this.config = {
			timeout: 5000,
			maxHeartbeatFailures: 3,
			...config,
		};

		this.http = axios.create({
			baseURL: this.config.baseUrl,
			timeout: this.config.timeout,
		});
	}

	/**
	 * Register a service instance with the discovery server
	 *
	 * @param instance - Service instance to register
	 * @throws Error if registration fails or validation fails
	 */
	async register(instance: Instance): Promise<void> {
		this.validateInstance(instance);

		try {
			await this.http.post("/register", instance);
		} catch (error) {
			throw this.handleError("Failed to register service", error);
		}
	}

	/**
	 * Send a heartbeat to keep the service instance alive
	 *
	 * @param serviceName - Name of the service
	 * @param id - Instance ID
	 * @throws Error if heartbeat fails
	 */
	async heartbeat(serviceName: string, id: string): Promise<void> {
		const request: HeartbeatRequest = { serviceName, id };

		try {
			await this.http.post("/heartbeat", request);
		} catch (error) {
			throw this.handleError("Failed to send heartbeat", error);
		}
	}

	/**
	 * Start automatic heartbeat sending
	 *
	 * @param serviceName - Name of the service
	 * @param id - Instance ID
	 * @param intervalMs - Heartbeat interval in milliseconds (default: 10000)
	 */
	startHeartbeat(serviceName: string, id: string, intervalMs = 10000): void {
		this.stopHeartbeat(); // Stop any existing heartbeat

		this.heartbeatIntervalId = setInterval(async () => {
			try {
				await this.heartbeat(serviceName, id);
				this.heartbeatFailures = 0;
			} catch (error) {
				this.heartbeatFailures++;
				console.error(`Heartbeat error (${this.heartbeatFailures}/${this.config.maxHeartbeatFailures}):`, error);

				if (this.heartbeatFailures >= this.config.maxHeartbeatFailures) {
					console.warn("⚠ Stopping heartbeat due to repeated failures.");
					this.stopHeartbeat();
				}
			}
		}, intervalMs);
	}

	/**
	 * Stop automatic heartbeat sending
	 */
	stopHeartbeat(): void {
		if (this.heartbeatIntervalId) {
			clearInterval(this.heartbeatIntervalId);
			this.heartbeatIntervalId = undefined;
			this.heartbeatFailures = 0;
		}
	}

	/**
	 * Lookup service instances
	 *
	 * @param filter - Lookup filter criteria
	 * @returns Array of matching service instances
	 */
	async lookup(filter: LookupFilter = {}): Promise<Instance[]> {
		try {
			const params: Record<string, string | number | boolean> = {};

			if (filter.service) {
				params.service = filter.service;
			}

			if (filter.metadata) {
				Object.entries(filter.metadata).forEach(([key, value]) => {
					if (value !== undefined) {
						params[key] = value as any;
					}
				});
			}

			const response = await this.http.get<Instance[]>("/lookup", { params });
			return response.data;
		} catch (error) {
			console.error("Lookup error:", error);
			return [];
		}
	}

	/**
	 * Automatically register a service and start heartbeat
	 *
	 * @param service - Service instance to register
	 * @param heartbeatMs - Heartbeat interval in milliseconds
	 */
	async autoRegister(service: Instance, heartbeatMs = 10_000): Promise<void> {
		console.log("📡 Registering with service discovery...");
		await this.register(service);

		console.log("❤️ Starting heartbeat...");
		this.startHeartbeat(service.serviceName, service.id, heartbeatMs);

		console.log(`🚀 Service Discovery active → ${service.serviceName} (${service.id})`);
	}

	/**
	 * Get the current heartbeat status
	 *
	 * @returns Object with heartbeat status information
	 */
	getHeartbeatStatus(): { isRunning: boolean; failureCount: number } {
		return {
			isRunning: this.heartbeatIntervalId !== undefined,
			failureCount: this.heartbeatFailures,
		};
	}

	/**
	 * Validate service instance data
	 *
	 * @private
	 * @param instance - Instance to validate
	 * @throws Error if validation fails
	 */
	private validateInstance(instance: Instance): void {
		if (!instance.serviceName?.trim()) {
			throw new Error("Service name is required");
		}

		if (!instance.id?.trim()) {
			throw new Error("Instance ID is required");
		}

		if (!instance.host?.trim()) {
			throw new Error("Host is required");
		}

		if (!instance.port || instance.port <= 0 || instance.port > 65535) {
			throw new Error("Valid port number is required (1-65535)");
		}

		if (!instance.mode || !["dev", "staging", "prod"].includes(instance.mode)) {
			throw new Error("Mode must be one of: dev, staging, prod");
		}

		this.validateMetadata(instance.metadata);
	}

	/**
	 * Validate metadata
	 *
	 * @private
	 * @param metadata - Metadata to validate
	 * @throws Error if validation fails
	 */
	private validateMetadata(metadata: Metadata): void {
		if (!metadata.environment || !["dev", "staging", "prod"].includes(metadata.environment)) {
			throw new Error("Environment must be one of: dev, staging, prod");
		}

		if (!metadata.region?.trim()) {
			throw new Error("Region is required");
		}

		if (typeof metadata.version !== "number" || metadata.version < 0) {
			throw new Error("Version must be a non-negative number");
		}
	}

	/**
	 * Handle and transform errors
	 *
	 * @private
	 * @param message - Error message
	 * @param error - Original error
	 * @returns Transformed error
	 */
	private handleError(message: string, error: unknown): Error {
		if (axios.isAxiosError(error)) {
			const axiosError = error as AxiosError;
			const status = axiosError.response?.status;
			const data = axiosError.response?.data as any;

			if (status === 400) {
				return new Error(`${message}: ${data?.error || "Bad request"}`);
			} else if (status === 404) {
				return new Error(`${message}: Service not found`);
			} else if (status === 500) {
				return new Error(`${message}: Server error`);
			}
		}

		return new Error(`${message}: ${error instanceof Error ? error.message : "Unknown error"}`);
	}
}
