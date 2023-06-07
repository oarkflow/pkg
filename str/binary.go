package str

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/klauspost/compress/gzip"
)

func ToCompressedString(data []byte) string {
	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	_, _ = w.Write(data)
	_ = w.Close()
	return base64.StdEncoding.EncodeToString(compressed.Bytes())
}

func FromCompressedString(data string) ([]byte, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	zipReader, err := gzip.NewReader(bytes.NewReader(decodedBytes))
	if err != nil {
		return nil, err
	}
	rawBytes, err := io.ReadAll(zipReader)
	if err != nil {
		return nil, err
	}
	return rawBytes, nil
}

func GenerateBinaryContent(packageName, varName string, data []byte, fileName ...string) []byte {
	encoded := ToCompressedString(data)
	output := &bytes.Buffer{}
	output.WriteString("package " + packageName + "\n\nvar " + varName + " = " + strconv.Quote(encoded) + "\n")
	bt := output.Bytes()
	if len(fileName) > 0 {
		writeFile(fileName[0], bt)
	}
	return bt
}

func DecodeBinaryString(data string) ([]byte, error) {
	return FromCompressedString(data)
}

func writeFile(filePath string, data []byte) {
	err := os.WriteFile(filePath, data, os.FileMode(0o664))
	if err != nil {
		log.Fatalf("Error writing '%s': %s", filePath, err)
	}
}
