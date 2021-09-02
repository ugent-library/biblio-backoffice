package engine

import (
	"fmt"
	"net/url"
)

type User map[string]interface{}

func (u User) ID() string {
	return u["_id"].(string)
}

func (e *Engine) GetUser(id string) (User, error) {
	user := make(User)
	if _, err := e.get(fmt.Sprintf("/user/%s", url.PathEscape(id)), nil, &user); err != nil {
		return nil, err
	}
	return user, nil
}

func (e *Engine) GetUserByUsername(username string) (User, error) {
	user := make(User)
	if _, err := e.get(fmt.Sprintf("/user/username/%s", url.PathEscape(username)), nil, &user); err != nil {
		return nil, err
	}
	return user, nil
}
