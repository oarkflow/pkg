package main

import (
	"log"

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
	// TestDepthStructToParams()
}

func TestPlaceholders() {
	var user = User{
		Name:      "Alice",
		Age:       27,
		Nicknames: []string{"amber", "", "AL", "ice", "", "", "", "", "", "", "", ""},
		Friends: []*User{
			{Name: "Bob", Age: 28, ImageLocal: &docx.Image{Path: "images/avatar-4.png"}},
			{Name: "Cecilia", Age: 29, ImageLocal: &docx.Image{Path: "images/avatar-5.png"}},
			{Name: "", Age: 999},
			{Name: "", Age: 999},
			{Name: "Den", Age: 30},
			{Name: "", Age: 999},
			{Name: "Edgar", Age: 31, ImageLocal: &docx.Image{Path: "images/avatar-6.png"}},
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
	fname := "user.template.docx"
	tdoc, _ := docx.OpenTemplate("test-data/" + fname)
	tdoc.Params(user)
	if err := tdoc.ExportDocx("test-data/~test-" + fname + ".docx"); err != nil {
		log.Fatalf("ExportDocx: %s", err)
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
