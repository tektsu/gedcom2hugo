package cmd

import "github.com/tektsu/gedcom"

type eventResponse struct {
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	Value     string `json:"value"`
	Type      string `json:"type"`
	Date      string `json:"date"`
	Place     string `json:"place"`
	Citations []int  `json:"citations"`
}

func newEventResponse(event *gedcom.EventRecord, ccb citationSubCallback) (*eventResponse, error) {

	response := &eventResponse{
		Name:      event.Tag,
		Tag:       event.Tag,
		Value:     event.Value,
		Type:      event.Type,
		Date:      event.Date,
		Place:     event.Place.Name,
		Citations: ccb(event.Citation),
	}
	if response.Tag == "EVEN" && response.Type != "" {
		response.Name = response.Type
	}
	name, exists := tagTable[response.Name]
	if exists {
		response.Name = name
	}

	response.Citations = append(response.Citations, ccb(event.Place.Citation)...) // Append place citations to the event

	return response, nil
}
