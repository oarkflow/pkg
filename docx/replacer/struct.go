package replacer

type Block struct {
	Head   Font
	Body   string
	Footer string
}
type Font struct {
	FontSize int
	FontName string
	Bold     bool
	Italic   bool
	Strike   bool
	Color    string
	Another  string
}

var (
	Black   = "000000"
	White   = "FFFFFF"
	Red     = "FF0000"
	Lime    = "00FF00"
	Blue    = "0000FF"
	Yellow  = "FFFF00"
	Cyan    = "00FFFF"
	Magenta = "FF00FF"
	Silver  = "C0C0C0"
	Gray    = "808080"
	Maroon  = "800000"
	Olive   = "808000"
	Green   = "008000"
	Purple  = "800080"
	Teal    = "008080"
	Navy    = "000080"
)

type WP struct {
	Tag  string
	Body []WPTokens
}
type WpPr struct {
	Tag  string
	Body string
}
type RPr struct {
	Tag  string
	Body string
}
type WR struct {
	Tag  string
	Body string
}
type WPTokens struct {
	Tag    string
	Body   string
	Attr   string
	Status int
}
type SectPr struct {
	Tag  string
	Body string
}
type Document struct {
	WP     []WP
	SectPr SectPr
}
