package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sujayai/backend-go/internal/handlers"
	"github.com/sujayai/backend-go/internal/middleware"
	"github.com/sujayai/backend-go/internal/storage"
	"github.com/sujayai/backend-go/pkg/config"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Initialize storage
	store, err := storage.New(cfg)
	if err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}

	// Initialize handlers
	postsHandler := handlers.NewPostsHandler(store)

	// Setup Gin
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())

	// Serve uploads
	r.Static("/uploads", cfg.UploadsDir)

	// API routes
	api := r.Group("/api")
	{
		api.GET("/posts", postsHandler.GetPosts)
		api.GET("/posts/:slug", postsHandler.GetPost)
		api.POST("/posts", middleware.AdminAuth(cfg.AdminKey), postsHandler.CreatePost)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}
