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
