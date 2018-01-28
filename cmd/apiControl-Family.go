package cmd

import (
	"html/template"
	"strings"

	"encoding/json"
	"os"
	"path/filepath"

	"github.com/tektsu/gedcom"
)

func (api *apiControl) addFamilies() error {

	for _, family := range api.gc.Family {
		err := api.addFamily(family)
		if err != nil {
			return err
		}
	}

	return nil
}

func (api *apiControl) addFamily(family *gedcom.FamilyRecord) error {

	fc := newFamilyControl(api)
	fc.family = family
	fc.response = &familyResponse{
		ID:        strings.ToLower(family.Xref),
		Citations: make(citationResponses),
	}
	ref, err := api.getFamilyIndexEntry(fc.response.ID)
	if err != nil {
		return err
	}
	fc.response.Ref = ref

	if family.Husband != nil {
		father, err := api.getIndividualIndexEntry(family.Husband.Xref)
		if err != nil {
			return err
		}
		fc.response.Ref.Husband = father
		fc.response.Ref.Title = father.LastNames[0]
	} else {
		fc.response.Ref.Title = "Unknown"
	}

	if family.Wife != nil {
		mother, err := api.getIndividualIndexEntry(family.Wife.Xref)
		if err != nil {
			return err
		}
		fc.response.Ref.Wife = mother
		fc.response.Ref.Title += "/" + mother.LastNames[0]
	} else {
		fc.response.Ref.Title += "/Unknown"
	}

	for _, i := range family.Child {
		child, err := api.getIndividualIndexEntry(i.Person.Xref)
		if err != nil {
			return err
		}
		fc.response.Children = append(fc.response.Children, child)
	}

	err = fc.addEvents()
	if err != nil {
		return err
	}

	// Get notes.
	if len(family.Note) > 0 {
		for _, n := range family.Note {
			fc.response.Notes = append(fc.response.Notes, n.Note)
		}
	}

	err = fc.addPhotos()
	if err != nil {
		return err
	}

	api.families[fc.response.ID] = fc.response

	return nil
}

func (api *apiControl) getFamilyIndexEntry(id string) (*familyReferenceResponse, error) {
	id = strings.ToLower(id)
	var entry *familyReferenceResponse
	var ok bool
	if entry, ok = api.famIndex[id]; !ok {
		entry = &familyReferenceResponse{ID: id}
		api.famIndex[id] = entry
	}

	return entry, nil
}

func (api *apiControl) exportFamilyAPI() error {
	familyAPIDir := filepath.Join(api.cx.String("project"), "static", "api", "family")
	err := os.MkdirAll(familyAPIDir, 0777)
	if err != nil {
		return err
	}
	for id, family := range api.families {
		file := filepath.Join(familyAPIDir, strings.ToLower(id+".json"))
		fh, err := os.Create(file)
		if err != nil {
			return err
		}

		j, err := json.Marshal(family)
		if err != nil {
			fh.Close()
			return err
		}
		_, err = fh.Write(j)
		if err != nil {
			fh.Close()
			return err
		}
		fh.Close()
	}

	return nil
}

func (api *apiControl) exportFamilyPages() error {

	const familyPageTemplate = `---
url: "/{{ .ID }}/"
categories:
  - Family
---
<script src="/js/jquery.min.js"></script>
<script src="/js/idrisutil.js"></script>
<script src="/js/familydisplay.js"></script>

<link rel="stylesheet" href="/js/photoswipe.css">
<link rel="stylesheet" href="/js/default-skin/default-skin.css">
<script src="/js/photoswipe.min.js"></script>
<script src="/js/photoswipe-ui-default.min.js"></script>

<script>
$(document).ready(function(){
    familydisplay("{{ .ID }}")
});
</script>

<div id="display"></div>

<div id="raw"></div>
`

	familyDir := filepath.Join(api.cx.String("project"), "content", "family")
	err := os.MkdirAll(familyDir, 0777)
	if err != nil {
		return err
	}

	for _, family := range api.families {
		file := filepath.Join(familyDir, family.ID+".md")

		fh, err := os.Create(file)
		if err != nil {
			return err
		}
		defer fh.Close()

		tpl := template.New("family")
		tpl, err = tpl.Parse(familyPageTemplate)
		if err != nil {
			return err
		}
		err = tpl.Execute(fh, family)
	}

	return nil
}
