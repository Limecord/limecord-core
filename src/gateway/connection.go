package gateway

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"strconv"
	"time"
)

const (
	closeGracePeriod = 2 * time.Second

	// Gateway encodings
	GATEWAY_ENCODING_JSON = 0
	GATEWAY_ENCODING_ETF = 1
	// default encoder
	GATEWAY_ENCODING_NONE = GATEWAY_ENCODING_JSON

	// Gateway compression
	GATEWAY_COMPRESS_ZLIB = 0
	GATEWAY_COMPRESS_ZSTD = 1
	// no compression
	GATEWAY_COMPRESS_NONE = GATEWAY_COMPRESS_ZSTD + 1
)

var (
	gatewayEncodings = []string { "json", "etf" }
	gatewayCompressions = []string { "zlib-stream",  "zstd-stream" }
)

type Connection struct {
	webSocket *websocket.Conn

	version int
	encoder Encoder

	compressor Compressor
	// do we compress the blocks
	compressBlocks bool
	compression int
}

func NewConnection(ws* websocket.Conn) *Connection {
	return &Connection {
		ws,
		// default to 0 or nil
		0,nil, nil, false, 0,
	}
}

func (c *Connection) ValidSetting(settings []string, setting string) (bool, int) {
	for idx, settingVal := range settings {
		if settingVal == setting {
			return true, idx
		}
	}
	return false, 0
}

func (c *Connection) ReadSettings(ctx *gin.Context) error {
	valid := false
	encodingIdx := GATEWAY_ENCODING_NONE
	compressionIdx := GATEWAY_COMPRESS_NONE

	encoding := ctx.Query("encoding")
	compression := ctx.Query("compress")
	version, err := strconv.ParseUint(ctx.Query("v"), 10, 9)
	if err != nil || version < 6 || version > 8 {
		fmt.Printf("Version: %v\n", version)
		if err != nil {
			return fmt.Errorf("Error reading gateway version: %v", err)
		}

		return fmt.Errorf("Invalid gateway version")
	}

	if encoding != "" {
		if valid, encodingIdx = c.ValidSetting(gatewayEncodings, encoding); !valid {
			return fmt.Errorf("Invalid gateway encoding")
		}
	}

	if compression != "" {
		if valid, compressionIdx = c.ValidSetting(gatewayCompressions, compression); !valid {
			return fmt.Errorf("Invalid gateway compress")
		}
	}

	fmt.Printf("Gateway Connection (version=%v, compression=%v, encoding=%v)\n",
		version, encoding, compression)

	c.version = int(version)
	c.encoder = NewEncoder(encodingIdx)
	c.compressor = NewCompressor(compressionIdx)
	c.compression = compressionIdx

	return nil
}

func (c *Connection) SendHeartbeatAck() error {
	packet := &GatewayPacket {
		Op: OP_HEARTBEAT_ACK,
	}

	return packet.WritePacket(c)
}

func (c *Connection) SendHello() error {
	packet := &GatewayPacket {
		Op: OP_HELLO,
		Data: gin.H {
			"heartbeat_interval": 45000,
			"_trace": []string{"limecord"},
		},
	}

	return packet.WritePacket(c)
}

func (c *Connection) ProcessPacket(packet *GatewayPacket) error {
	fmt.Printf("process: %v\n", packet.Op)
	switch packet.Op {
	case OP_HEARTBEAT:
		_ = c.SendHeartbeatAck()
		break

	case OP_IDENTIFY:
		// this should be determined by the identify packet
		c.compressBlocks = false
		break

	default:
		break
	}

	return nil
}

func (c *Connection) Close(err error) {
	message := ""
	if err != nil {
		message = err.Error()
	}

	_ = c.webSocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, message))
	time.Sleep(closeGracePeriod)
	_ = c.webSocket.Close()
}

func (c *Connection) Start(ctx *gin.Context) error {
	var (
		err error
	)

	defer c.Close(err)
	if err := c.ReadSettings(ctx); err != nil {
		return err
	}

	if err := c.SendHello(); err != nil {
		return err
	}

	for {
		err, packet := NewPacketFromConnection(c)
		if err != nil {
			return err
		} else if err = c.ProcessPacket(packet); err != nil {
			return err
		}
	}
}