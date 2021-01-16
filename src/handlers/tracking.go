// this registers the tracking endpoints but as we don't
// want to track our users  these are only here for dummy purposes

package handlers

import "github.com/gin-gonic/gin"

// dummy science endpoint
func science(c *gin.Context) {
	c.AbortWithStatus(204)
}

// register the routers for the tracking module
func RegisterTracking(router *gin.RouterGroup) {
	router.POST("/science", science)
}