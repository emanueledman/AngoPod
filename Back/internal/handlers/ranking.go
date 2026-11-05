package handlers

import (
	"net/http"

	"angopod-back/internal/data"
	"angopod-back/internal/db"
	"angopod-back/internal/models"
)

// GetRanking retorna os dados completos de ranking da semana a partir da SQLite
func GetRanking(w http.ResponseWriter, r *http.Request) {
	podcasts, err := db.GetPodcasts("", "", "", "", "")
	if err != nil || len(podcasts) == 0 {
		podcasts = data.Podcasts
	}

	episodes, err := db.GetEpisodes("")
	if err != nil || len(episodes) == 0 {
		episodes = data.Episodes
	}

	resp := models.RankingResponse{
		TopShows:     podcasts,
		TopEpisodes:  episodes,
		Movers:       data.Movers,
		Weeks:        []string{"24–30 JUN", "17–23 JUN", "10–16 JUN", "03–09 JUN"},
		UpdatedLabel: "Atualizado recentemente",
	}

	writeJSON(w, http.StatusOK, resp)
}

// GetHighlights retorna os destaques da semana e movers
func GetHighlights(w http.ResponseWriter, r *http.Request) {
	resp := map[string]any{
		"pulseHighlights": data.PulseHighlights,
		"movers":          data.Movers,
	}

	writeJSON(w, http.StatusOK, resp)
}
