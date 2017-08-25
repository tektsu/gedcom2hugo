package cmd

type personTmplData struct {
	ID            string
	Name          *personName
	Aliases       []*personName
	LastNames     []string
	Sex           string
	Sources       []*sourceRef
	ParentsFamily []personFamily
}

type personFamily struct {
	ID        string
	Pedigree  string
	AdoptedBy string
	Mother    personRef
	Father    personRef
	Children  []personRef
}

type personRef struct {
	ID         string
	Name       string
	Sex        string
	SourcesInd []int
}

type sourceTmplData struct {
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
