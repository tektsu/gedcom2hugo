package cmd

import (
	"fmt"
	"image"
	_ "image/jpeg" // Register decoders
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

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

	key := getPhotoKeyFromObject(o)
	if _, ok := photos[key]; !ok {
		photos[key] = &photoRef{
			ID:      key,
			File:    filepath.Base(o.File.Name),
			Title:   o.File.Title,
			Persons: make(photoPersonIndex),
		}

		file, err := os.Open(o.File.Name)
		defer file.Close()
		if err != nil {
			fmt.Printf("%v\n", err)
			return photos[key]
		}

		image, _, err := image.DecodeConfig(file) // Image Struct
		if err != nil {
			fmt.Printf("%s: %v\n", o.File.Name, err)
			return photos[key]
		}

		photos[key].Width = image.Width
		photos[key].Height = image.Height
	}

	if _, ok := photos[key].Persons[person.Xref]; !ok {
		photos[key].Persons[person.Xref] = newPersonRef(person)
	}

	return photos[key]
}

func getPhotoKeyFromObject(o *gedcom.ObjectRecord) string {

	key := "p" + strings.ToLower(strings.Replace(filepath.Base(o.File.Name), ".", "", -1))

	return key
}
