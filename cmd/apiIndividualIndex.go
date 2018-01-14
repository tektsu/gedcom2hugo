package cmd

import "strings"

type individualIndex map[string]*individualReferenceResponse

func (api *apiResponse) getIndividualIndexEntry(id string) (*individualReferenceResponse, error) {
	id = strings.ToLower(id)
	var entry *individualReferenceResponse
	var ok bool
	if entry, ok = api.indIndex[id]; !ok {
		entry = &individualReferenceResponse{ID: id}
		api.indIndex[id] = entry
	}

	return entry, nil
}
