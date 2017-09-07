package cmd

import (
	"path/filepath"

	"github.com/tektsu/gedcom"
)

type photoRef struct {
	File  string
	Title string
	Notes []string
}

func newPhotoRef(o *gedcom.ObjectRecord) *photoRef {

	p := &photoRef{
		File:  filepath.Base(o.File),
		Title: o.Title,
	}

	return p
}
