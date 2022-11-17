package models

import "fmt"

/*
copy from handle-server-api
*/
type HandleData struct {
	Url    string `json:"url"`
	Format string `json:"format"`
}

type HandleValue struct {
	Timestamp string      `json:"timestamp"`
	Type      string      `json:"type"`
	Index     int         `json:"index"`
	Ttl       int         `json:"ttl"`
	Data      *HandleData `json:"data"`
}

type Handle struct {
	Handle       string         `json:"handle"`
	ResponseCode int            `json:"responseCode"`
	Values       []*HandleValue `json:"values,omitempty"`
	Message      string         `json:"message,omitempty"`
}

func (h *Handle) IsSuccess() bool {
	return h.ResponseCode == 1
}

func (h *Handle) GetFullHandleURL() string {
	if !h.IsSuccess() {
		return ""
	}
	if h.Handle == "" {
		return ""
	}
	//http://hdl.handle.net/<handle> where <handle> is <prefix>/LU-<localId>
	return fmt.Sprintf("http://hdl.handle.net/%s", h.Handle)
}
