package cmd

import "github.com/tektsu/gedcom"

type eventRef struct {
	Tag        string
	Value      string
	Type       string
	Date       string
	Place      string
	SourcesInd []int
}

func newEventRef(e *gedcom.EventRecord, handleSources sourceCB) *eventRef {

	event := &eventRef{
		Tag:   e.Tag,
		Value: e.Value,
		Type:  e.Type,
		Date:  e.Date,
		Place: e.Place.Name,
	}
	event.SourcesInd = handleSources(sourcesFromCitations(e.Citation))
	placeSources := handleSources(sourcesFromCitations(e.Place.Citation))
	for _, i := range placeSources {
		event.SourcesInd = append(event.SourcesInd, i)
	}

	return event
}
