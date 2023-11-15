package replacer

import (
	"fmt"
	"strings"

	"github.com/oarkflow/pkg/docx/replacer/ooxml"
)

const OpenStatus = 4
const Self = 5

func wpParser(s string) ([]WPTokens, error) {
	nodes, err := ooxml.GetParentNodes(s)
	if err != nil {
		fmt.Println(err)
		return []WPTokens{}, err
	}
	_ = nodes
	var d []WPTokens
	for _, item := range nodes {
		status := Self
		switch item.TagStatus {
		case ooxml.TagComplimentary:
			status = OpenStatus
		case ooxml.TagSelfClosed:
			status = Self
		}
		d = append(d, WPTokens{
			Tag:    item.Name,
			Body:   item.Body,
			Attr:   item.Args,
			Status: status,
		})
	}
	return d, nil
}
func (d *Docx) Parser() (Document, error) {
	_, body, _ := d.ParseNode()
	nodes, err := ooxml.GetParentNodes(body)
	if err != nil {
		return Document{}, err
	}
	var doc Document
	for _, item := range nodes {
		switch item.Name {
		case "w:p":
			body, err := wpParser(item.Body)
			if err != nil {
				return Document{}, err
			}
			doc.WP = append(doc.WP, WP{
				Tag:  item.Name,
				Body: body,
			})
		case "w:sectPr":
			doc.SectPr = SectPr{
				Tag:  item.Name,
				Body: item.Body,
			}
		}

	}
	return doc, nil
}

func (d *Document) AddNewBlock(s string) {
	d.WP = append(d.WP, WP{
		Tag: "w:p",
		Body: []WPTokens{
			{
				Tag:    "w:pPr",
				Body:   `<w:pStyle w:val="Normal"/><w:rPr><w:rFonts w:ascii="Calibri" w:hAnsi="Calibri" w:eastAsia="Calibri" w:cs="" w:asciiTheme="minorHAnsi" w:cstheme="minorBidi" w:eastAsiaTheme="minorHAnsi" w:hAnsiTheme="minorHAnsi"/><w:color w:val="00000A"/><w:sz w:val="24"/><w:szCs w:val="24"/><w:lang w:val="en-US" w:eastAsia="en-US" w:bidi="ar-SA"/></w:rPr>`,
				Status: OpenStatus,
			},
			{
				Tag:    "w:r",
				Body:   fmt.Sprintf("<w:rPr></w:rPr><w:t>%s</w:t>", s),
				Status: OpenStatus,
			},
		},
	})
}

func (d *Document) AppendWPBlockInToEnd(block WP) {
	d.WP = append(d.WP, block)
}

func (d *Document) EditBlockWithNewLine(oldTag, newString string) error {
	id, err := d.GetBlockIDByTag(oldTag)
	if err != nil {
		// t.Error(err)
		return err
	}
	for i, itemD := range d.WP[id].Body {
		if itemD.Tag == "w:r" {
			if strings.Contains(itemD.Body, oldTag) {
				ex := d.WP[id].Body[i]
				tags := strings.Split(newString, "\n")
				var tempArray []WPTokens
				for _, tagI := range tags {
					if tagI != "" {
						//	fmc.Printfln("tag:[%v]", tagI)
						z := ex
						z.Body = "<w:t>" + tagI + "</w:t><w:br/>"
						tempArray = append(tempArray, z)
					}
				}
				if len(d.WP[id].Body) == 1 {
					d.WP[id].Body = tempArray
					return nil
				}
				right := d.WP[id].Body[i:]
				left := d.WP[id].Body[:i]
				d.WP[id].Body = append(left, tempArray...)
				d.WP[id].Body = append(d.WP[id].Body, right...)
				return nil
			}
		}
	}
	return fmt.Errorf("tag not found")
}

func (d *Document) GetBlockByID(id int) WP {
	return d.WP[id]
}

func (d *Docx) Compile(path string, doc Document) error {
	head, _, footer := d.ParseNode()
	body := "<w:body>" + doc.BodyToString() + "</w:body>"
	xml := head + body + footer
	d.content = xml
	return d.WriteToFile(path)
}
func (d *Document) BodyToString() string {
	var body string
	for _, item := range d.WP {
		body += wpTokenToString(item)
	}
	body += "<" + d.SectPr.Tag + ">" + d.SectPr.Body + "</" + d.SectPr.Tag + ">"
	return body
}
func wpTokenToString(item WP) string {
	// item.Tag
	var body string
	for _, it := range item.Body {
		var attr string
		if 0 < len(it.Attr) && string(it.Attr[0]) != " " {
			attr = " " + it.Attr
		} else {
			attr = it.Attr
		}
		if it.Status == OpenStatus {
			body += "<" + it.Tag + attr + ">" + it.Body + "</" + it.Tag + ">"
		}
		if it.Status == Self {
			body += "<" + it.Tag + attr + "/>"
		}

	}
	return "<w:p>" + body + "</w:p>"
}
func Screening(s string) string {
	var retString string
	for _, item := range s {
		it := string(item)
		switch it {
		case "<":
			retString += "&lt;"
		case ">":
			retString += "&gt;"
		case "&":
			retString += "&amp;"
		case "'":
			retString += "&apos;"
		default:
			retString += it
		}
	}
	return retString
}
