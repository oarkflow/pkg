package search

import (
	"bytes"
	"fmt"
	"io"

	"github.com/klauspost/compress/gzip"
)

func Compress(data []byte) ([]byte, error) {
	var compressed bytes.Buffer
	compressor := gzip.NewWriter(&compressed)
	// Compress the string
	_, err := compressor.Write(data)
	if err != nil {
		fmt.Println("Error compressing string:", err)
		return nil, err
	}
	compressor.Close()
	return compressed.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {
	// Decompress and print the original string
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		fmt.Println("Error reading string:", err)
		return nil, err
	}
	d, err := io.ReadAll(r)
	if err != nil {
		fmt.Println("Error decompressing string:", err)
		return nil, err
	}
	r.Close()
	return d, nil
}
