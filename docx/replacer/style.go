package replacer

import (
	"fmt"
	"strings"
)

type MarkerParams struct {
}

func returnTemplateMarkerWP() WP {
	return WP{
		Tag: "w:p",
		Body: []WPTokens{
			{
				Tag: "w:pPr",
				Body: AtomicWPTokensToString(WPTokens{
					Tag:    "w:pStyle",
					Attr:   ` w:val="Normal"`,
					Body:   "",
					Status: Self,
				}) +
					AtomicWPTokensToString(WPTokens{
						Tag:    "w:numPr",
						Attr:   ``,
						Body:   `<w:ilvl w:val="0"/><w:numId w:val="1"/>`,
						Status: OpenStatus,
					}) +
					AtomicWPTokensToString(WPTokens{
						Tag:    "w:bidi",
						Attr:   ` w:val="0"`,
						Body:   "",
						Status: Self,
					}) +
					AtomicWPTokensToString(WPTokens{
						Tag:    "w:jc",
						Attr:   ` w:val="left"`,
						Body:   "",
						Status: Self,
					}) +
					AtomicWPTokensToString(WPTokens{
						Tag:    "w:rPr",
						Attr:   ``,
						Body:   "",
						Status: OpenStatus,
					}),
				Attr:   ``,
				Status: OpenStatus,
			},
			{
				Tag:    "w:r",
				Body:   "<w:rPr></w:rPr><w:t>%s</w:t>",
				Attr:   ``,
				Status: OpenStatus,
			},
		},
	}
}
func (d *Document) CreateMarkedStringList(mp MarkerParams, letter ...string) []WP {
	if mp != (MarkerParams{}) {
		// With Params
		return []WP{returnTemplateMarkerWP()}
	}
	var wpArray []WP
	for _, item := range letter {
		// fmc.Printfln("ranger: %s", item)
		if strings.Contains(item, "\n") {
			arr := strings.Split(item, "\n")
			for _, arrItem := range arr {
				tempBlock := returnTemplateMarkerWP()
				tempBlock.Body[1].Body = fmt.Sprintf(tempBlock.Body[1].Body, Screening(arrItem))
				wpArray = append(wpArray, tempBlock)
			}
		} else {
			// tempBlock := templateMarkerWP
			tempBlock := returnTemplateMarkerWP()
			tempBlock.Body[1].Body = fmt.Sprintf(tempBlock.Body[1].Body, Screening(item))
			wpArray = append(wpArray, tempBlock)
		}
	}
	return wpArray
}
func AtomicWPTokensToString(token WPTokens) string {
	var attr string
	var body string
	if 0 < len(token.Attr) && string(token.Attr[0]) != " " {
		attr = " " + token.Attr
	} else {
		attr = token.Attr
	}
	if token.Status == OpenStatus {
		body += "<" + token.Tag + attr + ">" + token.Body + "</" + token.Tag + ">"
	}
	if token.Status == Self {
		body += "<" + token.Tag + attr + "/>"
	}
	return body
}
func GetTextFromXML(src string) (string, error) {
	res, err := wpParser(src)
	if err != nil {
		return "", err
	}
	var text string
	for _, item := range res {
		if item.Tag == "w:r" {
			res, err := wpParser(item.Body)
			if err != nil {
				return "", err
			}
			for _, wtTag := range res {
				if wtTag.Tag == "w:t" {
					text += wtTag.Body
				}
			}
		}
	}
	return text, nil
}
