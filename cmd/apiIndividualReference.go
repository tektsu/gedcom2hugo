package cmd

type individualReferenceResponse struct {
	ID        string   `json:"id"`
	Sex       string   `json:"sex"`
	Name      string   `json:"name"`
	Birth     string   `json:"birth"`
	Death     string   `json:"death"`
	Photo     string   `json:"photo"`
	LastNames []string `json:"lastnames"`
}

type individualReferenceResponses []*individualReferenceResponse
