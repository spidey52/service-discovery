package servicediscovery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client represents a Service Discovery client
type Client struct {
	httpClient         *resty.Client
	config             *Config
	heartbeatTicker    *time.Ticker
	heartbeatStopChan  chan struct{}
	heartbeatFailures  int
	heartbeatMutex     sync.RWMutex
	currentServiceName string
	currentInstanceID  string
}

// NewClient creates a new Service Discovery client
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if config.BaseURL == "" {
		return nil, fmt.Errorf("baseURL is required")
	}

	httpClient := resty.New().
		SetBaseURL(config.BaseURL).
		SetTimeout(config.Timeout).
		SetHeader("Content-Type", "application/json")

	return &Client{
		httpClient:        httpClient,
		config:            config,
		heartbeatStopChan: make(chan struct{}),
	}, nil
}

// Register registers a service instance with the discovery server
func (c *Client) Register(ctx context.Context, instance Instance) error {
	if err := c.validateInstance(instance); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetBody(instance).
		SetResult(&Instance{}).
		Post("/register")

	if err != nil {
		return fmt.Errorf("register request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("register failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// Heartbeat sends a heartbeat to keep the service instance alive
func (c *Client) Heartbeat(ctx context.Context, serviceName, id string) error {
	req := HeartbeatRequest{
		ServiceName: serviceName,
		ID:          id,
	}

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		Post("/heartbeat")

	if err != nil {
		return fmt.Errorf("heartbeat request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("heartbeat failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// StartHeartbeat starts automatic heartbeat sending
func (c *Client) StartHeartbeat(serviceName, id string, interval time.Duration) {
	c.StopHeartbeat() // Stop any existing heartbeat

	c.heartbeatMutex.Lock()
	c.currentServiceName = serviceName
	c.currentInstanceID = id
	c.heartbeatFailures = 0
	c.heartbeatMutex.Unlock()

	c.heartbeatTicker = time.NewTicker(interval)
	c.heartbeatStopChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-c.heartbeatTicker.C:
				c.heartbeatMutex.RLock()
				sn := c.currentServiceName
				iid := c.currentInstanceID
				c.heartbeatMutex.RUnlock()

				if sn == "" || iid == "" {
					continue
				}

				ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
				err := c.Heartbeat(ctx, sn, iid)
				cancel()

				if err != nil {
					c.heartbeatMutex.Lock()
					c.heartbeatFailures++
					fmt.Printf("Heartbeat error (%d/%d): %v\n", c.heartbeatFailures, c.config.MaxHeartbeatFailures, err)

					if c.heartbeatFailures >= c.config.MaxHeartbeatFailures {
						fmt.Println("⚠ Stopping heartbeat due to repeated failures.")
						c.heartbeatMutex.Unlock()
						c.StopHeartbeat()
						return
					}
					c.heartbeatMutex.Unlock()
				} else {
					c.heartbeatMutex.Lock()
					c.heartbeatFailures = 0
					c.heartbeatMutex.Unlock()
				}

			case <-c.heartbeatStopChan:
				return
			}
		}
	}()
}

// StopHeartbeat stops automatic heartbeat sending
func (c *Client) StopHeartbeat() {
	if c.heartbeatTicker != nil {
		c.heartbeatTicker.Stop()
		c.heartbeatTicker = nil
	}

	select {
	case c.heartbeatStopChan <- struct{}{}:
	default:
	}

	c.heartbeatMutex.Lock()
	c.currentServiceName = ""
	c.currentInstanceID = ""
	c.heartbeatFailures = 0
	c.heartbeatMutex.Unlock()
}

// Lookup finds service instances matching the filter criteria
func (c *Client) Lookup(ctx context.Context, filter LookupFilter) ([]Instance, error) {
	req := c.httpClient.R().SetContext(ctx)

	// Add service filter
	if filter.Service != "" {
		req.SetQueryParam("service", filter.Service)
	}

	// Add metadata filters
	for key, value := range filter.Metadata {
		switch v := value.(type) {
		case string:
			req.SetQueryParam(key, v)
		case int:
			req.SetQueryParam(key, fmt.Sprintf("%d", v))
		case bool:
			req.SetQueryParam(key, fmt.Sprintf("%t", v))
		}
	}

	var instances []Instance
	resp, err := req.SetResult(&instances).Get("/lookup")

	if err != nil {
		return nil, fmt.Errorf("lookup request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("lookup failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return instances, nil
}

// AutoRegister registers a service and starts automatic heartbeat
func (c *Client) AutoRegister(ctx context.Context, instance Instance, heartbeatInterval time.Duration) error {
	fmt.Println("📡 Registering with service discovery...")
	if err := c.Register(ctx, instance); err != nil {
		return fmt.Errorf("auto registration failed: %w", err)
	}

	fmt.Println("❤️ Starting heartbeat...")
	c.StartHeartbeat(instance.ServiceName, instance.ID, heartbeatInterval)

	fmt.Printf("🚀 Service Discovery active → %s (%s)\n", instance.ServiceName, instance.ID)
	return nil
}

// GetHeartbeatStatus returns the current heartbeat status
func (c *Client) GetHeartbeatStatus() (isRunning bool, failureCount int) {
	c.heartbeatMutex.RLock()
	defer c.heartbeatMutex.RUnlock()

	isRunning = c.heartbeatTicker != nil
	failureCount = c.heartbeatFailures
	return
}

// Close gracefully shuts down the client
func (c *Client) Close() {
	c.StopHeartbeat()
}
