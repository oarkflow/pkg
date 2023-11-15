package replacer

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

func encodeT(s string) (string, error) {
	var b bytes.Buffer
	enc := xml.NewEncoder(bufio.NewWriter(&b))
	if err := enc.Encode(s); err != nil {
		return s, err
	}
	fmt.Println(b.String())
	output := strings.Replace(b.String(), "<string>", "", 1) // remove string tag
	output = strings.Replace(output, "</string>", "", 1)
	output = strings.Replace(output, "&#xD;&#xA;", "<w:br/>", -1) // \r\n => newline
	return output, nil
}
func (d *Docx) ParseNode() (string, string, string) {
	bodySplit := strings.Split(d.content, "<w:body>")
	header := bodySplit[0]
	bodySplitAfter := strings.Split(bodySplit[1], "</w:body>")
	body := bodySplitAfter[0]
	footer := bodySplitAfter[1]
	return header, body, footer
}
func (d *Docx) GlueNodes(header, body, footer string) {
	d.content = header + "<w:body>" + body + "</w:body>" + footer
	//	fmt.Println(d.content)
}
func (d *Docx) GetFirstElementContain(word string, body []string) (int, error) {
	for i, item := range body {
		if strings.Contains(item, word) {
			return i, nil
		}
	}
	return 0, fmt.Errorf("element '%s' not found", word)
}

func (d *Docx) AddBlockAtTheEnd(element string, body []string) []string {

	return d.AddBlockAfterElement(len(body), element, body)
}
func (d *Docx) AddBlockAtTheBeginning(element string, body []string) []string {

	return d.AddBlockBeforeElement(0, element, body)
}
func (d *Docx) AddBlockAfterElement(elemNum int, element string, body []string) []string {
	f := len(body)
	if f == 0 {
		var el []string
		el = append(el, element)
		return el
	}
	if f < elemNum {
		body = append(body, element)
		return body
	}
	var newBody []string
	for i := 0; i < f; i++ {

		newBody = append(newBody, body[i])
		if i == elemNum-1 {
			newBody = append(newBody, element)
		}
	}
	return newBody
}
func (d *Docx) AddBlockBeforeElement(elemNum int, element string, body []string) []string {
	f := len(body)
	if f == 0 {
		var el []string
		el = append(el, element)
		return el
	}
	if f < elemNum {
		body = append(body, element)
		return body
	}
	var newBody []string
	for i := 0; i < f; i++ {
		if i == elemNum {
			newBody = append(newBody, element)
		}
		newBody = append(newBody, body[i])
	}
	return newBody
}
func (d *Docx) BodyParse(body string) []string {
	var bodyElems []string
	bsplit := strings.Split(body, "<w:p>")
	for _, item := range bsplit {
		parag := strings.Split(item, "</w:p>")
		bodyElems = append(bodyElems, "<w:p>"+parag[0]+"</w:p>")
		if len(parag) < 1 {
			if len(parag[1]) != 0 {
				bodyElems = append(bodyElems, parag[1])
			}
		}
	}
	return bodyElems
}
func (d *Docx) BodyGlue(body []string) string {
	full := ""
	for _, item := range body {
		full = full + item
	}
	return full
}
func (d *Docx) ReplaceWithTag(oldString string, newString string) (err error) {
	oldString, err = encodeT(oldString)
	if err != nil {
		return err
	}
	newString, err = encodeT(newString)
	if err != nil {
		return err
	}
	d.content = strings.Replace(d.content, "<w:t>{div id=&apos;1&apos;}{/div}</w:t>", `
	<w:p>
            <w:pPr>
                <w:pStyle w:val="Normal"/>
                <w:rPr>
                    <w:rFonts w:ascii="Droid Sans Mono;monospace;monospace;Droid Sans Fallback;noto-fonts-emoji-apple;Droid Sans Mono;monospace;monospace;Droid Sans Fallback" w:hAnsi="Droid Sans Mono;monospace;monospace;Droid Sans Fallback;noto-fonts-emoji-apple;Droid Sans Mono;monospace;monospace;Droid Sans Fallback" w:cstheme="minorBidi"/>
                    <w:b w:val="false"/>
                    <w:color w:val="000000"/>
                    <w:sz w:val="21"/>
                    <w:shd w:fill="auto" w:val="clear"/>
                </w:rPr>
            </w:pPr>
            <w:r>
                <w:rPr>
                    <w:rFonts w:asciiTheme="minorHAnsi" w:cstheme="minorBidi" w:eastAsiaTheme="minorHAnsi" w:hAnsiTheme="minorHAnsi" w:ascii="Droid Sans Mono;monospace;monospace;Droid Sans Fallback;noto-fonts-emoji-apple;Droid Sans Mono;monospace;monospace;Droid Sans Fallback" w:hAnsi="Droid Sans Mono;monospace;monospace;Droid Sans Fallback;noto-fonts-emoji-apple;Droid Sans Mono;monospace;monospace;Droid Sans Fallback"/>
                    <w:b w:val="false"/>
                    <w:color w:val="000000"/>
                    <w:sz w:val="21"/>
                    <w:shd w:fill="auto" w:val="clear"/>
                </w:rPr>
                <w:t>{div id=&apos;1&apos;}{/div}</w:t>
            </w:r>
        </w:p>
	
	`, -1)

	return nil
}
