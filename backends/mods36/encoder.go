package mods36

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/ugent-library/biblio-backoffice/models"
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

	// TODO
	// [%- IF type == 'researchData' %]
	// <location>
	// 	[%- IF url %]
	// 	<url>[% url | xml_strict %]</url>
	// 	[%- END %]
	// 	<holdingExternal>
	// 		<dcterms:simpledc xmlns:dcterms="http://purl.org/dc/terms/">
	// 			[%- IF access_level %]
	// 			<dcterms:accessRights>[% access_level | xml_strict %]</dcterms:accessRights>
	// 			[%- END %]
	// 			[%- IF license == "LicenseNotListed" %]
	// 			<dcterms:license>[% other_license | xml_strict %]</dcterms:license>
	// 			[%- ELSIF license -%]
	// 			<dcterms:license>[% license | xml_strict %]</dcterms:license>
	// 			[%- END %]
	// 			[%- FOREACH f IN format %]
	// 			<dcterms:format>[% f | xml_strict %]</dcterms:format>
	// 			[%- END %]
	// 		</dcterms:simpledc>
	// 	</holdingExternal>
	// </location>

	// TODO
	// 	[%- IF access_level %]
	// <accessCondition type="accessRights">[% access_level | xml_strict %]</accessCondition>
	// [%- END %]
	// [%- IF embargo %]
	// <accessCondition type="embargoEnd">[% embargo | xml_strict %]</accessCondition>
	// [%- END %]
	// [%- END %]

	// TODO
	// [%- FOREACH rel_pub IN related_publication %]
	// <relatedItem otherType="hasRelatedObject" otherTypeAuth="pcdm">
	//   <identifier type="hdl">http://hdl.handle.net/1854/LU-[% rel_pub._id | xml_strict %]</identifier>
	//   <recordInfo lang="eng">
	// 	<recordIdentifier>[% rel_pub._id | xml_strict %]</recordIdentifier>
	//   </recordInfo>
	// </relatedItem>
	// [%- END %]

	for _, rel := range p.RelatedDataset {
		r.RelatedItem = append(r.RelatedItem, RelatedItem{
			OtherType:     "hasRelatedObject",
			OtherTypeAuth: "pcdm",
			Identifier:    []Identifier{{Type: "hdl", Value: fmt.Sprintf("http://hdl.handle.net/1854/LU-%s", rel.ID)}},
			RecordInfo: &RecordInfo{
				Lang:             "eng",
				RecordIdentifier: []RecordIdentifier{{Value: rel.ID}},
			},
		})
	}

	// TODO
	// [%- IF type == 'book' || type == 'bookEditor' %]
	// [% PROCESS mods_origin_info %]
	// [%- IF page.count || article_number  %]
	// <physicalDescription>
	// 	[%- IF page.count %]
	// 	<extent>[% page.count | xml_strict %] p.</extent>
	// 	[%- END %]
	// 	[%- IF article_number %]
	// 	<form type="epublication"/>
	// 	[%# TODO why? flag only these as e-only in non-standard way %]
	// 	[%- END %]
	// </physicalDescription>
	// [%- END %]

	// [%- IF parent.title || series_title || volume || issn || vabb_series_id %]
	// <relatedItem type="series">
	// 	[%- IF parent.title || series_title %]
	// 	<titleInfo>
	// 		<title>[% parent.title || series_title | xml_strict %]</title>
	// 	</titleInfo>
	// 	[%- END %]
	// 	[%- IF parent.short_title %]
	// 	<titleInfo type="abbreviated">
	// 		<title>[% parent.short_title | xml_strict %]</title>
	// 	</titleInfo>
	// 	[%- END %]
	// 	[%- IF volume %]
	// 	<part>
	// 		<detail type="volume">
	// 			<number>[% volume | xml_strict %]</number>
	// 		</detail>
	// 	</part>
	// 	[%- END %]
	// 	[%- FOREACH i IN issn %]
	// 	[%- IF i.search('^[\dxX-]+$') %]
	// 	<identifier type="issn">[% i | xml_strict %]</identifier>
	// 	[%- END %]
	// 	[%- END %]
	// 	[%- IF vabb_series_id %]
	// 	<identifier type="vabb-series">[% vabb_series_id | xml_strict %]</identifier>
	// 	[%- END %]
	// </relatedItem>
	// [%- END %]
	// [%- ELSIF parent %]
	// <relatedItem type="host">
	// 	[%- IF parent.title %]
	// 	<titleInfo>
	// 		<title>[% parent.title | xml_strict %]</title>
	// 	</titleInfo>
	// 	[%- END %]
	// 	[%- IF parent.short_title %]
	// 	<titleInfo type="abbreviated">
	// 		<title>[% parent.short_title | xml_strict %]</title>
	// 	</titleInfo>
	// 	[%- END %]
	// 	[%- IF type == 'bookChapter' %]
	// 	[%- FOREACH person IN editor %]
	// 	[% mods_person(person, 'editor') %]
	// 	[%- END %]
	// 	[%- IF page.count %]
	// 	<physicalDescription>
	// 		<extent>[% page.count | xml_strict %] p.</extent>
	// 	</physicalDescription>
	// 	[%- END %]
	// 	[%- END %]
	// 	[%- IF type != 'bookChapter' %]
	// 	[%- FOREACH i IN issn %]
	// 	[%- IF i.search('^[\dxX-]+$') %]
	// 	<identifier type="issn">[% i | xml_strict %]</identifier>
	// 	[%- END %]
	// 	[%- END %]
	// 	[%- END %]
	// 	[%- FOREACH i IN isbn %]
	// 	[%- IF i.search('^[\dxX-]+$') %]
	// 	<identifier type="isbn">[% i | xml_strict %]</identifier>
	// 	[%- END %]
	// 	[%- END %]
	// 	[% PROCESS mods_origin_info %]
	// 	<part>
	// 		[%- IF volume && type != 'bookChapter' %]
	// 		<detail type="volume">
	// 			<number>[% volume | xml_strict %]</number>
	// 		</detail>
	// 		[%- END %]
	// 		[%- IF issue || issue_title %]
	// 		<detail type="issue">
	// 			[%- IF issue %]
	// 			<number>[% issue | xml_strict %]</number>
	// 			[%- END %]
	// 			[%- IF issue_title %]
	// 			<title>[% issue_title | xml_strict %]</title>
	// 			[%- END %]
	// 		</detail>
	// 		[%- END %]
	// 		[%- IF article_number %]
	// 		<detail type="article-number">
	// 			<number>[% article_number | xml_strict %]</number>
	// 		</detail>
	// 		[%- END %]
	// 		[%- IF page.item('first') || page.item('last') %]
	// 		<extent unit="page">
	// 			[%- IF page.item('first') %]
	// 			<start>[% page.item('first') | xml_strict %]</start>
	// 			[%- END %]
	// 			[%- IF page.item('last') %]
	// 			<end>[% page.item('last') | xml_strict %]</end>
	// 			[%- END %]
	// 		</extent>
	// 		[%- END %]
	// 		[%- IF year %]
	// 		<date encoding="w3cdtf">[% year | xml_strict %]</date>
	// 		[%- END %]
	// 	</part>

	// 	[%- IF type == 'bookChapter' && series_title %]
	// 	<relatedItem type="series">
	// 	  <titleInfo>
	// 		  <title>[% series_title | xml_strict %]</title>
	// 	  </titleInfo>
	// 	  [%- IF volume %]
	// 	  <part>
	// 		  <detail type="volume">
	// 			  <number>[% volume | xml_strict %]</number>
	// 		  </detail>
	// 	  </part>
	// 	  [%- END %]
	// 	  [%- FOREACH i IN issn %]
	// 	  [%- IF i.search('^[\dxX-]+$') %]
	// 	  <identifier type="issn">[% i | xml_strict %]</identifier>
	// 	  [%- END %]
	// 	  [%- END %]
	// 	  [%- IF vabb_series_id %]
	// 	  <identifier type="vabb-series">[% vabb_series_id | xml_strict %]</identifier>
	// 	  [%- END %]
	// 	</relatedItem>
	// 	[%- END %]
	// </relatedItem>
	// [%- ELSE %]
	// [% PROCESS mods_origin_info %]

	if p.PageCount != "" {
		if r.PhysicalDescription == nil {
			r.PhysicalDescription = &PhysicalDescription{}
		}
		r.PhysicalDescription.Extent = append(r.PhysicalDescription.Extent, Extent{Value: p.PageCount})
	}
	if p.ArticleNumber != "" {
		if r.PhysicalDescription == nil {
			r.PhysicalDescription = &PhysicalDescription{}
		}
		r.PhysicalDescription.Form = append(r.PhysicalDescription.Form, Form{Authority: "marcform", Value: "electronic"})
	}

	for _, val := range p.ResearchField {
		r.Subject = append(r.Subject, Subject{Occupation: []Occupation{{Lang: "end", Value: val}}})
	}
	for _, val := range p.Keyword {
		r.Subject = append(r.Subject, Subject{Topic: []Topic{{Lang: "und", Value: val}}})
	}

	// TODO
	// [%- IF best_file %]
	//   [%- IF best_file.change %]
	//   <accessCondition type="accessRights">info:eu-repo/semantics/embargoedAccess</accessCondition>
	//   <accessCondition type="embargoEnd">[% best_file.change.on | xml_strict %]</accessCondition>
	//   [%- ELSIF best_file.access == 'open'  %]
	//   <accessCondition type="accessRights">info:eu-repo/semantics/openAccess</accessCondition>
	//   [%- ELSIF best_file.access == 'restricted'  %]
	//   <accessCondition type="accessRights">info:eu-repo/semantics/restrictedAccess</accessCondition>
	//   [%- ELSIF best_file.access == 'private'  %]
	//   <accessCondition type="accessRights">info:eu-repo/semantics/closedAccess</accessCondition>
	//   [%- END %]
	// [%- END %]

	// TODO
	// [%- FOREACH f IN file %]
	// <location>
	// 	<url displayLabel="[% f.name | xml_strict %]" access="raw object">[% _config.uri_base | xml_strict %]/publication/[% _id | xml_strict %]/file/[% f._id | xml_strict %]</url>
	// 	<holdingExternal>
	// 		<dcterms:simpledc xmlns:dcterms="http://purl.org/dc/terms/">
	// 			[%- IF f.change %]
	// 			<dcterms:accessRights>info:eu-repo/semantics/embargoedAccess</dcterms:accessRights>
	// 			<dcterms:accessRights>info:eu-repo/date/embargoEnd/[% f.change.on | xml_strict %]</dcterms:accessRights>
	// 			[%- ELSIF f.access == 'open'  %]
	// 			<dcterms:accessRights>info:eu-repo/semantics/openAccess</dcterms:accessRights>
	// 			[%- ELSIF f.access == 'restricted'  %]
	// 			<dcterms:accessRights>info:eu-repo/semantics/restrictedAccess</dcterms:accessRights>
	// 			[%- ELSIF f.access == 'private'  %]
	// 			<dcterms:accessRights>info:eu-repo/semantics/closedAccess</dcterms:accessRights>
	// 			[%- END %]
	// 			[%- IF date_approved %]
	// 			<dcterms:valid>[% date_approved | xml_strict %]</dcterms:valid>
	// 			[%- END %]
	// 			[%- IF f.content_type %]
	// 			<dcterms:format>https://www.iana.org/assignments/media-types/[% f.content_type | xml_strict %]</dcterms:format>
	// 			[%- END %]
	// 			<dcterms:coverage>[% f.kind | xml_strict %]</dcterms:coverage>
	// 			[%- IF f.kind == 'dataset' %]
	// 			<dcterms:type>http://purl.org/dc/dcmitype/Dataset</dcterms:type>
	// 			[%- ELSE %]
	// 			<dcterms:type>http://purl.org/dc/dcmitype/Text</dcterms:type>
	// 			[%- END %]
	// 			[%- IF f.size %]
	// 			<dcterms:extent>[% f.size | xml_strict %] bytes</dcterms:extent>
	// 			[%- END %]
	// 			<dcterms:title>[% f.name | xml_strict %]</dcterms:title>
	// 		</dcterms:simpledc>
	// 	</holdingExternal>
	// </location>
	// [%- END %]

	// TODO
	// [%- FOREACH a IN alternative_location %]
	// <location>
	// 	<url access="object in context">[% a.url.trim | url | xml_strict %]</url>
	// 	<holdingExternal>
	// 		<dcterms:simpledc xmlns:dcterms="http://purl.org/dc/terms/">
	// 			<dcterms:accessRights>[% a.access | xml_strict %]</dcterms:accessRights>
	// 			[%- IF date_submitted %]
	// 			<dcterms:valid>[% date_submitted | xml_strict %]</dcterms:valid>
	// 			[%- END %]
	// 			<dcterms:coverage>[% a.kind | xml_strict %]</dcterms:coverage>
	// 			<dcterms:type>http://purl.org/dc/dcmitype/InteractiveResource</dcterms:type>
	// 		</dcterms:simpledc>
	// 	</holdingExternal>
	// </location>
	// [%- END %]

	// TODO
	// [%- IF plain_text_cite.fwo %]
	// <note type="preferred citation" lang="eng">[% plain_text_cite.fwo | xml_strict %]</note>
	// [%- END %]

	if p.AdditionalInfo != "" {
		r.Note = append(r.Note, Note{Type: "content", Lang: "und", Value: p.AdditionalInfo})
	}

	var recordInfoNotes []RecordInfoNote
	if p.Status == "private" {
		recordInfoNotes = append(recordInfoNotes, RecordInfoNote{Type: "ugent-submission-status", Value: "unsubmitted"})
	} else if p.Status == "deleted" && p.HasBeenPublic {
		recordInfoNotes = append(recordInfoNotes, RecordInfoNote{Type: "ugent-submission-status", Value: "pdeleted"})
	} else {
		recordInfoNotes = append(recordInfoNotes, RecordInfoNote{Type: "ugent-submission-status", Value: p.Status})
	}
	if p.Creator != nil && len(p.Creator.UGentID) > 0 {
		recordInfoNotes = append(recordInfoNotes, RecordInfoNote{Type: "ugent-creator", Value: p.Creator.UGentID[0]})
	}
	if p.SourceRecord != "" {
		recordInfoNotes = append(recordInfoNotes, RecordInfoNote{Type: "source note", Value: p.SourceRecord})
	}
	if p.SourceDB != "" && p.SourceID != "" {
		recordInfoNotes = append(recordInfoNotes, RecordInfoNote{Type: "source identifier", Value: fmt.Sprintf("%s:%s", p.SourceDB, p.SourceID)})
	}
	r.RecordInfo = &RecordInfo{
		Lang:                "eng",
		RecordContentSource: []RecordContentSource{{Value: "Ghent University Library"}},
		RecordIdentifier:    []RecordIdentifier{{Value: fmt.Sprintf("pug01:%s", p.ID)}},
		RecordCreationDate: []RecordCreationDate{{
			Encoding: "w3cdtf",
			Value:    p.DateCreated.UTC().Format(time.RFC3339),
		}},
		RecordChangeDate: []RecordChangeDate{{
			Encoding: "w3cdtf",
			Value:    p.DateUpdated.UTC().Format(time.RFC3339),
		}},
		LanguageOfCataloging: []LanguageOfCataloging{{LanguageTerm: LanguageTerm{Authority: "iso639-2b", Type: "code", Value: "eng"}}},
		RecordInfoNote:       recordInfoNotes,
	}
	// TODO
	// <recordInfo lang="eng">
	// 	[%- FOREACH fund IN ecoom %]
	// 	  [% IF fund.value.weight %]
	// 	  <recordInfoNote type="ecoom-[% fund.key %]-weight">[% fund.value.weight %]</recordInfoNote>
	// 	  [% END %]
	// 	  [% IF fund.value.css %]
	// 	  <recordInfoNote type="ecoom-[% fund.key %]-css">[% fund.value.css %]</recordInfoNote>
	// 	  [% END %]
	// 	  [% IF fund.value.international_collaboration.defined %]
	// 	  <recordInfoNote type="ecoom-[% fund.key %]-international-collaboration">[% fund.value.international_collaboration == 1 ? 'true' : 'false' %]</recordInfoNote>
	// 	  [% END %]
	// 	  [% FOREACH sector IN fund.value.sector %]
	// 	  <recordInfoNote type="ecoom-[% fund.key %]-sector">[% sector %]</recordInfoNote>
	// 	  [% END %]
	// 	  <recordInfoNote type="ecoom-[% fund.key %]-validation">[% fund.value.first_year %]</recordInfoNote>
	// 	[%- END %]
	// </recordInfo>

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

func NewRecord() *Record {
	return &Record{
		XmlnsXsi:          XmlnsXsi,
		XmlnsXlink:        XmlnsXlink,
		XsiSchemaLocation: XsiSchemaLocation,
		Version:           Version,
	}
}

type Record struct {
	XMLName           xml.Name `xml:"http://www.loc.gov/mods/v3 mods"`
	XmlnsXsi          string   `xml:"xmlns:xsi,attr"`
	XmlnsXlink        string   `xml:"xmlns:xlink,attr"`
	XsiSchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	Version           string   `xml:"version,attr"`

	Abstract            []Abstract           `xml:",omitempty"`
	AccessCondition     []AccessCondition    `xml:",omitempty"`
	Classification      []Classification     `xml:",omitempty"`
	Genre               []Genre              `xml:",omitempty"`
	Identifier          []Identifier         `xml:",omitempty"`
	Language            []Language           `xml:",omitempty"`
	TitleInfo           []TitleInfo          `xml:",omitempty"`
	Name                []Name               `xml:",omitempty"`
	OriginInfo          []OriginInfo         `xml:",omitempty"`
	PhysicalDescription *PhysicalDescription `xml:",omitempty"`
	Subject             []Subject            `xml:",omitempty"`
	Note                []Note               `xml:",omitempty"`
	RelatedItem         []RelatedItem        `xml:",omitempty"`
	RecordInfo          *RecordInfo          `xml:",omitempty"`
}

type RecordInfo struct {
	XMLName              xml.Name               `xml:"recordInfo"`
	Lang                 string                 `xml:"lang,attr,omitempty"`
	RecordContentSource  []RecordContentSource  `xml:",omitempty"`
	RecordIdentifier     []RecordIdentifier     `xml:",omitempty"`
	RecordCreationDate   []RecordCreationDate   `xml:",omitempty"`
	RecordChangeDate     []RecordChangeDate     `xml:",omitempty"`
	LanguageOfCataloging []LanguageOfCataloging `xml:",omitempty"`
	RecordInfoNote       []RecordInfoNote       `xml:",omitempty"`
}

type RecordContentSource struct {
	XMLName xml.Name `xml:"recordContentSource"`
	Value   string   `xml:",chardata"`
}

type RecordIdentifier struct {
	XMLName xml.Name `xml:"recordIdentifier"`
	Value   string   `xml:",chardata"`
}

type RecordCreationDate struct {
	XMLName  xml.Name `xml:"recordCreationDate"`
	Encoding string   `xml:"encoding,attr,omitempty"`
	Value    string   `xml:",chardata"`
}

type RecordChangeDate struct {
	XMLName  xml.Name `xml:"recordChangeDate"`
	Encoding string   `xml:"encoding,attr,omitempty"`
	Value    string   `xml:",chardata"`
}

type LanguageOfCataloging struct {
	XMLName      xml.Name     `xml:"languageOfCataloging"`
	LanguageTerm LanguageTerm `xml:"languageTerm"`
}

type RecordInfoNote struct {
	XMLName xml.Name `xml:"recordInfoNote"`
	Type    string   `xml:"type,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

type RelatedItem struct {
	XMLName             xml.Name             `xml:"relatedItem"`
	OtherType           string               `xml:"otherType,attr"`
	OtherTypeAuth       string               `xml:"otherTypeAuth,attr"`
	Abstract            []Abstract           `xml:",omitempty"`
	AccessCondition     []AccessCondition    `xml:",omitempty"`
	Classification      []Classification     `xml:",omitempty"`
	Genre               []Genre              `xml:",omitempty"`
	Identifier          []Identifier         `xml:",omitempty"`
	Language            []Language           `xml:",omitempty"`
	TitleInfo           []TitleInfo          `xml:",omitempty"`
	Name                []Name               `xml:",omitempty"`
	OriginInfo          []OriginInfo         `xml:",omitempty"`
	PhysicalDescription *PhysicalDescription `xml:",omitempty"`
	Subject             []Subject            `xml:",omitempty"`
	Note                []Note               `xml:",omitempty"`
	RelatedItem         []RelatedItem        `xml:",omitempty"`
	RecordInfo          *RecordInfo          `xml:",omitempty"`
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

type PhysicalDescription struct {
	XMLName xml.Name `xml:"physicalDescription"`
	Extent  []Extent `xml:",omitempty"`
	Form    []Form   `xml:",omitempty"`
}

type Extent struct {
	XMLName xml.Name `xml:"extent"`
	Value   string   `xml:",chardata"`
}

type Form struct {
	XMLName   xml.Name `xml:"extent"`
	Authority string   `xml:"authority,attr,omitempty"`
	Value     string   `xml:",chardata"`
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

type Note struct {
	XMLName xml.Name `xml:"note"`
	Type    string   `xml:"type,attr,omitempty"`
	Lang    string   `xml:"lang,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

type Subject struct {
	XMLName    xml.Name     `xml:"subject"`
	Topic      []Topic      `xml:",omitempty"`
	Occupation []Occupation `xml:",omitempty"`
}

type Topic struct {
	XMLName xml.Name `xml:"topic"`
	Lang    string   `xml:"lang,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

type Occupation struct {
	XMLName xml.Name `xml:"occupation"`
	Lang    string   `xml:"lang,attr,omitempty"`
	Value   string   `xml:",chardata"`
}
