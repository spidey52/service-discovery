package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spidey52/service-discovery/models"
	"github.com/spidey52/service-discovery/repository"
)

// SetupRoutes wires all endpoints
func SetupRoutes(r *gin.Engine, repo *repository.MongoRepo, heartbeatTTL time.Duration) {
	r.POST("/register", func(c *gin.Context) {
		var inst models.Instance
		if err := c.ShouldBindJSON(&inst); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := repo.Register(c.Request.Context(), inst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, inst)
	})

	r.POST("/heartbeat", func(c *gin.Context) {
		var req struct {
			ServiceName string `json:"serviceName"`
			ID          string `json:"id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := repo.UpdateHeartbeat(c.Request.Context(), req.ServiceName, req.ID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "heartbeat ok"})
	})

	r.GET("/lookup", func(c *gin.Context) {
		service := c.Query("service")
		mode := c.Query("mode")
		metadata := map[string]interface{}{}
		for key, vals := range c.Request.URL.Query() {
			if key == "service" || key == "mode" {
				continue
			}
			metadata[key] = parseString(vals[0])
		}
		instances, err := repo.Find(c.Request.Context(), service, mode, metadata, true, heartbeatTTL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, instances)
	})
}

// helper to parse strings to bool/int/float if possible
func parseString(s string) interface{} {
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}
	return s
}
