package handlers

import (
	"net/http"

	"angopod-back/internal/data"
	"angopod-back/internal/models"
)

// GetAllPresenters retorna a lista de apresentadores
func GetAllPresenters(w http.ResponseWriter, r *http.Request) {
	var list []models.Presenter
	for _, pres := range data.Presenters {
		list = append(list, pres)
	}

	writeJSON(w, http.StatusOK, list)
}

// GetPresenterBySlug retorna um apresentador pelo seu slug
func GetPresenterBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	if presenter, found := data.Presenters[slug]; found {
		writeJSON(w, http.StatusOK, presenter)
		return
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "Apresentador não encontrado"})
}

