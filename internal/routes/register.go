package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/controllers"
	"github.com/ugent-library/biblio-backend/internal/middleware"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/go-locale/locale"
	"gopkg.in/cas.v2"
)

func Register(c controllers.Context) {
	router := c.Router
	basePath := c.BaseURL.Path

	router.Use(handlers.RecoveryHandler())

	// static files
	router.PathPrefix(basePath + "/static/").Handler(http.StripPrefix(basePath+"/static/", http.FileServer(http.Dir("./static"))))

	// requireUser := middleware.RequireUser(baseURL.Path + "/logout")
	// setUser := middleware.SetUser(e, sessionName, sessionStore)
	requireUser := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !cas.IsAuthenticated(r) {
				cas.RedirectToLogin(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
	setUser := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var user *models.User

			session, _ := c.Session(r)
			userID := session.Values["user_id"]
			if userID != nil {
				u, err := c.Engine.GetUser(userID.(string))
				if err != nil {
					log.Printf("get user error: %s", err)
					// TODO
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				user = u
			}

			if user == nil {
				u, err := c.Engine.GetUserByUsername(cas.Username(r))
				if err != nil {
					log.Printf("get user error: %s", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				user = u
			}

			originalUserID := session.Values["original_user_id"]
			if originalUserID != nil {
				originalUser, err := c.Engine.GetUser(originalUserID.(string))
				if err != nil {
					log.Printf("get user error: %s", err)
					// TODO
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				c := context.WithOriginalUser(r.Context(), originalUser)
				r = r.WithContext(c)
			}

			c := context.WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(c))
		})
	}

	// authController := controllers.NewAuth(e, sessionName, sessionStore, oidcClient, router)
	usersController := controllers.NewUsers(c)
	publicationsController := controllers.NewPublications(c)
	datasetsController := controllers.NewDatasets(c)
	publicationFilesController := controllers.NewPublicationFiles(c)
	publicationDetailsController := controllers.NewPublicationDetails(c)
	publicationConferenceController := controllers.NewPublicationConference(c)
	publicationProjectsController := controllers.NewPublicationProjects(c)
	publicationDepartmentsController := controllers.NewPublicationDepartments(c)
	publicationAbstractsController := controllers.NewPublicationAbstracts(c)
	publicationLinksController := controllers.NewPublicationLinks(c)
	publicationAuthorsController := controllers.NewPublicationAuthors(c)
	publicationDatasetsController := controllers.NewPublicationDatasets(c)
	publicationAdditionalInfoController := controllers.NewPublicationAdditionalInfo(c)
	datasetDetailsController := controllers.NewDatasetDetails(c)
	datasetProjectsController := controllers.NewDatasetProjects(c)
	datasetPublicationsController := controllers.NewDatasetPublications(c)

	// TODO fix absolute url generation
	// var schemes []string
	// if u.Scheme == "http" {
	// 	schemes = []string{"http", "https"}
	// } else {
	// 	schemes = []string{"https", "http"}
	// }
	// r = r.Schemes(schemes...).Host(u.Host).PathPrefix(u.Path).Subrouter()

	r := router.PathPrefix(basePath).Subrouter()

	// r.Use(handlers.HTTPMethodOverrideHandler)
	r.Use(locale.Detect(c.Localizer))

	// home
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, basePath+"/publication", http.StatusFound)
	}).Methods("GET").Name("home")

	// auth
	// r.HandleFunc("/login", authController.Login).
	// 	Methods("GET").
	// 	Name("login")
	// r.HandleFunc("/auth/openid-connect/callback", authController.Callback).
	// 	Methods("GET")
	// r.HandleFunc("/logout", authController.Logout).
	// 	Methods("GET").
	// 	Name("logout")
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		session, _ := c.SessionStore.Get(r, c.SessionName)
		delete(session.Values, "original_user_id")
		delete(session.Values, "user_id")
		session.Save(r, w)
		cas.RedirectToLogout(w, r)
	}).Methods("GET").Name("logout")

	// users
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.Use(requireUser)
	userRouter.Use(setUser)
	userRouter.HandleFunc("/htmx/impersonate/choose", usersController.ImpersonateChoose).
		Methods("GET").
		Name("user_impersonate_choose")
	userRouter.HandleFunc("/impersonate", usersController.Impersonate).
		Methods("POST").
		Name("user_impersonate")
	// TODO why doesn't a DELETE with methodoverride work with CAS?
	userRouter.HandleFunc("/impersonate/remove", usersController.ImpersonateRemove).
		Methods("POST").
		Name("user_impersonate_remove")

	// publications
	pubsRouter := r.PathPrefix("/publication").Subrouter()
	pubsRouter.Use(middleware.SetActiveMenu("publications"))
	pubsRouter.Use(requireUser)
	pubsRouter.Use(setUser)
	pubsRouter.HandleFunc("", publicationsController.List).
		Methods("GET").
		Name("publications")
	pubsRouter.HandleFunc("/add", publicationsController.Add).
		Methods("GET").
		Name("publication_add")
	pubsRouter.HandleFunc("/add-single", publicationsController.AddSingle).
		Methods("GET").
		Name("publication_add_single")
	pubsRouter.HandleFunc("/add-multiple", publicationsController.AddMultiple).
		Methods("GET").
		Name("publication_add_multiple")
	pubsRouter.HandleFunc("/add-single/import", publicationsController.AddSingleImport).
		Methods("POST").
		Name("publication_add_single_import")

	pubRouter := pubsRouter.PathPrefix("/{id}").Subrouter()
	pubRouter.Use(middleware.SetPublication(c.Engine))
	pubRouter.HandleFunc("", publicationsController.Show).
		Methods("GET").
		Name("publication")
	pubRouter.HandleFunc("/add-single/description", publicationsController.AddSingleDescription).
		Methods("GET").
		Name("publication_add_single_description")
	pubRouter.HandleFunc("/add-single/confirm", publicationsController.AddSingleConfirm).
		Methods("GET").
		Name("publication_add_single_confirm")
	pubRouter.HandleFunc("/add-single/publish", publicationsController.AddSinglePublish).
		Methods("POST").
		Name("publication_add_single_publish")
	pubRouter.HandleFunc("/thumbnail", publicationsController.Thumbnail).
		Methods("GET").
		Name("publication_thumbnail")
	// Publication files
	pubRouter.HandleFunc("/file/{file_id}", publicationFilesController.Download).
		Methods("GET").
		Name("publication_file")
	pubRouter.HandleFunc("/file/{file_id}/thumbnail", publicationFilesController.Thumbnail).
		Methods("GET").
		Name("publication_file_thumbnail")
	pubRouter.HandleFunc("/htmx/file", publicationFilesController.Upload).
		Methods("POST").
		Name("upload_publication_file")
	pubRouter.HandleFunc("/htmx/file/{file_id}/edit", publicationFilesController.Edit).
		Methods("GET").
		Name("publication_file_edit")
	pubRouter.HandleFunc("/htmx/file/{file_id}", publicationFilesController.Update).
		Methods("PUT").
		Name("publication_file_update")
	pubRouter.HandleFunc("/htmx/file/{file_id}/remove", publicationFilesController.Remove).
		Methods("PATCH").
		Name("publication_file_remove")
	// Publication HTMX fragments
	pubRouter.HandleFunc("/htmx/summary", publicationsController.Summary).
		Methods("GET").
		Name("publication_summary")
	// Publication details HTMX fragments
	pubRouter.HandleFunc("/htmx", publicationDetailsController.Show).
		Methods("GET").
		Name("publication_details")
	pubRouter.HandleFunc("/htmx/edit", publicationDetailsController.OpenForm).
		Methods("GET").
		Name("publication_details_edit_form")
	pubRouter.HandleFunc("/htmx/edit", publicationDetailsController.SaveForm).
		Methods("PATCH").
		Name("publication_details_save_form")
	// Publication conference HTMX fragments
	pubRouter.HandleFunc("/htmx/conference", publicationConferenceController.Show).
		Methods("GET").
		Name("publication_conference")
	pubRouter.HandleFunc("/htmx/conference/edit", publicationConferenceController.OpenForm).
		Methods("GET").
		Name("publication_conference_edit_form")
	pubRouter.HandleFunc("/htmx/conference/edit", publicationConferenceController.SaveForm).
		Methods("PATCH").
		Name("publication_conference_save_form")
	// Publication additional info HTMX fragments
	pubRouter.HandleFunc("/htmx/additional_info", publicationAdditionalInfoController.Show).
		Methods("GET").
		Name("publication_additional_info")
	pubRouter.HandleFunc("/htmx/additional_info/edit", publicationAdditionalInfoController.OpenForm).
		Methods("GET").
		Name("publication_additional_info_edit_form")
	pubRouter.HandleFunc("/htmx/additional_info/edit", publicationAdditionalInfoController.SaveForm).
		Methods("PATCH").
		Name("publication_additional_info_save_form")
	// Publication projects HTMX fragments
	pubRouter.HandleFunc("/htmx/projects/list", publicationProjectsController.ListProjects).
		Methods("GET").
		Name("publication_projects")
	pubRouter.HandleFunc("/htmx/projects/list/activesearch", publicationProjectsController.ActiveSearch).
		Methods("POST").
		Name("publication_projects_activesearch")
	pubRouter.HandleFunc("/htmx/projects/add/{project_id}", publicationProjectsController.AddToPublication).
		Methods("PATCH").
		Name("publication_projects_add_to_publication")
	pubRouter.HandleFunc("/htmx/projects/remove/{project_id}", publicationProjectsController.ConfirmRemoveFromPublication).
		Methods("GET").
		Name("publication_projects_confirm_remove_from_publication")
	pubRouter.HandleFunc("/htmx/projects/remove/{project_id}", publicationProjectsController.RemoveFromPublication).
		Methods("PATCH").
		Name("publication_projects_remove_from_publication")
	// Publication departments HTMX fragments
	pubRouter.HandleFunc("/htmx/departments/list", publicationDepartmentsController.ListDepartments).
		Methods("GET").
		Name("publicationDepartments")
	pubRouter.HandleFunc("/htmx/departments/list/activesearch", publicationDepartmentsController.ActiveSearch).
		Methods("POST").
		Name("publicationDepartments_activesearch")
	pubRouter.HandleFunc("/htmx/departments/add/{department_id}", publicationDepartmentsController.AddToPublication).
		Methods("PATCH").
		Name("publicationDepartments_add_to_publication")
	pubRouter.HandleFunc("/htmx/departments/remove/{department_id}", publicationDepartmentsController.ConfirmRemoveFromPublication).
		Methods("GET").
		Name("publicationDepartments_confirm_remove_from_publication")
	pubRouter.HandleFunc("/htmx/departments/remove/{department_id}", publicationDepartmentsController.RemoveFromPublication).
		Methods("PATCH").
		Name("publicationDepartments_remove_from_publication")
	// Publication abstracts HTMX fragments
	pubRouter.HandleFunc("/htmx/abstracts/add", publicationAbstractsController.AddAbstract).
		Methods("GET").
		Name("publication_abstracts_add_abstract")
	pubRouter.HandleFunc("/htmx/abstracts/create", publicationAbstractsController.CreateAbstract).
		Methods("POST").
		Name("publication_abstracts_create_abstract")
	pubRouter.HandleFunc("/htmx/abstracts/edit/{delta}", publicationAbstractsController.EditAbstract).
		Methods("GET").
		Name("publication_abstracts_edit_abstract")
	pubRouter.HandleFunc("/htmx/abstracts/update/{delta}", publicationAbstractsController.UpdateAbstract).
		Methods("PUT").
		Name("publication_abstracts_update_abstract")
	pubRouter.HandleFunc("/htmx/abstracts/remove/{delta}", publicationAbstractsController.ConfirmRemoveFromPublication).
		Methods("GET").
		Name("publication_abstracts_confirm_remove_from_publication")
	pubRouter.HandleFunc("/htmx/abstracts/remove/{delta}", publicationAbstractsController.RemoveAbstract).
		Methods("DELETE").
		Name("publication_abstracts_remove_abstract")
	// Publication links HTMX fragments
	pubRouter.HandleFunc("/htmx/links/add", publicationLinksController.AddLink).
		Methods("GET").
		Name("publication_links_add_link")
	pubRouter.HandleFunc("/htmx/links/create", publicationLinksController.CreateLink).
		Methods("POST").
		Name("publication_links_create_link")
	pubRouter.HandleFunc("/htmx/links/edit/{delta}", publicationLinksController.EditLink).
		Methods("GET").
		Name("publication_links_edit_link")
	pubRouter.HandleFunc("/htmx/links/update/{delta}", publicationLinksController.UpdateLink).
		Methods("PUT").
		Name("publication_links_update_link")
	pubRouter.HandleFunc("/htmx/links/remove/{delta}", publicationLinksController.ConfirmRemoveFromPublication).
		Methods("GET").
		Name("publication_links_confirm_remove_from_publication")
	pubRouter.HandleFunc("/htmx/links/remove/{delta}", publicationLinksController.RemoveLink).
		Methods("DELETE").
		Name("publication_links_remove_link")
	// Publication authors HTMX fragments
	pubRouter.HandleFunc("/htmx/authors/list", publicationAuthorsController.List).
		Methods("GET").
		Name("publication_authors_list")
	pubRouter.HandleFunc("/htmx/authors/add/{delta}", publicationAuthorsController.AddRow).
		Methods("GET").
		Name("publication_authors_add_row")
	pubRouter.HandleFunc("/htmx/authors/shift/{delta}", publicationAuthorsController.ShiftRow).
		Methods("GET").
		Name("publication_authors_shift_row")
	pubRouter.HandleFunc("/htmx/authors/cancel/add/{delta}", publicationAuthorsController.CancelAddRow).
		Methods("DELETE").
		Name("publication_authors_cancel_add_row")
	pubRouter.HandleFunc("/htmx/authors/create/{delta}", publicationAuthorsController.CreateAuthor).
		Methods("POST").
		Name("publication_authors_create_author")
	pubRouter.HandleFunc("/htmx/authors/edit/{delta}", publicationAuthorsController.EditRow).
		Methods("GET").
		Name("publication_authors_edit_row")
	pubRouter.HandleFunc("/htmx/authors/cancel/edit/{delta}", publicationAuthorsController.CancelEditRow).
		Methods("DELETE").
		Name("publication_authors_cancel_edit_row")
	pubRouter.HandleFunc("/htmx/authors/update/{delta}", publicationAuthorsController.UpdateAuthor).
		Methods("POST").
		Name("publication_authors_update_author")
	pubRouter.HandleFunc("/htmx/authors/remove/{delta}", publicationAuthorsController.ConfirmRemoveFromPublication).
		Methods("GET").
		Name("publication_authors_confirm_remove_from_publication")
	pubRouter.HandleFunc("/htmx/authors/remove/{delta}", publicationAuthorsController.RemoveAuthor).
		Methods("DELETE").
		Name("publication_authors_remove_author")
	pubRouter.HandleFunc("/htmx/authors/order/{start}/{end}", publicationAuthorsController.OrderAuthors).
		Methods("PUT").
		Name("publication_authors_order_authors")
	// Publication datasets HTMX fragments
	pubRouter.HandleFunc("/htmx/datasets/choose", publicationDatasetsController.Choose).
		Methods("GET").
		Name("publication_datasets_choose")
	pubRouter.HandleFunc("/htmx/datasets/activesearch", publicationDatasetsController.ActiveSearch).
		Methods("POST").
		Name("publication_datasets_activesearch")
	pubRouter.HandleFunc("/htmx/datasets/add/{dataset_id}", publicationDatasetsController.Add).
		Methods("PATCH").
		Name("publication_datasets_add")
	pubRouter.HandleFunc("/htmx/datasets/remove/{dataset_id}", publicationDatasetsController.ConfirmRemove).
		Methods("GET").
		Name("publication_datasets_confirm_remove")
	pubRouter.HandleFunc("/htmx/datasets/remove/{dataset_id}", publicationDatasetsController.Remove).
		Methods("PATCH").
		Name("publication_datasets_remove")

	// datasets
	datasetsRouter := r.PathPrefix("/dataset").Subrouter()
	datasetsRouter.Use(middleware.SetActiveMenu("datasets"))
	datasetsRouter.Use(requireUser)
	datasetsRouter.Use(setUser)
	datasetsRouter.HandleFunc("", datasetsController.List).
		Methods("GET").
		Name("datasets")
	datasetsRouter.HandleFunc("/add", datasetsController.Add).
		Methods("GET").
		Name("dataset_add")
	datasetsRouter.HandleFunc("/add/import", datasetsController.AddImport).
		Methods("POST").
		Name("dataset_add_import")

	datasetRouter := datasetsRouter.PathPrefix("/{id}").Subrouter()
	datasetRouter.Use(middleware.SetDataset(c.Engine))
	datasetRouter.HandleFunc("", datasetsController.Show).
		Methods("GET").
		Name("dataset")
	datasetRouter.HandleFunc("/add/description", datasetsController.AddDescription).
		Methods("GET").
		Name("dataset_add_description")
	datasetRouter.HandleFunc("/add/confirm", datasetsController.AddConfirm).
		Methods("GET").
		Name("dataset_add_confirm")
	datasetRouter.HandleFunc("/add/publish", datasetsController.AddPublish).
		Methods("POST").
		Name("dataset_add_publish")
	// Dataset details HTMX fragments
	datasetRouter.HandleFunc("/htmx/details", datasetDetailsController.Show).
		Methods("GET").
		Name("dataset_details")
	datasetRouter.HandleFunc("/htmx/details/edit", datasetDetailsController.OpenForm).
		Methods("GET").
		Name("dataset_details_edit_form")
	datasetRouter.HandleFunc("/htmx/details/edit", datasetDetailsController.SaveForm).
		Methods("PATCH").
		Name("dataset_details_save_form")
	// Dataset projects HTMX fragments
	datasetRouter.HandleFunc("/htmx/projects/list", datasetProjectsController.Choose).
		Methods("GET").
		Name("dataset_projects")
	datasetRouter.HandleFunc("/htmx/projects/list/activesearch", datasetProjectsController.ActiveSearch).
		Methods("POST").
		Name("dataset_projects_activesearch")
	datasetRouter.HandleFunc("/htmx/projects/add/{project_id}", datasetProjectsController.Add).
		Methods("PATCH").
		Name("dataset_projects_add")
	datasetRouter.HandleFunc("/htmx/projects/remove/{project_id}", datasetProjectsController.ConfirmRemove).
		Methods("GET").
		Name("dataset_projects_confirm_remove")
	datasetRouter.HandleFunc("/htmx/projects/remove/{project_id}", datasetProjectsController.Remove).
		Methods("PATCH").
		Name("dataset_projects_remove")
	// Dataset publications HTMX fragments
	datasetRouter.HandleFunc("/htmx/publications/choose", datasetPublicationsController.Choose).
		Methods("GET").
		Name("dataset_publications_choose")
	datasetRouter.HandleFunc("/htmx/publications/activesearch", datasetPublicationsController.ActiveSearch).
		Methods("POST").
		Name("dataset_publications_activesearch")
	datasetRouter.HandleFunc("/htmx/publications/add/{publication_id}", datasetPublicationsController.Add).
		Methods("PATCH").
		Name("dataset_publications_add")
	datasetRouter.HandleFunc("/htmx/publications/remove/{publication_id}", datasetPublicationsController.ConfirmRemove).
		Methods("GET").
		Name("dataset_publications_confirm_remove")
	datasetRouter.HandleFunc("/htmx/publications/remove/{publication_id}", datasetPublicationsController.Remove).
		Methods("PATCH").
		Name("dataset_publications_remove")
}
