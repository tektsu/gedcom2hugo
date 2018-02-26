package cmd

import (
	"fmt"
	"html/template"
	"image"
	"os"
	"path/filepath"
	"strings"

	"encoding/json"

	"github.com/tektsu/gedcom"
)

func (api *apiControl) addPhoto(o *gedcom.ObjectRecord) *photoResponse {
	key := getPhotoKeyFromObject(o)
	if _, ok := api.photos[key]; !ok {
		api.photos[key] = &photoResponse{
			ID:    key,
			File:  filepath.Base(o.File.Name),
			Title: o.File.Title,
		}
		if o.File.Description != nil {
			api.photos[key].Description = o.File.Description.Note
		}

		for _, note := range o.Note {
			api.photos[key].Notes = append(api.photos[key].Notes, note.Note)
		}

		file, err := os.Open(o.File.Name)
		defer file.Close()
		if err != nil {
			fmt.Printf("%v\n", err)
			return api.photos[key]
		}

		img, _, err := image.DecodeConfig(file) // Image Struct
		if err != nil {
			fmt.Printf("hsgdhajsgdjhas %s: %v\n", o.File.Name, err)
			return api.photos[key]
		}

		api.photos[key].Width = img.Width
		api.photos[key].Height = img.Height
	}

	return api.photos[key]
}

func (api *apiControl) addPhotoForIndividual(o *gedcom.ObjectRecord, i *individualResponse) *photoResponse {
	response := api.addPhoto(o)

	ir, err := api.getIndividualIndexEntry(i.ID)
	if err != nil {
		fmt.Printf("%s\n", err)
		return response
	}
	response.People = append(response.People, ir)

	return response
}

func (api *apiControl) addPhotoForFamily(o *gedcom.ObjectRecord, f *familyResponse) *photoResponse {
	response := api.addPhoto(o)

	fr, err := api.getFamilyIndexEntry(f.ID)
	if err != nil {
		fmt.Printf("%s\n", err)
		return response
	}
	response.Families = append(response.Families, fr)

	return response
}

func (api *apiControl) buildFromGedcom(g *gedcom.Gedcom) error {

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

func (api *apiControl) exportPhotoAPI() error {
	photoAPIDir := filepath.Join(api.cx.String("project"), "static", "api", "photo")
	err := os.MkdirAll(photoAPIDir, 0777)
	if err != nil {
		return err
	}

	var photoIDs []string
	for id, photo := range api.photos {
		photoIDs = append(photoIDs, id)
		file := filepath.Join(photoAPIDir, strings.ToLower(id+".json"))
		fh, err := os.Create(file)
		if err != nil {
			return err
		}

		j, err := json.Marshal(photo)
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
	file := filepath.Join(photoAPIDir, strings.ToLower("list.json"))
	fh, err := os.Create(file)
	if err != nil {
		return err
	}

	j, err := json.Marshal(photoIDs)
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

	return nil
}

func (api *apiControl) exportPhotoPages() error {

	const photoPageTemplate = `---
url: "/{{ .ID }}/"
categories:
  - Photo
lead_photo: {{ .File }}
photo_key: {{ .ID  }}
---
<script src="/js/jquery.min.js"></script>
<script src="/js/photodisplay.js"></script>
<script>
$(document).ready(function(){
    photodisplay("{{ .ID }}")
});
</script>

<div id="display"></div>

<div id="raw"></div>
`

	photoDir := filepath.Join(api.cx.String("project"), "content", "media")
	err := os.MkdirAll(photoDir, 0777)
	if err != nil {
		return err
	}

	for key, photo := range api.photos {

		file := filepath.Join(photoDir, key+".md")

		fh, err := os.Create(file)
		if err != nil {
			return err
		}
		defer fh.Close()

		tpl := template.New("photo")
		tpl, err = tpl.Parse(photoPageTemplate)
		if err != nil {
			return err
		}
		err = tpl.Execute(fh, photo)
	}

	return nil
}

// Local functions

func getPhotoKeyFromObject(o *gedcom.ObjectRecord) string {
	key := "p" + strings.ToLower(strings.Replace(filepath.Base(o.File.Name), ".", "", -1))
	return key
}
