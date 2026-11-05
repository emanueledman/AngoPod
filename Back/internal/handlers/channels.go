package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"angopod-back/internal/data"
	"angopod-back/internal/db"
	appSync "angopod-back/internal/sync"
	"angopod-back/internal/youtube"
)

// GetChannels retorna a lista de canais configurados no sistema
func GetChannels(w http.ResponseWriter, r *http.Request) {
	channels, err := db.GetConfiguredChannels()
	if err != nil || len(channels) == 0 {
		channels = data.ConfiguredChannels
	}
	writeJSON(w, http.StatusOK, channels)
}

// AddChannelPayload dados necessários para adicionar um canal
type AddChannelPayload struct {
	ChannelID string `json:"channelId"`
	Slug      string `json:"slug"`
	Host      string `json:"host"`
	Category  string `json:"category"`
	Province  string `json:"province"`
	Language  string `json:"language"`
	Color     string `json:"color"`
}

// AddChannel adiciona um canal do YouTube e dispara sincronização imediata
func AddChannel(w http.ResponseWriter, r *http.Request) {
	var payload AddChannelPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Payload JSON inválido"})
		return
	}

	payload.ChannelID = strings.TrimSpace(payload.ChannelID)
	if payload.ChannelID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "channelId é obrigatório"})
		return
	}

	// Normaliza o slug (substitui espaços por traços, remove caracteres estranhos)
	cleanSlug := strings.TrimPrefix(payload.Slug, "@")
	if cleanSlug == "" {
		cleanSlug = strings.TrimPrefix(payload.ChannelID, "@")
	}
	cleanSlug = strings.ToLower(cleanSlug)
	cleanSlug = strings.ReplaceAll(cleanSlug, " ", "-")
	payload.Slug = cleanSlug

	if payload.Category == "" {
		payload.Category = "Sociedade"
	}
	if payload.Province == "" {
		payload.Province = "Luanda"
	}
	if payload.Language == "" {
		payload.Language = "Português"
	}
	if payload.Color == "" {
		payload.Color = "coral"
	}

	cfg := youtube.ChannelConfig{
		ChannelID: payload.ChannelID,
		Slug:      payload.Slug,
		Host:      payload.Host,
		Category:  payload.Category,
		Province:  payload.Province,
		Language:  payload.Language,
		Color:     payload.Color,
	}

	if err := db.AddChannel(cfg); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Falha ao gravar canal na base de dados"})
		return
	}

	// Dispara sincronização com o YouTube em background se a chave estiver presente
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey != "" {
		go func() {
			_, _, _ = appSync.RunSync(apiKey)
		}()
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message":   "Canal adicionado com sucesso. Sincronização iniciada.",
		"channelId": payload.ChannelID,
		"slug":      payload.Slug,
	})
}

// DeleteChannel remove um canal e os seus episódios da base de dados
func DeleteChannel(w http.ResponseWriter, r *http.Request) {
	idOrSlug := r.PathValue("id")
	if idOrSlug == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID ou Slug é obrigatório"})
		return
	}

	if err := db.DeleteChannel(idOrSlug); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Falha ao eliminar canal"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Canal e episódios eliminados com sucesso da base de dados",
		"target":  idOrSlug,
	})
}

