package cmd

import (
	"strconv"

	"github.com/tektsu/gedcom"
)

type sourceRef struct {
	RefNum int
	Ref    string
}

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

func newSourceRefFromCitation(c *gedcom.CitationRecord) *sourceRef {
	r := newSourceRef(c.Source)
	return r
}
