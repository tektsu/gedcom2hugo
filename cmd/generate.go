package cmd

import (
	"bytes"
	"errors"
	"github.com/iand/gedcom"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path/filepath"
)

var tagTable map[string]string

// Generate reads the GEDCOM file and builds the Hugo input files.
//noinspection GoUnusedExportedFunction
func Generate(cx *cli.Context) error {
	tagTable = map[string]string{
		"BAPM":  "Baptism",
		"BIRT":  "Birth",
		"BURI":  "Buried",
		"CENS":  "Census",
		"CHR":   "Christening",
		"DEAT":  "Death",
		"DIV":   "Divorced",
		"DIVF":  "Divorce Filed",
		"EMIG":  "Emigrated",
		"ENGA":  "Engaged",
		"GRAD":  "Graduated",
		"MARB":  "Marriage Bann",
		"MARL":  "Marriage License",
		"MARR":  "Married",
		"NATU":  "Naturalized",
		"OCCU":  "Occupation",
		"RELI":  "Religion",
		"RESI":  "Residence",
		"_MILT": "Military Service",
	}

	gc, err := readGedcom(cx)
	if err != nil {
		return cli.Exit(err, 1)
	}

	api := newAPIControl(cx)

	err = api.buildFromGedcom(gc)
	if err != nil {
		return cli.Exit(err, 1)
	}

	err = api.exportSourceAPI()
	if err != nil {
		return cli.Exit(err, 1)
	}

	err = api.exportSourcePages()
	if err != nil {
		return cli.Exit(err, 1)
	}

	err = api.exportIndividualAPI()
	if err != nil {
		return cli.Exit(err, 1)
	}

	err = api.exportIndividualPages()
	if err != nil {
		return cli.Exit(err, 1)
	}

	err = api.exportFamilyAPI()
	if err != nil {
		return cli.Exit(err, 1)
	}

	err = api.exportFamilyPages()
	if err != nil {
		return cli.Exit(err, 1)
	}

	err = api.exportPhotoAPI()
	if err != nil {
		return cli.Exit(err, 1)
	}

	err = api.exportPhotoPages()
	if err != nil {
		return cli.Exit(err, 1)
	}

	err = configureForJsonHeaders(api, err)
	if err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

func configureForJsonHeaders(api *apiControl, err error) error {
	headers := filepath.Join(api.cx.String("project"), "/static/api/_headers")
	file, err := os.Create(headers)
	if file == nil {
		return nil
	}
	if err != nil {
		_ = file.Close()
		return err
	}
	_, err = file.Write([]byte("/*  Access-Control-Allow-Origin: *  content-type: application/json; charset=utf-8"))
	if err != nil {
		_ = file.Close()
		return err
	}
	_ = file.Close()
	return nil
}

// readGedcom reads the GEDCOM file specified in the context into memory.
func readGedcom(cx *cli.Context) (*gedcom.Gedcom, error) {
	var gc *gedcom.Gedcom

	if cx.String("gedcom") == "" {
		return gc, errors.New("no GEDCOM file specified for input")
	}

	data, err := ioutil.ReadFile(cx.String("gedcom"))
	if err != nil {
		return gc, err
	}

	decoder := gedcom.NewDecoder(bytes.NewReader(data))

	gc, err = decoder.Decode()
	if err != nil {
		return gc, err
	}
	return gc, nil
}
