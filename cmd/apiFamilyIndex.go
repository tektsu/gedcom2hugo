package cmd

import "strings"

type familyIndex map[string]*familyReferenceResponse

func (api *apiResponse) getFamilyIndexEntry(id string) (*familyReferenceResponse, error) {
	id = strings.ToLower(id)
	var entry *familyReferenceResponse
	var ok bool
	if entry, ok = api.famIndex[id]; !ok {
		entry = &familyReferenceResponse{ID: id}
		api.famIndex[id] = entry
	}

	return entry, nil
}
