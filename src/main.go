package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"limecord-core/handlers"
	"limecord-core/middleware"
	"limecord-core/utils"
	"log"
	"os"
	"strings"
)

var (
	// The port that the http server will listen on
	SERVER_PORT = GetConfigVar("PORT", "8080")
)

// Get configuration variables from the environment, or fallback to the default value
func GetConfigVar(name string, def string) string {
	// Get the environment variable
	res := os.Getenv(name)
	// If the variable doesn't exist, it'll return a blank value (we also check for whitespace, just in-case)
	if strings.TrimSpace(res) == "" {
		// fallback to the default value
		return def
	}
	return res
}

func main() {
	// create the gin server
	httpServer := gin.New()

	// on 404 not found
	httpServer.NoRoute(func(c *gin.Context) {
		c.JSON(404, utils.GetAPIError(utils.API_NOTFOUND, utils.API_NOTFOUND_MESSAGE))
	})

	apiGroup := httpServer.Group("/api/:version", middleware.APIVersionMiddleware)
	handlers.RegisterTracking(apiGroup)

	// Start the api server, on the port specified in the environment, if errors log the error
	if err := httpServer.Run(fmt.Sprintf(":%s", SERVER_PORT)); err != nil {
		log.Fatal(err)
	}
}