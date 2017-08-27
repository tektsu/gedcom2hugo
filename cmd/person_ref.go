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

	return person
}

// newPersonRef builds person reference from a gedcom.IndividualRecord, adding
// in citations from the name of the person.
// In addition to the IndividualRecord, it is passed a local citation counter
// and a callback function to handle source references.
// It returns the new value of the citation counter and a new personRef.
func newPersonRefWithCitations(count int, i *gedcom.IndividualRecord, cbSources func([]*sourceRef)) (int, *personRef) {

	if i == nil {
		return count, nil
	}

	person := newPersonRef(i)

	sources := sourcesFromCitations(i.Name[0].Citation)
	for _ = range sources {
		count++
		person.SourcesInd = append(person.SourcesInd, count)
	}

	cbSources(sources)
	return count, person
}
