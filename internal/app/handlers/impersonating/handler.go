package impersonating

import (
	"errors"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/form"
)

type Handler struct {
	handlers.BaseHandler
	UserSearchService backends.UserSearchService
}

type Context struct {
	handlers.BaseContext
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil {
			handlers.Unauthorized(w, r)
			return
		}

		fn(w, r, Context{
			BaseContext: ctx,
		})
	})
}

type BindAddCImpersonationSuggest struct {
	FirstName string `query:"first_name"`
	LastName  string `query:"last_name"`
}

type BindCreateImpersonation struct {
	ID string `form:"id"`
}

type YieldAddImpersonation struct {
	Context
	Form *form.Form
}

type YieldAddImpersonationSuggest struct {
	Context
	FirstName string
	LastName  string
	Hits      []models.Person
}

func (h *Handler) AddImpersonation(w http.ResponseWriter, r *http.Request, ctx Context) {
	if ctx.OriginalUser != nil {
		h.Logger.Warn("add impersonation: already impersonating", "user", ctx.OriginalUser.ID)
		handlers.BadRequest(w, r, errors.New("already impersonating"))
	}

	if !ctx.User.CanImpersonateUser() {
		h.Logger.Warn("add impersonation: user does not have permission to impersonate", "user", ctx.User.ID)
		handlers.Unauthorized(w, r)
		return
	}

	render.Layout(w, "show_modal", "impersonation/add", YieldAddImpersonation{
		Context: ctx,
		Form:    h.addImpersonationForm(),
	})
}

func (h *Handler) AddImpersonationSuggest(w http.ResponseWriter, r *http.Request, ctx Context) {
	if ctx.OriginalUser != nil {
		h.Logger.Warn("add impersonation: already impersonating", "user", ctx.OriginalUser.ID)
		handlers.BadRequest(w, r, errors.New("already impersonating"))
	}

	if !ctx.User.CanImpersonateUser() {
		h.Logger.Warn("add impersonation: user does not have permission to impersonate", "user", ctx.User.ID)
		handlers.Unauthorized(w, r)
		return
	}

	b := BindAddCImpersonationSuggest{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("suggest impersonation: could not bind request arguments:", "errors", err, "request", r, "user", ctx.User.ID)
		handlers.BadRequest(w, r, err)
		return
	}

	hits, err := h.UserSearchService.SuggestUsers(b.FirstName + " " + b.LastName)
	if err != nil {
		h.Logger.Errorw("suggest impersonation: could not suggest users:", "errors", err, "request", r, "user", ctx.User.ID)
		handlers.InternalServerError(w, r, err)
		return
	}

	// exclude the current user
	for i, hit := range hits {
		if hit.ID == ctx.User.ID {
			if i == 0 {
				hits = hits[1:]
			} else {
				hits = append(hits[:i], hits[i+1:]...)
			}
			break
		}
	}

	render.Partial(w, "impersonation/suggest", YieldAddImpersonationSuggest{
		Context:   ctx,
		FirstName: b.FirstName,
		LastName:  b.LastName,
		Hits:      hits,
	})
}

func (h *Handler) CreateImpersonation(w http.ResponseWriter, r *http.Request, ctx Context) {
	if ctx.OriginalUser != nil {
		h.Logger.Warn("create impersonation: already impersonating", "user", ctx.OriginalUser.ID)
		handlers.BadRequest(w, r, errors.New("already impersonating"))
	}

	if !ctx.User.CanImpersonateUser() {
		h.Logger.Warn("create impersonation: user does not have permission to impersonate", "user", ctx.User.ID)
		handlers.Unauthorized(w, r)
		return
	}

	b := BindCreateImpersonation{}
	if err := bind.Request(r, &b); err != nil {
		h.Logger.Warnw("create impersonation: could not bind request arguments", "errors", err, "request", r)
		handlers.BadRequest(w, r, err)
		return
	}

	user, err := h.UserService.GetUser(b.ID)
	if err != nil {
		handlers.InternalServerError(w, r, err)
		return
	}

	// TODO handle user not found

	session, err := h.SessionStore.Get(r, h.SessionName)
	if err != nil {
		h.Logger.Errorw("create impersonation: session could not be retrieved:", "errors", err, "user", ctx.User.ID)
		handlers.InternalServerError(w, r, err)
		return
	}

	session.Values[handlers.OriginalUserIDKey] = ctx.User.ID
	session.Values[handlers.OriginalUserRoleKey] = ctx.UserRole
	session.Values[handlers.UserIDKey] = user.ID
	session.Values[handlers.UserRoleKey] = "user"

	if err = session.Save(r, w); err != nil {
		h.Logger.Errorw("create impersonation: session could not be saved:", "errors", err, "user", ctx.User.ID)
		handlers.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, h.PathFor("home").String(), http.StatusFound)
}

func (h *Handler) DeleteImpersonation(w http.ResponseWriter, r *http.Request, ctx Context) {
	if ctx.OriginalUser == nil {
		handlers.BadRequest(w, r, errors.New("no impersonation"))
	}

	session, err := h.SessionStore.Get(r, h.SessionName)
	if err != nil {
		h.Logger.Errorw("delete impersonation: session could not be retrieved:", "errors", err, "user", ctx.User.ID)
		handlers.InternalServerError(w, r, err)
		return
	}

	session.Values[handlers.UserIDKey] = session.Values[handlers.OriginalUserIDKey]
	session.Values[handlers.UserRoleKey] = session.Values[handlers.OriginalUserRoleKey]
	delete(session.Values, handlers.OriginalUserIDKey)
	delete(session.Values, handlers.OriginalUserRoleKey)

	if err = session.Save(r, w); err != nil {
		h.Logger.Errorw("delete impersonation: session could not be saved:", "errors", err, "user", ctx.User.ID)
		handlers.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, h.PathFor("home").String(), http.StatusFound)
}

func (h *Handler) addImpersonationForm() *form.Form {
	suggestURL := h.PathFor("suggest_impersonations").String()

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
