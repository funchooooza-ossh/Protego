package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/funchooooza-ossh/protego/cmd/server/handlers"
	_ "github.com/funchooooza-ossh/protego/docs"
	"github.com/funchooooza-ossh/protego/internal/composition"
	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title ProteGO API
// @version 1.0
// @description This is the infrastructure auth service for your platform.
// @host localhost:8080
// @securityDefinitions.apikey AccessToken
// @in cookie
// @name access_token
// @securityDefinitions.apikey RefreshToken
// @in cookie
// @name refresh_token
// @BasePath /
// @schemes http
func main() {
	cfg := config.Load()
	fmt.Printf("Starting ProteGO on port %s...\n", cfg.AppPort)
	fmt.Printf("Database DSN: %s\n", cfg.DatabaseDsn())

	usecases := composition.ProvideDependencies(cfg)

	router := gin.Default()

	// Routes
	group := router.Group(cfg.MainRoute)
	group.POST("/login", handlers.MakeLoginHandler(usecases.Login, cfg))
	group.POST("/auth/logout", func(c *gin.Context) {
		handlers.LogoutHandler(c, usecases.Logout)
	})
	group.POST("/register", handlers.MakeRegisterHandler(usecases.Register))
	group.GET("/auth/forward", handlers.MakeAuthHandler(usecases.Auth, cfg))
	group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create custom HTTP server with timeouts
	srv := composition.NewHttpServer(cfg, router)

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	fmt.Println("Server stopped")
}
