package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/odobenus/etf"
)

type Encoder interface {
	Encode(gin.H) ([]byte, error)
	Decode([]byte) (gin.H, error)
}

type JSONEncoder struct { }
type ETFEncoder struct { }

func (e *ETFEncoder) Encode(h gin.H) ([]byte, error) {
	context := new(etf.Context)
	buffer := new(bytes.Buffer)
	if err := context.Write(buffer, h); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (e *ETFEncoder) Decode(data []byte)(gin.H, error) {
	context := new(etf.Context)
	resp, err := context.Read(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	switch resp.(type) {
	case gin.H:
		return resp.(gin.H), nil
	default:
		return nil, fmt.Errorf("invalid type, dictionary required for gateway")
	}
}

func (e *JSONEncoder) Encode(h gin.H) ([]byte, error) {
	return json.Marshal(h)
}

func (e *JSONEncoder) Decode(data []byte) (gin.H, error) {
	decoded := gin.H {}
	err := json.Unmarshal(data, &decoded)
	return decoded, err
}

func NewEncoder(encoderType int) Encoder {
	switch encoderType {
	case GATEWAY_COMPRESS_NONE:
		fallthrough
	case GATEWAY_ENCODING_JSON:
		return &JSONEncoder{}
	case GATEWAY_ENCODING_ETF:
		return &ETFEncoder{}
	default:
		panic(fmt.Errorf("invalid encoder type"))

	}
}