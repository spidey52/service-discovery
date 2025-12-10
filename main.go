package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spidey52/service-discovery/handlers"
	"github.com/spidey52/service-discovery/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURI := "mongodb://localhost:27017/?directConnection=true"
	dbName := "service_registry"
	collName := "registry"
	heartbeatTTL := 30 * time.Second
	cleanupInterval := 10 * time.Second

	// Mongo connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database(dbName).Collection(collName)
	repo := repository.NewMongoRepo(coll)

	// Gin setup
	r := gin.Default()

	// WebSocket endpoint for real-time updates
	r.GET("/ws", handlers.HandleWebSocket)

	handlers.SetupRoutes(r, repo, heartbeatTTL)

	// Serve SPA
	spaHandler := handlers.NewSPAHandler("./ui")
	r.NoRoute(spaHandler.Handle)

	// Cleanup goroutine
	stop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				_ = repo.CleanupDead(context.Background(), heartbeatTTL)
			}
		}
	}()

	// Run server
	go func() {
		fmt.Println("Service discovery running on :4000")
		if err := r.Run(":4000"); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	close(stop)
	_ = client.Disconnect(context.Background())
	fmt.Println("Shutdown complete")
}
