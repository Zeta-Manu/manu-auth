package application

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Zeta-Manu/manu-auth/config"
	"github.com/Zeta-Manu/manu-auth/internal/adapter/idp"
	"github.com/Zeta-Manu/manu-auth/internal/api/route"
	"github.com/Zeta-Manu/manu-auth/pkg/utils"
)

func NewApplication(cfg config.Config) {
	idpAdapter, err := idp.NewCognitoAdapter(cfg.AuthService.AWS.AccessKey, cfg.AuthService.AWS.SecretAccessKey, cfg.AuthService.Cognito.UserPoolId, cfg.AuthService.Cognito.ClientId, cfg.AuthService.Cognito.Region)
	if err != nil {
		panic("failed")
	}

	router := gin.Default()

	logger, _ := zap.NewProduction()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "healthy"})
	})

	r := utils.RouterWithLogger{
		Router: router,
		Logger: logger,
	}

	route.InitRoutes(r, *idpAdapter)

	startServer(cfg, router, logger)
}

func startServer(cfg config.Config, handler http.Handler, logger *zap.Logger) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.AuthService.HTTP.Port),
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to Start Server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", zap.Error(err))
	}

	select {
	case <-ctx.Done():
		logger.Info("timeout of 5 seconds.")
	}
	logger.Info("Server exiting")
}
