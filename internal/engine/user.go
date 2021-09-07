package engine

import (
	"fmt"
	"net/url"
)

// TODO add missing fields
type User struct {
	ID        string `json:"_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
}

func (e *Engine) GetUser(id string) (*User, error) {
	user := &User{}
	if _, err := e.get(fmt.Sprintf("/user/%s", url.PathEscape(id)), nil, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (e *Engine) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	if _, err := e.get(fmt.Sprintf("/user/username/%s", url.PathEscape(username)), nil, user); err != nil {
		return nil, err
	}
	return user, nil
}
