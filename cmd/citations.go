package cmd

import "github.com/tektsu/gedcom"

func sourcesFromCitations(citations []*gedcom.CitationRecord) []*sourceRef {

	var sources []*sourceRef
	for _, c := range citations {
		sources = append(sources, newSourceRefFromCitation(c))
	}

	return sources
}
