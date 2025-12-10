/**
 * Service Discovery TypeScript SDK
 *
 * Types and interfaces for the Service Discovery system
 */

export type Environment = "dev" | "staging" | "prod";

export interface Metadata {
	/** Environment where the service is running */
	environment: Environment;
	/** Geographic region */
	region: string;
	/** Service version number */
	version: number;
	/** Developer name (optional) */
	developer?: string;
	/** Whether this is an experimental version (optional) */
	experimental?: boolean;
}

export interface Instance {
	/** Name of the service */
	serviceName: string;
	/** Unique identifier for this instance */
	id: string;
	/** Host address */
	host: string;
	/** Port number */
	port: number;
	/** Environment mode */
	mode: Environment;
	/** Service metadata */
	metadata: Metadata;
	/** Health status (optional) */
	health?: string;
	/** Last heartbeat timestamp (optional) */
	lastHeartbeat?: Date;
}

export interface LookupFilter {
	/** Service name to filter by */
	service?: string;
	/** Metadata filters */
	metadata?: Partial<Metadata>;
}

export interface HeartbeatRequest {
	/** Service name */
	serviceName: string;
	/** Instance ID */
	id: string;
}

export interface ServiceDiscoveryConfig {
	/** Base URL of the service discovery server */
	baseUrl: string;
	/** Request timeout in milliseconds */
	timeout?: number;
	/** Maximum number of heartbeat failures before stopping */
	maxHeartbeatFailures?: number;
}
