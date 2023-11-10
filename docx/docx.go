package docx

import (
	"bytes"
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

type Argument struct {
	Type        string `json:"type"`
	Placeholder string `json:"placeholder"`
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
			ptypes = append(ptypes, parseFunction(placeholder))
		} else {
			ptypes = append(ptypes, PType{Type: getType(placeholder), Placeholder: placeholder})
		}
	}
	return ptypes, nil
}

func parseFunction(funcStr string) PType {
	p := PType{
		Type:        "function",
		Placeholder: funcStr,
	}
	funcParts := strings.Split(funcStr, "(")
	if len(funcParts) != 2 {
		return p
	}
	p.Placeholder = funcParts[0]
	arguments := strings.Split(funcParts[1], ")")
	if len(arguments) != 2 {
		return p
	}
	argStr := arguments[0]
	if argStr == "" {
		return p
	}
	args := strings.Split(argStr, ",")
	for _, arg := range args {
		p.Arguments = append(p.Arguments, Argument{
			Type:        getType(arg),
			Placeholder: arg,
		})
	}
	return p
}

func getType(arg string) string {
	if strings.HasPrefix(arg, "'") || strings.HasPrefix(arg, `"`) {
		return "string"
	} else if IsNumber(arg) {
		return "number"
	}
	return "variable"
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
