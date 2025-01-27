package models

import "fmt"

type HandleValue struct {
	Timestamp string `json:"timestamp,omitempty"`
	Type      string `json:"type"`
	Index     int    `json:"index"`
	Ttl       int    `json:"ttl,omitempty"`
	Data      any    `json:"data"`
}

type UpsertHandleRequest struct {
	Handle       string         `json:"handle"`
	ResponseCode int            `json:"responseCode"`
	Values       []*HandleValue `json:"values,omitempty"`
}

type UpsertHandleResponse struct {
	Handle       string `json:"handle"`
	ResponseCode int    `json:"responseCode"`
	Message      string `json:"message,omitempty"`
}

func (h *UpsertHandleResponse) IsSuccess() bool {
	return h.ResponseCode == 1
}

func (h *UpsertHandleResponse) GetFullHandleURL() string {
	if !h.IsSuccess() {
		return ""
	}
	if h.Handle == "" {
		return ""
	}
	//http://hdl.handle.net/<handle> where <handle> is <prefix>/LU-<localId>
	return fmt.Sprintf("http://hdl.handle.net/%s", h.Handle)
}
