package routes

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/jpillora/ipfilter"
	"github.com/leonelquinteros/gotext"
	"github.com/nics/ich"
	"github.com/swaggest/swgui/v5emb"
	"github.com/ugent-library/biblio-backoffice/backends"
	"github.com/ugent-library/biblio-backoffice/ctx"
	"github.com/ugent-library/biblio-backoffice/handlers"
	"github.com/ugent-library/biblio-backoffice/handlers/authenticating"
	"github.com/ugent-library/biblio-backoffice/handlers/candidaterecords"
	"github.com/ugent-library/biblio-backoffice/handlers/dashboard"
	"github.com/ugent-library/biblio-backoffice/handlers/datasetcreating"
	"github.com/ugent-library/biblio-backoffice/handlers/datasetediting"
	"github.com/ugent-library/biblio-backoffice/handlers/datasetexporting"
	"github.com/ugent-library/biblio-backoffice/handlers/datasetsearching"
	"github.com/ugent-library/biblio-backoffice/handlers/datasetviewing"
	"github.com/ugent-library/biblio-backoffice/handlers/frontoffice"
	"github.com/ugent-library/biblio-backoffice/handlers/impersonating"
	"github.com/ugent-library/biblio-backoffice/handlers/mediatypes"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationbatch"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationcreating"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationediting"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationexporting"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationsearching"
	"github.com/ugent-library/biblio-backoffice/handlers/publicationviewing"
	"github.com/ugent-library/httpx"
	"github.com/ugent-library/mix"
	"github.com/ugent-library/oidc"
	"github.com/ugent-library/zaphttp"
	"github.com/ugent-library/zaphttp/zapchi"
	"github.com/unrolled/secure"
	"github.com/unrolled/secure/cspbuilder"
	"go.uber.org/zap"
)

type Version struct {
	Branch string
	Commit string
	Image  string
}

type Config struct {
	Version          Version
	Env              string
	Services         *backends.Services
	BaseURL          *url.URL
	Router           *ich.Mux
	Assets           mix.Manifest
	SessionStore     sessions.Store
	SessionName      string
	Timezone         *time.Location
	Loc              *gotext.Locale
	Logger           *zap.SugaredLogger
	OIDCAuth         *oidc.Auth
	UsernameClaim    string
	FrontendURL      string
	FrontendUsername string
	FrontendPassword string
	IPRanges         string
	MaxFileSize      int
	CSRFName         string
	CSRFSecret       string
	ApiServer        http.Handler
}

func Register(c Config) {
	c.Router.Use(middleware.RequestID)
	if c.Env != "local" {
		c.Router.Use(middleware.RealIP)
	}
	c.Router.Use(zaphttp.SetLogger(c.Logger.Desugar(), zapchi.RequestID))
	c.Router.Use(middleware.RequestLogger(zapchi.LogFormatter()))
	c.Router.Use(middleware.Recoverer)
	c.Router.Use(middleware.StripSlashes)

	// static files
	c.Router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// mount health and info
	c.Router.Get("/status", health.NewHandler(health.NewChecker())) // TODO add checkers
	c.Router.Get("/info", func(w http.ResponseWriter, r *http.Request) {
		httpx.RenderJSON(w, http.StatusOK, c.Version)
	})

	// rest api (api/v2)
	c.Router.Mount("/api/v2", http.StripPrefix("/api/v2", c.ApiServer))
	c.Router.Get("/api/v2/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/v2/openapi.yaml")
	})
	c.Router.Mount("/api/v2/docs", v5emb.New(
		"Biblio Backoffice",
		"/api/v2/openapi.yaml",
		"/api/v2/docs",
	))

	// handlers
	baseHandler := handlers.BaseHandler{
		Logger:          c.Logger,
		Router:          c.Router,
		SessionStore:    c.SessionStore,
		SessionName:     c.SessionName,
		Timezone:        c.Timezone,
		Loc:             c.Loc,
		UserService:     c.Services.UserService,
		BaseURL:         c.BaseURL,
		FrontendBaseUrl: c.FrontendURL,
	}
	frontofficeHandler := &frontoffice.Handler{
		Log:           c.Logger,
		Repo:          c.Services.Repo,
		FileStore:     c.Services.FileStore,
		PeopleRepo:    c.Services.PeopleRepo,
		PeopleIndex:   c.Services.PeopleIndex,
		ProjectsIndex: c.Services.ProjectsIndex,
		IPRanges:      c.IPRanges,
		IPFilter: ipfilter.New(ipfilter.Options{
			AllowedIPs:     strings.Split(c.IPRanges, ","),
			BlockByDefault: true,
		}),
	}
	datasetViewingHandler := &datasetviewing.Handler{
		BaseHandler: baseHandler,
		Repo:        c.Services.Repo,
	}
	datasetCreatingHandler := &datasetcreating.Handler{
		BaseHandler:         baseHandler,
		Repo:                c.Services.Repo,
		DatasetSearchIndex:  c.Services.DatasetSearchIndex,
		DatasetSources:      c.Services.DatasetSources,
		OrganizationService: c.Services.OrganizationService,
	}
	datasetEditingHandler := &datasetediting.Handler{
		BaseHandler:               baseHandler,
		Repo:                      c.Services.Repo,
		ProjectService:            c.Services.ProjectService,
		ProjectSearchService:      c.Services.ProjectSearchService,
		OrganizationSearchService: c.Services.OrganizationSearchService,
		OrganizationService:       c.Services.OrganizationService,
		PersonSearchService:       c.Services.PersonSearchService,
		PersonService:             c.Services.PersonService,
		PublicationSearchIndex:    c.Services.PublicationSearchIndex,
	}
	publicationCreatingHandler := &publicationcreating.Handler{
		BaseHandler:            baseHandler,
		Repo:                   c.Services.Repo,
		PublicationSearchIndex: c.Services.PublicationSearchIndex,
		PublicationSources:     c.Services.PublicationSources,
		PublicationDecoders:    c.Services.PublicationDecoders,
		OrganizationService:    c.Services.OrganizationService,
	}
	publicationEditingHandler := &publicationediting.Handler{
		BaseHandler:               baseHandler,
		Repo:                      c.Services.Repo,
		ProjectService:            c.Services.ProjectService,
		ProjectSearchService:      c.Services.ProjectSearchService,
		OrganizationSearchService: c.Services.OrganizationSearchService,
		OrganizationService:       c.Services.OrganizationService,
		PersonSearchService:       c.Services.PersonSearchService,
		PersonService:             c.Services.PersonService,
		DatasetSearchIndex:        c.Services.DatasetSearchIndex,
		FileStore:                 c.Services.FileStore,
		MaxFileSize:               c.MaxFileSize,
	}

	// frontoffice data exchange api
	c.Router.Group(func(r *ich.Mux) {
		r.Use(httpx.BasicAuth(c.FrontendUsername, c.FrontendPassword))
		r.Get("/frontoffice/publication/{id}", frontofficeHandler.GetPublication)
		r.Get("/frontoffice/publication", frontofficeHandler.GetAllPublications)
		r.Get("/frontoffice/dataset/{id}", frontofficeHandler.GetDataset)
		r.Get("/frontoffice/dataset", frontofficeHandler.GetAllDatasets)
		r.Get("/frontoffice/organization/{id}", frontofficeHandler.GetOrganization)
		r.Get("/frontoffice/organization", frontofficeHandler.GetAllOrganizations)
		r.Get("/frontoffice/organization-trees", frontofficeHandler.GetAllOrganizationTrees)
		r.Get("/frontoffice/user/{id}", frontofficeHandler.GetUser)
		r.Get("/frontoffice/user/username/{username}", frontofficeHandler.GetUserByUsername)
		r.Get("/frontoffice/person/{id}", frontofficeHandler.GetPerson)
		r.Put("/frontoffice/person/{id}/preferred-name", frontofficeHandler.SetPersonPreferredName)
		r.Get("/frontoffice/person/list", frontofficeHandler.GetPeople)
		r.Get("/frontoffice/person", frontofficeHandler.SearchPeople)
		r.Get("/frontoffice/project/{id}", frontofficeHandler.GetProject)
		r.Get("/frontoffice/project/browse", frontofficeHandler.BrowseProjects)
	})

	// frontoffice file download
	c.Router.Get("/download/{id}/{file_id}", frontofficeHandler.DownloadFile)
	c.Router.Head("/download/{id}/{file_id}", frontofficeHandler.DownloadFile)

	c.Router.Group(func(r *ich.Mux) {
		r.Use(httpx.MethodOverride) // TODO eliminate need for method override with htmx
		r.Use(csrf.Protect(
			[]byte(c.CSRFSecret),
			csrf.CookieName(c.CSRFName),
			csrf.Path("/"),
			csrf.Secure(c.BaseURL.Scheme == "https"),
			csrf.SameSite(csrf.SameSiteStrictMode),
			csrf.FieldName("csrf-token"),
		))
		r.Use(secure.New(secure.Options{
			IsDevelopment: c.Env == "local",
			ContentSecurityPolicy: (&cspbuilder.Builder{
				Directives: map[string][]string{
					cspbuilder.DefaultSrc: {"'self'"},
					cspbuilder.ScriptSrc:  {"'self'", "$NONCE"},
					// TODO: htmx injects style
					cspbuilder.StyleSrc: {"'self'", "'unsafe-inline'"},
					cspbuilder.ImgSrc:   {"'self'", "data:"},
				},
			}).MustBuild(),
		}).Handler)

		// BEGIN NEW STYLE HANDLERS
		r.Group(func(r *ich.Mux) {
			r.Use(ctx.Set(ctx.Config{
				Services:    c.Services,
				Router:      c.Router,
				Assets:      c.Assets,
				MaxFileSize: c.MaxFileSize,
				Timezone:    c.Timezone,
				Loc:         c.Loc,
				Env:         c.Env,
				ErrorHandlers: map[int]http.HandlerFunc{
					http.StatusNotFound:            handlers.NotFound,
					http.StatusInternalServerError: handlers.InternalServerError,
				},
				SessionName:   c.SessionName,
				SessionStore:  c.SessionStore,
				BaseURL:       c.BaseURL,
				FrontendURL:   c.FrontendURL,
				CSRFName:      "csrf-token",
				OIDCAuth:      c.OIDCAuth,
				UsernameClaim: c.UsernameClaim,
			}))

			r.NotFound(handlers.NotFound)

			// authentication
			r.Get("/auth/openid-connect/callback", authenticating.Callback)
			r.Get("/login", authenticating.Login).Name("login")
			r.Get("/logout", authenticating.Logout).Name("logout")

			// home
			r.Get("/", handlers.Home).Name("home")

			r.Group(func(r *ich.Mux) {
				r.Use(ctx.RequireUser)

				r.With(ctx.SetNav("dashboard")).Get("/dashboard", handlers.DashBoard).Name("dashboard")
				r.Get("/dashboard-icon", handlers.DashBoardIcon).Name("dashboard_icon")
				// dashboard action required component
				r.Get("/action-required", handlers.ActionRequired).Name("action_required")
				// dashboard drafts to complete component
				r.Get("/drafts-to-complete", handlers.DraftsToComplete).Name("drafts_to_complete")
				// dashboard recent activity component
				r.Get("/recent-activity", handlers.RecentActivity).Name("recent_activity")

				// curator only routes
				r.Group(func(r *ich.Mux) {
					r.Use(ctx.RequireCurator)

					r.Group(func(r *ich.Mux) {
						r.Use(ctx.SetNav("dashboard"))
						r.Get("/dashboard/datasets/{type}", dashboard.CuratorDatasets).Name("dashboard_datasets")
						r.Get("/dashboard/publications/{type}", dashboard.CuratorPublications).Name("dashboard_publications")
					})
					r.Post("/dashboard/refresh-apublications/{type}", dashboard.RefreshAPublications).Name("dashboard_refresh_apublications")
					r.Post("/dashboard/refresh-upublications/{type}", dashboard.RefreshUPublications).Name("dashboard_refresh_upublications")
					r.With(ctx.SetNav("candidate_records")).Get("/candidate-records", candidaterecords.CandidateRecords).Name("candidate_records")
					r.Get("/candidate-records-icon", candidaterecords.CandidateRecordsIcon).Name("candidate_records_icon")
					r.Get("/candidate-records/{id}/preview", candidaterecords.CandidateRecordPreview).Name("candidate_records_preview")
					r.Get("/candidate-records/{id}/confirm-reject", candidaterecords.ConfirmRejectCandidateRecord).Name("confirm_reject_candidate_record")
					r.Put("/candidate-records/{id}/reject", candidaterecords.RejectCandidateRecord).Name("reject_candidate_record")
					r.Put("/candidate-records/{id}/import", candidaterecords.ImportCandidateRecord).Name("import_candidate_record")
					r.Get("/candidate-records/{id}/files/{file_id}", candidaterecords.DownloadFile).Name("candidate_record_download_file")

					// impersonate user
					r.Get("/impersonation/add", impersonating.AddImpersonation).Name("add_impersonation")
					r.Get("/impersonation/suggestions", impersonating.AddImpersonationSuggest).Name("suggest_impersonations")
					r.Post("/impersonation", impersonating.CreateImpersonation).Name("create_impersonation")

					// export datasets
					r.Get("/dataset.{format}", datasetexporting.ExportByCurationSearch).Name("export_datasets")

					// change user role
					r.Put("/role/{role}", authenticating.UpdateRole).Name("update_role")

					// export publications
					r.Get("/publication.{format}", publicationexporting.ExportByCurationSearch).Name("export_publications")

					// publication batch operations
					r.With(ctx.SetNav("batch")).Get("/publication/batch", publicationbatch.Show).Name("publication_batch")
					r.Post("/publication/batch", publicationbatch.Process).Name("publication_process_batch")
				})

				// delete impersonation
				// TODO why doesn't a DELETE with methodoverride work here?
				r.Post("/delete-impersonation", impersonating.DeleteImpersonation).Name("delete_impersonation")

				// publications
				r.With(ctx.SetNav("publications")).Get("/publication", publicationsearching.Search).Name("publications")

				r.Route("/publication/{id}", func(r *ich.Mux) {
					r.Use(ctx.SetPublication(c.Services.Repo))
					r.Use(ctx.RequireViewPublication)

					// view only functions
					r.Get("/", publicationviewing.Show).Name("publication")
					r.With(ctx.SetSubNav("description")).Get("/description", publicationviewing.ShowDescription).Name("publication_description")
					r.With(ctx.SetSubNav("contributors")).Get("/contributors", publicationviewing.ShowContributors).Name("publication_contributors")
					r.With(ctx.SetSubNav("files")).Get("/files", publicationviewing.ShowFiles).Name("publication_files")
					r.With(ctx.SetSubNav("datasets")).Get("/datasets", publicationviewing.ShowDatasets).Name("publication_datasets")
					r.With(ctx.SetSubNav("activity")).Get("/activity", publicationviewing.ShowActivity).Name("publication_activity")
					r.Get("/files/{file_id}", publicationviewing.DownloadFile).Name("publication_download_file")

					// edit only
					r.Group(func(r *ich.Mux) {
						r.Use(ctx.RequireEditPublication)

						// delete
						r.Get("/confirm-delete", publicationediting.ConfirmDelete).Name("publication_confirm_delete")
						r.Delete("/", publicationediting.Delete).Name("publication_delete")

						// edit publication type
						r.Get("/type/confirm", publicationediting.ConfirmUpdateType).Name("publication_confirm_update_type")
						r.Put("/type", publicationEditingHandler.Wrap(publicationediting.UpdateType)).Name("publication_update_type")

						// details
						r.Get("/details/edit", publicationediting.EditDetails).Name("publication_edit_details")
						r.Put("/details", publicationediting.UpdateDetails).Name("publication_update_details")

						// projects
						r.Get("/projects/add", publicationediting.AddProject).Name("publication_add_project")
						r.Get("/projects/suggestions", publicationediting.SuggestProjects).Name("publication_suggest_projects")
						// project_id is last part of url because some id's contain slashes
						r.Get("/{snapshot_id}/projects/confirm-delete/{project_id:.+}", publicationediting.ConfirmDeleteProject).Name("publication_confirm_delete_project")

						// abstracts
						r.Get("/{snapshot_id}/abstracts/{abstract_id}/confirm-delete", publicationediting.ConfirmDeleteAbstract).Name("publication_confirm_delete_abstract")

						// links
						r.Get("/{snapshot_id}/links/{link_id}/confirm-delete", publicationediting.ConfirmDeleteLink).Name("publication_confirm_delete_link")

						// lay summaries
						r.Get("/{snapshot_id}/lay_summaries/{lay_summary_id}/confirm-delete", publicationediting.ConfirmDeleteLaySummary).Name("publication_confirm_delete_lay_summary")

						// files
						r.Get("/{snapshot_id}/files/{file_id}/confirm-delete", publicationediting.ConfirmDeleteFile).Name("publication_confirm_delete_file")

						// contributors
						r.Get("/contributors/{role}/{position}/confirm-delete", publicationediting.ConfirmDeleteContributor).Name("publication_confirm_delete_contributor")

						// departments
						r.Get("/departments/add", publicationediting.AddDepartment).Name("publication_add_department")
						r.Get("/departments/suggestions", publicationediting.SuggestDepartments).Name("publication_suggest_departments")
						r.Get("/{snapshot_id}/departments/{department_id}/confirm-delete", publicationediting.ConfirmDeleteDepartment).Name("publication_confirm_delete_department")

						// datasets
						r.Get("/{snapshot_id}/datasets/{dataset_id}/confirm-delete", publicationediting.ConfirmDeleteDataset).Name("publication_confirm_delete_dataset")

						// publish
						r.Get("/publish/confirm", publicationediting.ConfirmPublish).Name("publication_confirm_publish")
						r.Post("/publish", publicationediting.Publish).Name("publication_publish")

						// withdraw
						r.Get("/withdraw/confirm", publicationediting.ConfirmWithdraw).Name("publication_confirm_withdraw")
						r.Post("/withdraw", publicationediting.Withdraw).Name("publication_withdraw")

						// re-publish
						r.Get("/republish/confirm", publicationediting.ConfirmRepublish).Name("publication_confirm_republish")
						r.Post("/republish", publicationediting.Republish).Name("publication_republish")
					})

					// curator actions
					r.Group(func(r *ich.Mux) {
						r.Use(ctx.RequireCurator)

						// (un)lock publication
						r.Post("/lock", publicationediting.Lock).Name("publication_lock")
						r.Post("/unlock", publicationediting.Unlock).Name("publication_unlock")
					})
				})

				// datasets
				r.With(ctx.SetNav("datasets")).Get("/dataset", datasetsearching.Search).Name("datasets")

				r.Route("/dataset/{id}", func(r *ich.Mux) {
					r.Use(ctx.SetDataset(c.Services.Repo))
					r.Use(ctx.RequireViewDataset)

					// view only functions

					// edit only
					r.Group(func(r *ich.Mux) {
						r.Use(ctx.RequireEditDataset)

						// delete
						r.Get("/confirm-delete", datasetediting.ConfirmDelete).Name("dataset_confirm_delete")
						r.Delete("/", datasetediting.Delete).Name("dataset_delete")

						// projects
						r.Get("/projects/add", datasetediting.AddProject).Name("dataset_add_project")
						r.Get("/projects/suggestions", datasetediting.SuggestProjects).Name("dataset_suggest_projects")
						r.Get("/{snapshot_id}/projects/confirm-delete/{project_id:.+}", datasetediting.ConfirmDeleteProject).Name("dataset_confirm_delete_project")

						// abstracts
						r.Get("/abstracts/{abstract_id}/edit", datasetEditingHandler.Wrap(datasetediting.EditAbstract)).Name("dataset_edit_abstract")
						r.Get("/{snapshot_id}/abstracts/{abstract_id}/confirm-delete", datasetediting.ConfirmDeleteAbstract).Name("dataset_confirm_delete_abstract")

						// links
						r.Get("/{snapshot_id}/links/{link_id}/confirm-delete", datasetediting.ConfirmDeleteLink).Name("dataset_confirm_delete_link")

						// contributors
						r.Get("/contributors/{role}/{position}/confirm-delete", datasetediting.ConfirmDeleteContributor).Name("dataset_confirm_delete_contributor")

						// departments
						r.Get("/departments/add", datasetediting.AddDepartment).Name("dataset_add_department")
						r.Get("/departments/suggestions", datasetediting.SuggestDepartments).Name("dataset_suggest_departments")
						r.Get("/{snapshot_id}/departments/{department_id}/confirm-delete", datasetediting.ConfirmDeleteDepartment).Name("dataset_confirm_delete_department")

						// publications
						r.Get("/{snapshot_id}/publications/{publication_id}/confirm-delete", datasetediting.ConfirmDeletePublication).Name("dataset_confirm_delete_publication")

						// publish
						r.Get("/publish/confirm", datasetediting.ConfirmPublish).Name("dataset_confirm_publish")
						r.Post("/publish", datasetediting.Publish).Name("dataset_publish")

						// withdraw
						r.Get("/withdraw/confirm", datasetediting.ConfirmWithdraw).Name("dataset_confirm_withdraw")
						r.Post("/withdraw", datasetediting.Withdraw).Name("dataset_withdraw")

						// re-publish
						r.Get("/republish/confirm", datasetediting.ConfirmRepublish).Name("dataset_confirm_republish")
						r.Post("/republish", datasetediting.Republish).Name("dataset_republish")

						// publications
						r.Get("/publications/add", datasetediting.AddPublication).Name("dataset_add_publication")
						r.Get("/publications/suggestions", datasetediting.SuggestPublications).Name("dataset_suggest_publications")

						r.Post("/publications", datasetediting.CreatePublication).Name("dataset_create_publication")
						r.Delete("/publications/{publication_id}", datasetediting.DeletePublication).Name("dataset_delete_publication")

						r.With(ctx.SetNav("dataset"), ctx.SetSubNav("publications")).Get("/publications", datasetviewing.ShowPublications).Name("dataset_publications")
					})

					// curator actions
					r.Group(func(r *ich.Mux) {
						r.Use(ctx.RequireCurator)

						// (un)lock dataset
						r.Post("/lock", datasetediting.Lock).Name("dataset_lock")
						r.Post("/unlock", datasetediting.Unlock).Name("dataset_unlock")
					})
				})

				// media types
				r.Get("/media_type/suggestions", mediatypes.Suggest).Name("suggest_media_types")
			})
		})
		// END NEW STYLE HANDLERS

		// add dataset
		r.Get("/dataset/add",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.Add)).
			Name("dataset_add")
		r.Post("/dataset/add",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.Add))
		r.Post("/dataset/import",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.AddImport)).
			Name("dataset_add_import")
		r.Post("/dataset/import/confirm",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.ConfirmImport)).
			Name("dataset_confirm_import")
		r.Get("/dataset/{id}/add/description",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.AddDescription)).
			Name("dataset_add_description")
		r.Get("/dataset/{id}/add/confirm",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.AddConfirm)).
			Name("dataset_add_confirm")
		r.Post("/dataset/{id}/save",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.AddSaveDraft)).
			Name("dataset_add_save_draft")
		r.Post("/dataset/{id}/add/publish",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.AddPublish)).
			Name("dataset_add_publish")
		r.Get("/dataset/{id}/add/finish",
			datasetCreatingHandler.Wrap(datasetCreatingHandler.AddFinish)).
			Name("dataset_add_finish")

		// view dataset
		r.Get("/dataset/{id}",
			datasetViewingHandler.Wrap(datasetViewingHandler.Show)).
			Name("dataset")
		r.Get("/dataset/{id}/description",
			datasetViewingHandler.Wrap(datasetViewingHandler.ShowDescription)).
			Name("dataset_description")
		r.Get("/dataset/{id}/contributors",
			datasetViewingHandler.Wrap(datasetViewingHandler.ShowContributors)).
			Name("dataset_contributors")

		r.Get("/dataset/{id}/activity",
			datasetViewingHandler.Wrap(datasetViewingHandler.ShowActivity)).
			Name("dataset_activity")

		// edit dataset activity
		r.Get("/dataset/{id}/message/edit",
			datasetEditingHandler.Wrap(datasetEditingHandler.EditMessage)).
			Name("dataset_edit_message")
		r.Put("/dataset/{id}/message",
			datasetEditingHandler.Wrap(datasetEditingHandler.UpdateMessage)).
			Name("dataset_update_message")
		r.Get("/dataset/{id}/reviewer-tags/edit",
			datasetEditingHandler.Wrap(datasetEditingHandler.EditReviewerTags)).
			Name("dataset_edit_reviewer_tags")
		r.Put("/dataset/{id}/reviewer-tags",
			datasetEditingHandler.Wrap(datasetEditingHandler.UpdateReviewerTags)).
			Name("dataset_update_reviewer_tags")
		r.Get("/dataset/{id}/reviewer-note/edit",
			datasetEditingHandler.Wrap(datasetEditingHandler.EditReviewerNote)).
			Name("dataset_edit_reviewer_note")
		r.Put("/dataset/{id}/reviewer-note",
			datasetEditingHandler.Wrap(datasetEditingHandler.UpdateReviewerNote)).
			Name("dataset_update_reviewer_note")

		// edit dataset details
		r.Get("/dataset/{id}/details/edit",
			datasetEditingHandler.Wrap(datasetEditingHandler.EditDetails)).
			Name("dataset_edit_details")
		r.Put("/dataset/{id}/details/edit/refresh-form",
			datasetEditingHandler.Wrap(datasetEditingHandler.RefreshEditFileForm)).
			Name("dataset_edit_file_refresh_form")
		r.Put("/dataset/{id}/details",
			datasetEditingHandler.Wrap(datasetEditingHandler.UpdateDetails)).
			Name("dataset_update_details")

		// edit dataset projects
		r.Post("/dataset/{id}/projects",
			datasetEditingHandler.Wrap(datasetEditingHandler.CreateProject)).
			Name("dataset_create_project")
		r.Delete("/dataset/{id}/projects/{project_id:.+}",
			datasetEditingHandler.Wrap(datasetEditingHandler.DeleteProject)).
			Name("dataset_delete_project")

		// edit dataset links
		r.Get("/dataset/{id}/links/add",
			datasetEditingHandler.Wrap(datasetEditingHandler.AddLink)).
			Name("dataset_add_link")
		r.Post("/dataset/{id}/links",
			datasetEditingHandler.Wrap(datasetEditingHandler.CreateLink)).
			Name("dataset_create_link")
		r.Get("/dataset/{id}/links/{link_id}/edit",
			datasetEditingHandler.Wrap(datasetEditingHandler.EditLink)).
			Name("dataset_edit_link")
		r.Put("/dataset/{id}/links/{link_id}",
			datasetEditingHandler.Wrap(datasetEditingHandler.UpdateLink)).
			Name("dataset_update_link")
		r.Delete("/dataset/{id}/links/{link_id}",
			datasetEditingHandler.Wrap(datasetEditingHandler.DeleteLink)).
			Name("dataset_delete_link")

		// edit dataset departments
		r.Post("/dataset/{id}/departments",
			datasetEditingHandler.Wrap(datasetEditingHandler.CreateDepartment)).
			Name("dataset_create_department")
		r.Delete("/dataset/{id}/departments/{department_id}",
			datasetEditingHandler.Wrap(datasetEditingHandler.DeleteDepartment)).
			Name("dataset_delete_department")

		// edit dataset abstracts
		r.Get("/dataset/{id}/abstracts/add",
			datasetEditingHandler.Wrap(datasetEditingHandler.AddAbstract)).
			Name("dataset_add_abstract")
		r.Post("/dataset/{id}/abstracts",
			datasetEditingHandler.Wrap(datasetEditingHandler.CreateAbstract)).
			Name("dataset_create_abstract")
		r.Put("/dataset/{id}/abstracts/{abstract_id}",
			datasetEditingHandler.Wrap(datasetEditingHandler.UpdateAbstract)).
			Name("dataset_update_abstract")
		r.Delete("/dataset/{id}/abstracts/{abstract_id}",
			datasetEditingHandler.Wrap(datasetEditingHandler.DeleteAbstract)).
			Name("dataset_delete_abstract")

		// edit dataset contributors
		r.Post("/dataset/{id}/contributors/{role}/order",
			datasetEditingHandler.Wrap(datasetEditingHandler.OrderContributors)).
			Name("dataset_order_contributors")
		r.Get("/dataset/{id}/contributors/{role}/add",
			datasetEditingHandler.Wrap(datasetEditingHandler.AddContributor)).
			Name("dataset_add_contributor")
		r.Get("/dataset/{id}/contributors/{role}/suggestions",
			datasetEditingHandler.Wrap(datasetEditingHandler.AddContributorSuggest)).
			Name("dataset_add_contributor_suggest")
		r.Get("/dataset/{id}/contributors/{role}/confirm-create",
			datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmCreateContributor)).
			Name("dataset_confirm_create_contributor")
		r.Post("/dataset/{id}/contributors/{role}",
			datasetEditingHandler.Wrap(datasetEditingHandler.CreateContributor)).
			Name("dataset_create_contributor")
		r.Get("/dataset/{id}/contributors/{role}/{position}/edit",
			datasetEditingHandler.Wrap(datasetEditingHandler.EditContributor)).
			Name("dataset_edit_contributor")
		r.Get("/dataset/{id}/contributors/{role}/{position}/suggestions",
			datasetEditingHandler.Wrap(datasetEditingHandler.EditContributorSuggest)).
			Name("dataset_edit_contributor_suggest")
		r.Get("/dataset/{id}/contributors/{role}/{position}/confirm-update",
			datasetEditingHandler.Wrap(datasetEditingHandler.ConfirmUpdateContributor)).
			Name("dataset_confirm_update_contributor")
		r.Put("/dataset/{id}/contributors/{role}/{position}",
			datasetEditingHandler.Wrap(datasetEditingHandler.UpdateContributor)).
			Name("dataset_update_contributor")
		r.Delete("/dataset/{id}/contributors/{role}/{position}",
			datasetEditingHandler.Wrap(datasetEditingHandler.DeleteContributor)).
			Name("dataset_delete_contributor")

		// add publication
		r.Get("/publication/add",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.Add)).
			Name("publication_add")
		r.Post("/publication/add-single/import",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSingleImport)).
			Name("publication_add_single_import")
		r.Post("/publication/add-single/import/confirm",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSingleImportConfirm)).
			Name("publication_add_single_import_confirm")
		r.Get("/publication/{id}/add/description",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSingleDescription)).
			Name("publication_add_single_description")
		r.Get("/publication/{id}/add/confirm",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSingleConfirm)).
			Name("publication_add_single_confirm")
		r.Post("/publication/{id}/add/publish",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSinglePublish)).
			Name("publication_add_single_publish")
		r.Get("/publication/{id}/add/finish",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddSingleFinish)).
			Name("publication_add_single_finish")
		r.Post("/publication/add-multiple/import",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleImport)).
			Name("publication_add_multiple_import")
		r.Get("/publication/add-multiple/{batch_id}/confirm",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleConfirm)).
			Name("publication_add_multiple_confirm")
		r.Get("/publication/add-multiple/{batch_id}/publication/{id}",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleShow)).
			Name("publication_add_multiple_show")
		r.Post("/publication/add-multiple/{batch_id}/save",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleSave)).
			Name("publication_add_multiple_save_draft")
		r.Post("/publication/add-multiple/{batch_id}/publish",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultiplePublish)).
			Name("publication_add_multiple_publish")
		r.Get("/publication/add-multiple/{batch_id}/finish",
			publicationCreatingHandler.Wrap(publicationCreatingHandler.AddMultipleFinish)).
			Name("publication_add_multiple_finish")

		// edit publication activity
		r.Get("/publication/{id}/message/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditMessage)).
			Name("publication_edit_message")
		r.Put("/publication/{id}/message",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateMessage)).
			Name("publication_update_message")
		r.Get("/publication/{id}/reviewer-tags/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditReviewerTags)).
			Name("publication_edit_reviewer_tags")
		r.Put("/publication/{id}/reviewer-tags",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateReviewerTags)).
			Name("publication_update_reviewer_tags")
		r.Get("/publication/{id}/reviewer-note/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditReviewerNote)).
			Name("publication_edit_reviewer_note")
		r.Put("/publication/{id}/reviewer-note",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateReviewerNote)).
			Name("publication_update_reviewer_note")

		// edit publication conference
		r.Get("/publication/{id}/conference/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditConference)).
			Name("publication_edit_conference")
		r.Put("/publication/{id}/conference",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateConference)).
			Name("publication_update_conference")

		// edit publication additional info
		r.Get("/publication/{id}/additional-info/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditAdditionalInfo)).
			Name("publication_edit_additional_info")
		r.Put("/publication/{id}/additional-info",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateAdditionalInfo)).
			Name("publication_update_additional_info")

		// edit publication projects
		r.Post("/publication/{id}/projects",
			publicationEditingHandler.Wrap(publicationEditingHandler.CreateProject)).
			Name("publication_create_project")
		// project_id is last part of url because some id's contain slashes
		r.Delete("/publication/{id}/projects/{project_id:.+}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteProject)).
			Name("publication_delete_project")

		// edit publication links
		r.Get("/publication/{id}/links/add",
			publicationEditingHandler.Wrap(publicationEditingHandler.AddLink)).
			Name("publication_add_link")
		r.Post("/publication/{id}/links",
			publicationEditingHandler.Wrap(publicationEditingHandler.CreateLink)).
			Name("publication_create_link")
		r.Get("/publication/{id}/links/{link_id}/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditLink)).
			Name("publication_edit_link")
		r.Put("/publication/{id}/links/{link_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateLink)).
			Name("publication_update_link")
		r.Delete("/publication/{id}/links/{link_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteLink)).
			Name("publication_delete_link")

			// edit publication departments
		r.Post("/publication/{id}/departments",
			publicationEditingHandler.Wrap(publicationEditingHandler.CreateDepartment)).
			Name("publication_create_department")
		r.Delete("/publication/{id}/departments/{department_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteDepartment)).
			Name("publication_delete_department")

		// edit publication abstracts
		r.Get("/publication/{id}/abstracts/add",
			publicationEditingHandler.Wrap(publicationEditingHandler.AddAbstract)).
			Name("publication_add_abstract")
		r.Post("/publication/{id}/abstracts",
			publicationEditingHandler.Wrap(publicationEditingHandler.CreateAbstract)).
			Name("publication_create_abstract")
		r.Get("/publication/{id}/abstracts/{abstract_id}/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditAbstract)).
			Name("publication_edit_abstract")
		r.Put("/publication/{id}/abstracts/{abstract_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateAbstract)).
			Name("publication_update_abstract")
		r.Delete("/publication/{id}/abstracts/{abstract_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteAbstract)).
			Name("publication_delete_abstract")

		// edit publication lay summaries
		r.Get("/publication/{id}/lay_summaries/add",
			publicationEditingHandler.Wrap(publicationEditingHandler.AddLaySummary)).
			Name("publication_add_lay_summary")
		r.Post("/publication/{id}/lay_summaries",
			publicationEditingHandler.Wrap(publicationEditingHandler.CreateLaySummary)).
			Name("publication_create_lay_summary")
		r.Get("/publication/{id}/lay_summaries/{lay_summary_id}/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditLaySummary)).
			Name("publication_edit_lay_summary")
		r.Put("/publication/{id}/lay_summaries/{lay_summary_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateLaySummary)).
			Name("publication_update_lay_summary")
		r.Delete("/publication/{id}/lay_summaries/{lay_summary_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteLaySummary)).
			Name("publication_delete_lay_summary")

		// edit publication datasets
		r.Get("/publication/{id}/datasets/add",
			publicationEditingHandler.Wrap(publicationEditingHandler.AddDataset)).
			Name("publication_add_dataset")
		r.Get("/publication/{id}/datasets/suggestions",
			publicationEditingHandler.Wrap(publicationEditingHandler.SuggestDatasets)).
			Name("publication_suggest_datasets")
		r.Post("/publication/{id}/datasets",
			publicationEditingHandler.Wrap(publicationEditingHandler.CreateDataset)).
			Name("publication_create_dataset")
		r.Delete("/publication/{id}/datasets/{dataset_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteDataset)).
			Name("publication_delete_dataset")

		// edit publication contributors
		r.Post("/publication/{id}/contributors/{role}/order",
			publicationEditingHandler.Wrap(publicationEditingHandler.OrderContributors)).
			Name("publication_order_contributors")
		r.Get("/publication/{id}/contributors/{role}/add",
			publicationEditingHandler.Wrap(publicationEditingHandler.AddContributor)).
			Name("publication_add_contributor")
		r.Get("/publication/{id}/contributors/{role}/suggestions",
			publicationEditingHandler.Wrap(publicationEditingHandler.AddContributorSuggest)).
			Name("publication_add_contributor_suggest")
		r.Get("/publication/{id}/contributors/{role}/confirm-create",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmCreateContributor)).
			Name("publication_confirm_create_contributor")
		r.Post("/publication/{id}/contributors/{role}",
			publicationEditingHandler.Wrap(publicationEditingHandler.CreateContributor)).
			Name("publication_create_contributor")
		r.Get("/publication/{id}/contributors/{role}/{position}/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditContributor)).
			Name("publication_edit_contributor")
		r.Get("/publication/{id}/contributors/{role}/{position}/suggestions",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditContributorSuggest)).
			Name("publication_edit_contributor_suggest")
		r.Get("/publication/{id}/contributors/{role}/{position}/confirm-update",
			publicationEditingHandler.Wrap(publicationEditingHandler.ConfirmUpdateContributor)).
			Name("publication_confirm_update_contributor")
		r.Put("/publication/{id}/contributors/{role}/{position}",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateContributor)).
			Name("publication_update_contributor")
		r.Delete("/publication/{id}/contributors/{role}/{position}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteContributor)).
			Name("publication_delete_contributor")

		// edit publication files
		r.Post("/publication/{id}/files",
			publicationEditingHandler.Wrap(publicationEditingHandler.UploadFile)).
			Name("publication_upload_file")
		r.Get("/publication/{id}/files/{file_id}/edit",
			publicationEditingHandler.Wrap(publicationEditingHandler.EditFile)).
			Name("publication_edit_file")
		r.Get("/publication/{id}/refresh-files",
			publicationEditingHandler.Wrap(publicationEditingHandler.RefreshFiles)).
			Name("publication_refresh_files")
		r.Get("/publication/{id}/files/{file_id}/refresh-form",
			publicationEditingHandler.Wrap(publicationEditingHandler.RefreshEditFileForm)).
			Name("publication_edit_file_refresh_form")
		r.Put("/publication/{id}/files/{file_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.UpdateFile)).
			Name("publication_update_file")
		r.Delete("/publication/{id}/files/{file_id}",
			publicationEditingHandler.Wrap(publicationEditingHandler.DeleteFile)).
			Name("publication_delete_file")
	})
}
