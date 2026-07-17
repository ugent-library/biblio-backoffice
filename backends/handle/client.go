package handle

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/models"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	BaseURL         string
	FrontEndBaseURL string
	Prefix          string
	ADMID           string
	ADMPrivateKey   string
}

type Client struct {
	config    Config
	http      *http.Client
	sessionId string
}

type authResponse struct {
	SessionID     string `json:"sessionId,omitempty"`
	Error         string `json:"error,omitempty"`
	ID            string `json:"id,omitempty"`
	Authenticated bool   `json:"authenticated"`
}

func NewClient(c Config) *Client {
	return &Client{
		config: c,
		http:   http.DefaultClient,
	}
}

func (c *Client) put(path string, requestBody io.Reader, responseData any) (*http.Response, error) {
	req, err := c.newRequest(http.MethodPut, path, requestBody)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req, responseData)
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	url := c.config.BaseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf(`handle version="0",sessionId="%s"`, c.sessionId))

	return req, nil
}

func (c *Client) doRequest(req *http.Request, responseData any) (*http.Response, error) {
	res, err := c.http.Do(req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(responseData); err != nil {
		return res, err
	}

	return res, nil
}

func (c *Client) UpsertHandle(localId string) (*models.UpsertHandleResponse, error) {
	if !c.authenticated() {
		if err := c.authenticate(); err != nil {
			return nil, err
		}
	}

	handle := fmt.Sprintf("%s/LU-%s", c.config.Prefix, localId)
	handleReq := &models.UpsertHandleRequest{
		ResponseCode: 1,
		Handle:       handle,
		Values: []*models.HandleValue{
			{
				Index: 1,
				Type:  "URL",
				Data: map[string]any{
					"format": "string",
					"value":  fmt.Sprintf("%s/%s", c.config.FrontEndBaseURL, localId),
				},
			},
			{
				Index: 100,
				Type:  "HS_ADMIN",
				Data: map[string]any{
					"format": "admin",
					"value": map[string]any{
						"handle":      c.config.ADMID,
						"index":       200,
						"permissions": "111111111111", //TODO
					},
				},
			},
		},
	}
	handleReqBytes, _ := json.MarshalIndent(handleReq, "", "  ")
	handleRes := &models.UpsertHandleResponse{}

	_, err := c.put("/api/handles/"+handle+"?overwrite=true", bytes.NewReader(handleReqBytes), handleRes)
	if err != nil {
		return nil, err
	}

	return handleRes, nil
}

func (c *Client) authenticated() bool {
	return c.sessionId != ""
}

func (c *Client) authenticate() error {
	rawPriv, err := ssh.ParseRawPrivateKey([]byte(c.config.ADMPrivateKey))
	if err != nil {
		return err
	}
	priv := rawPriv.(*rsa.PrivateKey)

	var nonce string
	var nonceBytes []byte
	var cnonce string
	var cnonceBytes []byte = make([]byte, 16)
	var sessionId string

	// generate session (without authorization)
	cnonce2Bytes := make([]byte, 16)
	rand.Read(cnonce2Bytes)
	reqBody, _ := json.Marshal(map[string]string{
		"version": "0",
		"cnonce":  base64.StdEncoding.EncodeToString(cnonce2Bytes),
	})
	req, _ := http.NewRequest(http.MethodPost, c.config.BaseURL+"/api/sessions", bytes.NewReader(reqBody))
	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("unable to create session: %w", err)
	}
	resH := make(map[string]string)
	if err := json.NewDecoder(res.Body).Decode(&resH); err != nil {
		return fmt.Errorf("failed to decode session response: %w", err)
	}
	nonce = resH["nonce"]
	sessionId = resH["sessionId"]
	nonceBytes, _ = base64.StdEncoding.DecodeString(nonce)

	//authenticate
	rand.Read(cnonceBytes)
	cnonce = base64.StdEncoding.EncodeToString(cnonceBytes)
	msg := make([]byte, 0, len(nonceBytes)+len(cnonceBytes))
	msg = append(msg, nonceBytes...)
	msg = append(msg, cnonceBytes...)
	hashGen := sha256.New()
	hashGen.Write(msg)
	hash := hashGen.Sum(nil)
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash)
	if err != nil {
		return fmt.Errorf("unable to sign message: %w", err)
	}
	body, _ := json.MarshalIndent(map[string]string{
		"version":   "0",
		"sessionId": sessionId,
		"id":        c.config.ADMID,
		"cnonce":    cnonce,
		"type":      "HS_PUBKEY",
		"alg":       "SHA256",
		"signature": base64.StdEncoding.EncodeToString(sig),
	}, "", " ")
	req, _ = http.NewRequest(http.MethodPost, c.config.BaseURL+"/api/sessions/this", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Length", fmt.Sprint(len(body)))
	res, err = c.http.Do(req)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	authRes := &authResponse{}
	if err := json.NewDecoder(res.Body).Decode(authRes); err != nil {
		return fmt.Errorf("unable to decode authentication response: %w", err)
	}

	c.sessionId = authRes.SessionID

	if authRes.Authenticated {
		return nil
	}
	return errors.New("authentication failed: " + authRes.Error)
}
