package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	sharederrorenvelope "github.com/unsia-erp/shared-errorenvelope"
	sharedobservability "github.com/unsia-erp/shared-observability"
	"github.com/unsia-erp/unsia-sso-token-service/internal/domain"
	"github.com/unsia-erp/unsia-sso-token-service/internal/handler"
	"github.com/unsia-erp/unsia-sso-token-service/internal/infrastructure/database"
	"github.com/unsia-erp/unsia-sso-token-service/internal/service"
)

func main() {
	_ = godotenv.Load()
	sharedobservability.InitLogger("sso-token-service")

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	if err := db.AutoMigrate(&domain.RefreshToken{}, &domain.TokenMetadata{}); err != nil {
		log.Printf("Warning: AutoMigrate failed: %v", err)
	}

	r := gin.New()
	r.Use(sharedobservability.CorrelationIDMiddleware())
	r.Use(sharedobservability.RequestLoggerMiddleware())
	r.Use(sharedobservability.MetricsMiddleware())
	r.GET("/metrics", sharedobservability.MetricsHandler())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "sso-token-service"})
	})

	tokenService := service.NewTokenService(db)
	tokenHandler := handler.NewTokenHandler(tokenService)

	r.POST("/api/v1/tokens/refresh", tokenHandler.RefreshToken)
	r.POST("/api/v1/tokens/revoke", tokenHandler.RevokeToken)

	r.Use(sharederrorenvelope.ErrorHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}

	go func() {
		sharedobservability.Logger.Info().Msgf("SSO Token Service started on port %s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Server failed to run: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	sharedobservability.Logger.Info().Msg("SSO Token Service shutting down...")
	time.Sleep(1 * time.Second)
}
