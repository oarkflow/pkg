package docx

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Contains functions to work with data from a zip file
type ZipData interface {
	files() []*zip.File
	close() error
}

// Type for in memory zip files
type ZipInMemory struct {
	data *zip.Reader
}

func (d ZipInMemory) files() []*zip.File {
	return d.data.File
}

// Since there is nothing to close for in memory, just nil the data and return nil
func (d ZipInMemory) close() error {
	d.data = nil
	return nil
}

// Type for zip files read from disk
type ZipFile struct {
	data *zip.ReadCloser
}

func (d ZipFile) files() []*zip.File {
	return d.data.File
}

func (d ZipFile) close() error {
	return d.data.Close()
}

type ReplaceDocx struct {
	zipReader ZipData
	content   string
	links     string
	headers   map[string]string
	footers   map[string]string
	images    map[string]string
	imgIndex  int32
	refIndex  int32
}

func (r *ReplaceDocx) Editable() *Docx {
	return &Docx{
		files:      r.zipReader.files(),
		content:    r.content,
		links:      r.links,
		headers:    r.headers,
		footers:    r.footers,
		images:     r.images,
		appendFile: make(map[string][]byte),
		imgIndex:   r.imgIndex,
		refIndex:   r.refIndex,
	}
}

func (r *ReplaceDocx) Close() error {
	return r.zipReader.close()
}

type Docx struct {
	files      []*zip.File
	content    string
	links      string
	headers    map[string]string
	footers    map[string]string
	images     map[string]string
	imgIndex   int32
	refIndex   int32
	appendFile map[string][]byte
}

func (d *Docx) GetContent() string {
	return d.content
}

func (d *Docx) SetContent(content string) {
	d.content = content
}

func (d *Docx) ReplaceRaw(oldString string, newString string, num int) {
	d.content = strings.Replace(d.content, oldString, newString, num)
}

func (d *Docx) ReplaceTagContains(oldTag string, oldstring string, newString string) error {
	reg, err := regexp.Compile(`<` + oldTag + `>.*?</\` + oldTag + ">")
	if err != nil {
		return err
	}
	for _, content := range reg.FindAllString(d.content, -1) {
		if strings.Contains(content, oldstring) {
			d.content = strings.Replace(d.content, content, newString, 1)
		}
	}
	return nil
}

func (d *Docx) ReplaceStartEnd(start, end string, new string, num int) {
	indexstart := strings.Index(d.content, start)
	indexend := strings.Index(d.content, end)
	if indexend == -1 || indexstart == -1 {
		return
	}
	indexend += len(end)
	d.content = strings.Replace(d.content, d.content[indexstart:indexend], new, num)
}

func (d *Docx) Replace(oldString string, newString string, num int) (err error) {
	oldString, err = encode(oldString)
	if err != nil {
		return err
	}
	newString, err = encode(newString)
	if err != nil {
		return err
	}
	d.content = strings.Replace(d.content, oldString, newString, num)

	return nil
}

func (d *Docx) ReplaceLink(oldString string, newString string, num int) (err error) {
	oldString, err = encode(oldString)
	if err != nil {
		return err
	}
	newString, err = encode(newString)
	if err != nil {
		return err
	}
	d.links = strings.Replace(d.links, oldString, newString, num)

	return nil
}

func (d *Docx) ReplaceHeader(oldString string, newString string) (err error) {
	return replaceHeaderFooter(d.headers, oldString, newString)
}

func (d *Docx) ReplaceFooter(oldString string, newString string) (err error) {
	return replaceHeaderFooter(d.footers, oldString, newString)
}

func (d *Docx) WriteToFile(path string) (err error) {
	var target *os.File
	target, err = os.Create(path)
	if err != nil {
		return
	}
	defer target.Close()
	err = d.Write(target)
	return
}

func (d *Docx) Write(ioWriter io.Writer) (err error) {
	w := zip.NewWriter(ioWriter)
	defer w.Close()
	for _, file := range d.files {
		var writer io.Writer
		var readCloser io.ReadCloser

		writer, err = w.Create(file.Name)
		if err != nil {
			return err
		}
		readCloser, err = file.Open()
		if err != nil {
			return err
		}
		if file.Name == "word/document.xml" {
			writer.Write([]byte(d.content))
		} else if file.Name == "word/_rels/document.xml.rels" {
			writer.Write([]byte(d.links))
		} else if strings.Contains(file.Name, "header") && d.headers[file.Name] != "" {
			writer.Write([]byte(d.headers[file.Name]))
		} else if strings.Contains(file.Name, "footer") && d.footers[file.Name] != "" {
			writer.Write([]byte(d.footers[file.Name]))
		} else if strings.HasPrefix(file.Name, "word/media/") && d.images[file.Name] != "" {
			newImage, err := os.Open(d.images[file.Name])
			if err != nil {
				return err
			}
			writer.Write(streamToByte(newImage))
			newImage.Close()
		} else {
			writer.Write(streamToByte(readCloser))
		}
	}
	for fileName, appendFile := range d.appendFile {
		writer, err := w.Create(fileName)
		if err != nil {
			return err
		}
		writer.Write([]byte(appendFile))
	}
	return
}

func (d *Docx) AddPic(dir string) (string, string, error) {
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		return "", "", err
	}
	file, err := ioutil.ReadFile(dir)
	if err != nil {
		return "", "", err
	}
	d.imgIndex++
	d.refIndex++
	idIndx := strconv.Itoa(int(d.imgIndex))
	refIndx := strconv.Itoa(int(d.refIndex))
	fileName := d.fileName(dir)
	imgName := "image" + idIndx + d.fileSuffix(fileName)
	id := "rId" + refIndx
	sb := strings.Builder{}
	sb.WriteString(`<Relationship Id="rId`)
	sb.WriteString(idIndx)
	sb.WriteString(`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/`)
	sb.WriteString(imgName)
	sb.WriteString(`"/>`)
	d.links += sb.String()
	d.appendFile["word/media/"+imgName] = file
	return refIndx, id, err
}

func (d *Docx) AddPicStream(buf []byte, suffix string) (string, string) {
	d.imgIndex++
	d.refIndex++
	idIndx := strconv.Itoa(int(d.imgIndex))
	refIndx := strconv.Itoa(int(d.refIndex))
	imgName := "image" + idIndx + suffix
	id := "rId" + refIndx

	sb := strings.Builder{}
	sb.WriteString(`<Relationship Id="rId`)
	sb.WriteString(refIndx)
	sb.WriteString(`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/`)
	sb.WriteString(imgName)
	sb.WriteString(`"/></Relationships>`)
	d.links = strings.Replace(d.links, "</Relationships>", sb.String(), 1)
	d.appendFile["word/media/"+imgName] = buf
	return refIndx, id
}

func (d *Docx) fileName(dir string) string {
	list := strings.Split(dir, string(os.PathSeparator))
	return list[len(list)-1]
}

func (d *Docx) fileSuffix(file string) string {
	list := strings.Split(file, ".")
	if len(list) > 1 {
		return "." + list[len(list)-1]
	}
	return ""
}

func replaceHeaderFooter(headerFooter map[string]string, oldString string, newString string) (err error) {
	oldString, err = encode(oldString)
	if err != nil {
		return err
	}
	newString, err = encode(newString)
	if err != nil {
		return err
	}

	for k := range headerFooter {
		headerFooter[k] = strings.Replace(headerFooter[k], oldString, newString, -1)
	}

	return nil
}

// ReadDocxFromFS opens a docx file from the file system
func ReadDocxFromFS(file string, fs fs.FS) (*ReplaceDocx, error) {
	f, err := fs.Open(file)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buff, f)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(buff.Bytes())
	return ReadDocxFromMemory(reader, size)
}

func ReadDocxFromMemory(data io.ReaderAt, size int64) (*ReplaceDocx, error) {
	reader, err := zip.NewReader(data, size)
	if err != nil {
		return nil, err
	}
	zipData := ZipInMemory{data: reader}
	return ReadDocx(zipData)
}

func ReadDocxFile(path string) (*ReplaceDocx, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	zipData := ZipFile{data: reader}
	return ReadDocx(zipData)
}

func ReadDocx(reader ZipData) (*ReplaceDocx, error) {
	content, err := readText(reader.files())
	if err != nil {
		return nil, err
	}

	refid, imgid, links, err := readLinks(reader.files())
	if err != nil {
		return nil, err
	}

	headers, footers, _ := readHeaderFooter(reader.files())
	images, _ := retrieveImageFilenames(reader.files())
	return &ReplaceDocx{zipReader: reader, content: content, links: links, headers: headers, footers: footers, images: images, imgIndex: int32(imgid), refIndex: int32(refid)}, nil
}

func retrieveImageFilenames(files []*zip.File) (map[string]string, error) {
	images := make(map[string]string)
	for _, f := range files {
		if strings.HasPrefix(f.Name, "word/media/") {
			images[f.Name] = ""
		}
	}
	return images, nil
}

func readHeaderFooter(files []*zip.File) (headerText map[string]string, footerText map[string]string, err error) {

	h, f, err := retrieveHeaderFooterDoc(files)

	if err != nil {
		return map[string]string{}, map[string]string{}, err
	}

	headerText, err = buildHeaderFooter(h)
	if err != nil {
		return map[string]string{}, map[string]string{}, err
	}

	footerText, err = buildHeaderFooter(f)
	if err != nil {
		return map[string]string{}, map[string]string{}, err
	}

	return headerText, footerText, err
}

func buildHeaderFooter(headerFooter []*zip.File) (map[string]string, error) {

	headerFooterText := make(map[string]string)
	for _, element := range headerFooter {
		documentReader, err := element.Open()
		if err != nil {
			return map[string]string{}, err
		}

		text, err := wordDocToString(documentReader)
		if err != nil {
			return map[string]string{}, err
		}

		headerFooterText[element.Name] = text
	}

	return headerFooterText, nil
}

func readText(files []*zip.File) (text string, err error) {
	var documentFile *zip.File
	documentFile, err = retrieveWordDoc(files)
	if err != nil {
		return text, err
	}
	var documentReader io.ReadCloser
	documentReader, err = documentFile.Open()
	if err != nil {
		return text, err
	}

	text, err = wordDocToString(documentReader)
	return
}

var imgReg = regexp.MustCompile(`"media/image\d+`)
var refReg = regexp.MustCompile(`"rId\d+`)

func readLinks(files []*zip.File) (refid, imgid int, text string, err error) {
	var documentFile *zip.File
	documentFile, err = retrieveLinkDoc(files)
	if err != nil {
		return
	}
	var documentReader io.ReadCloser
	documentReader, err = documentFile.Open()
	if err != nil {
		return
	}

	text, err = wordDocToString(documentReader)
	for _, v := range imgReg.FindAllString(text, -1) {
		id := strings.TrimPrefix(v, `"media/image`)
		mid, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		if mid > imgid {
			imgid = mid
		}
	}
	for _, v := range refReg.FindAllString(text, -1) {
		id := strings.TrimPrefix(v, `"rId`)
		mid, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		if mid > refid {
			refid = mid
		}
	}
	return
}

func wordDocToString(reader io.Reader) (string, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	data := string(b)

	return data, nil
}

func retrieveWordDoc(files []*zip.File) (file *zip.File, err error) {
	for _, f := range files {
		if f.Name == "word/document.xml" {
			file = f
		}
	}
	if file == nil {
		err = errors.New("document.xml file not found")
	}
	return
}

func retrieveLinkDoc(files []*zip.File) (file *zip.File, err error) {
	for _, f := range files {
		if f.Name == "word/_rels/document.xml.rels" {
			file = f
		}
	}
	if file == nil {
		err = errors.New("document.xml.rels file not found")
	}
	return
}

func retrieveHeaderFooterDoc(files []*zip.File) (headers []*zip.File, footers []*zip.File, err error) {
	for _, f := range files {

		if strings.Contains(f.Name, "header") {
			headers = append(headers, f)
		}
		if strings.Contains(f.Name, "footer") {
			footers = append(footers, f)
		}
	}
	if len(headers) == 0 && len(footers) == 0 {
		err = errors.New("headers[1-3].xml file not found and footers[1-3].xml file not found")
	}
	return
}

func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

// To get Word to recognize a tab character, we have to first close off the previous
// text element.  This means if there are multiple consecutive tabs, there are empty <w:t></w:t>
// in between but it still seems to work correctly in the output document, certainly better
// than other combinations I tried.
const TAB = "</w:t><w:tab/><w:t>"
const NEWLINE = "<w:br/>"

func encode(s string) (string, error) {
	var b bytes.Buffer
	enc := xml.NewEncoder(bufio.NewWriter(&b))
	if err := enc.Encode(s); err != nil {
		return s, err
	}
	output := strings.Replace(b.String(), "<string>", "", 1) // remove string tag
	output = strings.Replace(output, "</string>", "", 1)
	output = strings.Replace(output, "&#xD;&#xA;", NEWLINE, -1) // \r\n (Windows newline)
	output = strings.Replace(output, "&#xD;", NEWLINE, -1)      // \r (earlier Mac newline)
	output = strings.Replace(output, "&#xA;", NEWLINE, -1)      // \n (unix/linux/OS X newline)
	output = strings.Replace(output, "&#x9;", TAB, -1)          // \t (tab)
	return output, nil
}

func (d *Docx) ReplaceImage(oldImage string, newImage string) (err error) {
	if _, ok := d.images[oldImage]; ok {
		d.images[oldImage] = newImage
		return nil
	}
	return fmt.Errorf("old image: %q, file not found", oldImage)
}

func (d *Docx) ImagesLen() int {
	return len(d.images)
}
