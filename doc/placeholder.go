package doc

import (
	"bytes"

	"github.com/klauspost/compress/zip"

	"github.com/oarkflow/pkg/jet"
	"github.com/oarkflow/pkg/str"
)

func Replace(filename string, data map[string]any) error {
	zipReader, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	var byteBuffer bytes.Buffer

	zipWriter := zip.NewWriter(&byteBuffer)
	defer zipWriter.Close()
	for _, file := range zipReader.File {

		writer, err := zipWriter.Create(file.Name)
		if err != nil {
			return err
		}

		readCloser, err := file.Open()
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		buf.ReadFrom(readCloser)
		newContent, err := jet.Parse(str.FromByte(buf.Bytes()), data)
		if err != nil {
			panic(err)
			return err
		}
		writer.Write(str.ToByte(newContent))
	}
	return nil
}
