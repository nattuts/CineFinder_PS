package handler

import (
	"encoding/json"
	"net/http"

	"cinefinder/internal/model"
	"cinefinder/internal/service"
)

type MovieHandler struct {
	service service.MovieServiceInterface
}

func NewMovieHandler(s service.MovieServiceInterface) *MovieHandler {
	return &MovieHandler{service: s}
}

func (h *MovieHandler) Create(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	createdMovie, err := h.service.Create(movie)
	if err != nil {
		http.Error(w, "Erro ao salvar filme", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdMovie)
}