package servicediscovery

import (
	"testing"
)

func TestValidateInstance(t *testing.T) {
	client, _ := NewClient(DefaultConfig("http://localhost:4000"))

	tests := []struct {
		name     string
		instance Instance
		wantErr  bool
	}{
		{
			name: "valid instance",
			instance: Instance{
				ServiceName: "test-service",
				ID:          "test-001",
				Host:        "127.0.0.1",
				Port:        8080,
				Mode:        EnvDev,
				Metadata: Metadata{
					Environment: EnvDev,
					Region:      "us-east",
					Version:     1,
				},
			},
			wantErr: false,
		},
		{
			name: "missing service name",
			instance: Instance{
				ServiceName: "",
				ID:          "test-001",
				Host:        "127.0.0.1",
				Port:        8080,
				Mode:        EnvDev,
				Metadata: Metadata{
					Environment: EnvDev,
					Region:      "us-east",
					Version:     1,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			instance: Instance{
				ServiceName: "test-service",
				ID:          "test-001",
				Host:        "127.0.0.1",
				Port:        0,
				Mode:        EnvDev,
				Metadata: Metadata{
					Environment: EnvDev,
					Region:      "us-east",
					Version:     1,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid environment",
			instance: Instance{
				ServiceName: "test-service",
				ID:          "test-001",
				Host:        "127.0.0.1",
				Port:        8080,
				Mode:        Environment("invalid"),
				Metadata: Metadata{
					Environment: EnvDev,
					Region:      "us-east",
					Version:     1,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.validateInstance(tt.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInstance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMetadata(t *testing.T) {
	client, _ := NewClient(DefaultConfig("http://localhost:4000"))

	tests := []struct {
		name     string
		metadata Metadata
		wantErr  bool
	}{
		{
			name: "valid metadata",
			metadata: Metadata{
				Environment: EnvDev,
				Region:      "us-east",
				Version:     1,
			},
			wantErr: false,
		},
		{
			name: "missing region",
			metadata: Metadata{
				Environment: EnvDev,
				Region:      "",
				Version:     1,
			},
			wantErr: true,
		},
		{
			name: "negative version",
			metadata: Metadata{
				Environment: EnvDev,
				Region:      "us-east",
				Version:     -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.validateMetadata(tt.metadata)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMetadata() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetHeartbeatStatus(t *testing.T) {
	client, _ := NewClient(DefaultConfig("http://localhost:4000"))

	// Initially should not be running
	isRunning, failureCount := client.GetHeartbeatStatus()
	if isRunning {
		t.Error("Expected heartbeat to not be running initially")
	}
	if failureCount != 0 {
		t.Errorf("Expected failure count to be 0, got %d", failureCount)
	}
}
