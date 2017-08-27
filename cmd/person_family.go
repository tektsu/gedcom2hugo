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
// It is passed the local citation counter, and returns the new citation
// counter value, an array of source summaries, and a new peersonFamily record.
func newPersonFamily(count int, flr *gedcom.FamilyLinkRecord) (int, []*sourceRef, *personFamily) {
	var sources []*sourceRef

	if flr.Family == nil {
		return count, sources, nil
	}

	family := &personFamily{
		ID:        flr.Family.Xref,
		Pedigree:  flr.Pedigree,
		AdoptedBy: flr.AdoptedBy,
	}

	createPersonRef := func(c int, i *gedcom.IndividualRecord) (int, *personRef) {

		c, s, person := newPersonRefWithCitations(count, i)
		for _, source := range s {
			sources = append(sources, source)
		}

		return c, person
	}

	count, family.Father = createPersonRef(count, flr.Family.Husband)
	count, family.Mother = createPersonRef(count, flr.Family.Wife)
	for _, i := range flr.Family.Child {
		var child *personRef
		count, child = createPersonRef(count, i)
		family.Children = append(family.Children, child)
	}

	return count, sources, family
}
