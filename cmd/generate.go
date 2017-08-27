package cmd

import (
	"bytes"
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

// Generate reads the GEDCOM file and builds the Hugo input files.
func Generate(cx *cli.Context) error {

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
		if people[id].Living {
			continue
		}
		//if id == "I126" {
		//	fmt.Printf("%s\n", spew.Sdump(person))
		//}
		file := filepath.Join(personDir, strings.ToLower(id+".md"))
		fh, err := os.Create(file)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		defer fh.Close()

		data := newPersonTmplData(person)

		tpl := template.New("person")
		funcs := template.FuncMap{
			"add":     func(x, y int) int { return x + y },
			"ToLower": strings.ToLower,
		}
		tpl.Funcs(funcs)
		tpl, err = tpl.Parse(personPageTemplate)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		err = tpl.Execute(fh, data)
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
