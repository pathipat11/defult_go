package cmd

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/app/routes"
	"app/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// HttpCmd serves the application over HTTP with graceful shutdown.
func HttpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "Run server on HTTP protocol",
		Run: func(cmd *cobra.Command, args []string) {
			if viper.GetString("APP_ENV") == "production" {
				gin.SetMode(gin.ReleaseMode)
			}

			r := gin.Default()
			routes.Router(r)

			srv := &http.Server{
				Addr:              ":" + viper.GetString("APP_PORT"),
				Handler:           r,
				ReadHeaderTimeout: 10 * time.Second,
			}

			// Start the server in a goroutine so it doesn't block shutdown handling.
			go func() {
				logger.Infof("HTTP server listening on %s", srv.Addr)
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Err("http server error: ", err)
					os.Exit(1)
				}
			}()

			// Wait for an interrupt/terminate signal.
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			logger.Infof("shutting down server...")

			// Give in-flight requests up to 10 seconds to complete.
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				logger.Err("forced shutdown: ", err)
			}
			logger.Infof("server stopped")
		},
	}
}
