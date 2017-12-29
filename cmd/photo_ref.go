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

type photoPersonIndex map[string]*personRef

type photoRef struct {
	ID      string
	File    string
	Title   string
	Height  int
	Width   int
	Persons photoPersonIndex
	Notes   []string
}

type photoIndex map[string]*photoRef

func newPhotoRef(o *gedcom.ObjectRecord, person *gedcom.IndividualRecord) *photoRef {

	if _, ok := photos[o.Xref]; !ok {
		photos[o.Xref] = &photoRef{
			ID:      o.Xref,
			File:    filepath.Base(o.File.Name),
			Title:   o.File.Title,
			Persons: make(photoPersonIndex),
		}

		file, err := os.Open(o.File.Name)
		defer file.Close()
		if err != nil {
			fmt.Printf("%v\n", err)
			return photos[o.Xref]
		}

		image, _, err := image.DecodeConfig(file) // Image Struct
		if err != nil {
			fmt.Printf("%s: %v\n", o.File.Name, err)
			return photos[o.Xref]
		}

		photos[o.Xref].Width = image.Width
		photos[o.Xref].Height = image.Height
	}
	p := photos[o.Xref]

	if _, ok := photos[o.Xref].Persons[person.Xref]; !ok {
		photos[o.Xref].Persons[person.Xref] = newPersonRef(person)
	}

	return p
}
