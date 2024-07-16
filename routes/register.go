package routes

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/jpillora/ipfilter"
	"github.com/leonelquinteros/gotext"
	"github.com/nics/ich"
	"github.com/samber/lo"
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
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/httpx"
	"github.com/ugent-library/oidc"
	"github.com/unrolled/secure"
	"github.com/unrolled/secure/cspbuilder"
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
	Assets           map[string]string
	SessionStore     sessions.Store
	SessionName      string
	Timezone         *time.Location
	Loc              *gotext.Locale
	Logger           *slog.Logger
	OIDCAuth         *oidc.Auth
	UsernameClaim    string
	FrontendURL      string
	FrontendUsername string
	FrontendPassword string
	IPRanges         string
	MaxFileSize      int
	CSRFName         string
	CSRFSecret       string
}

func Register(c Config) {
	if c.Env != "local" {
		c.Router.Use(middleware.RealIP)
	}
	c.Router.Use(httplog.RequestLogger(httplog.NewLogger("biblio-backoffice-http", httplog.Options{
		JSON:             c.Env != "local",
		LogLevel:         lo.Ternary(c.Env == "local", slog.LevelDebug, slog.LevelInfo),
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		QuietDownRoutes: []string{
			"/dashboard-icon",
			"/candidate-records-icon",
		},
		QuietDownPeriod: 1 * time.Minute,
	})))
	c.Router.Use(middleware.StripSlashes)

	// static files
	c.Router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// mount health and info
	c.Router.Get("/status", health.NewHandler(health.NewChecker())) // TODO add checkers
	c.Router.Get("/info", func(w http.ResponseWriter, r *http.Request) {
		httpx.RenderJSON(w, http.StatusOK, c.Version)
	})

	// frontoffice data exchange api
	frontofficeHandler := &frontoffice.Handler{
		Log:       c.Logger,
		Repo:      c.Services.Repo,
		FileStore: c.Services.FileStore,
		IPRanges:  c.IPRanges,
		IPFilter: ipfilter.New(ipfilter.Options{
			AllowedIPs:     strings.Split(c.IPRanges, ","),
			BlockByDefault: true,
		}),
	}

	c.Router.Group(func(r *ich.Mux) {
		r.Use(httpx.BasicAuth(c.FrontendUsername, c.FrontendPassword))
		r.Get("/frontoffice/publication/{id}", frontofficeHandler.GetPublication)
		r.Get("/frontoffice/publication", frontofficeHandler.GetAllPublications)
		r.Get("/frontoffice/dataset/{id}", frontofficeHandler.GetDataset)
		r.Get("/frontoffice/dataset", frontofficeHandler.GetAllDatasets)
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
			csrf.SameSite(csrf.SameSiteStrictMode),
			csrf.FieldName("csrf-token"),
			csrf.Secure(c.Env != "local"),
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

		r.Group(func(r *ich.Mux) {
			r.Use(ctx.Set(ctx.Config{
				Services:    c.Services,
				Logger:      c.Logger,
				Router:      c.Router,
				Assets:      c.Assets,
				MaxFileSize: c.MaxFileSize,
				Timezone:    c.Timezone,
				Loc:         c.Loc,
				Env:         c.Env,
				StatusErrorHandlers: map[int]http.HandlerFunc{
					http.StatusNotFound:            handlers.NotFound,
					http.StatusInternalServerError: handlers.InternalServerError,
				},
				ErrorHandlers: map[error]http.HandlerFunc{
					models.ErrUserNotFound: handlers.UserNotFound,
					models.ErrNotFound:     handlers.NotFound,
				},
				SessionName:  c.SessionName,
				SessionStore: c.SessionStore,
				BaseURL:      c.BaseURL,
				FrontendURL:  c.FrontendURL,
				CSRFName:     "csrf-token",
			}))

			r.NotFound(handlers.NotFound)

			// authentication
			authHandler := authenticating.NewAuthHandler(c.OIDCAuth, c.UsernameClaim)
			r.Get("/auth/openid-connect/callback", authHandler.Callback)
			r.Get("/login", authHandler.Login).Name("login")
			r.Get("/logout", authHandler.Logout).Name("logout")

			// home
			r.Get("/", handlers.Home).Name("home")

			r.Group(func(r *ich.Mux) {
				r.Use(ctx.RequireUser)

				r.With(ctx.SetNav("dashboard")).With(ctx.SetBreadcrumbs("dashboard")).Get("/dashboard", handlers.DashBoard).Name("dashboard")
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

						r.With(ctx.SetBreadcrumbs("dashboard_datasets")).
							Get("/dashboard/datasets/{type}", dashboard.CuratorDatasets).Name("dashboard_datasets")
						r.With(ctx.SetBreadcrumbs("dashboard_publications")).
							Get("/dashboard/publications/{type}", dashboard.CuratorPublications).Name("dashboard_publications")
					})

					r.Post("/dashboard/refresh-apublications/{type}", dashboard.RefreshAPublications).Name("dashboard_refresh_apublications")
					r.Post("/dashboard/refresh-upublications/{type}", dashboard.RefreshUPublications).Name("dashboard_refresh_upublications")

					r.With(ctx.SetNav("candidate_records")).
						With(ctx.SetBreadcrumbs("candidate_records")).
						Get("/candidate-records", candidaterecords.CandidateRecords).Name("candidate_records")
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
					r.With(ctx.SetNav("batch")).
						With(ctx.SetBreadcrumbs("publications", "publication_batch")).
						Get("/publication/batch", publicationbatch.Show).Name("publication_batch")
					r.Post("/publication/batch", publicationbatch.Process).Name("publication_process_batch")
				})

				// delete impersonation
				// TODO why doesn't a DELETE with methodoverride work here?
				r.Post("/delete-impersonation", impersonating.DeleteImpersonation).Name("delete_impersonation")

				// publications
				r.Group(func(r *ich.Mux) {
					r.Use(ctx.SetNav("publications"))
					r.Use(ctx.SetBreadcrumbs("publications"))

					// search
					r.Get("/publication", publicationsearching.Search).Name("publications")

					// import (wizard part 1 - before save)
					r.Route("/add-publication", func(r *ich.Mux) {
						r.Use(ctx.AddBreadcrumb("publication_add"))

						r.Get("/", publicationcreating.Add).Name("publication_add")
						r.Post("/import/single", publicationcreating.AddSingleImport).Name("publication_add_single_import")
						r.Post("/import/single/confirm", publicationcreating.AddSingleImportConfirm).Name("publication_add_single_import_confirm")
						r.Post("/import/multiple", publicationcreating.AddMultipleImport).Name("publication_add_multiple_import")
						r.Get("/import/multiple/{batch_id}/confirm", publicationcreating.AddMultipleConfirm).Name("publication_add_multiple_confirm")
						r.With(ctx.SetPublication(c.Services.Repo)).
							Get("/import/multiple/{batch_id}/publication/{id}", publicationcreating.AddMultipleShow).Name("publication_add_multiple_show")
						r.Post("/import/multiple/{batch_id}/save", publicationcreating.AddMultipleSave).Name("publication_add_multiple_save_draft")
						r.Post("/import/multiple/{batch_id}/publish", publicationcreating.AddMultiplePublish).Name("publication_add_multiple_publish")
						r.Get("/import/multiple/{batch_id}/finish", publicationcreating.AddMultipleFinish).Name("publication_add_multiple_finish")
					})

					r.Route("/publication/{id}", func(r *ich.Mux) {
						r.Use(ctx.SetPublication(c.Services.Repo))
						r.Use(ctx.RequireViewPublication)
						r.Use(ctx.AddBreadcrumb("publication"))

						// view only functions
						r.Group(func(r *ich.Mux) {
							r.Use(ctx.RequireViewPublication)

							r.Get("/", publicationviewing.Show).Name("publication")
							r.With(ctx.SetSubNav("description")).Get("/description", publicationviewing.ShowDescription).Name("publication_description")
							r.With(ctx.SetSubNav("contributors")).Get("/contributors", publicationviewing.ShowContributors).Name("publication_contributors")
							r.With(ctx.SetSubNav("files")).Get("/files", publicationviewing.ShowFiles).Name("publication_files")
							r.With(ctx.SetSubNav("datasets")).Get("/datasets", publicationviewing.ShowDatasets).Name("publication_datasets")
							r.With(ctx.SetSubNav("activity")).Get("/activity", publicationviewing.ShowActivity).Name("publication_activity")
							r.Get("/files/{file_id}", publicationviewing.DownloadFile).Name("publication_download_file")
						})

						// edit only
						r.Group(func(r *ich.Mux) {
							r.Use(ctx.RequireEditPublication)

							// add (wizard part 2 - after save)
							r.Group(func(r *ich.Mux) {
								r.Use(ctx.SetBreadcrumbs("publications", "publication_add"))

								r.Get("/add/description", publicationcreating.AddSingleDescription).Name("publication_add_single_description")
								r.Get("/add/confirm", publicationcreating.AddSingleConfirm).Name("publication_add_single_confirm")
								r.Post("/add/publish", publicationcreating.AddSinglePublish).Name("publication_add_single_publish")
								r.Get("/add/finish", publicationcreating.AddSingleFinish).Name("publication_add_single_finish")
							})

							// delete
							r.Get("/confirm-delete", publicationediting.ConfirmDelete).Name("publication_confirm_delete")
							r.Delete("/", publicationediting.Delete).Name("publication_delete")

							// details
							r.Get("/details/edit", publicationediting.EditDetails).Name("publication_edit_details")
							r.Put("/details", publicationediting.UpdateDetails).Name("publication_update_details")

							// edit publication type
							r.Get("/type/confirm", publicationediting.ConfirmUpdateType).Name("publication_confirm_update_type")
							r.Put("/type", publicationediting.UpdateType).Name("publication_update_type")

							// projects
							r.Get("/projects/add", publicationediting.AddProject).Name("publication_add_project")
							r.Get("/projects/suggestions", publicationediting.SuggestProjects).Name("publication_suggest_projects")
							r.Post("/projects", publicationediting.CreateProject).Name("publication_create_project")
							// project_id is last part of url because some id's contain slashes
							r.Get("/{snapshot_id}/projects/confirm-delete/{project_id:.+}", publicationediting.ConfirmDeleteProject).Name("publication_confirm_delete_project")
							r.Delete("/projects/{project_id:.+}", publicationediting.DeleteProject).Name("publication_delete_project")

							// conference
							r.Get("/conference/edit", publicationediting.EditConference).Name("publication_edit_conference")
							r.Put("/conference", publicationediting.UpdateConference).Name("publication_update_conference")

							// abstracts
							r.Get("/abstracts/add", publicationediting.AddAbstract).Name("publication_add_abstract")
							r.Post("/abstracts", publicationediting.CreateAbstract).Name("publication_create_abstract")
							r.Get("/abstracts/{abstract_id}/edit", publicationediting.EditAbstract).Name("publication_edit_abstract")
							r.Put("/abstracts/{abstract_id}", publicationediting.UpdateAbstract).Name("publication_update_abstract")
							r.Get("/{snapshot_id}/abstracts/{abstract_id}/confirm-delete", publicationediting.ConfirmDeleteAbstract).Name("publication_confirm_delete_abstract")
							r.Delete("/abstracts/{abstract_id}", publicationediting.DeleteAbstract).Name("publication_delete_abstract")

							// links
							r.Get("/links/add", publicationediting.AddLink).Name("publication_add_link")
							r.Post("/links", publicationediting.CreateLink).Name("publication_create_link")
							r.Get("/links/{link_id}/edit", publicationediting.EditLink).Name("publication_edit_link")
							r.Put("/links/{link_id}", publicationediting.UpdateLink).Name("publication_update_link")
							r.Get("/{snapshot_id}/links/{link_id}/confirm-delete", publicationediting.ConfirmDeleteLink).Name("publication_confirm_delete_link")
							r.Delete("/links/{link_id}", publicationediting.DeleteLink).Name("publication_delete_link")

							// lay summaries
							r.Get("/lay_summaries/add", publicationediting.AddLaySummary).Name("publication_add_lay_summary")
							r.Post("/lay_summaries", publicationediting.CreateLaySummary).Name("publication_create_lay_summary")
							r.Get("/lay_summaries/{lay_summary_id}/edit", publicationediting.EditLaySummary).Name("publication_edit_lay_summary")
							r.Put("/lay_summaries/{lay_summary_id}", publicationediting.UpdateLaySummary).Name("publication_update_lay_summary")
							r.Get("/{snapshot_id}/lay_summaries/{lay_summary_id}/confirm-delete", publicationediting.ConfirmDeleteLaySummary).Name("publication_confirm_delete_lay_summary")
							r.Delete("/lay_summaries/{lay_summary_id}", publicationediting.DeleteLaySummary).Name("publication_delete_lay_summary")

							// additional info
							r.Get("/additional-info/edit", publicationediting.EditAdditionalInfo).Name("publication_edit_additional_info")
							r.Put("/additional-info", publicationediting.UpdateAdditionalInfo).Name("publication_update_additional_info")

							// files
							r.Post("/files", publicationediting.UploadFile).Name("publication_upload_file")
							r.Get("/refresh-files", publicationediting.RefreshFiles).Name("publication_refresh_files")
							r.Get("/files/{file_id}/edit", publicationediting.EditFile).Name("publication_edit_file")
							r.Get("/files/{file_id}/refresh-form", publicationediting.RefreshEditFileForm).Name("publication_edit_file_refresh_form")
							r.Put("/files/{file_id}", publicationediting.UpdateFile).Name("publication_update_file")
							r.Get("/{snapshot_id}/files/{file_id}/confirm-delete", publicationediting.ConfirmDeleteFile).Name("publication_confirm_delete_file")
							r.Delete("/files/{file_id}", publicationediting.DeleteFile).Name("publication_delete_file")

							// contributors
							r.Post("/contributors/{role}/order", publicationediting.OrderContributors).Name("publication_order_contributors")
							r.Get("/contributors/{role}/add", publicationediting.AddContributor).Name("publication_add_contributor")
							r.Get("/contributors/{role}/suggestions", publicationediting.AddContributorSuggest).Name("publication_add_contributor_suggest")
							r.Get("/contributors/{role}/confirm-create", publicationediting.ConfirmCreateContributor).Name("publication_confirm_create_contributor")
							r.Post("/contributors/{role}", publicationediting.CreateContributor).Name("publication_create_contributor")
							r.Get("/contributors/{role}/{position}/edit", publicationediting.EditContributor).Name("publication_edit_contributor")
							r.Get("/contributors/{role}/{position}/suggestions", publicationediting.EditContributorSuggest).Name("publication_edit_contributor_suggest")
							r.Get("/contributors/{role}/{position}/confirm-update", publicationediting.ConfirmUpdateContributor).Name("publication_confirm_update_contributor")
							r.Put("/contributors/{role}/{position}", publicationediting.UpdateContributor).Name("publication_update_contributor")
							r.Get("/contributors/{role}/{position}/confirm-delete", publicationediting.ConfirmDeleteContributor).Name("publication_confirm_delete_contributor")
							r.Delete("/contributors/{role}/{position}", publicationediting.DeleteContributor).Name("publication_delete_contributor")

							// departments
							r.Get("/departments/add", publicationediting.AddDepartment).Name("publication_add_department")
							r.Post("/departments", publicationediting.CreateDepartment).Name("publication_create_department")
							r.Get("/departments/suggestions", publicationediting.SuggestDepartments).Name("publication_suggest_departments")
							r.Get("/{snapshot_id}/departments/{department_id}/confirm-delete", publicationediting.ConfirmDeleteDepartment).Name("publication_confirm_delete_department")
							r.Delete("/departments/{department_id}", publicationediting.DeleteDepartment).Name("publication_delete_department")

							// datasets
							r.Get("/datasets/add", publicationediting.AddDataset).Name("publication_add_dataset")
							r.Get("/datasets/suggestions", publicationediting.SuggestDatasets).Name("publication_suggest_datasets")
							r.Post("/datasets", publicationediting.CreateDataset).Name("publication_create_dataset")
							r.Get("/{snapshot_id}/datasets/{dataset_id}/confirm-delete", publicationediting.ConfirmDeleteDataset).Name("publication_confirm_delete_dataset")
							r.Delete("/datasets/{dataset_id}", publicationediting.DeleteDataset).Name("publication_delete_dataset")

							// activity
							r.Get("/message/edit", publicationediting.EditMessage).Name("publication_edit_message")
							r.Put("/message", publicationediting.UpdateMessage).Name("publication_update_message")

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

							// activity
							r.Get("/reviewer-tags/edit", publicationediting.EditReviewerTags).Name("publication_edit_reviewer_tags")
							r.Put("/reviewer-tags", publicationediting.UpdateReviewerTags).Name("publication_update_reviewer_tags")
							r.Get("/reviewer-note/edit", publicationediting.EditReviewerNote).Name("publication_edit_reviewer_note")
							r.Put("/reviewer-note", publicationediting.UpdateReviewerNote).Name("publication_update_reviewer_note")

							// (un)lock publication
							r.Post("/lock", publicationediting.Lock).Name("publication_lock")
							r.Post("/unlock", publicationediting.Unlock).Name("publication_unlock")
						})
					})
				})

				// datasets
				r.Group(func(r *ich.Mux) {
					r.Use(ctx.SetNav("datasets"))
					r.Use(ctx.SetBreadcrumbs("datasets"))

					r.Get("/dataset", datasetsearching.Search).Name("datasets")

					// dataset wizard (part 1)
					r.Route("/add-dataset", func(r *ich.Mux) {
						r.Use(ctx.AddBreadcrumb("dataset_add"))

						r.Get("/", datasetcreating.Add).Name("dataset_add")
						r.Post("/", datasetcreating.Add).Name("dataset_add")
						r.Post("/import/confirm", datasetcreating.ConfirmImport).Name("dataset_confirm_import")
						r.Post("/import", datasetcreating.AddImport).Name("dataset_add_import")
					})

					r.Route("/dataset/{id}", func(r *ich.Mux) {
						r.Use(ctx.SetDataset(c.Services.Repo))
						r.Use(ctx.RequireViewDataset)
						r.Use(ctx.AddBreadcrumb("dataset"))

						// view only functions
						r.Get("/", datasetviewing.Show).Name("dataset")
						r.With(ctx.SetSubNav("description")).Get("/description", datasetviewing.ShowDescription).Name("dataset_description")
						r.With(ctx.SetSubNav("contributors")).Get("/contributors", datasetviewing.ShowContributors).Name("dataset_contributors")
						r.With(ctx.SetSubNav("publications")).Get("/publications", datasetviewing.ShowPublications).Name("dataset_publications")
						r.With(ctx.SetSubNav("activity")).Get("/activity", datasetviewing.ShowActivity).Name("dataset_activity")

						// edit only
						r.Group(func(r *ich.Mux) {
							r.Use(ctx.RequireEditDataset)

							// wizard (part 2)
							r.Group(func(r *ich.Mux) {
								r.Use(ctx.SetBreadcrumbs("datasets", "dataset_add"))

								r.Post("/save", datasetcreating.AddSaveDraft).Name("dataset_add_save_draft")
								r.Post("/add/publish", datasetcreating.AddPublish).Name("dataset_add_publish")
								r.Get("/add/finish", datasetcreating.AddFinish).Name("dataset_add_finish")
								r.Get("/add/confirm", datasetcreating.AddConfirm).Name("dataset_add_confirm")
								r.Get("/add/description", datasetcreating.AddDescription).Name("dataset_add_description")
							})

							// delete
							r.Get("/confirm-delete", datasetediting.ConfirmDelete).Name("dataset_confirm_delete")
							r.Delete("/", datasetediting.Delete).Name("dataset_delete")

							// projects
							r.Get("/projects/add", datasetediting.AddProject).Name("dataset_add_project")
							r.Get("/projects/suggestions", datasetediting.SuggestProjects).Name("dataset_suggest_projects")
							r.Post("/projects", datasetediting.CreateProject).Name("dataset_create_project")
							r.Get("/{snapshot_id}/projects/confirm-delete/{project_id:.+}", datasetediting.ConfirmDeleteProject).Name("dataset_confirm_delete_project")
							r.Delete("/projects/{project_id:.+}", datasetediting.DeleteProject).Name("dataset_delete_project")

							// abstracts
							r.Get("/abstracts/add", datasetediting.AddAbstract).Name("dataset_add_abstract")
							r.Post("/abstracts", datasetediting.CreateAbstract).Name("dataset_create_abstract")
							r.Get("/abstracts/{abstract_id}/edit", datasetediting.EditAbstract).Name("dataset_edit_abstract")
							r.Put("/abstracts/{abstract_id}", datasetediting.UpdateAbstract).Name("dataset_update_abstract")
							r.Get("/{snapshot_id}/abstracts/{abstract_id}/confirm-delete", datasetediting.ConfirmDeleteAbstract).Name("dataset_confirm_delete_abstract")
							r.Delete("/abstracts/{abstract_id}", datasetediting.DeleteAbstract).Name("dataset_delete_abstract")

							// links
							r.Get("/{snapshot_id}/links/{link_id}/confirm-delete", datasetediting.ConfirmDeleteLink).Name("dataset_confirm_delete_link")

							// departments
							r.Get("/departments/add", datasetediting.AddDepartment).Name("dataset_add_department")
							r.Get("/departments/suggestions", datasetediting.SuggestDepartments).Name("dataset_suggest_departments")
							r.Post("/departments", datasetediting.CreateDepartment).Name("dataset_create_department")
							r.Get("/{snapshot_id}/departments/{department_id}/confirm-delete", datasetediting.ConfirmDeleteDepartment).Name("dataset_confirm_delete_department")
							r.Delete("/departments/{department_id}", datasetediting.DeleteDepartment).Name("dataset_delete_department")

							// publications
							r.Get("/publications/add", datasetediting.AddPublication).Name("dataset_add_publication")
							r.Get("/publications/suggestions", datasetediting.SuggestPublications).Name("dataset_suggest_publications")
							r.Post("/publications", datasetediting.CreatePublication).Name("dataset_create_publication")
							r.Get("/{snapshot_id}/publications/{publication_id}/confirm-delete", datasetediting.ConfirmDeletePublication).Name("dataset_confirm_delete_publication")
							r.Delete("/publications/{publication_id}", datasetediting.DeletePublication).Name("dataset_delete_publication")

							// activity
							r.Get("/message/edit", datasetediting.EditMessage).Name("dataset_edit_message")
							r.Put("/message", datasetediting.UpdateMessage).Name("dataset_update_message")

							// publish
							r.Get("/publish/confirm", datasetediting.ConfirmPublish).Name("dataset_confirm_publish")
							r.Post("/publish", datasetediting.Publish).Name("dataset_publish")

							// withdraw
							r.Get("/withdraw/confirm", datasetediting.ConfirmWithdraw).Name("dataset_confirm_withdraw")
							r.Post("/withdraw", datasetediting.Withdraw).Name("dataset_withdraw")

							// re-publish
							r.Get("/republish/confirm", datasetediting.ConfirmRepublish).Name("dataset_confirm_republish")
							r.Post("/republish", datasetediting.Republish).Name("dataset_republish")

							// edit links
							r.Get("/links/add", datasetediting.AddLink).Name("dataset_add_link")
							r.Post("/links", datasetediting.CreateLink).Name("dataset_create_link")
							r.Get("/links/{link_id}/edit", datasetediting.EditLink).Name("dataset_edit_link")
							r.Put("/links/{link_id}", datasetediting.UpdateLink).Name("dataset_update_link")
							r.Delete("/links/{link_id}", datasetediting.DeleteLink).Name("dataset_delete_link")

							// edit details
							r.Get("/details/edit", datasetediting.EditDetails).Name("dataset_edit_details")
							r.Put("/details/edit/refresh", datasetediting.RefreshEditDetails).Name("dataset_refresh_edit_details")
							r.Put("/details", datasetediting.UpdateDetails).Name("dataset_update_details")

							// edit contributors
							r.Post("/contributors/{role}/order", datasetediting.OrderContributors).Name("dataset_order_contributors")
							r.Get("/contributors/{role}/add", datasetediting.AddContributor).Name("dataset_add_contributor")
							r.Get("/contributors/{role}/suggestions", datasetediting.AddContributorSuggest).Name("dataset_add_contributor_suggest")
							r.Get("/contributors/{role}/confirm-create", datasetediting.ConfirmCreateContributor).Name("dataset_confirm_create_contributor")
							r.Post("/contributors/{role}", datasetediting.CreateContributor).Name("dataset_create_contributor")
							r.Get("/contributors/{role}/{position}/edit", datasetediting.EditContributor).Name("dataset_edit_contributor")
							r.Get("/contributors/{role}/{position}/suggestions", datasetediting.EditContributorSuggest).Name("dataset_edit_contributor_suggest")
							r.Get("/contributors/{role}/{position}/confirm-update", datasetediting.ConfirmUpdateContributor).Name("dataset_confirm_update_contributor")
							r.Put("/contributors/{role}/{position}", datasetediting.UpdateContributor).Name("dataset_update_contributor")
							r.Get("/contributors/{role}/{position}/confirm-delete", datasetediting.ConfirmDeleteContributor).Name("dataset_confirm_delete_contributor")
							r.Delete("/contributors/{role}/{position}", datasetediting.DeleteContributor).Name("dataset_delete_contributor")
						})

						// curator actions
						r.Group(func(r *ich.Mux) {
							r.Use(ctx.RequireCurator)

							// activity
							r.Get("/reviewer-tags/edit", datasetediting.EditReviewerTags).Name("dataset_edit_reviewer_tags")
							r.Put("/reviewer-tags", datasetediting.UpdateReviewerTags).Name("dataset_update_reviewer_tags")
							r.Get("/reviewer-note/edit", datasetediting.EditReviewerNote).Name("dataset_edit_reviewer_note")
							r.Put("/reviewer-note", datasetediting.UpdateReviewerNote).Name("dataset_update_reviewer_note")

							// (un)lock dataset
							r.Post("/lock", datasetediting.Lock).Name("dataset_lock")
							r.Post("/unlock", datasetediting.Unlock).Name("dataset_unlock")
						})
					})
				})

				// media types
				r.Get("/media_type/suggestions", mediatypes.Suggest).Name("suggest_media_types")
			})
		})
	})
}
