package handlers

import (
	"github.com/gin-gonic/gin"
	"limecord-core/middleware"
)

func getUser(ctx *gin.Context) {

}

// register the routers for the tracking module
func RegisterUser(router *gin.RouterGroup) {
	router.GET("/users/:user", middleware.AuthMiddleware, getUser)
}