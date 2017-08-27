package cmd

import "github.com/tektsu/gedcom"

// sourcesFromCitations builds a list of sources from a gedcom.CitatiionRecord.
func sourcesFromCitations(citations []*gedcom.CitationRecord) []*sourceRef {

	var sourceRefs []*sourceRef
	for _, c := range citations {
		sourceRefs = append(sourceRefs, newSourceRefFromCitation(c))
	}

	return sourceRefs
}
