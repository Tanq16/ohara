package server

import (
	"encoding/json"
	"net/http"

	"github.com/tanishqrupaal/ohara/internal/model"
)

func (s *Server) listTouchpoints(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	tag := r.URL.Query().Get("tag")
	startDate := r.URL.Query().Get("start_date")

	tps, err := s.store.ListTouchpoints(category, tag, startDate)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tps)
}

func (s *Server) createTouchpoint(w http.ResponseWriter, r *http.Request) {
	var input model.TouchpointInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	tp, err := s.store.CreateTouchpoint(input)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, tp)
}

func (s *Server) updateTouchpoint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var input model.TouchpointInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	tp, err := s.store.UpdateTouchpoint(id, input)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tp)
}

func (s *Server) deleteTouchpoint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := s.store.DeleteTouchpoint(id)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
