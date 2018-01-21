package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/tektsu/gedcom"
	"github.com/urfave/cli"
)

var tagTable map[string]string

// Generate reads the GEDCOM file and builds the Hugo input files.
func Generate(cx *cli.Context) error {

	tagTable = map[string]string{
		"BAPM": "Baptism",
		"BIRT": "Birth",
		"BURI": "Buried",
		"CENS": "Census",
		"CHR":  "Christening",
		"DEAT": "Death",
		"DIV":  "Divorced",
		"DIVF": "Divorce Filed",
		"EMIG": "Emigrated",
		"ENGA": "Engaged",
		"GRAD": "Graduated",
		"MARB": "Marriage Bann",
		"MARL": "Marriage License",
		"MARR": "Married",
		"NATU": "Naturalized",
		"OCCU": "Occupation",
		"RELI": "Religion",
		"RESI": "Residence",
	}

	gc, err := readGedcom(cx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	api := newAPIResponse(cx)

	err = api.buildFromGedcom(gc)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = api.exportSourceAPI()
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = api.exportSourcePages()
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = api.exportIndividualAPI()
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = api.exportIndividualPages()
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = api.exportFamilyAPI()
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = api.exportFamilyPages()
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = api.exportPhotoAPI()
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = api.exportPhotoPages()
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// Configure for JSON headers.
	err = ioutil.WriteFile("static/api/_headers", []byte("/*  Access-Control-Allow-Origin: *  content-type: application/json; charset=utf-8"), 0644)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}

// readGedcom reads the GEDCOM file specified in the context into memory.
func readGedcom(cx *cli.Context) (*gedcom.Gedcom, error) {
	var gc *gedcom.Gedcom

	if cx.String("gedcom") == "" {
		return gc, errors.New("No GEDCOM file specified for input")
	}

	data, err := ioutil.ReadFile(cx.String("gedcom"))
	if err != nil {
		return gc, err
	}

	decoder := gedcom.NewDecoder(bytes.NewReader(data))
	decoder.SetUnrecTagFunc(func(l int, t, v, x string) {
		if t[0:1] == "_" {
			return
		}
		fmt.Printf("Unrecognized tag: %d %s %s", l, t, v)
		if x != "" {
			fmt.Printf(" (%s)", x)
		}
		fmt.Println("")
	})
	gc, err = decoder.Decode()
	if err != nil {
		return gc, err
	}
	return gc, nil
}
