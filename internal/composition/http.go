package composition

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/gin-gonic/gin"
)

func NewHttpServer(cfg *config.Config, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           router,
		ReadTimeout:       time.Duration(cfg.HttpReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.HttpWriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(cfg.HttpIdleTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(cfg.HttpReadHeaderTimeout) * time.Second,
	}
}

func RunHttpServer(ctx context.Context, cfg *config.Config, server *http.Server) error {
	lc := net.ListenConfig{
		KeepAlive: 30 * time.Second,
	}

	listener, err := lc.Listen(ctx, "tcp", ":"+cfg.AppPort)
	if err != nil {
		return err
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}
