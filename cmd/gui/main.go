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

func getFile(w http.ResponseWriter, r *http.Request, project string, file string, db *persistent_storage.Db) {
	t, found, err := db.GetFileContent(project, file)
	if err != nil {
		log.Panic(err)
	}
	log.Print(t)

	if !found {
		http.Error(w, http.StatusText(404), 404)
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

func (k *DefinitionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (k *SearchResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getDef(w http.ResponseWriter, r *http.Request, project string, path_limit string, db *persistent_storage.Db) {
	t, found, err := db.GetDefinition(project, path_limit, r.URL.Query().Get("name"))
	if err != nil {
		log.Panic(err)
	}

	if !found {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	list := []render.Renderer{}
	for _, def := range t {
		list = append(list, &DefinitionResponse{Definition: &def})
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
		list = append(list, &SearchResponse{FileContent: &val})
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
