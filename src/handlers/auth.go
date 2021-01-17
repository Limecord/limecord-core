package handlers

import "github.com/gin-gonic/gin"

func location_metadata(c *gin.Context) {
	c.JSON(200, gin.H{
		"consent_required": false,
		"country_code":     "US",
	})
}

func login(c *gin.Context) {

}

func register(c *gin.Context) {

}

// register the routers for the tracking module
func RegisterAuth(router *gin.RouterGroup) {
	router = router.Group("/auth")
	router.GET("/location-metadata", location_metadata)
}
