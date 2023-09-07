package mods36

import (
	"encoding/xml"

	"github.com/ugent-library/biblio-backoffice/internal/models"
)

// TODO copied from frontoffice handler, DRY this
var licenses = map[string]string{
	"CC0-1.0":          "Creative Commons Public Domain Dedication (CC0 1.0)",
	"CC-BY-4.0":        "Creative Commons Attribution 4.0 International Public License (CC-BY 4.0)",
	"CC-BY-SA-4.0":     "Creative Commons Attribution-ShareAlike 4.0 International Public License (CC BY-SA 4.0)",
	"CC-BY-NC-4.0":     "Creative Commons Attribution-NonCommercial 4.0 International Public License (CC BY-NC 4.0)",
	"CC-BY-ND-4.0":     "Creative Commons Attribution-NoDerivatives 4.0 International Public License (CC BY-ND 4.0)",
	"CC-BY-NC-SA-4.0":  "Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International Public License (CC BY-NC-SA 4.0)",
	"CC-BY-NC-ND-4.0":  "Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International Public License (CC BY-NC-ND 4.0)",
	"InCopyright":      "No license (in copyright)",
	"LicenseNotListed": "A specific license has been chosen by the rights holder. Get in touch with the rights holder for reuse rights.",
	"CopyrightUnknown": "Information pending",
	"":                 "No license (in copyright)",
}

var openLicenses = map[string]struct{}{
	"CC0-1.0":         {},
	"CC-BY-4.0":       {},
	"CC-BY-SA-4.0":    {},
	"CC-BY-NC-4.0":    {},
	"CC-BY-ND-4.0":    {},
	"CC-BY-NC-SA-4.0": {},
	"CC-BY-NC-ND-4.0": {},
}

type Encoder struct {
	baseURL string
}

func New(baseURL string) *Encoder {
	return &Encoder{
		baseURL: baseURL,
	}
}

func (e *Encoder) EncodePublication(p *models.Publication) ([]byte, error) {
	r := NewRecord()

	if p.Handle != "" {
		r.Identifier = append(r.Identifier, Identifier{Type: "hdl", Value: p.Handle})
	}
	if p.DOI != "" {
		r.Identifier = append(r.Identifier, Identifier{Type: "doi", Value: p.DOI})
	}
	if p.WOSID != "" {
		r.Identifier = append(r.Identifier, Identifier{Type: "wos", Value: p.WOSID})
	}
	if p.VABBID != "" {
		r.Identifier = append(r.Identifier, Identifier{Type: "vabb", Value: p.VABBID})
	}
	for _, val := range p.ISBN {
		r.Identifier = append(r.Identifier, Identifier{Type: "isbn", Value: val})
	}
	for _, val := range p.EISBN {
		r.Identifier = append(r.Identifier, Identifier{Type: "isbn", Value: val})
	}

	r.Classification = append(r.Classification, Classification{Type: "ugent-classification", Value: p.Classification})
	if p.Status == "public" {
		r.Classification = append(r.Classification, Classification{Type: "biblio-review-status", Value: p.Status})
	}
	if p.WOSType != "" {
		r.Classification = append(r.Classification, Classification{Type: "wos", Value: p.WOSType})
	}
	if p.VABBType != "" {
		r.Classification = append(r.Classification, Classification{Type: "vabb-type", Value: p.VABBType})
		if p.VABBApproved {
			r.Classification = append(r.Classification, Classification{Type: "vabb-status", Value: "approved"})
		} else {
			r.Classification = append(r.Classification, Classification{Type: "vabb-status", Value: "not-approved"})
		}
	}

	switch p.PublicationStatus {
	case "unpublished":
		r.Classification = append(r.Classification, Classification{Type: "ugent-classification-text", Value: "unpublished"})
		r.Classification = append(r.Classification, Classification{Type: "pso", Value: "unpublished"})
	case "accepted":
		// TODO use the same vocabulary everywhere, see also frontoffice handler
		r.Classification = append(r.Classification, Classification{Type: "ugent-classification-text", Value: "inpress"})
		r.Classification = append(r.Classification, Classification{Type: "pso", Value: "accepted-for-publication"})
	default:
		r.Classification = append(r.Classification, Classification{Type: "ugent-classification-text", Value: "published"})
		r.Classification = append(r.Classification, Classification{Type: "pso", Value: "published"})
	}

	if !p.Extern {
		r.Classification = append(r.Classification, Classification{Type: "ugent-publication-credit", Value: "ugent"})
	}

	// TODO move jcr info to backoffice
	// [%- IF jcr.category_vigintile %]
	// <classification authority="jcr-category-vigintile">[% jcr.category_vigintile | xml_strict %]</classification>
	// [%- END %]

	// TODO type in frontend oai is the camelcased form
	r.Genre = append(r.Genre, Genre{Authority: "ugent", Type: "biblio", Value: p.Type})
	if p.JournalArticleType != "" {
		// TOOD inconsistent type value
		r.Genre = append(r.Genre, Genre{Authority: "ugent", Type: "article", Value: p.JournalArticleType})
	}
	if p.ConferenceType != "" {
		r.Genre = append(r.Genre, Genre{Authority: "ugent", Type: "conference", Value: p.ConferenceType})
	}
	if p.MiscellaneousType != "" {
		// TOOD inconsistent type value
		r.Genre = append(r.Genre, Genre{Authority: "ugent", Type: "misc", Value: p.MiscellaneousType})
	}

	for _, val := range p.Language {
		r.Language = append(r.Language, Language{LanguageTerm: LanguageTerm{Type: "code", Authority: "iso639-2b", Value: val}})
	}

	r.TitleInfo = append(r.TitleInfo, TitleInfo{
		Lang:  "und",
		Title: &Title{Value: p.Title},
	})
	for _, val := range p.AlternativeTitle {
		r.TitleInfo = append(r.TitleInfo, TitleInfo{
			Lang:  "und",
			Type:  "alternative",
			Title: &Title{Value: val},
		})
	}

	for _, val := range p.Abstract {
		r.Abstract = append(r.Abstract, Abstract{Lang: val.Lang, Value: val.Text})
	}

	if len(p.File) > 0 {
		bestLicense := ""
		for _, f := range p.File {
			if bestLicense == "" {
				if _, isLicense := licenses[f.License]; isLicense {
					bestLicense = f.License
				}
			}
			if _, isOpenLicense := openLicenses[f.License]; isOpenLicense {
				bestLicense = f.License
				break
			}
		}

		r.AccessCondition = append(r.AccessCondition, AccessCondition{
			Type:  "useAndReproduction",
			Lang:  "eng",
			Value: licenses[bestLicense],
		})
	}

	for _, c := range p.Author {
		addContributor(r, "author", c)
	}
	for _, c := range p.Author {
		addContributor(r, "editor", c)
	}
	// TODO use promoter terminology everywhere?
	for _, c := range p.Supervisor {
		addContributor(r, "promoter", c)
	}

	for _, rel := range p.RelatedOrganizations {
		r.Name = append(r.Name, Name{
			Type: "corporate",
			DisplayForm: []DisplayForm{
				{Value: rel.Organization.Name},
			},
			NameIdentifier: []NameIdentifier{
				{Type: "ugent", Value: rel.OrganizationID},
			},
			Role: []Role{
				{RoleTerm: RoleTerm{Authority: "marcrelator", AuthorityURI: "http://id.loc.gov/vocabulary/relators", Type: "code", Value: "sht"}},
				{RoleTerm: RoleTerm{Type: "text", Lang: "eng", Value: "department"}},
			},
		})
	}

	if p.ConferenceName != "" {
		ri := RelatedItem{
			OtherType: "conference",
			TitleInfo: []TitleInfo{
				{Title: &Title{Value: p.ConferenceName}},
			},
		}
		if p.ConferenceOrganizer != "" || p.ConferenceLocation != "" || p.ConferenceStartDate != "" || p.ConferenceEndDate != "" {
			oi := OriginInfo{}
			if p.ConferenceOrganizer != "" {
				oi.Publisher = &Publisher{Value: p.ConferenceOrganizer}
			}
			if p.ConferenceLocation != "" {
				oi.Place = []Place{{PlaceTerm: []PlaceTerm{{Value: p.ConferenceLocation}}}}
			}
			if p.ConferenceStartDate != "" {
				oi.DateOther = append(oi.DateOther, DateOther{
					Type:     "conference",
					Encoding: "w3cdtf",
					Point:    "start",
					Value:    p.ConferenceStartDate,
				})
			}
			if p.ConferenceEndDate != "" {
				oi.DateOther = append(oi.DateOther, DateOther{
					Type:     "conference",
					Encoding: "w3cdtf",
					Point:    "end",
					Value:    p.ConferenceEndDate,
				})
			}
			ri.OriginInfo = []OriginInfo{oi}
		}
		r.RelatedItem = append(r.RelatedItem, ri)
	}

	if p.DefensePlace != "" || p.DefenseDate != "" {
		oi := OriginInfo{EventType: "promotion"}
		if p.DefensePlace != "" {
			oi.Place = []Place{{PlaceTerm: []PlaceTerm{{Value: p.DefensePlace}}}}
		}
		if p.DefenseDate != "" {
			oi.DateOther = []DateOther{{Type: "promotion", Encoding: "w3cdtf", Value: p.DefenseDate}}
		}
		r.OriginInfo = append(r.OriginInfo, oi)
	}

	for _, rp := range p.RelatedProjects {
		ri := RelatedItem{
			OtherType:  "project",
			Identifier: []Identifier{{Type: "iweto", Value: rp.ProjectID}},
		}
		// TODO map gismo id
		// [%- IF pro.gismo_id %]
		// <identifier type="gismo-uuid">[% pro.gismo_id | xml_strict %]</identifier>
		// [%- END %]
		if rp.Project.Title != "" {
			ri.TitleInfo = append(ri.TitleInfo, TitleInfo{Title: &Title{Value: rp.Project.Title}})
		}
		if rp.Project.StartDate != "" || rp.Project.EndDate != "" {
			oi := OriginInfo{}
			if rp.Project.StartDate != "" {
				oi.DateOther = append(oi.DateOther, DateOther{Type: "project", Encoding: "w3cdtf", Point: "start", Value: rp.Project.StartDate})
			}
			if rp.Project.EndDate != "" {
				oi.DateOther = append(oi.DateOther, DateOther{Type: "project", Encoding: "w3cdtf", Point: "end", Value: rp.Project.EndDate})
			}
			ri.OriginInfo = []OriginInfo{oi}
		}
		r.RelatedItem = append(r.RelatedItem, ri)
	}

	return xml.Marshal(r)
}

func addContributor(r *Record, role string, c *models.Contributor) {
	var marcRelator string
	switch role {
	case "author":
		marcRelator = "aut"
	case "editor":
		marcRelator = "edt"
	case "promoter":
		marcRelator = "ths"
	}

	name := Name{
		Type: "personal",
		NamePart: []NamePart{
			{Type: "given", Value: c.FirstName()},
			{Type: "family", Value: c.LastName()},
		},
		Role: []Role{
			{RoleTerm: RoleTerm{Type: "code", Authority: "marcrelator", AuthorityURI: "http://id.loc.gov/vocabulary/relators", Value: marcRelator}},
			{RoleTerm: RoleTerm{Type: "text", Lang: "eng", Value: role}},
		},
	}
	if c.PersonID != "" {
		name.Authority = "ugent"

		for _, val := range c.Person.UGentID {
			name.NameIdentifier = append(name.NameIdentifier, NameIdentifier{Type: "ugent", Value: val})
		}

		if c.Person.ORCID != "" {
			name.NameIdentifier = append(name.NameIdentifier, NameIdentifier{Type: "orcid", Value: c.Person.ORCID})
		}
		for _, val := range c.Person.Affiliations {
			name.Affiliation = append(name.Affiliation, Affiliation{Value: val.OrganizationID})
		}
	}

	r.Name = append(r.Name, name)
}

const (
	XmlnsXsi          = "http://www.w3.org/2001/XMLSchema-instance"
	XmlnsXlink        = "http://www.w3.org/1999/xlink"
	XsiSchemaLocation = "http://www.loc.gov/mods/v3 http://www.loc.gov/standards/mods/v3/mods-3-6.xsd"
	Version           = "3.6"
)

type Record struct {
	XMLName           xml.Name `xml:"http://www.loc.gov/mods/v3 mods"`
	XmlnsXsi          string   `xml:"xmlns:xsi,attr"`
	XmlnsXlink        string   `xml:"xmlns:xlink,attr"`
	XsiSchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	Version           string   `xml:"version,attr"`

	Abstract        []Abstract        `xml:",omitempty"`
	AccessCondition []AccessCondition `xml:",omitempty"`
	Classification  []Classification  `xml:",omitempty"`
	Genre           []Genre           `xml:",omitempty"`
	Identifier      []Identifier      `xml:",omitempty"`
	Language        []Language        `xml:",omitempty"`
	TitleInfo       []TitleInfo       `xml:",omitempty"`
	Name            []Name            `xml:",omitempty"`
	OriginInfo      []OriginInfo      `xml:",omitempty"`
	RelatedItem     []RelatedItem     `xml:",omitempty"`
}

type RelatedItem struct {
	XMLName   xml.Name `xml:"relatedItem"`
	OtherType string   `xml:"otherType,attr"`

	Abstract        []Abstract        `xml:",omitempty"`
	AccessCondition []AccessCondition `xml:",omitempty"`
	Classification  []Classification  `xml:",omitempty"`
	Genre           []Genre           `xml:",omitempty"`
	Identifier      []Identifier      `xml:",omitempty"`
	Language        []Language        `xml:",omitempty"`
	TitleInfo       []TitleInfo       `xml:",omitempty"`
	Name            []Name            `xml:",omitempty"`
	OriginInfo      []OriginInfo      `xml:",omitempty"`
	RelatedItem     []RelatedItem     `xml:",omitempty"`
}

func NewRecord() *Record {
	return &Record{
		XmlnsXsi:          XmlnsXsi,
		XmlnsXlink:        XmlnsXlink,
		XsiSchemaLocation: XsiSchemaLocation,
		Version:           Version,
	}
}

type Classification struct {
	XMLName xml.Name `xml:"classification"`
	Type    string   `xml:"type,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

type Genre struct {
	XMLName   xml.Name `xml:"classification"`
	Type      string   `xml:"type,attr,omitempty"`
	Authority string   `xml:"authority,attr,omitempty"`
	Value     string   `xml:",chardata"`
}

type Identifier struct {
	XMLName xml.Name `xml:"identifier"`
	Type    string   `xml:"type,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

type Language struct {
	XMLName      xml.Name     `xml:"language"`
	LanguageTerm LanguageTerm `xml:"languageTerm"`
}

type LanguageTerm struct {
	Type      string `xml:"type,attr,omitempty"`
	Authority string `xml:"authority,attr,omitempty"`
	Value     string `xml:",chardata"`
}

type AccessCondition struct {
	XMLName xml.Name `xml:"accessCondition"`
	Type    string   `xml:"type,attr,omitempty"`
	Lang    string   `xml:"lang,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

type TitleInfo struct {
	XMLName xml.Name `xml:"titleInfo"`
	Type    string   `xml:"type,attr,omitempty"`
	Lang    string   `xml:"lang,attr,omitempty"`
	Title   *Title   `xml:"title,omitempty"`
}

type Title struct {
	Value string `xml:",chardata"`
}

type Abstract struct {
	XMLName xml.Name `xml:"abstract"`
	Type    string   `xml:"type,attr,omitempty"`
	Lang    string   `xml:"lang,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

type Name struct {
	XMLName        xml.Name         `xml:"name"`
	Type           string           `xml:"type,attr,omitempty"`
	Authority      string           `xml:"authority,attr,omitempty"`
	NamePart       []NamePart       `xml:",omitempty"`
	DisplayForm    []DisplayForm    `xml:",omitempty"`
	NameIdentifier []NameIdentifier `xml:",omitempty"`
	Role           []Role           `xml:",omitempty"`
	Affiliation    []Affiliation    `xml:",omitempty"`
}

type NamePart struct {
	XMLName xml.Name `xml:"namePart"`
	Type    string   `xml:"type,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

type DisplayForm struct {
	XMLName xml.Name `xml:"displayForm"`
	Value   string   `xml:",chardata"`
}

type Role struct {
	XMLName  xml.Name `xml:"role"`
	RoleTerm RoleTerm `xml:"roleTerm"`
}

type RoleTerm struct {
	Type         string `xml:"type,attr,omitempty"`
	Lang         string `xml:"lang,attr,omitempty"`
	Authority    string `xml:"authority,attr,omitempty"`
	AuthorityURI string `xml:"authorityURI,attr,omitempty"`
	Value        string `xml:",chardata"`
}

type NameIdentifier struct {
	XMLName xml.Name `xml:"nameIdentifier"`
	Type    string   `xml:"type,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

type Affiliation struct {
	XMLName xml.Name `xml:"affiliation"`
	Value   string   `xml:",chardata"`
}

type OriginInfo struct {
	XMLName   xml.Name    `xml:"originInfo"`
	EventType string      `xml:"eventType,attr,omitempty"`
	Publisher *Publisher  `xml:",omitempty"`
	Place     []Place     `xml:",omitempty"`
	DateOther []DateOther `xml:",omitempty"`
}

type Publisher struct {
	Value string `xml:",chardata"`
}
type Place struct {
	XMLName   xml.Name    `xml:"place"`
	PlaceTerm []PlaceTerm `xml:",omitempty"`
}

type PlaceTerm struct {
	XMLName xml.Name `xml:"placeTerm"`
	Value   string   `xml:",chardata"`
}

type DateOther struct {
	XMLName  xml.Name `xml:"dateOther"`
	Type     string   `xml:"type,attr,omitempty"`
	Encoding string   `xml:"encoding,attr,omitempty"`
	Point    string   `xml:"point,attr,omitempty"`
	Value    string   `xml:",chardata"`
}
