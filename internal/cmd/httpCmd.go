package cmd

import (
	"app/app/routes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// HTTP is serve http ot https
func HttpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "Run server on HTTP protocol",
		Run: func(cmd *cobra.Command, args []string) {
			r := gin.Default()
			routes.Router(r)
			r.Run(":" + viper.GetString("APP_PORT")) // Start server on the configured port
		},
	}
}
