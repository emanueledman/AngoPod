package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"angopod-back/internal/data"
	"angopod-back/internal/db"
	"angopod-back/internal/models"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// GetAllPodcasts retorna a lista de podcasts da SQLite com filtros via query params
func GetAllPodcasts(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	province := r.URL.Query().Get("province")
	language := r.URL.Query().Get("language")
	platform := r.URL.Query().Get("platform")
	search := strings.TrimSpace(r.URL.Query().Get("search"))

	podcasts, err := db.GetPodcasts(category, province, language, platform, search)
	if err != nil {
		// Fallback para memória em caso de exceção de leitura
		podcasts = fallbackFilterPodcasts(category, province, language, platform, search)
	}

	if podcasts == nil {
		podcasts = []models.Podcast{}
	}

	writeJSON(w, http.StatusOK, podcasts)
}

// GetPodcastBySlug retorna um podcast específico pelo slug a partir da SQLite
func GetPodcastBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	pod, err := db.GetPodcastBySlug(slug)
	if err == nil && pod != nil {
		writeJSON(w, http.StatusOK, pod)
		return
	}

	// Fallback para dados em memória se não encontrado na base
	for _, p := range data.Podcasts {
		if p.Slug == slug {
			writeJSON(w, http.StatusOK, p)
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "Podcast não encontrado"})
}

// GetPodcastEpisodes retorna os episódios de um podcast a partir da SQLite
func GetPodcastEpisodes(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	episodes, err := db.GetEpisodes(slug)
	if err != nil || len(episodes) == 0 {
		var fallbackEps []models.Episode
		for _, ep := range data.Episodes {
			if ep.PodcastSlug == slug {
				fallbackEps = append(fallbackEps, ep)
			}
		}
		episodes = fallbackEps
	}

	if episodes == nil {
		episodes = []models.Episode{}
	}

	writeJSON(w, http.StatusOK, episodes)
}

func fallbackFilterPodcasts(category, province, language, platform, search string) []models.Podcast {
	var results []models.Podcast
	searchLower := strings.ToLower(search)

	for _, pod := range data.Podcasts {
		if category != "" && category != "Todas" && category != "Todos" && pod.Category != category {
			continue
		}
		if province != "" && province != "Todas" && pod.Province != province {
			continue
		}
		if language != "" && language != "Todos" && pod.Language != language {
			continue
		}
		if platform != "" && platform != "Todas" && pod.Platform != platform {
			continue
		}
		if search != "" {
			fullText := strings.ToLower(pod.Title + " " + pod.Host + " " + pod.Category + " " + pod.Province)
			if !strings.Contains(fullText, searchLower) {
				continue
			}
		}
		results = append(results, pod)
	}
	return results
}
