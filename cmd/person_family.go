package cmd

import (
	"github.com/tektsu/gedcom"
)

// personFamily describes a person's family, either their own or their
// parents.
// ID is the Gedcom family Xref.
// Mother, Father and Children are personRefs representing the members
// of the family.
type personFamily struct {
	ID        string
	Pedigree  string
	AdoptedBy string
	Events    []*eventRef
	Mother    *personRef
	Father    *personRef
	Children  []*personRef
}

// newPersonFamily builds a personFamily record from a
// gedcom.FamilyLinkRecord.
// In addition to the FamilyLinkRecord, it is passed a callback function to
// handle source references.
func newPersonFamily(flr *gedcom.FamilyLinkRecord, handleSources sourceCB) *personFamily {

	if flr.Family == nil {
		return nil
	}

	family := &personFamily{
		ID:        flr.Family.Xref,
		Pedigree:  flr.Pedigree,
		AdoptedBy: flr.AdoptedBy,
	}

	family.Father = newPersonRefWithCitations(flr.Family.Husband, handleSources)
	family.Mother = newPersonRefWithCitations(flr.Family.Wife, handleSources)
	for _, e := range flr.Family.Event {
		event := newEventRef(e, handleSources)
		family.Events = append(family.Events, event)
	}
	for _, i := range flr.Family.Child {
		var child *personRef
		child = newPersonRefWithCitationsAsChild(i, handleSources)
		family.Children = append(family.Children, child)
	}

	return family
}
