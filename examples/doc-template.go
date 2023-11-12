package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/oarkflow/pkg/docx"
)

// Using this data for all examples below
type User struct {
	Name                   string
	Age                    int
	Nicknames              []string
	Friends                []*User
	BrokenStylePlaceholder string
	TriggerRemove          string
	ImageLocal             *docx.Image
	ImageURL               *docx.Image
	Images                 []*docx.Image
}

func main() {
	TestPlaceholders()
	TestDepthStructToParams()
}

func TestPlaceholders() {
	var user = User{
		Name:      "Alice",
		Age:       27,
		Nicknames: []string{"amber", "", "AL", "ice", "", "", "", "", "", "", "", ""},
		Friends: []*User{
			{Name: "Bob", Age: 28, ImageLocal: &docx.Image{Path: "images/avatar-4.png", Width: 25, Height: 25}},
			{Name: "Cecilia", Age: 29, ImageLocal: &docx.Image{Path: "images/avatar-5.png", Width: 25, Height: 25}},
			{Name: "", Age: 999},
			{Name: "", Age: 999},
			{Name: "Den", Age: 30},
			{Name: "", Age: 999},
			{Name: "Edgar", Age: 31, ImageLocal: &docx.Image{Path: "images/avatar-6.png", Width: 25, Height: 25}},
			{Name: "", Age: 999},
			{Name: "", Age: 999},
		},
		BrokenStylePlaceholder: "(NOT ANYMORE)",
		ImageLocal: &docx.Image{
			Path:   "images/avatar-1.png",
			Width:  25,
			Height: 25,
		},
		ImageURL: &docx.Image{
			URL:    "https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png",
			Width:  25,
			Height: 25,
		},
		Images: []*docx.Image{
			{
				Path:   "images/avatar-2.png",
				Width:  25,
				Height: 25,
			},
			{
				Path:   "images/avatar-3.png",
				Width:  25,
				Height: 25,
			},
		},
	}

	filenames := []string{
		"user.template.docx",
		"tables.docx",
		"lists.docx",
	}

	for _, fname := range filenames {

		inputs := []string{
			"struct",
			"json",
		}

		// Test param setup byu different input types
		for _, inType := range inputs {
			tdoc, _ := docx.OpenTemplate("test-data/" + fname)

			plaintext := tdoc.Plaintext()
			params := []string{
				"{{Name}}",
				"{{Friends.1.Name}}",
				"{{Friends.Name :empty",
				"{{Friends.Age}}",
			}
			for _, p := range params {
				if !strings.Contains(plaintext, p) {
					log.Fatalf("Param `%s` should be found in plaintext: \n\n%s", p, plaintext)
				}
			}

			// Run
			switch inType {
			case "struct":
				tdoc.Params(user)
			case "json":
				buf, _ := json.Marshal(user)
				tdoc.Params(buf)
			}

			plaintext = tdoc.Plaintext()
			if err := tdoc.ExportDocx("test-data/~test-" + fname + "-" + inType + ".docx"); err != nil {
				log.Fatalf("[%s] ExportDocx: %s", inType, err)
			}

			// Check for "must remove" text
			removedTexts := []string{
				":empty",
				":empty:remove",
				"must be removed",
				"999 y/o",
			}
			for _, s := range removedTexts {
				if strings.Contains(plaintext, s) {
					log.Fatalf("Text `%s` must be removed: \n%s", s, tdoc.Plaintext())
				}
			}

			// Test for know lines
			for _, u := range user.Friends {
				if u.Name == "" {
					continue
				}

				if !strings.Contains(plaintext, u.Name) {
					log.Fatalf("User[%s] friends Name must be found: \n\n%s", u.Name, tdoc.Plaintext())
				}

				years := fmt.Sprintf("%d y/o", u.Age)
				if !strings.Contains(plaintext, years) {
					log.Fatalf("User[%s] friends Age[%d] must be found: \n\n%s", u.Name, u.Age, tdoc.Plaintext())
				}
			}

			// Test empty struct
			// non-empty-trigger placeholders must stay as is
			tdoc, _ = docx.OpenTemplate("test-data/" + fname)
			tdoc.Params(struct{ Dummy string }{Dummy: "never"})
			if err := tdoc.ExportDocx("test-data/~test-" + inType + ".docx"); err != nil {
				log.Fatalf("[%s] ExportDocx: %s", inType, err)
			}
			notreplacedPlaceholdersStr := strings.Join(tdoc.Placeholders(), ", ")

			plaintext = tdoc.Plaintext()
			mustBeMissing := []string{
				"{{Name}}",
				"{{Friends.1.Name}}",
				"{{NotReplacable}}",
				"{{NotReplacable , }}",
			}
			for _, s := range mustBeMissing {
				if !strings.Contains(notreplacedPlaceholdersStr, s) {
					log.Fatalf("Placeholder [%s] must be inPlaceholders() slice. Found: %s\n\n%s\n\n", s, notreplacedPlaceholdersStr, plaintext)
				}
				if !strings.Contains(plaintext, s) {
					log.Fatalf("Text [%s] must be found: \n\n%s", s, plaintext)
				}
			}
		}
	}

}

func TestDepthStructToParams() {
	var user = User{
		Name: "Alice",
		Age:  27,
		Friends: []*User{
			{Name: "Bob", Age: 28, Friends: []*User{
				{Name: "Cecilia", Age: 29},
				{Name: "Sun", Age: 999},
				{Name: "Tony", Age: 999},
			}},
			{Name: "Den", Age: 30, Friends: []*User{
				{Name: "Ben", Age: 999},
				{Name: "Edgar", Age: 31},
				{Name: "Jouny", Age: 999},
				{Name: "Carrzy", Age: 999},
			}},
		},
	}

	tdoc, _ := docx.OpenTemplate("test-data/depth.docx")
	tdoc.Params(user)
	if err := tdoc.ExportDocx("test-data/~test-depth.docx"); err != nil {
		log.Fatal(err)
	}
}
