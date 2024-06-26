package mods36

import (
	"bytes"
	"encoding/xml"
	"strings"
	"text/template"

	"github.com/ugent-library/biblio-backoffice/frontoffice"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/biblio-backoffice/repositories"
)

type personWithRole struct {
	Person *frontoffice.Contributor
	Role   string
}

var funcs = template.FuncMap{
	"xml": func(s string) string {
		b := bytes.Buffer{}
		xml.EscapeText(&b, []byte(s))
		return b.String()
	},
	"personWithRole": func(p *frontoffice.Contributor, r string) personWithRole {
		return personWithRole{p, r}
	},
	"hasPrefix": strings.HasPrefix,
}

var tmpl = template.Must(template.New("").Funcs(funcs).Parse(`
{{define "person"}}
<name type="personal"{{if .Person.ID}} authority="ugent"{{end}}>
	{{if .Person.FirstName}}
    <namePart type="given">{{.Person.FirstName | xml}}</namePart>
    {{end}}
	{{if .Person.LastName}}
    <namePart type="family">{{.Person.LastName | xml}}</namePart>
    {{end}}
	{{if .Person.Name}}
    <displayForm>{{.Person.Name | xml}}</displayForm>
    {{end}}
    <role>
		{{if eq .Role "author"}}
        <roleTerm authority="marcrelator" authorityURI="http://id.loc.gov/vocabulary/relators" type="code">aut</roleTerm>
		{{else if eq .Role "editor"}}
        <roleTerm authority="marcrelator" authorityURI="http://id.loc.gov/vocabulary/relators" type="code">edt</roleTerm>
		{{else if eq .Role "promoter"}}
        <roleTerm authority="marcrelator" authorityURI="http://id.loc.gov/vocabulary/relators" type="code">ths</roleTerm>
    	{{end}}
        <roleTerm type="text" lang="eng">{{.Role | xml}}</roleTerm>
    </role>
	{{range .Person.UGentID}}
    <nameIdentifier type="ugent">{{. | xml}}</nameIdentifier>
	{{end}}
	{{if .Person.ORCID}}
    <nameIdentifier type="orcid">{{.Person.ORCID | xml}}</nameIdentifier>
	{{end}}
	{{range .Person.Affiliation}}
    {{with .UGentID}}<affiliation>{{. | xml}}</affiliation>{{end}}
	{{end}}
</name>
{{end}}

{{define "record"}}
<mods version="3.6"
    xmlns:xlink="http://www.w3.org/1999/xlink"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns="http://www.loc.gov/mods/v3"
    xsi:schemaLocation="http://www.loc.gov/mods/v3 http://www.loc.gov/standards/mods/v3/mods-3-6.xsd"
>
	{{range .Rec.DOI}}
	<identifier type="doi">{{. | xml }}</identifier>
	{{end}}
	{{if .Rec.WOSID}}
	<identifier type="wos">{{.Rec.WOSID | xml }}</identifier>
	{{end}}
	{{if .Rec.VABBID}}
	<identifier type="vabb">{{.Rec.VABBID | xml }}</identifier>
	{{end}}
	{{range .Rec.ISBN}}
	<identifier type="isbn">{{. | xml }}</identifier>
	{{end}}

	<titleInfo lang="und">
 	   <title>{{.Rec.Title | xml }}</title>
	</titleInfo>
	{{range .Rec.AlternativeTitle}}
	<titleInfo type="alternative" lang="und">
 	   <title>{{. | xml }}</title>
	</titleInfo>
	{{end}}

	{{range .Rec.Language}}
	<language>
 	   <languageTerm authority="iso639-2b" type="code">{{. | xml }}</languageTerm>
	</language>
	{{end}}

	{{if .Rec.CopyrightStatement}}
	<accessCondition lang="eng" type="useAndReproduction">{{.Rec.CopyrightStatement | xml }}</accessCondition>
	{{end}}

	{{range .Rec.AbstractFull}}
	<abstract lang="{{.Lang | xml}}">{{.Text | xml }}</abstract>
	{{end}}

	<genre authority="ugent" type="biblio">{{.Rec.Type | xml }}</genre>
	{{if .Rec.ArticleType}}
	<genre authority="ugent" type="article">{{.Rec.ArticleType | xml}}</genre>
	{{end}}
	{{if .Rec.MiscType}}
	<genre authority="ugent" type="misc">{{.Rec.MiscType | xml}}</genre>
	{{end}}
	{{if .Rec.ConferenceType}}
	<genre authority="ugent" type="conference">{{.Rec.ConferenceType | xml}}</genre>
	{{end}}

	<classification authority="biblio-review-status">public</classification>
	{{if .Rec.WOSType}}
	<classification authority="wos">{{.Rec.WOSType | xml}}</classification>
	{{end}}
	{{if .Rec.VABBType}}
	<classification authority="vabb-type">{{.Rec.VABBType | xml}}</classification>
	<classification authority="vabb-status">{{if not .Rec.IsVABBApproved}}not-{{end}}approved</classification>
	{{end}}
	<classification authority="ugent-classification">{{.Rec.Classification | xml}}</classification>
	{{if .Rec.PublicationStatus}}
	<classification authority="ugent-publication-status">{{.Rec.PublicationStatus | xml}}</classification>
	{{end}}
	{{if eq .Rec.PublicationStatus "published"}}
	<classification authority="pso">published</classification>
	{{else if eq .Rec.PublicationStatus "unpublished"}}
	<classification authority="pso">unpublished</classification>
	{{else if eq .Rec.PublicationStatus "inpress"}}
	<classification authority="pso">accepted-for-publication</classification>
	{{end}}
	{{if not .Rec.IsExternal}}
	<classification authority="ugent-publication-credit">ugent</classification>
	{{end}}
	{{if and .Rec.JCR .Rec.JCR.CategoryVigintile}}
	<classification authority="jcr-category-vigintile">{{print .Rec.JCR.CategoryVigintile | xml}}</classification>
	{{end}}

	{{range .Rec.Author}}
	{{template "person" (personWithRole . "author")}}
	{{end}}
	{{range .Rec.Editor}}
	{{template "person" (personWithRole . "editor")}}
	{{end}}
	{{range .Rec.Promoter}}
	{{template "person" (personWithRole . "promoter")}}
	{{end}}
	{{range .Rec.Affiliation}}
	<name type="corporate">
		{{if .Name}}
		<displayForm>{{.Name | xml}}</displayForm>
		{{else}}
		<displayForm>{{.UGentID | xml}}</displayForm>
		{{end}}
		<nameIdentifier type="ugent">{{.UGentID | xml}}</nameIdentifier>
		<role>
			<roleTerm authority="marcrelator" authorityURI="http://id.loc.gov/vocabulary/relators" type="code">sht</roleTerm>
			<roleTerm type="text" lang="eng">department</roleTerm>
		</role>
	</name>
	{{end}}

	{{with .Rec.Conference}}
	<relatedItem otherType="conference">
		<titleInfo>
			<title>{{.Name | xml}}</title>
		</titleInfo>
		{{if or .Organizer .Location .StartDate .EndDate}}
		<originInfo>
			{{if .Organizer}}
			<publisher>{{.Organizer | xml}}</publisher>
			{{end}}
			{{if .Location}}
			<place>
				<placeTerm>{{.Location | xml}}</placeTerm>
			</place>
			{{end}}
			{{if .StartDate}}
			<dateOther type="conference" encoding="w3cdtf" point="start">{{.StartDate | xml}}</dateOther>
			{{end}}
			{{if .EndDate}}
			<dateOther type="conference" encoding="w3cdtf" point="end">{{.EndDate | xml}}</dateOther>
			{{end}}
		</originInfo>
		{{end}}
	</relatedItem>
	{{end}}

	{{with .Rec.Defense}}
	<originInfo eventType="promotion">
		{{if .Location}}
		<place>
			<placeTerm>{{.Location | xml}}</placeTerm>
		</place>
		{{end}}
		{{if .Date}}
		<dateOther type="promotion" encoding="w3cdtf">{{.Date | xml}}</dateOther>
		{{end}}
	</originInfo>
	{{end}}

	{{range .Rec.Project}}
	<relatedItem otherType="project">
		{{if .IWETOID}}
		<identifier type="iweto">{{.IWETOID | xml}}</identifier>
		{{end}}
		{{if .GISMOID}}
		<identifier type="gismo-uuid">{{.GISMOID | xml}}</identifier>
		{{end}}
		{{if .Title}}
		<titleInfo>
			<title>{{.Title | xml}}</title>
		</titleInfo>
		{{end}}
		{{if or .StartDate .EndDate}}
		<originInfo>
			{{if .StartDate}}
			<dateOther type="project" encoding="w3cdtf" point="start">{{.StartDate | xml}}</dateOther>
			{{end}}
			{{if .EndDate}}
			<dateOther type="project" encoding="w3cdtf" point="end">{{.EndDate | xml}}</dateOther>
			{{end}}
		</originInfo>
		{{end}}
	</relatedItem>
	{{end}}

	{{if eq .Rec.Type "researchData"}}
	<location>
		{{if .Rec.URL}}
		<url>{{.Rec.URL | xml}}</url>
		{{end}}
		<holdingExternal>
			<dcterms:simpledc xmlns:dcterms="http://purl.org/dc/terms/">
				{{if .Rec.AccessLevel}}
				<dcterms:accessRights>{{.Rec.AccessLevel | xml}}</dcterms:accessRights>
				{{end}}
				{{if eq .Rec.License "LicenseNotListed"}}
				<dcterms:license>{{.Rec.OtherLicense | xml}}</dcterms:license>
				{{else if .Rec.License}}
				<dcterms:license>{{.Rec.License | xml}}</dcterms:license>
				{{end}}
				{{range .Rec.Format}}
				<dcterms:format>{{. | xml}}</dcterms:format>
				{{end}}
			</dcterms:simpledc>
		</holdingExternal>
	</location>
	{{if .Rec.AccessLevel}}
	<accessCondition type="accessRights">{{.Rec.AccessLevel | xml}}</accessCondition>
	{{end}}
	{{if .Rec.Embargo}}
	<accessCondition type="embargoEnd">{{.Rec.Embargo | xml}}</accessCondition>
	{{end}}
	{{end}}

	{{range .Rec.RelatedPublication}}
	<relatedItem otherType="hasRelatedObject" otherTypeAuth="pcdm">
		<identifier type="hdl">http://hdl.handle.net/1854/LU-{{.ID | xml}}</identifier>
		<recordInfo lang="eng">
			<recordIdentifier>{{.ID | xml}}</recordIdentifier>
		</recordInfo>
	</relatedItem>
	{{end}}

	{{range .Rec.RelatedDataset}}
	<relatedItem otherType="hasRelatedObject" otherTypeAuth="pcdm">
		<identifier type="hdl">http://hdl.handle.net/1854/LU-{{.ID | xml}}</identifier>
		<recordInfo lang="eng">
			<recordIdentifier>{{.ID | xml}}</recordIdentifier>
		</recordInfo>
	</relatedItem>
	{{end}}

	{{if or (eq .Rec.Type "book") (eq .Rec.Type "bookEditor")}}
		<originInfo eventType="publication">
			{{if and .Rec.Publisher .Rec.Publisher.Location}}
			<place>
				<placeTerm>{{.Rec.Publisher.Location | xml}}</placeTerm>
			</place>
			{{end}}
			{{if and .Rec.Publisher .Rec.Publisher.Name}}
			<publisher>{{.Rec.Publisher.Name | xml}}</publisher>
			{{end}}
			<dateIssued encoding="w3cdtf">{{.Rec.Year | xml}}</dateIssued>
			{{if .Rec.Edition}}
			<edition>{{.Rec.Edition | xml}}</edition>
			{{end}}
		</originInfo>

		{{if or (and .Rec.Page .Rec.Page.Count) .Rec.ArticleNumber}}
		<physicalDescription>
			{{if and .Rec.Page .Rec.Page.Count}}
			<extent>{{.Rec.Page.Count | xml}} p.</extent>
			{{end}}
			{{if .Rec.ArticleNumber}}
			<form authority="marcform">electronic</form>
			{{end}}
		</physicalDescription>
		{{end}}

		{{if or .Rec.Parent .Rec.SeriesTitle .Rec.Volume .Rec.ISSN}}
		<relatedItem type="series">
			{{if and .Rec.Parent .Rec.Parent.Title}}
			<titleInfo>
				<title>{{.Rec.Parent.Title | xml}}</title>
			</titleInfo>
			{{else if .Rec.SeriesTitle}}
			<titleInfo>
				<title>{{.Rec.SeriesTitle | xml}}</title>
			</titleInfo>
			{{end}}
			{{if and .Rec.Parent .Rec.Parent.ShortTitle}}
			<titleInfo type="abbreviated">
				<title>{{.Rec.Parent.ShortTitle | xml}}</title>
			</titleInfo>
			{{end}}
			{{if .Rec.Volume}}
			<part>
				<detail type="volume">
					<number>{{.Rec.Volume | xml}}</number>
				</detail>
			</part>
			{{end}}
			{{range .Rec.ISSN}}
			<identifier type="issn">{{. | xml}}</identifier>
			{{end}}
		</relatedItem>
		{{end}}
	{{else if .Rec.Parent}}
		<relatedItem type="host">
			{{if .Rec.Parent.Title}}
			<titleInfo>
				<title>{{.Rec.Parent.Title | xml}}</title>
			</titleInfo>
			{{end}}
			{{if .Rec.Parent.ShortTitle}}
			<titleInfo type="abbreviated">
				<title>{{.Rec.Parent.ShortTitle | xml}}</title>
			</titleInfo>
			{{end}}

			{{if eq .Rec.Type "bookChapter"}}
				{{range .Rec.Editor}}
				{{template "person" (personWithRole . "editor")}}
				{{end}}
				{{if and .Rec.Page .Rec.Page.Count}}
				<physicalDescription>
					<extent>{{.Rec.Page.Count | xml}} p.</extent>
				</physicalDescription>
				{{end}}
			{{else}}
			 	{{range .Rec.ISSN}}
				<identifier type="issn">{{. | xml}}</identifier>
				{{end}}
			{{end}}

			{{range .Rec.ISBN}}
			<identifier type="isbn">{{. | xml}}</identifier>
			{{end}}

			<originInfo eventType="publication">
				{{if and .Rec.Publisher .Rec.Publisher.Location}}
				<place>
					<placeTerm>{{.Rec.Publisher.Location | xml}}</placeTerm>
				</place>
				{{end}}
				{{if and .Rec.Publisher .Rec.Publisher.Name}}
				<publisher>{{.Rec.Publisher.Name | xml}}</publisher>
				{{end}}
				<dateIssued encoding="w3cdtf">{{.Rec.Year | xml}}</dateIssued>
				{{if .Rec.Edition}}
				<edition>{{.Rec.Edition | xml}}</edition>
				{{end}}
			</originInfo>

			<part>
				{{if and .Rec.Volume (ne .Rec.Type "bookChapter")}}
				<detail type="volume">
					<number>{{.Rec.Volume | xml}}</number>
				</detail>
				{{end}}
				{{if or .Rec.Issue .Rec.IssueTitle}}
				<detail type="issue">
					{{if .Rec.Issue}}
					<number>{{.Rec.Issue | xml}}</number>
					{{end}}
					{{if .Rec.IssueTitle}}
					<title>{{.Rec.IssueTitle | xml}}</title>
					{{end}}
				</detail>
				{{end}}
				{{if .Rec.ArticleNumber}}
				<detail type="article-number">
					<number>{{.Rec.ArticleNumber | xml}}</number>
				</detail>
				{{end}}
				{{if and .Rec.Page (or .Rec.Page.First .Rec.Page.Last)}}
				<extent unit="page">
					{{if .Rec.Page.First}}
					<start>{{.Rec.Page.First | xml}}</start>
					{{end}}
					{{if .Rec.Page.Last}}
					<end>{{.Rec.Page.Last | xml}}</end>
					{{end}}
				</extent>
				{{end}}
				<date encoding="w3cdtf">{{.Rec.Year | xml}}</date>
			</part>

			{{if and (eq .Rec.Type "bookChapter") .Rec.SeriesTitle}}
			<relatedItem type="series">
				<titleInfo>
					<title>{{.Rec.SeriesTitle | xml}}</title>
				</titleInfo>
				{{if .Rec.Volume}}
				<part>
					<detail type="volume">
						<number>{{.Rec.Volume | xml}}</number>
					</detail>
				</part>
				{{end}}
				{{range .Rec.ISSN}}
				<identifier type="issn">{{. | xml}}</identifier>
				{{end}}
			</relatedItem>
			{{end}}
		</relatedItem>
	{{else}}
		<originInfo eventType="publication">
			{{if and .Rec.Publisher .Rec.Publisher.Location}}
			<place>
				<placeTerm>{{.Rec.Publisher.Location | xml}}</placeTerm>
			</place>
			{{end}}
			{{if and .Rec.Publisher .Rec.Publisher.Name}}
			<publisher>{{.Rec.Publisher.Name | xml}}</publisher>
			{{end}}
			<dateIssued encoding="w3cdtf">{{.Rec.Year | xml}}</dateIssued>
			{{if .Rec.Edition}}
			<edition>{{.Rec.Edition | xml}}</edition>
			{{end}}
		</originInfo>

		{{if or (and .Rec.Page .Rec.Page.Count) .Rec.ArticleNumber}}
		<physicalDescription>
			{{if and .Rec.Page .Rec.Page.Count}}
			<extent>{{.Rec.Page.Count | xml}} p.</extent>
			{{end}}
			{{if .Rec.ArticleNumber}}
			<form authority="marcform">electronic</form>
			{{end}}
		</physicalDescription>
		{{end}}
	{{end}}

    {{range .Rec.Subject}}
	<subject>
		<occupation lang="eng">{{. | xml}}</occupation>
	</subject>
	{{end}}
    {{range .Rec.Keyword}}
	<subject>
		<topic lang="und">{{. | xml}}</topic>
	</subject>
	{{end}}

	{{$bestFile := .Rec.BestFile}}
	{{if $bestFile}}
		{{if $bestFile.Change}}
		<accessCondition type="accessRights">info:eu-repo/semantics/embargoedAccess</accessCondition>
		<accessCondition type="embargoEnd">{{$bestFile.Change.On | xml}}</accessCondition>
		{{else if eq $bestFile.Access "open"}}
		<accessCondition type="accessRights">info:eu-repo/semantics/openAccess</accessCondition>
		{{else if eq $bestFile.Access "restricted"}}
		<accessCondition type="accessRights">info:eu-repo/semantics/restrictedAccess</accessCondition>
		{{else if eq $bestFile.Access "private"}}
		<accessCondition type="accessRights">info:eu-repo/semantics/closedAccess</accessCondition>
		{{end}}
	{{end}}

	{{range .Rec.File}}
	<location>
		<url displayLabel="{{.Name | xml}}" access="raw object">{{$.BaseURL | xml}}/publication/{{$.Rec.ID | xml}}/file/{{.ID | xml}}</url>
		<holdingExternal>
			<dcterms:simpledc xmlns:dcterms="http://purl.org/dc/terms/">
				{{if .Change}}
				<dcterms:accessRights>info:eu-repo/semantics/embargoedAccess</dcterms:accessRights>
				<dcterms:accessRights>info:eu-repo/date/embargoEnd/{{.Change.On | xml}}</dcterms:accessRights>
				{{else if eq .Access "open"}}
				<accessCondition type="accessRights">info:eu-repo/semantics/openAccess</accessCondition>
				{{else if eq .Access "restricted"}}
				<accessCondition type="accessRights">info:eu-repo/semantics/restrictedAccess</accessCondition>
				{{else if eq .Access "private"}}
				<accessCondition type="accessRights">info:eu-repo/semantics/closedAccess</accessCondition>
				{{end}}
				{{if .ContentType}}
				<dcterms:format>https://www.iana.org/assignments/media-types/{{.ContentType | xml}}</dcterms:format>
				{{end}}
				<dcterms:coverage>{{.Kind | xml}}</dcterms:coverage>
				<dcterms:type>http://purl.org/dc/dcmitype/Text</dcterms:type>
				<dcterms:extent>{{.Size | xml}} bytes</dcterms:extent>
				<dcterms:title>{{.Name | xml}}</dcterms:title>
			</dcterms:simpledc>
		</holdingExternal>
	</location>
	{{end}}

	{{range .Rec.AlternativeLocation}}
	<location>
		<url access="object in context">{{.URL | xml}}</url>
		<holdingExternal>
			<dcterms:simpledc xmlns:dcterms="http://purl.org/dc/terms/">
				<dcterms:accessRights>{{.Access | xml}}</dcterms:accessRights>
				<dcterms:coverage>{{.Kind | xml}}</dcterms:coverage>
				<dcterms:type>http://purl.org/dc/dcmitype/InteractiveResource</dcterms:type>
			</dcterms:simpledc>
		</holdingExternal>
	</location>
	{{end}}

	{{/* TODO
	[%- IF plain_text_cite.fwo %]
	<note type="preferred citation" lang="eng">[% plain_text_cite.fwo | xml_strict %]</note>
	[%- END %]
	*/}}

	{{if .Rec.AdditionalInfo}}
	<note type="content" lang="und">{{.Rec.AdditionalInfo | xml}}</note>
	{{end}}

	<recordInfo lang="eng">
		<recordContentSource>Ghent University Library</recordContentSource>
		<recordIdentifier>pug01:{{.Rec.ID}}</recordIdentifier>
		<recordCreationDate encoding="w3cdtf">{{slice .Rec.DateCreated 0 10 | xml}}T{{slice .Rec.DateCreated 11 | xml}}Z</recordCreationDate>
		<recordChangeDate encoding="w3cdtf">{{slice .Rec.DateUpdated 0 10 | xml}}T{{slice .Rec.DateUpdated 11 | xml}}Z</recordChangeDate>
		<languageOfCataloging>
			<languageTerm authority="iso639-2b" type="code">eng</languageTerm>
		</languageOfCataloging>
		<recordInfoNote type="ugent-submission-status">{{.Rec.Status | xml}}</recordInfoNote>
		{{if and .Rec.CreatedBy .Rec.CreatedBy.UGentID}}
		<recordInfoNote type="ugent-creator">{{index .Rec.CreatedBy.UGentID 0 | xml}}</recordInfoNote>
		{{end}}
		{{if and .Rec.Source .Rec.Source.Record}}
		<recordInfoNote type="source note">{{.Rec.Source.Record | xml}}</recordInfoNote>
		{{end}}
		{{if and .Rec.Source .Rec.Source.DB .Rec.Source.ID}}
		<recordInfoNote type="source identifier">{{.Rec.Source.DB | xml}}:{{.Rec.Source.ID | xml}}</recordInfoNote>
		{{end}}
		{{range $fundName, $fund := .Rec.ECOOM}}
			{{if $fund.Weight}}
			<recordInfoNote type="ecoom-{{$fundName | xml}}-weight">{{$fund.Weight | xml}}</recordInfoNote>
			{{end}}
			{{if $fund.CSS}}
			<recordInfoNote type="ecoom-{{$fundName | xml}}-css">{{$fund.CSS | xml}}</recordInfoNote>
			{{end}}
			{{if $fund.InternationalCollaboration}}
			<recordInfoNote type="ecoom-{{$fundName | xml}}-international-collaboration">{{$fund.InternationalCollaboration | xml}}</recordInfoNote>
			{{end}}
			{{range $fund.Sector}}
			<recordInfoNote type="ecoom-{{$fundName | xml}}-sector">{{. | xml}}</recordInfoNote>
			{{end}}
			{{if $fund.Validation}}
			<recordInfoNote type="ecoom-{{$fundName | xml}}-validation">{{$fund.Validation | xml}}</recordInfoNote>
			{{end}}
		{{end}}
	</recordInfo>
</mods>
{{end}}
`))

type Encoder struct {
	repo    *repositories.Repo
	baseURL string
}

func New(repo *repositories.Repo, baseURL string) *Encoder {
	return &Encoder{
		repo:    repo,
		baseURL: baseURL,
	}
}

func (e *Encoder) encode(r *frontoffice.Record) ([]byte, error) {
	b := bytes.Buffer{}
	err := tmpl.ExecuteTemplate(&b, "record", struct {
		BaseURL string
		Rec     *frontoffice.Record
	}{
		BaseURL: e.baseURL,
		Rec:     r,
	})
	return b.Bytes(), err
}

func (e *Encoder) EncodePublication(p *models.Publication) ([]byte, error) {
	return e.encode(frontoffice.MapPublication(p, e.repo))
}

func (e *Encoder) EncodeDataset(d *models.Dataset) ([]byte, error) {
	return e.encode(frontoffice.MapDataset(d, e.repo))
}
