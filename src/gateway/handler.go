package gateway

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	WsUpgrader = websocket.Upgrader{
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		// we handle compression ourself
		EnableCompression: false,
		CheckOrigin:       func(r *http.Request) bool {
			return true
		},
	}
)

func startGateway(ctx *gin.Context) {
	ws, err := WsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %v\n", err)
		return
	}

	conn := NewConnection(ws)
	if err = conn.Start(ctx); err != nil {
		fmt.Printf("Error occured in connection: %v\n", err)
	}
}

func redirectGateway(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"url": "wss://localhost:8080/api/v8/gateway_ws",
	})
}

func RegisterGateway(router *gin.RouterGroup) {
	router.GET("/gateway", redirectGateway)
	router.GET("/gateway_ws/", startGateway)
}
