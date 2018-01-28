package cmd

type familyReferenceResponse struct {
	ID      string                       `json:"id"`
	Married string                       `json:"married"`
	Title   string                       `json:"title"`
	Husband *individualReferenceResponse `json:"husband"`
	Wife    *individualReferenceResponse `json:"wife"`
}

type familyReferenceResponses []*familyReferenceResponse