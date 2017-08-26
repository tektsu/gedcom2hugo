package cmd

import (
	"fmt"

	"github.com/tektsu/gedcom"
)

type personName struct {
	Full       string
	Last       string
	LastFirst  string
	SourcesInd []int
}

func newPersonName(n *gedcom.NameRecord) *personName {

	given, family := extractNames(n.Name)
	r := &personName{
		Last:      family,
		Full:      fmt.Sprintf("%s %s", given, family),
		LastFirst: fmt.Sprintf("%s, %s", family, given),
	}

	return r
}

func (r *personName) citations(cc int, c []*gedcom.CitationRecord) (int, []*sourceRef) {

	sources := sourcesFromCitations(c)
	for _ = range sources {
		cc++
		r.SourcesInd = append(r.SourcesInd, cc)
	}
	return cc, sources
}
