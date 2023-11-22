package projects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/ugent-library/biblio-backoffice/models"
)

type SuggestQuery struct {
	Query string `json:"query"`
}

type GetProjectQuery struct {
	ID string `json:"id"`
}

type Config struct {
	APIUrl string
	APIKey string
}

type Client struct {
	config Config
	http   *http.Client
}

func New(c Config) *Client {
	return &Client{
		config: c,
		http: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *Client) GetProject(id string) (*models.Project, error) {
	getProjectURL := fmt.Sprintf("%s/%s", c.config.APIUrl, "get-project")

	buf := &bytes.Buffer{}
	// json.NewEncoder(buf).Encode(&GetProjectQuery{ID: id})
	json.NewEncoder(buf).Encode(&GetProjectQuery{ID: fmt.Sprintf("urn:iweto:%s", id)})

	req, err := http.NewRequest(http.MethodPost, getProjectURL, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.config.APIKey)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	src, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Project service: get-project call yielded an error: %s", src)
	}

	data := gjson.ParseBytes(src)
	p := &models.Project{
		EUProject: &models.EUProject{},
	}

	if id := data.Get("id"); id.Exists() {
		strId := id.String()
		if idx := strings.LastIndex(strId, ":"); idx != -1 {
			p.ID = strId[idx+1:]
		} else {
			p.ID = ""
		}
		// p.ID = id.String()
	}

	if name := data.Get("name.#[language==\"und\"].value"); name.Exists() {
		p.Title = name.String()
	}

	if foundingDate := data.Get("foundingDate"); foundingDate.Exists() {
		p.StartDate = foundingDate.String()
	}

	if dissolutionDate := data.Get("dissolutionDate"); dissolutionDate.Exists() {
		p.StartDate = dissolutionDate.String()
	}

	if euID := data.Get("identifier.#[propertyID=\"CORDIS\"].value"); euID.Exists() {
		p.EUProject.ID = euID.String()
	}

	if callID := data.Get("isFundedBy.hasCallNumber"); callID.Exists() {
		p.EUProject.CallID = callID.String()
	}

	if acronym := data.Get("hasAcronym.0"); acronym.Exists() {
		p.EUProject.Acronym = acronym.String()
	}

	if fp := data.Get("isFundedBy.isAwardedBy.name"); fp.Exists() {
		p.EUProject.FrameworkProgramme = fp.String()
	}

	if gismoID := data.Get("identifier.#[propertyID=\"GISMO\"].value"); gismoID.Exists() {
		p.GISMOID = gismoID.String()
	}

	if iwetoID := data.Get("identifier.#[propertyID=\"CORDIS\"].value"); iwetoID.Exists() {
		p.IWETOID = iwetoID.String()
	}

	// iweto_id not filled in everywhere, but should be same as id for now
	if p.IWETOID == "" {
		p.IWETOID = p.ID
	}

	return p, nil
}

func (c *Client) SuggestProjects(q string) ([]models.Completion, error) {
	suggestUrl := fmt.Sprintf("%s/%s", c.config.APIUrl, "suggest-projects")

	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(&SuggestQuery{Query: q})

	req, err := http.NewRequest(http.MethodPost, suggestUrl, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.config.APIKey)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	src, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Project service: suggest-projects call yielded an error: %s", src)
	}

	completions := make([]models.Completion, 0)

	if data := gjson.ParseBytes(src).Get("data"); data.Exists() {
		if data.IsArray() {
			for _, hit := range data.Array() {
				c := models.Completion{}

				if id := hit.Get("id"); id.Exists() {
					strId := id.String()
					if idx := strings.LastIndex(strId, ":"); idx != -1 {
						c.ID = strId[idx+1:]
					} else {
						c.ID = ""
					}
					// c.ID = id.String()
				}

				if name := hit.Get("name.#[language==\"und\"].value"); name.Exists() {
					c.Heading = name.String()
				}

				if acronym := hit.Get("hasAcronym.0"); acronym.Exists() {
					c.Description = acronym.String()
				}

				completions = append(completions, c)
			}
		}
	}

	return completions, nil
}
