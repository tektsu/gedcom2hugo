package cmd

import "github.com/tektsu/gedcom"

type personFamily struct {
	ID        string
	Pedigree  string
	AdoptedBy string
	Mother    *personRef
	Father    *personRef
	Children  []*personRef
}

func newPersonFamily(flr *gedcom.FamilyLinkRecord) *personFamily {
	f := &personFamily{
		ID:        flr.Family.Xref,
		Pedigree:  flr.Pedigree,
		AdoptedBy: flr.AdoptedBy,
	}

	return f
}

func (r *personFamily) processParents(cc int, fr *gedcom.FamilyRecord) (int, []*sourceRef) {
	var sources []*sourceRef
	var tmpSrc []*sourceRef

	if fr == nil {
		return cc, sources
	}

	if fr.Husband != nil {
		r.Father = newPersonRef(fr.Husband)
		cc, tmpSrc = r.Father.citations(cc, fr.Husband.Name[0].Citation)
		for _, s := range tmpSrc {
			sources = append(sources, s)
		}
	}

	if fr.Wife != nil {
		r.Mother = newPersonRef(fr.Wife)
		cc, tmpSrc = r.Mother.citations(cc, fr.Wife.Name[0].Citation)
		for _, s := range tmpSrc {
			sources = append(sources, s)
		}
	}

	for _, cr := range fr.Child {
		child := newPersonRef(cr)
		cc, tmpSrc = child.citations(cc, cr.Name[0].Citation)
		for _, s := range tmpSrc {
			sources = append(sources, s)
		}
		r.Children = append(r.Children, child)
	}

	return cc, sources
}
