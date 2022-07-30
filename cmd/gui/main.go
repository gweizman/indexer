package main

import (
	"internal/persistent_storage"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

func getFile(w http.ResponseWriter, r *http.Request, project string, file string, db *persistent_storage.Db) {
	t, found, err := db.GetFileContent(project, file)
	if err != nil {
		log.Panic(err)
	}

	if !found {
		render.Render(w, r, ErrNotFound)
		return
	}
	w.Write([]byte(t))
}

type DefinitionResponse struct {
	*persistent_storage.Definition
}

type SearchResponse struct {
	*persistent_storage.FileContent
}

type DirResponse struct {
	*persistent_storage.DirOrFile
}

func (k DefinitionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (k SearchResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (k DirResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getDef(w http.ResponseWriter, r *http.Request, project string, path_limit string, db *persistent_storage.Db) {
	t, found, err := db.GetDefinition(project, path_limit, r.URL.Query().Get("name"))
	if err != nil {
		log.Panic(err)
	}

	if !found {
		render.Render(w, r, ErrNotFound)
		return
	}

	list := []render.Renderer{}
	for _, def := range t {
		list = append(list, &DefinitionResponse{Definition: def})
	}

	render.RenderList(w, r, list)
}

func getDir(w http.ResponseWriter, r *http.Request, project string, path string, db *persistent_storage.Db) {
	t, found, err := db.GetDirChildren(project, path)
	if err != nil {
		log.Panic(err)
	}

	if !found {
		render.Render(w, r, ErrNotFound)
		return
	}

	list := []render.Renderer{}
	for _, def := range t {
		list = append(list, &DirResponse{DirOrFile: def})
	}

	render.RenderList(w, r, list)
}

func textSearch(w http.ResponseWriter, r *http.Request, project string, path_limit string, db *persistent_storage.Db) {
	t, found, err := db.SearchFileContent(project, path_limit, r.URL.Query().Get("query"))
	if err != nil {
		log.Panic(err)
	}

	if !found {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	list := []render.Renderer{}
	for _, val := range t {
		list = append(list, &SearchResponse{FileContent: val})
	}

	render.RenderList(w, r, list)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	db, err := persistent_storage.CreateDb()
	if err != nil {
		log.Panic(err)
	}

	project := chi.URLParam(r, "project")
	action := chi.URLParam(r, "action")
	file := chi.URLParam(r, "*")

	switch strings.ToLower(action) {
	case "file":
		getFile(w, r, project, file, db)
		return
	case "definition":
		getDef(w, r, project, file, db)
		return
	case "search":
		textSearch(w, r, project, file, db)
		return
	case "dir":
		getDir(w, r, project, file, db)
		return
	default:
		http.Error(w, http.StatusText(404), 404)
		return
	}
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/api/{project}/{action}/*", apiHandler)
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root"))
	})

	http.ListenAndServe(":8080", r)
}
