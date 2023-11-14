package ooxml

import (
	"fmt"
	"strconv"
	"strings"
)

type NewToken struct {
	Name           string
	TagStatus      int
	TokenSymbOpen  int
	TokenSymbClose int
	Args           string
}

const TagComplimentary = 4
const TagOpenComplimentary = 1
const TagClosingComplimentary = 2
const TagSelfClosed = 3

const Creating = true
const Empty = false

func getTokens(d string) []NewToken {
	var Tokens []NewToken
	// token flags
	name := ""
	tagStatus := TagSelfClosed // check Tag open/closed/selfclosed
	tokenStatus := Empty       // token opening( true) or not(false)
	ArgStatus := Empty
	arg := ""
	lastChar := ""
	TokenSymbOpen := 0
	TokenSymbClose := 0
	for i, l := range d {
		let := string(l)
		switch let {
		case "<":
			switch tokenStatus {
			case Creating:
				//  if lastChar!="\\"{

				//  }
			case Empty:
				tokenStatus = Creating
				tagStatus = TagOpenComplimentary
				TokenSymbOpen = i
			}
		case ">":
			switch tokenStatus {
			case Creating:
				if lastChar == "/" {
					tagStatus = TagSelfClosed
				}
				TokenSymbClose = i
				Tokens = append(Tokens, NewToken{
					Name:           name,
					TagStatus:      tagStatus,
					TokenSymbOpen:  TokenSymbOpen,
					TokenSymbClose: TokenSymbClose,
					Args:           arg,
				})
				name = ""
				tagStatus = TagSelfClosed // check Tag open/closed/selfclosed
				ArgStatus = Empty
				arg = ""
				tokenStatus = Empty
			case Empty:

			}
		case "/":
			if lastChar == "<" {
				tagStatus = TagClosingComplimentary
			}
		case "\\":
		case " ":
			if tokenStatus {
				if ArgStatus == Creating {
					arg += " "
				} else {
					ArgStatus = Creating
				}

			}
		case "\t":
		case "\n":
		default:
			if tokenStatus {
				if ArgStatus {
					arg += let
				} else {
					name += let
				}
			}
		}
		lastChar = let
	}
	return Tokens
}

func getFirstNodes(tokens []NewToken) ([]NewToken, error) {
	// fmt.Printf("debug run %v\n", tokens)
	if len(tokens) == 0 {
		return []NewToken{}, fmt.Errorf("null tokens")
	}
	token := tokens[0]
	if len(tokens) == 1 {
		return []NewToken{token}, nil
	}
	if token.TagStatus == TagSelfClosed {
		var f []NewToken
		f = append(f, token)
		//	fmt.Printf("runnest [%d]\n", len(tokens[1:]))
		another, err := getFirstNodes(tokens[1:])
		if err != nil {
			return []NewToken{}, err
		}
		f = append(f, another...)
		return f, nil
	}
	tokenStatus := Creating
	iter := 0
	var tk []NewToken
	for i := 1; i < len(tokens); i++ {
		let := tokens[i]
		//	fmt.Printf("let : %v |token %s\n", let, token.Name)
		if let.Name == token.Name {
			if let.TagStatus == TagOpenComplimentary {
				//			fmt.Printf("iter += 1\n")
				iter += 1
			}
			if let.TagStatus == TagClosingComplimentary {
				if iter != 0 {
					//				fmt.Printf("iter -= 1\n")
					iter -= 1
				} else {
					//				fmt.Printf("tag %s closed\n", token.Name)
					tokenStatus = Empty
					tk = append(tk, token, let)
					if len(tokens) <= i {
						return tk, nil
					}
					//				fmt.Printf("another\n")
					//	fmt.Printf("another runnest [%d]\n", len(tokens[i+1:]))
					if len(tokens[i+1:]) != 0 {
						another, err := getFirstNodes(tokens[i+1:])
						if err != nil {
							return nil, err
						}
						tk = append(tk, another...)
						return tk, nil
					}
					return tk, nil
				}

			}
		}
	}
	if tokenStatus != Empty {
		return []NewToken{}, fmt.Errorf("tag [%v] wasn't closed", token)
	}
	return []NewToken{}, nil
}

type EmbeddedToken struct {
	Name      string
	TagStatus int
	Body      string
	Args      string
}

func GetParentNodes(d string) ([]EmbeddedToken, error) {
	if len(d) == 0 {
		return nil, fmt.Errorf("len 'xml string' == 0")
	}
	tokens := getTokens(d)
	//	fmt.Printf("tokens [%v]\n", tokens)
	items, err := getFirstNodes(tokens)
	if err != nil {
		return nil, err
	}
	var env NewToken
	envStatus := Empty
	var eTokens []EmbeddedToken
	// fmt.Printf("len items [%d] items[%v] d[%s]\n", len(items), items, d)
	for _, item := range items {
		//	fmt.Printf("item: %v\n", item)
		if item.TagStatus == TagSelfClosed {
			eTokens = append(eTokens, EmbeddedToken{
				Name:      item.Name,
				TagStatus: TagSelfClosed,
				Args:      item.Args,
				Body:      "",
			})
		}
		if item.TagStatus == TagOpenComplimentary {
			if envStatus == Empty {
				envStatus = Creating
				env = item
			} else {
				//	fmt.Printf("error adding element %v\n", item)
				return nil, fmt.Errorf("error adding element %v", item)
			}
		}
		if item.TagStatus == TagClosingComplimentary {
			if envStatus == Creating {
				//		fmt.Printf("creating\n")
				eTokens = append(eTokens, EmbeddedToken{
					Name:      item.Name,
					TagStatus: TagComplimentary,
					Args:      env.Args,
					Body:      d[env.TokenSymbClose+1 : item.TokenSymbOpen],
				})

				envStatus = Empty
				env = NewToken{}
			} else {
				return nil, fmt.Errorf("error adding element %v", item)
			}
		}
	}
	if envStatus != Empty {
		return nil, fmt.Errorf("broken xml.tag [%s] not closed", env.Name)
	}
	return eTokens, nil
}

type Font struct {
	FontSize           int // <w:sz w:val="36"/>& <w:szCs w:val="36"/>
	FontName           string
	Bold               bool // <w:b/> <w:bCs/>
	Italic             bool // <w:i/> <w:iCs/>
	Strike             bool
	Color              string // <w:color w:val="F10D0C"/>
	Underline          string // <w:u w:val="single"/> non or single or double
	Another            []EmbeddedToken
	AnotherFontTagAttr []string
}

func GetWVal(s string) string {
	if len(s) == 0 {
		return ""
	}
	z := strings.Split(s, "w:val=")
	if 2 <= len(z) {
		size := strings.Split(z[1], "\"")[1]
		// fmt.Printf("size: %s\n", size)
		return size
	}
	return ""
}
func ParseRPR(xml string) (Font, error) {
	nodes, err := GetParentNodes(xml)
	if err != nil {
		return Font{}, err
	}
	var rpr Font
	for _, item := range nodes {
		//	fmt.Printf("item: %s\n", item.Name)
		switch item.Name {
		case "w:b", "w:bCs":
			rpr.Bold = true

		case "w:i", "w:iCs":
			rpr.Italic = true

		case "w:sz", "w:szCs":
			z, err := strconv.Atoi(GetWVal(item.Args))
			if err != nil {
				return Font{}, fmt.Errorf("error parsing node [%v]. error [%s]", item, err)
			}
			rpr.FontSize = z
		case "w:u":
			rpr.Underline = GetWVal(item.Args)
		case "w:color":
			rpr.Color = GetWVal(item.Args)
		default:
			rpr.Another = append(rpr.Another, item)
		}
	}
	return rpr, nil
}
