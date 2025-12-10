package handlers

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// SPAHandler serves a single-page React app with fallback to index.html
type SPAHandler struct {
	BuildDir string // Path to React build folder
}

// NewSPAHandler returns a Gin handler for serving SPA
func NewSPAHandler(buildDir string) *SPAHandler {
	return &SPAHandler{BuildDir: buildDir}
}

// Handle serves static files if they exist, otherwise serves index.html
func (h *SPAHandler) Handle(c *gin.Context) {
	// Full path to the requested file
	reqPath := filepath.Join(h.BuildDir, c.Request.URL.Path)

	// Check if file exists and is not a directory
	if info, err := os.Stat(reqPath); err == nil && !info.IsDir() {
		c.File(reqPath)
		return
	}

	// Fallback to index.html
	c.File(filepath.Join(h.BuildDir, "index.html"))
}
