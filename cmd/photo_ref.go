package cmd

import (
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
		File:   filepath.Base(o.File),
		Title:  o.Title,
		Height: o.Height,
		Width:  o.Width,
	}

	return p
}
