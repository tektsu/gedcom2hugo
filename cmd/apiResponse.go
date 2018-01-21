package cmd

import (
	"fmt"
	"image"
	_ "image/jpeg" // Register decoders
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/tektsu/gedcom"
	"github.com/urfave/cli"
)

type apiResponse struct {
	cx          *cli.Context
	gc          *gedcom.Gedcom
	indIndex    individualIndex
	famIndex    familyIndex
	sources     sourceResponses
	individuals individualResponses
	families    familyResponses
	photos      photoResponses
}

type citationCallback func(string, []*gedcom.CitationRecord)
type citationSubCallback func([]*gedcom.CitationRecord) []int
type photoCallback func(*gedcom.ObjectRecord, *individualResponse) *photoResponse

func newAPIResponse(c *cli.Context) *apiResponse {
	response := &apiResponse{
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

func (api *apiResponse) addIndividualCitations(individualID string, citations []*gedcom.CitationRecord) error {
	for _, citation := range citations {
		sourceID := strings.ToLower(citation.Source.Xref)
		var c *sourceCitationResponse
		if _, ok := api.sources[sourceID].Citations[citation.Page]; ok {
			c = api.sources[sourceID].Citations[citation.Page]
		} else {
			c = newCitationResponse()
			api.sources[sourceID].Citations[citation.Page] = c
		}
		c.Individuals[individualID], _ = api.getIndividualIndexEntry(individualID)
	}

	return nil
}

func (api *apiResponse) addFamilyCitations(familyID string, citations []*gedcom.CitationRecord) error {
	for _, citation := range citations {
		sourceID := strings.ToLower(citation.Source.Xref)
		var c *sourceCitationResponse
		if _, ok := api.sources[sourceID].Citations[citation.Page]; ok {
			c = api.sources[sourceID].Citations[citation.Page]
		} else {
			c = newCitationResponse()
			api.sources[sourceID].Citations[citation.Page] = c
		}
		c.Families[familyID], _ = api.getFamilyIndexEntry(familyID)
	}

	return nil
}

func (api *apiResponse) addPhoto(o *gedcom.ObjectRecord) *photoResponse {
	key := getPhotoKeyFromObject(o)
	if _, ok := api.photos[key]; !ok {
		api.photos[key] = &photoResponse{
			ID:    key,
			File:  filepath.Base(o.File.Name),
			Title: o.File.Title,
		}

		file, err := os.Open(o.File.Name)
		defer file.Close()
		if err != nil {
			fmt.Printf("%v\n", err)
			return api.photos[key]
		}

		image, _, err := image.DecodeConfig(file) // Image Struct
		if err != nil {
			fmt.Printf("hsgdhajsgdjhas %s: %v\n", o.File.Name, err)
			return api.photos[key]
		}

		api.photos[key].Width = image.Width
		api.photos[key].Height = image.Height
	}

	return api.photos[key]
}

func (api *apiResponse) addPhotoForIndividual(o *gedcom.ObjectRecord, i *individualResponse) *photoResponse {
	response := api.addPhoto(o)

	ir, err := api.getIndividualIndexEntry(i.ID)
	if err != nil {
		fmt.Printf("%s\n", err)
		return response
	}
	response.People = append(response.People, ir)

	return response
}

func (api *apiResponse) addPhotoForFamily(o *gedcom.ObjectRecord, f *familyResponse) *photoResponse {
	response := api.addPhoto(o)

	fr, err := api.getFamilyIndexEntry(f.ID)
	if err != nil {
		fmt.Printf("%s\n", err)
		return response
	}
	response.Families = append(response.Families, fr)

	return response
}

func (api *apiResponse) buildFromGedcom(g *gedcom.Gedcom) error {

	api.gc = g

	var err error

	err = api.addSources()
	if err != nil {
		return err
	}

	err = api.addIndividuals()
	if err != nil {
		return err
	}

	err = api.addFamilies()
	if err != nil {
		return err
	}

	return nil
}

func getPhotoKeyFromObject(o *gedcom.ObjectRecord) string {

	key := "p" + strings.ToLower(strings.Replace(filepath.Base(o.File.Name), ".", "", -1))

	return key
}
