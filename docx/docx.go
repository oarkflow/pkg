package docx

import (
	"bytes"
	"go/parser"
	"path/filepath"
	"strings"
)

func PrepareDocxToFile(file string, data map[string]interface{}, outputFile ...string) error {
	var output string
	if len(outputFile) == 0 {
		output = strings.Replace(filepath.Base(file), filepath.Ext(file), "-Filled", -1) + filepath.Ext(file)
	} else {
		output = outputFile[0]
	}
	doc, err := Open(file)
	if err != nil {
		return err
	}
	err = doc.ReplaceAll(data)
	if err != nil {
		return err
	}

	return doc.WriteToFile(output)
}

type PType struct {
	Type        string     `json:"type"`
	Placeholder string     `json:"placeholder"`
	Arguments   []Argument `json:"arguments"`
}

func Placeholders(file string) ([]PType, error) {
	var ptypes []PType
	doc, err := Open(file)
	if err != nil {
		return nil, err
	}
	placeholders, err := doc.GetPlaceHoldersList()
	if err != nil {
		return nil, err
	}
	for _, placeholder := range placeholders {
		placeholder = strReplacer.Replace(placeholder)
		if IsFunction(placeholder) {
			functions, err := ParseExpr(placeholder)
			if err != nil {
				return nil, err
			}
			for _, f := range functions {
				ptypes = append(ptypes, PType{Type: "function", Placeholder: f.Name, Arguments: f.Arguments})
			}
		} else {
			node, err := parser.ParseExpr(placeholder)
			if err != nil {
				return nil, err
			}
			arg := extractArg(node)
			ptypes = append(ptypes, PType{Type: arg.Type, Placeholder: placeholder})
		}
	}
	return ptypes, nil
}

func PrepareDocx(file string, data map[string]interface{}) (*bytes.Buffer, error) {
	var byteBuffer bytes.Buffer
	doc, err := Open(file)
	if err != nil {
		return nil, err
	}
	err = doc.ReplaceAll(data)
	if err != nil {
		return nil, err
	}
	err = doc.Write(&byteBuffer)
	if err != nil {
		return nil, err
	}
	return &byteBuffer, nil
}
