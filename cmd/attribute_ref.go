package cmd

import "github.com/tektsu/gedcom"

type eventRef struct {
	Name       string
	Tag        string
	Value      string
	Type       string
	Date       string
	Place      string
	SourcesInd []int
}

func newEventRef(e *gedcom.EventRecord, handleSources sourceCB) *eventRef {

	event := &eventRef{
		Name:  e.Tag,
		Tag:   e.Tag,
		Value: e.Value,
		Type:  e.Type,
		Date:  e.Date,
		Place: e.Place.Name,
	}
	if event.Tag == "EVEN" && event.Type != "" {
		event.Name = event.Type
	}
	name, exists := tagTable[event.Name]
	if exists {
		event.Name = name
	}
	event.SourcesInd = handleSources(sourcesFromCitations(e.Citation))
	placeSources := handleSources(sourcesFromCitations(e.Place.Citation))
	for _, i := range placeSources {
		event.SourcesInd = append(event.SourcesInd, i)
	}

	return event
}
