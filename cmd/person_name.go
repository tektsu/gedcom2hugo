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
	name := &personName{
		Last:      family,
		Full:      fmt.Sprintf("%s %s", given, family),
		LastFirst: fmt.Sprintf("%s, %s", family, given),
	}

	return name
}

func newPersonNameWithCitations(count int, n *gedcom.NameRecord) (int, []*sourceRef, *personName) {
	var sources []*sourceRef

	name := newPersonName(n)

	sources = sourcesFromCitations(n.Citation)
	for _ = range sources {
		count++
		name.SourcesInd = append(name.SourcesInd, count)
	}

	return count, sources, name
}

func (name *personName) citations(count int, c []*gedcom.CitationRecord) (int, []*sourceRef) {

	sources := sourcesFromCitations(c)
	for _ = range sources {
		count++
		name.SourcesInd = append(name.SourcesInd, count)
	}
	return count, sources
}
