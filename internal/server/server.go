package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/tanq16/ohara/internal/model"
	"github.com/tanq16/ohara/internal/store"
)

type Storer interface {
	ListTouchpoints(category, tag, startDate string) ([]model.Touchpoint, error)
	CreateTouchpoint(input model.TouchpointInput) (model.Touchpoint, error)
	UpdateTouchpoint(id string, input model.TouchpointInput) (model.Touchpoint, error)
	DeleteTouchpoint(id string) error
	GetMetadata() (model.Metadata, error)
	AddCategory(name string) error
	RemoveCategory(name string) error
	AddTag(name string) error
	RemoveTag(name string) error
	ListReports() ([]string, error)
	GetReport(filename string) (string, error)
	CreateReport(filename, content string) error
}

type Config struct {
	Port int
}

type Server struct {
	config Config
	store  Storer
	mux    *http.ServeMux
}

func New(cfg Config, st Storer) *Server {
	s := &Server{
		config: cfg,
		store:  st,
		mux:    http.NewServeMux(),
	}
	s.setup()
	return s
}

func (s *Server) setup() {
	s.mux.HandleFunc("GET /api/touchpoints", s.listTouchpoints)
	s.mux.HandleFunc("POST /api/touchpoints", s.createTouchpoint)
	s.mux.HandleFunc("PUT /api/touchpoints/{id}", s.updateTouchpoint)
	s.mux.HandleFunc("DELETE /api/touchpoints/{id}", s.deleteTouchpoint)

	s.mux.HandleFunc("GET /api/metadata", s.getMetadata)
	s.mux.HandleFunc("POST /api/metadata/categories", s.addCategory)
	s.mux.HandleFunc("DELETE /api/metadata/categories/{name}", s.removeCategory)
	s.mux.HandleFunc("POST /api/metadata/tags", s.addTag)
	s.mux.HandleFunc("DELETE /api/metadata/tags/{name}", s.removeTag)

	s.mux.HandleFunc("GET /api/reports", s.listReports)
	s.mux.HandleFunc("GET /api/reports/{filename}", s.getReport)
	s.mux.HandleFunc("POST /api/reports", s.createReport)

	sub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("ERROR [server] failed to create static sub-filesystem: %v", err)
	}
	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(sub)))

	s.mux.HandleFunc("GET /{$}", s.serveIndex)
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {
	data, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("INFO [server] %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	log.Printf("INFO [server] Starting on %s", addr)
	return http.ListenAndServe(addr, withLogging(s.mux))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, store.ErrValidation),
		errors.Is(err, store.ErrAlreadyExists),
		errors.Is(err, store.ErrInvalidFilename):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		log.Printf("ERROR [server] %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}
