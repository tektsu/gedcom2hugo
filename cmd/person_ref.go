package cmd

import "github.com/tektsu/gedcom"

// personRef describes a reference to a person.
// It is used to reference other people on a person page.
// ID is the Gedcom Xref of the person.
// Name is the name of the person.
// Sex is the sex of the person.
// SourcesInd is an array of local references to sources for the person's name.
type personRef struct {
	ID         string
	Name       string
	Sex        string
	SourcesInd []int
}

// newPersonRef builds person reference from a gedcom.IndividualRecord.
func newPersonRef(i *gedcom.IndividualRecord) *personRef {

	person := &personRef{
		ID:   i.Xref,
		Sex:  i.Sex,
		Name: people[i.Xref].FullName,
	}
	if person.Sex != "M" && person.Sex != "F" {
		person.Sex = "U"
	}

	return person
}

// newPersonRefwithCitations builds person reference from a
// gedcom.IndividualRecord, adding in citations from the name of the person.
// In addition to the IndividualRecord, it is passed callback function to
// handle source references.
func newPersonRefWithCitations(i *gedcom.IndividualRecord, handleSources sourceCB) *personRef {

	if i == nil {
		return nil
	}

	person := newPersonRef(i)
	person.SourcesInd = handleSources(sourcesFromCitations(i.Name[0].Citation))

	return person
}

// newPersonRefWithCitationsAsChild is the same as mewPersonRefWithCitations,
// except that it omits the last name of the person.
func newPersonRefWithCitationsAsChild(i *gedcom.IndividualRecord, handleSources sourceCB) *personRef {

	if i == nil {
		return nil
	}

	person := newPersonRefWithCitations(i, handleSources)
	person.Name = people[i.Xref].GivenName

	return person
}
