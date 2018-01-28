package cmd

import (
	"github.com/tektsu/gedcom"
	"github.com/urfave/cli"
)

type apiControl struct {
	cx          *cli.Context
	gc          *gedcom.Gedcom
	indIndex    individualIndex
	famIndex    familyIndex
	sources     sourceResponses
	individuals individualResponses
	families    familyResponses
	photos      photoResponses
}

type citationResponse struct {
	ID        int    `json:"id"`
	SourceID  string `json:"sourceid"`
	SourceRef string `json:"sourceref"`
	Detail    string `json:"detail"`
}

type citationResponses map[int]*citationResponse

type eventResponse struct {
	Name      string   `json:"name"`
	Tag       string   `json:"tag"`
	Value     string   `json:"value"`
	Type      string   `json:"type"`
	Date      string   `json:"date"`
	Place     string   `json:"place"`
	Notes     []string `json:"notes"`
	Citations []int    `json:"citations"`
}

type familyControl struct {
	api           *apiControl
	citationCount int
	citationIndex map[string]int
	family        *gedcom.FamilyRecord
	response      *familyResponse
}

type familyIndex map[string]*familyReferenceResponse

type familyLinkResponse struct {
	ID        string                       `json:"id"`
	Pedigree  string                       `json:"pedigree"`
	AdoptedBy string                       `json:"adoptedby"`
	Events    []*eventResponse             `json:"events"`
	Mother    *individualReferenceResponse `json:"mother"`
	Father    *individualReferenceResponse `json:"father"`
	Children  individualReferenceResponses `json:"children"`
}

type familyReferenceResponse struct {
	ID      string                       `json:"id"`
	Married string                       `json:"married"`
	Title   string                       `json:"title"`
	Husband *individualReferenceResponse `json:"husband"`
	Wife    *individualReferenceResponse `json:"wife"`
}

type familyReferenceResponses []*familyReferenceResponse

type familyResponse struct {
	ID        string                       `json:"id"`
	Notes     []string                     `json:"note"`
	Ref       *familyReferenceResponse     `json:"ref"`
	Children  individualReferenceResponses `json:"children"`
	Photos    []*photoResponse             `json:"photos"`
	Events    []*eventResponse             `json:"events"`
	Citations citationResponses            `json:"citations"`
}

type familyResponses map[string]*familyResponse

type photoResponse struct {
	ID          string                       `json:"id"`
	File        string                       `json:"file"`
	Title       string                       `json:"title"`
	Description string                       `json:"description"`
	Height      int                          `json:"height"`
	Width       int                          `json:"width"`
	People      individualReferenceResponses `json:"people"`
	Families    familyReferenceResponses     `json:"families"`
	Notes       []string                     `json:"notes"`
}

type individualControl struct {
	api           *apiControl
	citationCount int
	citationIndex map[string]int
	individual    *gedcom.IndividualRecord
	response      *individualResponse
}

type individualIndex map[string]*individualReferenceResponse

type individualNameResponse struct {
	First     string `json:"first"`
	Last      string `json:"last"`
	Full      string `json:"full"`
	LastFirst string `json:"lastfirst"`
	Citations []int  `json:"citations"`
}

type individualReferenceResponse struct {
	ID        string   `json:"id"`
	Sex       string   `json:"sex"`
	Name      string   `json:"name"`
	Birth     string   `json:"birth"`
	Death     string   `json:"death"`
	Photo     string   `json:"photo"`
	LastNames []string `json:"lastnames"`
}

type individualReferenceResponses []*individualReferenceResponse

type individualResponse struct {
	ID            string                       `json:"id"`
	Ref           *individualReferenceResponse `json:"ref"`
	Name          *individualNameResponse      `json:"name"`
	Aliases       []*individualNameResponse    `json:"aliases"`
	Events        []*eventResponse             `json:"events"`
	Attributes    []*eventResponse             `json:"attributes"`
	ParentsFamily []*familyLinkResponse        `json:"parentsfamily"`
	Family        []*familyLinkResponse        `json:"family"`
	TopPhoto      *photoResponse               `json:"topphoto"`
	Photos        []*photoResponse             `json:"photos"`
	Citations     citationResponses            `json:"citations"`
	Notes         []string                     `json:"notes"`
}

type individualResponses map[string]*individualResponse

type photoResponses map[string]*photoResponse

type sourceCitationResponse struct {
	Individuals map[string]*individualReferenceResponse `json:"individuals"`
	Families    map[string]*familyReferenceResponse     `json:"families"`
}

type sourceCitationResponses map[string]*sourceCitationResponse

type sourceResponse struct {
	ID          string                  `json:"id"`
	Author      string                  `json:"author"`
	Title       string                  `json:"title"`
	Publication string                  `json:"publication"`
	File        []string                `json:"file"`
	RefNum      int                     `json:"refnum"`
	Ref         string                  `json:"ref"`
	Note        string                  `json:"note"`
	Citations   sourceCitationResponses `json:"citations"`
}

type sourceResponses map[string]*sourceResponse
