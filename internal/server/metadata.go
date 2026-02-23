package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) getMetadata(w http.ResponseWriter, r *http.Request) {
	md, err := s.store.GetMetadata()
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, md)
}

type namePayload struct {
	Name string `json:"name"`
}

func (s *Server) addCategory(w http.ResponseWriter, r *http.Request) {
	var p namePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := s.store.AddCategory(p.Name); err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"name": p.Name})
}

func (s *Server) removeCategory(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if err := s.store.RemoveCategory(name); err != nil {
		writeStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) addTag(w http.ResponseWriter, r *http.Request) {
	var p namePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := s.store.AddTag(p.Name); err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"name": p.Name})
}

func (s *Server) removeTag(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if err := s.store.RemoveTag(name); err != nil {
		writeStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
