package cmd

import (
	"fmt"
	"image"
	_ "image/jpeg" // Register decoders
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/tektsu/gedcom"
)

type photoRef struct {
	File   string
	Title  string
	Height int
	Width  int
	Notes  []string
}

func newPhotoRef(o *gedcom.ObjectRecord) *photoRef {

	p := &photoRef{
		File:  filepath.Base(o.File.Name),
		Title: o.File.Title,
	}

	file, err := os.Open(o.File.Name)
	defer file.Close()
	if err != nil {
		fmt.Printf("%v\n", err)
		return p
	}

	image, _, err := image.DecodeConfig(file) // Image Struct
	if err != nil {
		fmt.Printf("%s: %v\n", o.File.Name, err)
		return p
	}

	p.Width = image.Width
	p.Height = image.Height

	return p
}
