// Code generated by ogen, DO NOT EDIT.

package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ogen-go/ogen/uri"
)

func (s *Server) cutPrefix(path string) (string, bool) {
	prefix := s.cfg.Prefix
	if prefix == "" {
		return path, true
	}
	if !strings.HasPrefix(path, prefix) {
		// Prefix doesn't match.
		return "", false
	}
	// Cut prefix from the path.
	return strings.TrimPrefix(path, prefix), true
}

// ServeHTTP serves http request as defined by OpenAPI v3 specification,
// calling handler that matches the path or returning not found error.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	elem := r.URL.Path
	elemIsEscaped := false
	if rawPath := r.URL.RawPath; rawPath != "" {
		if normalized, ok := uri.NormalizeEscapedPath(rawPath); ok {
			elem = normalized
			elemIsEscaped = strings.ContainsRune(elem, '%')
		}
	}

	elem, ok := s.cutPrefix(elem)
	if !ok || len(elem) == 0 {
		s.notFound(w, r)
		return
	}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/"
			origElem := elem
			if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'a': // Prefix: "add-p"
				origElem := elem
				if l := len("add-p"); len(elem) >= l && elem[0:l] == "add-p" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'e': // Prefix: "erson"
					origElem := elem
					if l := len("erson"); len(elem) >= l && elem[0:l] == "erson" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handleAddPersonRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}

					elem = origElem
				case 'r': // Prefix: "roject"
					origElem := elem
					if l := len("roject"); len(elem) >= l && elem[0:l] == "roject" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handleAddProjectRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}

					elem = origElem
				}

				elem = origElem
			case 'g': // Prefix: "get-"
				origElem := elem
				if l := len("get-"); len(elem) >= l && elem[0:l] == "get-" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'o': // Prefix: "organization"
					origElem := elem
					if l := len("organization"); len(elem) >= l && elem[0:l] == "organization" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handleGetOrganizationRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}

					elem = origElem
				case 'p': // Prefix: "p"
					origElem := elem
					if l := len("p"); len(elem) >= l && elem[0:l] == "p" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'e': // Prefix: "erson"
						origElem := elem
						if l := len("erson"); len(elem) >= l && elem[0:l] == "erson" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "POST":
								s.handleGetPersonRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "POST")
							}

							return
						}

						elem = origElem
					case 'r': // Prefix: "roject"
						origElem := elem
						if l := len("roject"); len(elem) >= l && elem[0:l] == "roject" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "POST":
								s.handleGetProjectRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "POST")
							}

							return
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			case 'i': // Prefix: "import-"
				origElem := elem
				if l := len("import-"); len(elem) >= l && elem[0:l] == "import-" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'o': // Prefix: "organizations"
					origElem := elem
					if l := len("organizations"); len(elem) >= l && elem[0:l] == "organizations" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handleImportOrganizationsRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}

					elem = origElem
				case 'p': // Prefix: "p"
					origElem := elem
					if l := len("p"); len(elem) >= l && elem[0:l] == "p" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'e': // Prefix: "erson"
						origElem := elem
						if l := len("erson"); len(elem) >= l && elem[0:l] == "erson" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "POST":
								s.handleImportPersonRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "POST")
							}

							return
						}

						elem = origElem
					case 'r': // Prefix: "roject"
						origElem := elem
						if l := len("roject"); len(elem) >= l && elem[0:l] == "roject" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "POST":
								s.handleImportProjectRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "POST")
							}

							return
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			case 's': // Prefix: "search-"
				origElem := elem
				if l := len("search-"); len(elem) >= l && elem[0:l] == "search-" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'o': // Prefix: "organizations"
					origElem := elem
					if l := len("organizations"); len(elem) >= l && elem[0:l] == "organizations" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handleSearchOrganizationsRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}

					elem = origElem
				case 'p': // Prefix: "p"
					origElem := elem
					if l := len("p"); len(elem) >= l && elem[0:l] == "p" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'e': // Prefix: "eople"
						origElem := elem
						if l := len("eople"); len(elem) >= l && elem[0:l] == "eople" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "POST":
								s.handleSearchPeopleRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "POST")
							}

							return
						}

						elem = origElem
					case 'r': // Prefix: "rojects"
						origElem := elem
						if l := len("rojects"); len(elem) >= l && elem[0:l] == "rojects" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "POST":
								s.handleSearchProjectsRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "POST")
							}

							return
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			}

			elem = origElem
		}
	}
	s.notFound(w, r)
}

// Route is route object.
type Route struct {
	name        string
	summary     string
	operationID string
	pathPattern string
	count       int
	args        [0]string
}

// Name returns ogen operation name.
//
// It is guaranteed to be unique and not empty.
func (r Route) Name() string {
	return r.name
}

// Summary returns OpenAPI summary.
func (r Route) Summary() string {
	return r.summary
}

// OperationID returns OpenAPI operationId.
func (r Route) OperationID() string {
	return r.operationID
}

// PathPattern returns OpenAPI path.
func (r Route) PathPattern() string {
	return r.pathPattern
}

// Args returns parsed arguments.
func (r Route) Args() []string {
	return r.args[:r.count]
}

// FindRoute finds Route for given method and path.
//
// Note: this method does not unescape path or handle reserved characters in path properly. Use FindPath instead.
func (s *Server) FindRoute(method, path string) (Route, bool) {
	return s.FindPath(method, &url.URL{Path: path})
}

// FindPath finds Route for given method and URL.
func (s *Server) FindPath(method string, u *url.URL) (r Route, _ bool) {
	var (
		elem = u.Path
		args = r.args
	)
	if rawPath := u.RawPath; rawPath != "" {
		if normalized, ok := uri.NormalizeEscapedPath(rawPath); ok {
			elem = normalized
		}
		defer func() {
			for i, arg := range r.args[:r.count] {
				if unescaped, err := url.PathUnescape(arg); err == nil {
					r.args[i] = unescaped
				}
			}
		}()
	}

	elem, ok := s.cutPrefix(elem)
	if !ok {
		return r, false
	}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/"
			origElem := elem
			if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'a': // Prefix: "add-p"
				origElem := elem
				if l := len("add-p"); len(elem) >= l && elem[0:l] == "add-p" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'e': // Prefix: "erson"
					origElem := elem
					if l := len("erson"); len(elem) >= l && elem[0:l] == "erson" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "POST":
							// Leaf: AddPerson
							r.name = "AddPerson"
							r.summary = "Upsert a person"
							r.operationID = "addPerson"
							r.pathPattern = "/add-person"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}

					elem = origElem
				case 'r': // Prefix: "roject"
					origElem := elem
					if l := len("roject"); len(elem) >= l && elem[0:l] == "roject" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "POST":
							// Leaf: AddProject
							r.name = "AddProject"
							r.summary = "Upsert a project"
							r.operationID = "addProject"
							r.pathPattern = "/add-project"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}

					elem = origElem
				}

				elem = origElem
			case 'g': // Prefix: "get-"
				origElem := elem
				if l := len("get-"); len(elem) >= l && elem[0:l] == "get-" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'o': // Prefix: "organization"
					origElem := elem
					if l := len("organization"); len(elem) >= l && elem[0:l] == "organization" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "POST":
							// Leaf: GetOrganization
							r.name = "GetOrganization"
							r.summary = "Get organization"
							r.operationID = "getOrganization"
							r.pathPattern = "/get-organization"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}

					elem = origElem
				case 'p': // Prefix: "p"
					origElem := elem
					if l := len("p"); len(elem) >= l && elem[0:l] == "p" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'e': // Prefix: "erson"
						origElem := elem
						if l := len("erson"); len(elem) >= l && elem[0:l] == "erson" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "POST":
								// Leaf: GetPerson
								r.name = "GetPerson"
								r.summary = "Get person"
								r.operationID = "getPerson"
								r.pathPattern = "/get-person"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}

						elem = origElem
					case 'r': // Prefix: "roject"
						origElem := elem
						if l := len("roject"); len(elem) >= l && elem[0:l] == "roject" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "POST":
								// Leaf: GetProject
								r.name = "GetProject"
								r.summary = "Get project"
								r.operationID = "getProject"
								r.pathPattern = "/get-project"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			case 'i': // Prefix: "import-"
				origElem := elem
				if l := len("import-"); len(elem) >= l && elem[0:l] == "import-" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'o': // Prefix: "organizations"
					origElem := elem
					if l := len("organizations"); len(elem) >= l && elem[0:l] == "organizations" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "POST":
							// Leaf: ImportOrganizations
							r.name = "ImportOrganizations"
							r.summary = "Import organization hierarchy"
							r.operationID = "importOrganizations"
							r.pathPattern = "/import-organizations"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}

					elem = origElem
				case 'p': // Prefix: "p"
					origElem := elem
					if l := len("p"); len(elem) >= l && elem[0:l] == "p" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'e': // Prefix: "erson"
						origElem := elem
						if l := len("erson"); len(elem) >= l && elem[0:l] == "erson" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "POST":
								// Leaf: ImportPerson
								r.name = "ImportPerson"
								r.summary = "Import a person"
								r.operationID = "importPerson"
								r.pathPattern = "/import-person"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}

						elem = origElem
					case 'r': // Prefix: "roject"
						origElem := elem
						if l := len("roject"); len(elem) >= l && elem[0:l] == "roject" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "POST":
								// Leaf: ImportProject
								r.name = "ImportProject"
								r.summary = "Import a project"
								r.operationID = "importProject"
								r.pathPattern = "/import-project"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			case 's': // Prefix: "search-"
				origElem := elem
				if l := len("search-"); len(elem) >= l && elem[0:l] == "search-" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'o': // Prefix: "organizations"
					origElem := elem
					if l := len("organizations"); len(elem) >= l && elem[0:l] == "organizations" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "POST":
							// Leaf: SearchOrganizations
							r.name = "SearchOrganizations"
							r.summary = "Search organizations"
							r.operationID = "searchOrganizations"
							r.pathPattern = "/search-organizations"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}

					elem = origElem
				case 'p': // Prefix: "p"
					origElem := elem
					if l := len("p"); len(elem) >= l && elem[0:l] == "p" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'e': // Prefix: "eople"
						origElem := elem
						if l := len("eople"); len(elem) >= l && elem[0:l] == "eople" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "POST":
								// Leaf: SearchPeople
								r.name = "SearchPeople"
								r.summary = "Search people"
								r.operationID = "searchPeople"
								r.pathPattern = "/search-people"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}

						elem = origElem
					case 'r': // Prefix: "rojects"
						origElem := elem
						if l := len("rojects"); len(elem) >= l && elem[0:l] == "rojects" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "POST":
								// Leaf: SearchProjects
								r.name = "SearchProjects"
								r.summary = "Search projects"
								r.operationID = "searchProjects"
								r.pathPattern = "/search-projects"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			}

			elem = origElem
		}
	}
	return r, false
}
