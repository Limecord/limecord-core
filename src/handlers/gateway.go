package handlers

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func gateway(c *gin.Context) {
	c.JSON(200, gin.H{
		"url": "wss://localhost:8080/api/v8/gateway_ws",
	})
}

func checkOrigin(r *http.Request) bool {
	return true
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	CheckOrigin:       checkOrigin,
	EnableCompression: true,
}

type GatewayPacket struct {
	Op          int             `json:"op"`
	Data        json.RawMessage `json:"d"`
	SequenceNum int             `json:"s"`
	EventName   string          `json:"t"`
}

func writeCompressedJson(conn *websocket.Conn, data interface{}) error {
	var buf bytes.Buffer
	zl := zlib.NewWriter(&buf)
	if err := json.NewEncoder(zl).Encode(data); err != nil {
		return err
	}
	zl.Flush()
	conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
	return zl.Close()
}

func gatewayHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}

	zlib.NewWriter()

	writeCompressedJson(conn, gin.H{
		"op": 10,
		"d": gin.H{
			"heartbeat_interval": 45000,
		},
	})

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var packet GatewayPacket
		if t == websocket.BinaryMessage {
			zlibReader, err := zlib.NewReader(bytes.NewBuffer(msg))
			if err != nil {
				fmt.Println("zlib reader creation failed")
				return
			}

			decoder := json.NewDecoder(zlibReader)
			if err = decoder.Decode(&packet); err != nil {
				fmt.Println("Gateway json unmarshal failed", err)
				return
			}
		} else {
			if err = json.Unmarshal(msg, &packet); err != nil {
				fmt.Println("Gateway json unmarshal failed", err)
				return
			}
		}

		if packet.Op == 1 {
			writeCompressedJson(conn, gin.H{
				"op": 11,
			})
		}

	}
}

// register the routers for the tracking module
func RegisterGateway(router *gin.RouterGroup) {
	router.GET("/gateway", gateway)
	router.GET("/gateway_ws/", func(c *gin.Context) {
		gatewayHandler(c.Writer, c.Request)
	})
}
