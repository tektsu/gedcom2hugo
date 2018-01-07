package cmd

type photoResponse struct {
	ID     string `json:"id"`
	File   string `json:"file"`
	Title  string `json:"title"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
	//People photoPersonIndex `json:"people"`
	Notes []string `json:"notes"`
}

type photoResponses map[string]*photoResponse

func (ic *individualControl) addPhotos() error {
	for _, o := range ic.individual.Object {
		if o.File.Form != "jpg" && o.File.Form != "png" {
			continue
		}
		p := ic.api.addPhoto(o, ic.response)
		ic.response.Photos = append(ic.response.Photos, p)
	}

	return nil
}

func (ic *individualControl) addTopPhoto() error {
	if ic.individual.Photo != nil {
		p := ic.api.addPhoto(ic.individual.Photo, ic.response)
		ic.response.TopPhoto = p
	}

	return nil
}
