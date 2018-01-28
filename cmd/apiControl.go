package cmd

import (
	_ "image/jpeg" // Register decoders
	_ "image/png"
	"strings"

	"github.com/tektsu/gedcom"
	"github.com/urfave/cli"
)

func newAPIControl(c *cli.Context) *apiControl {
	response := &apiControl{
		cx:          c,
		indIndex:    make(individualIndex),
		famIndex:    make(familyIndex),
		sources:     make(sourceResponses),
		individuals: make(individualResponses),
		families:    make(familyResponses),
		photos:      make(photoResponses),
	}

	return response
}

func (api *apiControl) addFamilyCitations(familyID string, citations []*gedcom.CitationRecord) error {
	for _, citation := range citations {
		sourceID := strings.ToLower(citation.Source.Xref)
		var c *sourceCitationResponse
		if _, ok := api.sources[sourceID].Citations[citation.Page]; ok {
			c = api.sources[sourceID].Citations[citation.Page]
		} else {
			c = newSourceCitationResponse()
			api.sources[sourceID].Citations[citation.Page] = c
		}
		c.Families[familyID], _ = api.getFamilyIndexEntry(familyID)
	}

	return nil
}
