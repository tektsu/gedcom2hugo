package cmd

import (
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

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

type photoResponses map[string]*photoResponse

func (ic *individualControl) addPhotos() error {
	for _, o := range ic.individual.Object {
		if o.File.Form != "jpg" && o.File.Form != "png" {
			continue
		}
		p := ic.api.addPhotoForIndividual(o, ic.response)
		ic.response.Photos = append(ic.response.Photos, p)
	}

	return nil
}

func (ic *individualControl) addTopPhoto() error {
	if ic.individual.Photo != nil {
		p := ic.api.addPhoto(ic.individual.Photo) // Don't use addPhotoForIndividual() here of there will be a duplicate on the photo page
		ic.response.TopPhoto = p
	}

	return nil
}

func (fc *familyControl) addPhotos() error {
	for _, o := range fc.family.Object {
		if o.File.Form != "jpg" && o.File.Form != "png" {
			continue
		}
		p := fc.api.addPhotoForFamily(o, fc.response)
		fc.response.Photos = append(fc.response.Photos, p)
	}

	return nil
}

func (api *apiResponse) exportPhotoAPI() error {
	photoAPIDir := filepath.Join(api.cx.String("project"), "static", "api", "photo")
	err := os.MkdirAll(photoAPIDir, 0777)
	if err != nil {
		return err
	}
	for id, photo := range api.photos {
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

	return nil
}

func (api *apiResponse) exportPhotoPages() error {

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
