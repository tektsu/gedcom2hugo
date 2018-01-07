package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/tektsu/gedcom"
	"github.com/urfave/cli"
)

// sourceIndex is a cache of information about each source
type sourceIndex map[int]string

// Global caches
var sources sourceIndex
var people personIndex
var tagTable map[string]string
var photos photoIndex

func shortcode(c string) string {
	return fmt.Sprintf("{{< %s >}}", c)
}

func openShortcode() string {
	return fmt.Sprintf("{{< ")
}

func closeShortcode() string {
	return fmt.Sprintf(" >}}")
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

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

	photos = make(photoIndex)

	project := cx.String("project")

	gc, err := readGedcom(cx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	people = newPersonIndex(gc)

	// Generate Source Pages.
	sources = make(sourceIndex)
	sourceDir := filepath.Join(project, "content", "source")
	err = os.MkdirAll(sourceDir, 0777)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	for _, source := range gc.Source {
		id := source.Xref
		file := filepath.Join(sourceDir, strings.ToLower(id+".md"))
		fh, err := os.Create(file)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		defer fh.Close()

		data := newSourceTmplData(source)
		sources[data.RefNum] = data.Ref

		tpl := template.New("source")
		tpl, err = tpl.Parse(sourcePageTemplate)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		err = tpl.Execute(fh, data)
	}

	// Generate Person Pages.
	personDir := filepath.Join(project, "content", "person")
	err = os.MkdirAll(personDir, 0777)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	for _, person := range gc.Individual {
		id := person.Xref
		file := filepath.Join(personDir, strings.ToLower(id+".md"))

		fh, err := os.Create(file)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		defer fh.Close()

		data := newPersonTmplData(person)

		tpl := template.New("person")
		funcs := template.FuncMap{
			"add":            func(x, y int) int { return x + y },
			"min":            min,
			"ToLower":        strings.ToLower,
			"shortcode":      shortcode,
			"openShortcode":  openShortcode,
			"closeShortcode": closeShortcode,
		}
		tpl.Funcs(funcs)
		tpl, err = tpl.Parse(personPageTemplate)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		err = tpl.Execute(fh, data)
	}

	// Generate Media Pages.
	mediaDir := filepath.Join(project, "content", "media")
	err = os.MkdirAll(mediaDir, 0777)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	for key, photo := range photos {

		file := filepath.Join(mediaDir, strings.ToLower(key+".md"))

		fh, err := os.Create(file)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		defer fh.Close()

		tpl := template.New("photo")
		funcs := template.FuncMap{
			"add":            func(x, y int) int { return x + y },
			"min":            min,
			"ToLower":        strings.ToLower,
			"shortcode":      shortcode,
			"openShortcode":  openShortcode,
			"closeShortcode": closeShortcode,
		}
		tpl.Funcs(funcs)
		tpl, err = tpl.Parse(photoPageTemplate)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		err = tpl.Execute(fh, newPhotoTmplData(photo))
	}

	err = generateAPI(cx, gc)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}

func generateAPI(cx *cli.Context, gc *gedcom.Gedcom) error {

	project := cx.String("project")

	api := newAPIResponse(cx)

	err := api.buildFromGedcom(gc)
	if err != nil {
		return err
	}

	// Generate source api responses.
	sourceAPIDir := filepath.Join(project, "static", "api", "source")
	err = os.MkdirAll(sourceAPIDir, 0777)
	if err != nil {
		return err
	}
	for id, source := range api.sources {
		file := filepath.Join(sourceAPIDir, strings.ToLower(id+".json"))
		fh, err := os.Create(file)
		if err != nil {
			return err
		}
		defer fh.Close()

		j, err := json.Marshal(source)
		if err != nil {
			return err
		}
		_, err = fh.Write(j)
		if err != nil {
			return err
		}
	}

	// Generate individual api responses.
	individualAPIDir := filepath.Join(project, "static", "api", "individual")
	err = os.MkdirAll(individualAPIDir, 0777)
	if err != nil {
		return err
	}
	for id, individual := range api.individuals {
		file := filepath.Join(individualAPIDir, strings.ToLower(id+".json"))
		fh, err := os.Create(file)
		if err != nil {
			return err
		}
		defer fh.Close()

		j, err := json.Marshal(individual)
		if err != nil {
			return err
		}
		_, err = fh.Write(j)
		if err != nil {
			return err
		}
	}

	// Generate photo api responses.
	photoAPIDir := filepath.Join(project, "static", "api", "photo")
	err = os.MkdirAll(photoAPIDir, 0777)
	if err != nil {
		return err
	}
	for id, photo := range api.photos {
		file := filepath.Join(photoAPIDir, strings.ToLower(id+".json"))
		fh, err := os.Create(file)
		if err != nil {
			return err
		}
		defer fh.Close()

		j, err := json.Marshal(photo)
		if err != nil {
			return err
		}
		_, err = fh.Write(j)
		if err != nil {
			return err
		}
	}

	// Configure for JSON headers.
	err = ioutil.WriteFile("static/api/_headers", []byte("/*  Access-Control-Allow-Origin: *  content-type: application/json; charset=utf-8"), 0644)
	if err != nil {
		return err
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
