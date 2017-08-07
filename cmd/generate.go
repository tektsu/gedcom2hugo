package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/tektsu/gedcom"
	"github.com/urfave/cli"
)

var sourceList SourceList

func add(x, y int) int {
	return x + y
}

// Generate reads the GEDCOM file and builds the Hugo input files.
func Generate(cx *cli.Context) error {

	project := cx.String("project")

	gc, err := readGedcom(cx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	people, err := createPersonIndex(cx, gc)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// Generate Source Pages.
	sourceList = make(SourceList)
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

		data, err := newSourceData(cx, source)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		sourceList[data.RefNum] = data.Ref

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

	spew.Config.MaxDepth = 5
	for _, person := range gc.Individual {
		id := person.Xref
		//if id == "I126" {
		//	fmt.Printf("%s\n", spew.Sdump(person))
		//}
		file := filepath.Join(personDir, strings.ToLower(id+".md"))
		fh, err := os.Create(file)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		defer fh.Close()

		data, err := newPersonData(cx, people, person)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		tpl := template.New("person")
		funcs := template.FuncMap{"add": add}
		tpl.Funcs(funcs)
		tpl, err = tpl.Parse(personPageTemplate)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		err = tpl.Execute(fh, data)
	}
	return nil
}

// readGedcom reads the GEDCOM file specified ib the context into memory.
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
	gc, err = decoder.Decode()
	if err != nil {
		return gc, err
	}

	return gc, nil
}

type indSortable struct {
	ID, Name string
}

type indSortableList []indSortable

func (l indSortableList) Len() int           { return len(l) }
func (l indSortableList) Less(i, j int) bool { return l[i].Name < l[j].Name }
func (l indSortableList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

// individualIndex creates a map information about individuals keyed to Individual ID.
func createPersonIndex(cx *cli.Context, gc *gedcom.Gedcom) (*personIndex, error) {
	index := make(personIndex)

	//Build the index.
	for _, i := range gc.Individual {
		index[i.Xref] = &individual{}

		if len(i.Name) > 0 {
			given, family := extractNames(i.Name[0].Name)
			index[i.Xref].FamilyName = family
			index[i.Xref].FullName = fmt.Sprintf("%s %s", given, family)
			index[i.Xref].LastNameFirst = fmt.Sprintf("%s, %s", family, given)
		}
	}

	// Assign weights
	l := make(indSortableList, len(index))
	i := 0
	for id, ind := range index {
		l[i] = indSortable{
			ID:   id,
			Name: ind.LastNameFirst,
		}
		i++
	}
	sort.Sort(l)

	var weight int64 = 1
	for _, entry := range l {
		index[entry.ID].AlphaWeight = weight
		weight++
	}

	return &index, nil
}

// extractNames splits a full name into a given name and a family name.
func extractNames(name string) (string, string) {
	var given, family string

	re := regexp.MustCompile("^([^/]+) +/(.+)/$")
	names := re.FindStringSubmatch(name)
	given = names[1]
	family = names[2]

	return given, family
}
