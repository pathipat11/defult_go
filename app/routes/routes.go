// app/routes/routes.go
package routes

import (
	"context"
	"net/http"
	"strings"
	"time"

	"app/app/controller"
	"app/config"
	"app/internal/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Router sets up all the routes for the application
func Router(app *gin.Engine) {

	// Middleware
	app.Use(otelgin.Middleware(viper.GetString("APP_NAME")))
	app.Use(corsMiddleware())

	// Liveness: process is up.
	app.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Readiness: dependencies (database) are reachable.
	app.GET("/readyz", func(ctx *gin.Context) {
		db := config.GetDB()
		if db == nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": "database not initialized"})
			return
		}
		c, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(c); err != nil {
			logger.Err("readiness check failed: ", err)
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": "database unreachable"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// Create a new group for /api/v1
	apiV1 := app.Group("/api/v1")

	// Build the controllers once and share them across route groups.
	ctl := controller.New()

	// Define groups of routes under /api/v1
	Auth(apiV1.Group("/auth"), ctl)
	Product(apiV1.Group("/products"), ctl)
	User(apiV1.Group("/users"), ctl)
}

// corsMiddleware builds the CORS config from CORS_ALLOW_ORIGINS. A value of "*"
// allows all origins; otherwise provide a comma-separated allow-list.
func corsMiddleware() gin.HandlerFunc {
	origins := viper.GetString("CORS_ALLOW_ORIGINS")

	conf := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	if strings.TrimSpace(origins) == "*" || origins == "" {
		// Credentials cannot be combined with a wildcard origin per the CORS spec,
		// so reflect any origin instead.
		conf.AllowAllOrigins = false
		conf.AllowOriginFunc = func(origin string) bool { return true }
	} else {
		list := strings.Split(origins, ",")
		for i := range list {
			list[i] = strings.TrimSpace(list[i])
		}
		conf.AllowOrigins = list
	}

	return cors.New(conf)
}
