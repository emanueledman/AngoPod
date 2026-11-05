package handlers

import (
	"net/http"

	"angopod-back/internal/data"
	"angopod-back/internal/db"
)

// GetAllProvinces retorna todas as províncias com contagens atualizadas
func GetAllProvinces(w http.ResponseWriter, r *http.Request) {
	provinces, err := db.GetProvinces()
	if err != nil || len(provinces) == 0 {
		provinces = data.Provinces
	}
	writeJSON(w, http.StatusOK, provinces)
}
