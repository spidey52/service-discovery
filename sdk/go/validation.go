package servicediscovery

import (
	"fmt"
	"strings"
)

// validateInstance validates a service instance
func (c *Client) validateInstance(instance Instance) error {
	if strings.TrimSpace(instance.ServiceName) == "" {
		return fmt.Errorf("serviceName is required")
	}

	if strings.TrimSpace(instance.ID) == "" {
		return fmt.Errorf("id is required")
	}

	if strings.TrimSpace(instance.Host) == "" {
		return fmt.Errorf("host is required")
	}

	if instance.Port < 1 || instance.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	if instance.Mode != EnvDev && instance.Mode != EnvStaging && instance.Mode != EnvProd {
		return fmt.Errorf("mode must be one of: dev, staging, prod")
	}

	return c.validateMetadata(instance.Metadata)
}

// validateMetadata validates service metadata
func (c *Client) validateMetadata(metadata Metadata) error {
	if metadata.Environment != EnvDev && metadata.Environment != EnvStaging && metadata.Environment != EnvProd {
		return fmt.Errorf("environment must be one of: dev, staging, prod")
	}

	if strings.TrimSpace(metadata.Region) == "" {
		return fmt.Errorf("region is required")
	}

	if metadata.Version < 0 {
		return fmt.Errorf("version must be non-negative")
	}

	return nil
}
