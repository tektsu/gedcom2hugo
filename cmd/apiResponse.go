package cmd

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/tektsu/gedcom"
	"github.com/urfave/cli"
)

type apiResponse struct {
	cx          *cli.Context
	gc          *gedcom.Gedcom
	sources     sourceResponses
	individuals individualResponses
	photos      photoResponses
}

type citationCallback func(string, []*gedcom.CitationRecord)
type citationSubCallback func([]*gedcom.CitationRecord) []int
type photoCallback func(*gedcom.ObjectRecord, *individualResponse) *photoResponse

func buildAPIResponseFromGedcom(c *cli.Context, g *gedcom.Gedcom) (*apiResponse, error) {

	response := &apiResponse{
		cx:     c,
		gc:     g,
		photos: make(photoResponses),
	}

	var err error

	response.sources, err = newSourceResponsesFromGedcom(g)
	if err != nil {
		return response, err
	}

	// Callback for individual citations.
	iccb := func(individualID string, citations []*gedcom.CitationRecord) {
		for _, citation := range citations {
			sourceID := strings.ToLower(citation.Source.Xref)
			var c *sourceCitationResponse
			if _, ok := response.sources[sourceID].Citations[citation.Page]; ok {
				c = response.sources[sourceID].Citations[citation.Page]
			} else {
				c = newCitationResponse()
				response.sources[sourceID].Citations[citation.Page] = c
			}
			c.Individuals[individualID] = true
		}
	}

	// Callback for photos.
	photocb := func(o *gedcom.ObjectRecord, i *individualResponse) *photoResponse {
		key := getPhotoKeyFromObject(o)
		if _, ok := response.photos[key]; !ok {
			response.photos[key] = &photoResponse{
				ID:    key,
				File:  filepath.Base(o.File.Name),
				Title: o.File.Title,
				//People: make(photoPersonIndex),
			}

			file, err := os.Open(o.File.Name)
			defer file.Close()
			if err != nil {
				fmt.Printf("%v\n", err)
				return response.photos[key]
			}

			image, _, err := image.DecodeConfig(file) // Image Struct
			if err != nil {
				fmt.Printf("%s: %v\n", o.File.Name, err)
				return response.photos[key]
			}

			response.photos[key].Width = image.Width
			response.photos[key].Height = image.Height
		}

		//if _, ok := response.photos[key].Persons[person.Xref]; !ok {
		//response.photos[key].Persons[person.Xref] = newPersonRef(person)
		//}

		return response.photos[key]

	}

	response.individuals, err = newIndividualResponses()
	err = response.individuals.addAll(g.Individual, iccb, photocb)
	if err != nil {
		return response, err
	}

	return response, nil
}
