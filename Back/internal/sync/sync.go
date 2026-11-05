package sync

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"angopod-back/internal/data"
	"angopod-back/internal/db"
	"angopod-back/internal/youtube"
)

var (
	isSyncing bool
	syncMu    sync.Mutex
)

// SyncStatus representa o estado da sincronização para consulta da API
type SyncStatus struct {
	IsSyncing       bool      `json:"isSyncing"`
	LastSyncAt      time.Time `json:"lastSyncAt"`
	ChannelsSynced  int       `json:"channelsSynced"`
	EpisodesSynced  int       `json:"episodesSynced"`
	Status          string    `json:"status"`
	Message         string    `json:"message"`
	YouTubeConfigured bool    `json:"youtubeConfigured"`
}

var lastStatus = SyncStatus{
	Status:  "idle",
	Message: "Nenhuma sincronização executada ainda",
}

// GetStatus retorna o estado atual da sincronização
func GetStatus(apiKey string) SyncStatus {
	syncMu.Lock()
	defer syncMu.Unlock()
	s := lastStatus
	s.IsSyncing = isSyncing
	s.YouTubeConfigured = apiKey != ""
	return s
}

// StartPeriodicSync inicializa a rotina de sincronização em segundo plano
func StartPeriodicSync(apiKey string, intervalHours int) {
	if intervalHours <= 0 {
		intervalHours = 6
	}

	go func() {
		// Executa na inicialização se a chave estiver configurada
		if apiKey != "" {
			log.Println("🚀 Chave do YouTube detetada. A iniciar sincronização inicial com a base de dados...")
			_, _, _ = RunSync(apiKey)
		} else {
			log.Println("ℹ️ YOUTUBE_API_KEY vazia no .env. A operar com dados locais da base de dados SQLite. Adiciona a chave para ativar a sincronização com o YouTube.")
		}

		ticker := time.NewTicker(time.Duration(intervalHours) * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			if apiKey != "" {
				log.Println("⏰ Intervalo de sincronização periódica atingido. A sincronizar com o YouTube...")
				_, _, _ = RunSync(apiKey)
			}
		}
	}()
}

// RunSync sincroniza canais e episódios com o YouTube Data API v3 e persiste na SQLite
func RunSync(apiKey string) (int, int, error) {
	if apiKey == "" {
		return 0, 0, fmt.Errorf("YOUTUBE_API_KEY não configurada no .env")
	}

	syncMu.Lock()
	if isSyncing {
		syncMu.Unlock()
		return 0, 0, fmt.Errorf("já existe uma sincronização em curso")
	}
	isSyncing = true
	syncMu.Unlock()

	defer func() {
		syncMu.Lock()
		isSyncing = false
		syncMu.Unlock()
	}()

	startTime := time.Now()
	client := youtube.NewClient(apiKey)

	// Registra início no log do banco
	res, _ := db.DB.Exec(
		"INSERT INTO sync_log (started_at, status) VALUES (?, ?)",
		startTime, "running",
	)
	var logID int64
	if res != nil {
		logID, _ = res.LastInsertId()
	}

	totalChannels := 0
	totalEpisodes := 0
	var syncErrors []string

	channelsToSync, err := db.GetConfiguredChannels()
	if err != nil || len(channelsToSync) == 0 {
		channelsToSync = data.ConfiguredChannels
	}

	// Duração mínima para considerar episódio de podcast completo (ignorar cortes/shorts)
	minMinutes := 30
	if envMin := os.Getenv("MIN_EPISODE_MINUTES"); envMin != "" {
		if parsed, err := strconv.Atoi(envMin); err == nil && parsed > 0 {
			minMinutes = parsed
		}
	}

	for _, cfg := range channelsToSync {
		chResp, err := client.FetchChannel(cfg.ChannelID)
		if err != nil {
			msg := fmt.Sprintf("Canal %s (%s): %v", cfg.Slug, cfg.ChannelID, err)
			log.Printf("⚠️ Erro ao sincronizar canal: %s", msg)
			syncErrors = append(syncErrors, msg)
			continue
		}

		pod, uploadsPlaylistID := youtube.ChannelToPodcast(chResp, cfg)
		if pod == nil {
			continue
		}

		// Persiste podcast na SQLite
		if err := db.UpsertPodcast(*pod, cfg.ChannelID); err != nil {
			log.Printf("⚠️ Erro ao persistir podcast %s: %v", pod.Slug, err)
			continue
		}
		totalChannels++

		// Se tiver playlist de uploads, busca os últimos vídeos (máx 50)
		if uploadsPlaylistID != "" {
			videoIDs, err := client.FetchPlaylistVideoIDs(uploadsPlaylistID, 50)
			if err != nil {
				log.Printf("⚠️ Erro ao buscar vídeos de %s: %v", pod.Slug, err)
				continue
			}

			if len(videoIDs) > 0 {
				vResp, err := client.FetchVideosDetails(videoIDs)
				if err != nil {
					log.Printf("⚠️ Erro ao buscar detalhes dos vídeos de %s: %v", pod.Slug, err)
					continue
				}

				for _, vItem := range vResp.Items {
					// Filtra cortes e shorts: só aceita vídeos com duração >= minMinutes
					if !youtube.IsFullEpisode(vItem.ContentDetails.Duration, minMinutes) {
						continue
					}

					ep := youtube.VideoToEpisode(vItem, pod)
					if err := db.UpsertEpisode(ep, pod.ID); err != nil {
						log.Printf("⚠️ Erro ao persistir episódio %s: %v", ep.ID, err)
						continue
					}
					totalEpisodes++
				}
			}
		}
	}

	finishTime := time.Now()
	statusText := "success"
	errMsg := ""
	if len(syncErrors) > 0 {
		if totalChannels > 0 {
			statusText = "partial"
		} else {
			statusText = "error"
		}
		errMsg = fmt.Sprintf("%d erros durante a sincronização", len(syncErrors))
	}

	// Atualiza log no banco
	if logID > 0 {
		_, _ = db.DB.Exec(
			"UPDATE sync_log SET finished_at=?, channels_synced=?, episodes_synced=?, status=?, error_message=? WHERE id=?",
			finishTime, totalChannels, totalEpisodes, statusText, errMsg, logID,
		)
	}

	syncMu.Lock()
	lastStatus = SyncStatus{
		IsSyncing:       false,
		LastSyncAt:      finishTime,
		ChannelsSynced:  totalChannels,
		EpisodesSynced:  totalEpisodes,
		Status:          statusText,
		Message:         fmt.Sprintf("Sincronizados %d canais e %d episódios em %s", totalChannels, totalEpisodes, finishTime.Sub(startTime).Round(time.Millisecond)),
		YouTubeConfigured: true,
	}
	syncMu.Unlock()

	log.Printf("🎉 Sincronização concluída! %d canais e %d episódios atualizados na SQLite.", totalChannels, totalEpisodes)
	return totalChannels, totalEpisodes, nil
}

