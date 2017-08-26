package cmd

import (
	"strconv"

	"github.com/tektsu/gedcom"
)

// sourceRef describes a source reference.
// RefNum is the numeric id of the source from the Gedcom file.
// Ref is a string representation of the source, and may contain additional
// information when a sourceRef is built from a citation.
type sourceRef struct {
	RefNum int
	Ref    string
}

// newSourceRef build a sourceRef from a gedcom.SourceRecord.
// A poorly-formed Gedcam file can produce a panic.
func newSourceRef(s *gedcom.SourceRecord) *sourceRef {
	refNum, err := strconv.Atoi(s.Xref[1:len(s.Xref)])
	if err != nil {
		panic(err)
	}

	return &sourceRef{
		RefNum: refNum,
		Ref:    sl[refNum],
	}
}

// newSourceRefFromCitation build a sourceRef from a gedcom.CitationRecord.
// A poorly-formed Gedcam file can produce a panic.
func newSourceRefFromCitation(c *gedcom.CitationRecord) *sourceRef {
	r := newSourceRef(c.Source)
	return r
}