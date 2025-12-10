// Package servicediscovery provides a Go client for the Service Discovery system
package servicediscovery

import (
	"time"
)

// Environment represents the deployment environment
type Environment string

const (
	EnvDev     Environment = "dev"
	EnvStaging Environment = "staging"
	EnvProd    Environment = "prod"
)

// Metadata contains service metadata information
type Metadata struct {
	Environment  Environment `json:"environment" validate:"required,oneof=dev staging prod"`
	Region       string      `json:"region" validate:"required"`
	Version      int         `json:"version" validate:"required,min=0"`
	Developer    string      `json:"developer,omitempty"`
	Experimental bool        `json:"experimental,omitempty"`
}

// Instance represents a service instance
type Instance struct {
	ServiceName   string      `json:"serviceName" validate:"required"`
	ID            string      `json:"id" validate:"required"`
	Host          string      `json:"host" validate:"required"`
	Port          int         `json:"port" validate:"required,min=1,max=65535"`
	Mode          Environment `json:"mode" validate:"required,oneof=dev staging prod"`
	Metadata      Metadata    `json:"metadata" validate:"required"`
	Health        string      `json:"health,omitempty"`
	LastHeartbeat time.Time   `json:"lastHeartbeat,omitempty"`
}

// LookupFilter contains filters for service lookup
type LookupFilter struct {
	Service  string                 `json:"service,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// HeartbeatRequest represents a heartbeat request
type HeartbeatRequest struct {
	ServiceName string `json:"serviceName" validate:"required"`
	ID          string `json:"id" validate:"required"`
}

// Config contains client configuration
type Config struct {
	BaseURL              string        `validate:"required,url"`
	Timeout              time.Duration `validate:"min=1s"`
	MaxHeartbeatFailures int           `validate:"min=1"`
}

// DefaultConfig returns a default client configuration
func DefaultConfig(baseURL string) *Config {
	return &Config{
		BaseURL:              baseURL,
		Timeout:              5 * time.Second,
		MaxHeartbeatFailures: 3,
	}
}
