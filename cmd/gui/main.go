package main

import (
	"internal/persistent_storage"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func apiHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("a")
	db, err := persistent_storage.CreateDb()
	if err != nil {
		log.Panic(err)
	}

	project := chi.URLParam(r, "project")
	action := chi.URLParam(r, "action")
	file := chi.URLParam(r, "*")

	switch strings.ToLower(action) {
	case "file":
	case "get":
	case "content": // TODO: Decide what the canonical name should be
		getFile(w, r, project, file, db)
		return
	case "definition":
		// TODO: do
	case "text":
		// TODO: do
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

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/api/{project}/{action}/*", apiHandler)
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root"))
	})

	http.ListenAndServe(":8080", r)
}
