package impersonating

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/render/form"
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
		c.Log.Warn("add impersonation: already impersonating", "user", c.OriginalUser.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	views.AddImpersonation(c, addImpersonationForm(c)).Render(r.Context(), w)
}

func AddImpersonationSuggest(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	if c.OriginalUser != nil {
		c.Log.Warn("add impersonation: already impersonating", "user", c.OriginalUser.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	b := BindAddCImpersonationSuggest{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("suggest impersonation: could not bind request arguments:", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	hits, err := c.UserSearchService.SuggestUsers(b.FirstName + " " + b.LastName)
	if err != nil {
		c.Log.Errorw("suggest impersonation: could not suggest users:", "errors", err, "request", r, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
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
		c.Log.Warn("create impersonation: already impersonating", "user", c.OriginalUser.ID)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	b := BindCreateImpersonation{}
	if err := bind.Request(r, &b); err != nil {
		c.Log.Warnw("create impersonation: could not bind request arguments", "errors", err, "request", r)
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	user, err := c.UserService.GetUser(b.ID)
	if err != nil {
		c.Log.Errorf("create impersonation: unable to fetch user %s: %w", b.ID, err)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	// TODO handle user not found

	session, err := c.SessionStore.Get(r, c.SessionName)
	if err != nil {
		c.Log.Errorw("create impersonation: session could not be retrieved:", "errors", err, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	session.Values[ctx.OriginalUserIDKey] = c.User.ID
	session.Values[ctx.OriginalUserRoleKey] = c.UserRole
	session.Values[ctx.UserIDKey] = user.ID
	session.Values[ctx.UserRoleKey] = "user"

	if err = session.Save(r, w); err != nil {
		c.Log.Errorw("create impersonation: session could not be saved:", "errors", err, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	http.Redirect(w, r, c.PathTo("home").String(), http.StatusFound)
}

func DeleteImpersonation(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	if c.OriginalUser == nil {
		c.Log.Warnf("delete impersonation: %w", errors.New("no impersonation"))
		c.HandleError(w, r, httperror.BadRequest)
		return
	}

	session, err := c.SessionStore.Get(r, c.SessionName)
	if err != nil {
		c.Log.Errorw("delete impersonation: session could not be retrieved:", "errors", err, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	session.Values[ctx.UserIDKey] = session.Values[ctx.OriginalUserIDKey]
	session.Values[ctx.UserRoleKey] = session.Values[ctx.OriginalUserRoleKey]
	delete(session.Values, ctx.OriginalUserIDKey)
	delete(session.Values, ctx.OriginalUserRoleKey)

	if err = session.Save(r, w); err != nil {
		c.Log.Errorw("delete impersonation: session could not be saved:", "errors", err, "user", c.User.ID)
		c.HandleError(w, r, httperror.InternalServerError)
		return
	}

	http.Redirect(w, r, c.PathTo("home").String(), http.StatusFound)
}

func addImpersonationForm(c *ctx.Ctx) *form.Form {
	suggestURL := c.PathTo("suggest_impersonations").String()

	return form.New().
		WithTheme("cols").
		AddSection(
			&form.Text{
				Template: "contributor_name",
				Name:     "first_name",
				Label:    "First name",
				Vars: struct {
					SuggestURL string
				}{
					SuggestURL: suggestURL,
				},
			},
			&form.Text{
				Template: "contributor_name",
				Name:     "last_name",
				Label:    "Last name",
				Vars: struct {
					SuggestURL string
				}{
					SuggestURL: suggestURL,
				},
			},
		)
}
