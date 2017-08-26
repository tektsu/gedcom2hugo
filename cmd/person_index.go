package cmd

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/tektsu/gedcom"
)

type personIndexEntry struct {
	AlphaWeight int64 // Weight of index entry based on aphabetical order
	GivenName,
	FamilyName,
	FullName,
	LastNameFirst string // Names in different forms
}

type personIndex map[string]*personIndexEntry

func newIndividual(i *gedcom.IndividualRecord) *personIndexEntry {
	r := &personIndexEntry{}

	if len(i.Name) > 0 {
		given, family := extractNames(i.Name[0].Name)
		r.GivenName = given
		r.FamilyName = family
		r.FullName = fmt.Sprintf("%s %s", given, family)
		r.LastNameFirst = fmt.Sprintf("%s, %s", family, given)
	}

	return r
}

type indSortable struct {
	ID, Name string
}

type indSortableList []indSortable

func (l indSortableList) Len() int           { return len(l) }
func (l indSortableList) Less(i, j int) bool { return l[i].Name < l[j].Name }
func (l indSortableList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func newPersonIndex(gc *gedcom.Gedcom) personIndex {
	idx := make(personIndex)

	//Build the index.
	for _, i := range gc.Individual {
		idx[i.Xref] = newIndividual(i)
	}

	// Assign the weights
	l := make(indSortableList, len(idx))
	i := 0
	for id, ind := range idx {
		l[i] = indSortable{
			ID:   id,
			Name: ind.LastNameFirst,
		}
		i++
	}
	sort.Sort(l)

	var weight int64 = 1
	for _, entry := range l {
		idx[entry.ID].AlphaWeight = weight
		weight++
	}

	return idx
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
