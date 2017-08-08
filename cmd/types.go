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
	Sources   []SourceRef
}

type personIndex map[string]*individual

type personName struct {
	Full       string
	Last       string
	LastFirst  string
	SourcesInd []int
}

type SourceList map[int]string

type SourceRef struct {
	RefNum int
	Ref    string
}

type sourceData struct {
	ID          string
	Author      string
	Abbr        string
	Publication string
	Text        string
	Title       string
	Type        string
	Form        string
	File        []string
	FileNumber  []string
	Place       []string
	Date        []string
	DateViewed  []string
	URL         []string
	DocLocation []string
	RefNum      int
	Ref         string
}
