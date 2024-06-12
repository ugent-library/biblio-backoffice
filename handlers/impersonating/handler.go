package impersonating

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/views"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/httperror"
)

type BindAddCImpersonationSuggest struct {
	FirstName string `query:"first_name"`
	LastName  string `query:"last_name"`
}

type BindCreateImpersonation struct {
	ID string `form:"id"`
}

func AddImpersonation(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	if c.OriginalUser != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(fmt.Errorf("already impersonating user %s", c.OriginalUser.ID)))
		return
	}

	views.ShowModal(views.AddImpersonation(c)).Render(r.Context(), w)
}

func AddImpersonationSuggest(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	if c.OriginalUser != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(fmt.Errorf("already impersonating user %s", c.OriginalUser.ID)))
		return
	}

	b := BindAddCImpersonationSuggest{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	hits, err := c.UserSearchService.SuggestUsers(b.FirstName + " " + b.LastName)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	// exclude the current user
	for i, hit := range hits {
		if hit.ID == c.User.ID {
			if i == 0 {
				hits = hits[1:]
			} else {
				hits = append(hits[:i], hits[i+1:]...)
			}
			break
		}
	}

	views.AddImpersonationSuggest(c, b.FirstName, b.LastName, hits).Render(r.Context(), w)
}

func CreateImpersonation(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	if c.OriginalUser != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(fmt.Errorf("already impersonating user %s", c.OriginalUser.ID)))
		return
	}

	b := BindCreateImpersonation{}
	if err := bind.Request(r, &b); err != nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(err))
		return
	}

	user, err := c.UserService.GetUser(b.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	// TODO handle user not found

	session, err := c.SessionStore.Get(r, c.SessionName)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	session.Values[ctx.OriginalUserIDKey] = c.User.ID
	session.Values[ctx.OriginalUserRoleKey] = c.UserRole
	session.Values[ctx.UserIDKey] = user.ID
	session.Values[ctx.UserRoleKey] = "user"

	if err = session.Save(r, w); err != nil {
		c.HandleError(w, r, err)
		return
	}

	http.Redirect(w, r, c.PathTo("home").String(), http.StatusFound)
}

func DeleteImpersonation(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	if c.OriginalUser == nil {
		c.HandleError(w, r, httperror.BadRequest.Wrap(errors.New("missing impersonation")))
		return
	}

	session, err := c.SessionStore.Get(r, c.SessionName)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	session.Values[ctx.UserIDKey] = session.Values[ctx.OriginalUserIDKey]
	session.Values[ctx.UserRoleKey] = session.Values[ctx.OriginalUserRoleKey]
	delete(session.Values, ctx.OriginalUserIDKey)
	delete(session.Values, ctx.OriginalUserRoleKey)

	if err = session.Save(r, w); err != nil {
		c.HandleError(w, r, err)
		return
	}

	http.Redirect(w, r, c.PathTo("home").String(), http.StatusFound)
}
