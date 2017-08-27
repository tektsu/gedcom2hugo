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

// newPersonName builds a personName from a gedcom.NameRecord.
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
// In addition to a NameRecord, it is passed a callback function to handle
// source references.
func newPersonNameWithCitations(n *gedcom.NameRecord, handleSources sourceCB) *personName {

	name := newPersonName(n)
	name.SourcesInd = handleSources(sourcesFromCitations(n.Citation))

	return name
}
