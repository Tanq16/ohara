package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) listReports(w http.ResponseWriter, r *http.Request) {
	names, err := s.store.ListReports()
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, names)
}

func (s *Server) getReport(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")

	content, err := s.store.GetReport(filename)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Write([]byte(content))
}

type reportPayload struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

func (s *Server) createReport(w http.ResponseWriter, r *http.Request) {
	var p reportPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := s.store.CreateReport(p.Filename, p.Content); err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"filename": p.Filename})
}
