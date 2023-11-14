package replacer

import (
	"fmt"
	"strings"
)

func (d *Docx) CreateStructedBlock(simpleword string, arr ...Font) Block {
	block := Block{
		Head: Font{
			FontSize: 15,
			FontName: "Calibri",
			Bold:     false,
			Italic:   false,
			Strike:   false,
			Color:    Black,
			Another:  "",
		},
		Body:   simpleword,
		Footer: "</w:t></w:r></w:p>",
	}
	arg := arr[0]
	if len(arr) != 0 {
		if arg.FontSize != 0 {
			block.Head.FontSize = arg.FontSize
		}
		if arg.FontName != "" {
			block.Head.FontName = arg.FontName
		}
		if arg.Bold {
			block.Head.Bold = true
		}
		if arg.Italic {
			block.Head.Italic = true
		}
		if arg.Strike {
			block.Head.Strike = true
		}
		if arg.Color != "" { // add check is color in allowed list
			block.Head.Color = arg.Color
		}
	}
	return block
}

//	BlockToString func (d *Docx) BlockToString(block Block) string {
//		return fmt.Sprintf("%s%s%s", block.Head, block.Body, block.Footer)
//	}
func (d *Docx) BlockToString(block Block) string {
	header := ""
	header += fmt.Sprintf("<w:sz w:val=\"%d\"/><w:szCs w:val=\"%d\"/>", block.Head.FontSize*2, block.Head.FontSize*2)
	header += fmt.Sprintf(`<w:rFonts w:eastAsia="%s"/>`, block.Head.FontName)
	header += fmt.Sprintf(` <w:color w:val="%s"/>`, block.Head.Color)
	header += block.Head.Another
	if block.Head.Bold {
		header += "<w:b/><w:bCs/>"
	}
	if block.Head.Italic {
		header += "<w:i/><w:iCs/>"
	}
	if block.Head.Strike {
		header += `<w:u w:val="single"/>`
	}
	return fmt.Sprintf("<w:p><w:r><w:rPr>%s </w:rPr><w:t>%s%s", header, block.Body, block.Footer)
}

// CreateNewBlock <w:b/>
//
//	<w:bCs/>
//	<w:i/>
//	<w:iCs/>
//	<w:color w:val="F10D0C"/>
//	<w:sz w:val="36"/>
//	<w:szCs w:val="36"/>
//	<w:u w:val="single"/>
func (d *Docx) CreateNewBlock(simpleword string) string {
	block := fmt.Sprintf("<w:p><w:pPr><w:pStyle w:val=\"Normal\"/><w:rPr></w:rPr></w:pPr><w:r><w:rPr></w:rPr><w:t>%s</w:t></w:r></w:p>", simpleword)
	return block
}
func (d *Docx) parseHeader(p string) Font {
	fmt.Printf("header: %s\n", p)
	return Font{}
}
func (d *Docx) ParseBlockToStruct(s string) Block {
	f := strings.Split(s, "<w:t>")
	bf := strings.Split(f[1], "</w:t>")
	//	header := f[0]
	header := d.parseHeader(strings.Split(f[0], "<w:p>")[1])
	body := bf[0]
	footer := bf[1]
	return Block{
		Head:   header,
		Body:   body,
		Footer: "</w:t>" + footer,
	}
}
func (d *Docx) EditTextInStringBlock(block string, newText string) string {
	p := d.ParseBlockToStruct(block)
	p.Body = newText
	return d.BlockToString(p)
}
func (d *Docx) GetTextInBlock(s string) string {
	return d.ParseBlockToStruct(s).Body
}
func (d *Docx) ReplaceTextInBlock(old, new string, body []string) []string {
	id, err := d.GetFirstElementContain(old, body)
	if err == nil {
		body[id] = d.EditTextInStringBlock(body[id], new)
	}
	return body
}
func (d *Docx) EditStructedBlockParams(block Block) Block {

	return block
}

func (d *Docx) ReplaceBlockToBlocks(bodypart []string, id int, blocks []Block) []string {
	return []string{}
}
