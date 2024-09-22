package main

import (
	"app/app/routes"
	"app/config"
	"app/internal/cmd"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func main() {

	config.Init() // Initialize configuration

	// Start the HTTP server
	if err := command(); err != nil {
		log.Fatalf("Error running command: %v", err)
	}
}

// Command setup for running the HTTP server
func command() error {
	cmda := &cobra.Command{
		Use:  "app",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmda.AddCommand(httpCmd())
	cmda.AddCommand(cmd.Migrate())

	return cmda.Execute()

}
// Command for starting the HTTP server
func httpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "Run server on HTTP protocol",
		Run: func(cmd *cobra.Command, args []string) {
			r := gin.Default()
			routes.Router(r)
			r.Run(":8080") // Start server on port 8080
		},
	}
}
