package cmd

import "github.com/tektsu/gedcom"

// personFamily describes a person's family, either their own or their
// parents.
// ID is the Gedcom family Xref.
// Mother, Father and Children are personRefs representing the members
// of the family.
type personFamily struct {
	ID        string
	Pedigree  string
	AdoptedBy string
	Mother    *personRef
	Father    *personRef
	Children  []*personRef
}

// newPersonFamily builds a personFamily record from a
// gedcom.FamilyLinkRecord.
// It is passed the local citation counter and a callback function to handle
// source references.
// It returns the new citation counter value, an array of source summaries,
// and a new peersonFamily record.
func newPersonFamily(count int, flr *gedcom.FamilyLinkRecord, cbSources func([]*sourceRef)) (int, *personFamily) {
	var sources []*sourceRef

	if flr.Family == nil {
		return count, nil
	}

	family := &personFamily{
		ID:        flr.Family.Xref,
		Pedigree:  flr.Pedigree,
		AdoptedBy: flr.AdoptedBy,
	}

	count, family.Father = newPersonRefWithCitations(count, flr.Family.Husband, cbSources)
	count, family.Mother = newPersonRefWithCitations(count, flr.Family.Wife, cbSources)
	for _, i := range flr.Family.Child {
		var child *personRef
		count, child = newPersonRefWithCitations(count, i, cbSources)
		family.Children = append(family.Children, child)
	}

	cbSources(sources)
	return count, family
}
