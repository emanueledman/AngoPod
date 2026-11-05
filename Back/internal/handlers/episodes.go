package handlers

import (
	"net/http"

	"angopod-back/internal/data"
	"angopod-back/internal/db"
	"angopod-back/internal/models"
)

// GetAllEpisodes retorna a lista de todos os episódios recentes da SQLite
func GetAllEpisodes(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")

	episodes, err := db.GetEpisodes("")
	if err != nil || len(episodes) == 0 {
		episodes = data.Episodes
	}

	if tag != "" {
		var filtered []models.Episode
		for _, ep := range episodes {
			if ep.Tag == tag {
				filtered = append(filtered, ep)
			}
		}
		episodes = filtered
	}

	if episodes == nil {
		episodes = []models.Episode{}
	}

	writeJSON(w, http.StatusOK, episodes)
}

// GetEpisodeByID retorna um episódio pelo seu ID a partir da SQLite
func GetEpisodeByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ep, err := db.GetEpisodeByID(id)
	if err == nil && ep != nil {
		writeJSON(w, http.StatusOK, ep)
		return
	}

	for _, e := range data.Episodes {
		if e.ID == id {
			writeJSON(w, http.StatusOK, e)
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "Episódio não encontrado"})
}
