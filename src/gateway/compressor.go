package gateway

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

type Compressor interface {
	Compress([]byte) ([]byte, int64, error)
	Decompress([]byte) ([]byte, error)
}

type ZlibCompressor struct {}
type ZstdCompressor struct {}
type NoCompressor struct {}

func (z *ZlibCompressor) Compress(data []byte) ([]byte, int64, error) {
	var (
		err error
		written int64
	)

	buffer := new(bytes.Buffer)
	zlibWriter := zlib.NewWriter(buffer)

	if written, err = io.Copy(zlibWriter, bytes.NewBuffer(data)); err != nil {
		return nil, 0, fmt.Errorf("failed to copy encoded bytes: %v", err)
	}

	if err = zlibWriter.Flush(); err != nil {
		return nil, 0, fmt.Errorf("failed to flush zlib buffer: %v", err)
	}

	return buffer.Bytes(), written, nil
}

func (z *ZlibCompressor) Decompress(data []byte) ([]byte, error) {
	buffer := new(bytes.Buffer)
	zlibReader, err := zlib.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("zlib reader creation failed: %v", err)
	}

	if _, err = io.Copy(buffer, zlibReader); err != nil {
		return nil, fmt.Errorf("failed to copy decoded bytes: %v", err)
	}

	return buffer.Bytes(), nil
}

func (z *ZstdCompressor) Compress(data []byte) ([]byte, int64, error) {
	//return gozstd.Compress(nil, data), nil
	return nil, 0, nil
}

func (z *ZstdCompressor) Decompress(data []byte) ([]byte, error) {
	//return gozstd.Decompress(nil, data)
	return nil, nil
}

func (z *NoCompressor) Compress(data []byte) ([]byte, int64, error) {
	return data, int64(len(data)), nil
}

func (z *NoCompressor) Decompress(data []byte) ([]byte, error) {
	return data, nil
}

func NewCompressor(compressorType int) Compressor {
	switch compressorType {
	case GATEWAY_COMPRESS_NONE:
		return &NoCompressor{}
	case GATEWAY_COMPRESS_ZLIB:
		return &ZlibCompressor{}
	case GATEWAY_COMPRESS_ZSTD:
		return &ZstdCompressor{}
	default:
		panic(fmt.Errorf("invalid compression type"))
	}
}