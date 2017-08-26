package cmd

import "github.com/tektsu/gedcom"

type personRef struct {
	ID         string
	Name       string
	Sex        string
	SourcesInd []int
}

func newPersonRef(i *gedcom.IndividualRecord) *personRef {

	r := &personRef{
		ID:   i.Xref,
		Sex:  i.Sex,
		Name: people[i.Xref].FullName,
	}

	return r
}

func (r *personRef) citations(cc int, c []*gedcom.CitationRecord) (int, []*sourceRef) {

	sources := sourcesFromCitations(c)
	for _ = range sources {
		cc++
		r.SourcesInd = append(r.SourcesInd, cc)
	}
	return cc, sources
}
