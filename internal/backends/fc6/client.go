package fc6

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	URL      string
	Username string
	Password string
}

type Client struct {
	config Config
	http   *http.Client
}

func New(conf Config) *Client {
	if conf.URL == "" {
		conf.URL = "http://localhost:8080/fcrepo/rest"
	}
	if conf.Username == "" {
		conf.Username = "fedoraAdmin"
	}
	if conf.Password == "" {
		conf.Password = "fedoraAdmin"
	}
	c := &Client{
		config: conf,
		http: &http.Client{
			Timeout: 3 * time.Second,
		},
	}

	c.init()

	return c
}

func (c *Client) init() {
	ok, err := c.resourceExists("/biblio-objects")
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Print("creating /biblio-objects")
		url := c.config.URL + "/biblio-objects"
		req, err := http.NewRequest(http.MethodPut, url, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth("fedoraAdmin", "fedoraAdmin")
		req.Header.Set("Content-Type", "text/turtle")
		if _, err = c.http.Do(req); err != nil {
			log.Fatal(err)
		}
	}
}

func (c *Client) resourceExists(p string) (bool, error) {
	url := c.config.URL + p
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	req.SetBasicAuth("fedoraAdmin", "fedoraAdmin")
	req.Header.Set("Content-Type", "text/turtle")

	res, err := c.http.Do(req)
	if err != nil {
		log.Print(err)
		return false, err
	}
	return res.StatusCode == http.StatusOK, nil
}

func (c *Client) markResourceAsFile(p string) error {
	url := c.config.URL + p
	body := strings.NewReader(`  
		PREFIX pcdm: <http://pcdm.org/models#>
		INSERT {
		<> a pcdm:File
		} WHERE {
		}
	`)
	req, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return err
	}
	req.SetBasicAuth("fedoraAdmin", "fedoraAdmin")
	req.Header.Set("Content-Type", "application/sparql-update")

	_, err = c.http.Do(req)
	return err
}
