package cmd

type citationResponse struct {
	ID        int    `json:"id"`
	SourceID  string `json:"sourceid"`
	SourceRef string `json:"sourceref"`
	Detail    string `json:"detail"`
}

type citationResponses map[int]*citationResponse
