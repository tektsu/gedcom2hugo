package cmd

type photoResponse struct {
	ID       string                       `json:"id"`
	File     string                       `json:"file"`
	Title    string                       `json:"title"`
	Height   int                          `json:"height"`
	Width    int                          `json:"width"`
	People   individualReferenceResponses `json:"people"`
	Families familyReferenceResponses     `json:"families"`
	Notes    []string                     `json:"notes"`
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
		p := ic.api.addPhotoForIndividual(ic.individual.Photo, ic.response)
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
