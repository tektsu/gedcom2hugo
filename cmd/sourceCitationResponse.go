package cmd

func newSourceCitationResponse() *sourceCitationResponse {
	return &sourceCitationResponse{
		Individuals: make(map[string]*individualReferenceResponse),
		Families:    make(map[string]*familyReferenceResponse),
	}
}
