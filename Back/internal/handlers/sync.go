package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	appSync "angopod-back/internal/sync"
)

// GetSyncStatus retorna o estado da sincronização
func GetSyncStatus(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	status := appSync.GetStatus(apiKey)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

// TriggerManualSync aciona uma sincronização manual imediata
func TriggerManualSync(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error":   "YOUTUBE_API_KEY não está configurada no .env",
			"status":  http.StatusBadRequest,
			"details": "Adiciona a tua chave no ficheiro Back/.env para ativar a sincronização com o YouTube.",
		})
		return
	}

	// Executa sincronização em background ou síncrona
	go func() {
		_, _, _ = appSync.RunSync(apiKey)
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Sincronização com o YouTube iniciada em segundo plano. Podes verificar o progresso em /api/sync/status",
		"status":  "processing",
	})
}

