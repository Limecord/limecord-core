package main

import (
	"fmt"
	"limecord-core/gateway"
	"limecord-core/handlers"
	"limecord-core/middleware"
	"limecord-core/utils"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	// The port that the http server will listen on
	SERVER_PORT = GetConfigVar("PORT", "8080")
	// the domain name of the discord instance
	SERVER_HOSTNAME = GetConfigVar("HOSTNAME", "localhost")
	//
	SERVER_HTTPS_CERT = GetConfigVar("HTTPS_CERT", "./test/server.crt")
	//
	SERVER_HTTPS_KEY = GetConfigVar("HTTPS_KEY", "./test/server.pem")
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
	// allow all through cors by default
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowWildcard = true
	corsConfig.AllowWebSockets = true
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("*")
	httpServer.Use(cors.New(corsConfig))

	// on 404 not found
	httpServer.NoRoute(func(c *gin.Context) {
		c.JSON(404, utils.GetAPIError(utils.API_NOTFOUND, utils.API_NOTFOUND_MESSAGE))
	})

	apiGroup := httpServer.Group("/api/:version", middleware.APIVersionMiddleware)

	// register our tracking handlers
	handlers.RegisterTracking(apiGroup)

	// register our user handlers
	handlers.RegisterUser(apiGroup)

	handlers.RegisterAuth(apiGroup)

	handlers.RegisterInvites(apiGroup)

	// register our gateway
	gateway.RegisterGateway(apiGroup)

	// Start the api server, on the port specified in the environment, if errors log the error
	err := httpServer.RunTLS(fmt.Sprintf(":%s", SERVER_PORT),
		SERVER_HTTPS_CERT, SERVER_HTTPS_KEY)
	// this is disabled as it gives us an invalid response
	// err := autotls.RunWithManager(httpServer, &autoCertMgr);
	if err != nil {
		log.Fatal(err)
	}
}
