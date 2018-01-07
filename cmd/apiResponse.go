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

func newAPIResponse(c *cli.Context) *apiResponse {
	response := &apiResponse{
		cx:          c,
		sources:     make(sourceResponses),
		individuals: make(individualResponses),
		photos:      make(photoResponses),
	}

	return response
}

func (api *apiResponse) addCitations(individualID string, citations []*gedcom.CitationRecord) {
	for _, citation := range citations {
		sourceID := strings.ToLower(citation.Source.Xref)
		var c *sourceCitationResponse
		if _, ok := api.sources[sourceID].Citations[citation.Page]; ok {
			c = api.sources[sourceID].Citations[citation.Page]
		} else {
			c = newCitationResponse()
			api.sources[sourceID].Citations[citation.Page] = c
		}
		c.Individuals[individualID] = true
	}
}

func (api *apiResponse) addPhoto(o *gedcom.ObjectRecord, i *individualResponse) *photoResponse {
	key := getPhotoKeyFromObject(o)
	if _, ok := api.photos[key]; !ok {
		api.photos[key] = &photoResponse{
			ID:    key,
			File:  filepath.Base(o.File.Name),
			Title: o.File.Title,
			//People: make(photoPersonIndex),
		}

		file, err := os.Open(o.File.Name)
		defer file.Close()
		if err != nil {
			fmt.Printf("%v\n", err)
			return api.photos[key]
		}

		image, _, err := image.DecodeConfig(file) // Image Struct
		if err != nil {
			fmt.Printf("%s: %v\n", o.File.Name, err)
			return api.photos[key]
		}

		api.photos[key].Width = image.Width
		api.photos[key].Height = image.Height
	}

	//if _, ok := response.photos[key].Persons[person.Xref]; !ok {
	//response.photos[key].Persons[person.Xref] = newPersonRef(person)
	//}

	return api.photos[key]

}

func (api *apiResponse) buildFromGedcom(g *gedcom.Gedcom) error {

	api.gc = g

	var err error

	err = api.addSources()
	if err != nil {
		return err
	}

	// Callback for individual citations.
	//iccb := func(individualID string, citations []*gedcom.CitationRecord) {
	//for _, citation := range citations {
	//sourceID := strings.ToLower(citation.Source.Xref)
	//var c *sourceCitationResponse
	//if _, ok := api.sources[sourceID].Citations[citation.Page]; ok {
	//c = api.sources[sourceID].Citations[citation.Page]
	//} else {
	//c = newCitationResponse()
	//api.sources[sourceID].Citations[citation.Page] = c
	//}
	//c.Individuals[individualID] = true
	//}
	//}

	// Callback for photos.
	//photocb := func(o *gedcom.ObjectRecord, i *individualResponse) *photoResponse {
	//key := getPhotoKeyFromObject(o)
	//if _, ok := api.photos[key]; !ok {
	//api.photos[key] = &photoResponse{
	//ID:    key,
	//File:  filepath.Base(o.File.Name),
	//Title: o.File.Title,
	////People: make(photoPersonIndex),
	//}

	//file, err := os.Open(o.File.Name)
	//defer file.Close()
	//if err != nil {
	//fmt.Printf("%v\n", err)
	//return api.photos[key]
	//}

	//image, _, err := image.DecodeConfig(file) // Image Struct
	//if err != nil {
	//fmt.Printf("%s: %v\n", o.File.Name, err)
	//return api.photos[key]
	//}

	//api.photos[key].Width = image.Width
	//api.photos[key].Height = image.Height
	//}

	////if _, ok := response.photos[key].Persons[person.Xref]; !ok {
	////response.photos[key].Persons[person.Xref] = newPersonRef(person)
	////}

	//return api.photos[key]

	//}

	//api.individuals, err = newIndividualResponses()
	//err = api.individuals.addAll(g.Individual, iccb, photocb)
	err = api.addIndividuals()
	if err != nil {
		return err
	}

	return nil
}
