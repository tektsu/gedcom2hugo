package cmd

type individual struct {
	AlphaWeight int64 // Weight of individual entry based on aphabetical order
	FamilyName,
	FullName,
	LastNameFirst string // Names in different forms
}

type personData struct {
	ID        string
	Name      personName
	Aliases   []personName
	LastNames []string
	Sex       string
	Sources   []sourceRef
}

type personIndex map[string]*individual

type personName struct {
	Full       string
	Last       string
	LastFirst  string
	SourcesInd []int
}

type sourceData struct {
	ID          string
	Author      string
	Abbr        string
	Publication string
	Text        string
	Title       string
	Type        string
	File        []string
	FileNumber  []string
	Place       []string
	Date        []string
	DateViewed  []string
	URL         []string
	DocLocation []string
	RefNum      int
	Ref         string
	Periodical  string
	Volume      string
	MediaType   string
	Repository  []string
	Submitter   []string
	Page        []string
	Film        []string
}

type sourceList map[int]string

type sourceRef struct {
	RefNum int
	Ref    string
}
