package cmd

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/tektsu/gedcom"
)

// personIndex is a cache of information about each individual/
type personIndex map[string]*personIndexEntry

// personIndexEntry contains cached information about a person.
// It lists names in various forms, and calculates a weight based on
// alphabetical order of the names.
type personIndexEntry struct {
	AlphaWeight int64
	GivenName,
	FamilyName,
	FullName,
	LastNameFirst string
}

// newPersonIndexEntry builds a personIndexEntry from a
// gedcom.IndividualRecord.
func newPersonIndexEntry(i *gedcom.IndividualRecord) *personIndexEntry {
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

// indSortableList is a sortable list used to alphabetize the personIndex.
type indSortableList []indSortable

func (l indSortableList) Len() int           { return len(l) }
func (l indSortableList) Less(i, j int) bool { return l[i].Name < l[j].Name }
func (l indSortableList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

type indSortable struct {
	ID, Name string
}

// newPersonIndex builds a personIndex from a gedcom.Gedcom.
func newPersonIndex(gc *gedcom.Gedcom) personIndex {
	idx := make(personIndex)

	//Build the index.
	for _, i := range gc.Individual {
		idx[i.Xref] = newPersonIndexEntry(i)
	}

	// Assign the weights.
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
