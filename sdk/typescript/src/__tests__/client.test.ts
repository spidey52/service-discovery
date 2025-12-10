import { ServiceDiscoveryClient } from "../client";
import { Instance } from "../types";

describe("ServiceDiscoveryClient", () => {
	const baseUrl = "http://localhost:4000";
	let client: ServiceDiscoveryClient;

	beforeEach(() => {
		client = new ServiceDiscoveryClient({ baseUrl });
	});

	describe("validation", () => {
		it("should validate instance correctly", () => {
			const validInstance: Instance = {
				serviceName: "test-service",
				id: "test-001",
				host: "127.0.0.1",
				port: 8080,
				mode: "dev",
				metadata: {
					environment: "dev",
					region: "us-east",
					version: 1,
				},
			};

			expect(() => client["validateInstance"](validInstance)).not.toThrow();
		});

		it("should throw error for invalid service name", () => {
			const invalidInstance: Instance = {
				serviceName: "",
				id: "test-001",
				host: "127.0.0.1",
				port: 8080,
				mode: "dev",
				metadata: {
					environment: "dev",
					region: "us-east",
					version: 1,
				},
			};

			expect(() => client["validateInstance"](invalidInstance)).toThrow("Service name is required");
		});

		it("should throw error for invalid port", () => {
			const invalidInstance: Instance = {
				serviceName: "test-service",
				id: "test-001",
				host: "127.0.0.1",
				port: 0,
				mode: "dev",
				metadata: {
					environment: "dev",
					region: "us-east",
					version: 1,
				},
			};

			expect(() => client["validateInstance"](invalidInstance)).toThrow("Valid port number is required");
		});

		it("should throw error for invalid environment", () => {
			const invalidInstance: Instance = {
				serviceName: "test-service",
				id: "test-001",
				host: "127.0.0.1",
				port: 8080,
				mode: "dev",
				metadata: {
					environment: "invalid" as any,
					region: "us-east",
					version: 1,
				},
			};

			expect(() => client["validateInstance"](invalidInstance)).toThrow("Environment must be one of: dev, staging, prod");
		});
	});

	describe("heartbeat status", () => {
		it("should return correct initial status", () => {
			const status = client.getHeartbeatStatus();
			expect(status.isRunning).toBe(false);
			expect(status.failureCount).toBe(0);
		});
	});
});
