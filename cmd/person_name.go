package cmd

import (
	"fmt"

	"github.com/tektsu/gedcom"
)

// personName describes a person's name in several forms, and may include a
// list of local source references.
type personName struct {
	Full       string
	Last       string
	LastFirst  string
	SourcesInd []int
}

// newPersonName builds a personRef from a gedcom.NameRecord.
func newPersonName(n *gedcom.NameRecord) *personName {

	given, family := extractNames(n.Name)
	name := &personName{
		Last:      family,
		Full:      fmt.Sprintf("%s %s", given, family),
		LastFirst: fmt.Sprintf("%s, %s", family, given),
	}

	return name
}

// newPersonNameWithCitations builds a personRef from a gedcom.NameRecord and
// processes citations.
// In addition to a NameRecord, it is passed a local citation counter
// and a callback function to handle source references.
// It returns the new value of the citation counter and a new personName.
func newPersonNameWithCitations(count int, n *gedcom.NameRecord, cbSources func([]*sourceRef)) (int, *personName) {
	var sources []*sourceRef

	name := newPersonName(n)

	sources = sourcesFromCitations(n.Citation)
	for _ = range sources {
		count++
		name.SourcesInd = append(name.SourcesInd, count)
	}
	cbSources(sources)

	return count, name
}
