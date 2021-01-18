package gateway

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

type GatewayPacket struct {
	Op          int             `json:"op,omitempty" mapstructure:"op,omitempty"`
	Data        interface{}		`json:"d,omitempty"  mapstructure:"d,omitempty"`
	SequenceNum int             `json:"s,omitempty"  mapstructure:"s,omitempty"`
	EventName   string          `json:"t,omitempty"  mapstructure:"t,omitempty"`
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (p *GatewayPacket) Write(conn *Connection, data []byte) error {
	return conn.webSocket.WriteMessage(websocket.BinaryMessage, data)
}

func (p *GatewayPacket) WritePacketChunked(conn *Connection, data []byte, chunkSize int) error {
	totalSize := len(data)

	fmt.Printf("block: %v\n", data)

	for i := 0; i < totalSize; {
		chunk := data[i:min(i + chunkSize, totalSize)]
		fmt.Printf("write: %v\n", chunk)
		if err := p.Write(conn, chunk); err != nil {
			return err
		}

		i += chunkSize
	}

	fmt.Printf("totalSize: %d, chunkSize: %d\n", totalSize, chunkSize)
	return nil
}
func (p *GatewayPacket) WritePacket(conn *Connection) error {
	mapped := new(gin.H)
	if err := mapstructure.Decode(p, mapped); err != nil {
		return err
	}

	encoded, err := conn.encoder.Encode(*mapped)
	if err != nil {
		return err
	}

	fmt.Printf("encoded: %v\n", string(encoded))

	compressed, written, err := conn.compressor.Compress(encoded)
	if err != nil {
		return err
	}

	// i've never seen such a jank way of doing it
	switch conn.compression {
	case GATEWAY_COMPRESS_ZLIB:
		if !conn.compressBlocks {
			err = p.Write(conn, compressed)
			break
		}

		if written > 0 {
			chunk0 := compressed[:2]
			chunk1 := compressed[2:]

			err = p.WritePacketChunked(conn, chunk0, 1024)
			if err == nil {
				err = p.WritePacketChunked(conn, chunk1, 1024)
			}
		} else {
			err = p.WritePacketChunked(conn, compressed, 1024)
		}

		break
	default:
		err = p.Write(conn, compressed)
	}

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}

	return nil
}

func (p *GatewayPacket) ReadPacket(conn *Connection, messageType int, message []byte) (error) {
	decoded, err := conn.encoder.Decode(message)
	if err != nil {
		return err
	}

	if err = mapstructure.Decode(decoded, p); err != nil {
		return err
	}

	return nil
}

func NewPacketFromConnection(conn *Connection) (error, *GatewayPacket) {
	packet := &GatewayPacket {}
	msgType, msg, err := conn.webSocket.ReadMessage()
	if err != nil {
		return err, nil
	} else if len(msg) > 4096 {
		return fmt.Errorf("payload exceeded length"), nil
	} else if err := packet.ReadPacket(conn, msgType, msg); err != nil {
		return err, nil
	}

	return nil,packet
}