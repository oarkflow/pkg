package docx

import (
	"bytes"
	"fmt"
	"regexp"
	"sync"
)

type TableReplacer struct {
	document     []byte
	placeholders []*Placeholder
	distinctRuns []*Run // slice of all distinct runs extracted from the placeholders used for validation
	ReplaceCount int
	BytesChanged int64
	mu           sync.Mutex
}

type TablePlaceholder struct {
	Key   string
	Value string
}

func NewTableReplacer(docBytes []byte) *TableReplacer {
	r := &TableReplacer{
		document:     docBytes,
		ReplaceCount: 0,
	}

	return r
}

func (r *TableReplacer) Replace(tablePrefix string, placeholders [][]TablePlaceholder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	found := false
	foundRowTemplate := false
	firstPlaceholder := placeholders[0]
	// Matches any XML table containing the first tag.
	// Table must have two rows, where one is header and the other is a template for the rest of the rows.
	// Header will remain, template will be replaced with the rest of the rows.
	re := regexp.MustCompile(fmt.Sprintf(`(?s)<w:tbl>.*?<w:tr>.*?</w:tr>.*?(<w:tr>.*?\[%s\.%s\].*?</w:tr>).*?</w:tbl>`, tablePrefix, firstPlaceholder[0].Key))
	row := re.FindSubmatch(r.document)
	if len(row) == 2 {
		foundRowTemplate = true
		found = true
	}

	rowTemplate := row[1]

	if !foundRowTemplate {
		return fmt.Errorf("row template not found")
	}

	outputRows := []byte{}
	for i := 0; i < len(placeholders); i++ {
		outputRow := rowTemplate
		for j := 0; j < len(placeholders[i]); j++ {
			outputRow = bytes.Replace(outputRow, []byte(fmt.Sprintf("[%s.%s]", tablePrefix, placeholders[i][j].Key)), []byte(placeholders[i][j].Value), -1)
		}
		outputRows = append(outputRows, outputRow...)
	}

	r.document = bytes.Replace(r.document, rowTemplate, outputRows, -1)

	// all replacing actions might potentially screw up the XML structure
	// in order to capture this, all tags are re-validated after replacing a value
	if err := ValidatePositions(r.document, r.distinctRuns); err != nil {
		return fmt.Errorf("replace produced invalid result: %w", err)
	}

	if !found {
		return ErrPlaceholderNotFound
	}
	return nil
}

// Bytes returns the document bytes.
// If called after Replace(), the bytes will be modified.
func (r *TableReplacer) Bytes() []byte {
	return r.document
}
