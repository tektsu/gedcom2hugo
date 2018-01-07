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

func (ic *individualControl) newEventResponse(event *gedcom.EventRecord) (*eventResponse, error) {

	response := &eventResponse{
		Name:      event.Tag,
		Tag:       event.Tag,
		Value:     event.Value,
		Type:      event.Type,
		Date:      event.Date,
		Place:     event.Place.Name,
		Citations: ic.addCitations(event.Citation),
	}
	if response.Tag == "EVEN" && response.Type != "" {
		response.Name = response.Type
	}
	name, exists := tagTable[response.Name]
	if exists {
		response.Name = name
	}

	response.Citations = append(response.Citations, ic.addCitations(event.Place.Citation)...) // Append place citations to the event

	return response, nil
}

func (ic *individualControl) addAttributes() error {
	for _, r := range ic.individual.Attribute {

		if r.Tag == "SSN" { // Skip social security number.
			continue
		}
		eventResponse, err := ic.newEventResponse(r)
		if err != nil {
			return err
		}
		ic.response.Attributes = append(ic.response.Attributes, eventResponse)
	}

	return nil
}

func (ic *individualControl) addEvents() error {
	for _, r := range ic.individual.Event {

		if r.Tag == "Photo" {
			continue
		}
		eventResponse, err := ic.newEventResponse(r)
		if err != nil {
			return err
		}
		if eventResponse.Name == "Birth" {
			ic.response.Birth = eventResponse.Date
		}
		if eventResponse.Name == "Death" {
			ic.response.Death = eventResponse.Date
		}

		ic.response.Events = append(ic.response.Events, eventResponse)
	}

	return nil
}
