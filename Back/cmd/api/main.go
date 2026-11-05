package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"angopod-back/internal/db"
	"angopod-back/internal/handlers"
	"angopod-back/internal/middleware"
	appSync "angopod-back/internal/sync"
)

// loadEnv carrega automaticamente as variáveis de um ficheiro .env se presente
func loadEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if os.Getenv(key) == "" {
				_ = os.Setenv(key, val)
			}
		}
	}
}

func main() {
	// Carrega .env do diretório atual se existir
	loadEnv(".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "angopod.db"
	}

	// 1. Inicializa a Base de Dados SQLite
	database, err := db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("🚨 Erro fatal ao inicializar base de dados SQLite: %v", err)
	}
	defer database.Close()
	log.Printf("📦 Base de dados SQLite ativa em: %s", dbPath)

	// 2. Inicia sincronização com o YouTube (goroutine em segundo plano)
	youtubeKey := os.Getenv("YOUTUBE_API_KEY")
	syncIntervalHours := 6
	if h := os.Getenv("SYNC_INTERVAL_HOURS"); h != "" {
		if parsed, err := strconv.Atoi(h); err == nil && parsed > 0 {
			syncIntervalHours = parsed
		}
	}
	appSync.StartPeriodicSync(youtubeKey, syncIntervalHours)

	// Origens permitidas para CORS (whitelist configurável via .env)
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsEnv != "" {
		for _, o := range strings.Split(allowedOriginsEnv, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	} else {
		allowedOrigins = []string{
			"http://localhost:5173",
			"http://localhost:5174",
			"http://localhost:3000",
			"http://127.0.0.1:5173",
			"http://127.0.0.1:5174",
		}
	}

	// Rate limit configurável via .env
	rateLimit := 120
	if rl := os.Getenv("RATE_LIMIT_PER_MINUTE"); rl != "" {
		if parsed, err := strconv.Atoi(rl); err == nil && parsed > 0 {
			rateLimit = parsed
		}
	}

	mux := http.NewServeMux()

	// Health check seguro com estatísticas do SQLite e YouTube
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":            "healthy",
			"service":           "angopod-backend-go",
			"database":          "sqlite",
			"db_path":           dbPath,
			"youtube_connected": youtubeKey != "",
			"security":          "maximum",
			"timestamp":         time.Now().UTC(),
		})
	})

	// Rotas de Podcasts e Episódios (servidas via SQLite)
	mux.HandleFunc("GET /api/podcasts", handlers.GetAllPodcasts)
	mux.HandleFunc("GET /api/podcasts/{slug}", handlers.GetPodcastBySlug)
	mux.HandleFunc("GET /api/podcasts/{slug}/episodes", handlers.GetPodcastEpisodes)

	mux.HandleFunc("GET /api/episodes", handlers.GetAllEpisodes)
	mux.HandleFunc("GET /api/episodes/{id}", handlers.GetEpisodeByID)

	mux.HandleFunc("GET /api/presenters", handlers.GetAllPresenters)
	mux.HandleFunc("GET /api/presenters/{slug}", handlers.GetPresenterBySlug)

	mux.HandleFunc("GET /api/provinces", handlers.GetAllProvinces)

	mux.HandleFunc("GET /api/ranking", handlers.GetRanking)
	mux.HandleFunc("GET /api/highlights", handlers.GetHighlights)

	// Rotas de Sincronização com o YouTube
	mux.HandleFunc("GET /api/sync/status", handlers.GetSyncStatus)
	mux.HandleFunc("POST /api/sync", handlers.TriggerManualSync)

	// Rotas de Gestão de Canais (Adicionar / Eliminar)
	mux.HandleFunc("GET /api/channels", handlers.GetChannels)
	mux.HandleFunc("POST /api/channels", handlers.AddChannel)
	mux.HandleFunc("DELETE /api/channels/{id}", handlers.DeleteChannel)

	// Pipeline de Segurança Máxima
	var handler http.Handler = mux
	handler = middleware.LimitBodySize(1 << 20)(handler)            // 1 MB
	handler = middleware.RateLimit(rateLimit, time.Minute)(handler) // Configurado via .env
	handler = middleware.SanitizeRequest(handler)                   // Sanitização de entradas
	handler = middleware.SecurityHeaders(handler)                   // Cabeçalhos OWASP
	handler = middleware.CORSStrict(allowedOrigins)(handler)        // CORS Restrito (intercepta OPTIONS e anexa headers)
	handler = middleware.Logger(handler)                            // Logger de auditoria
	handler = middleware.RecoverPanic(handler)                      // Anti-crash / Panic Recovery

	// Servidor HTTP reforçado contra Slowloris
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,  // Mitigação Slowloris
		ReadTimeout:       5 * time.Second,  // Timeout leitura payload
		WriteTimeout:      10 * time.Second, // Timeout envio resposta
		IdleTimeout:       60 * time.Second, // Timeout conexões inativas
		MaxHeaderBytes:    1 << 20,          // 1 MB limite no cabeçalho
	}

	fmt.Println("==========================================================")
	fmt.Printf("🛡️  AngoPod Back-end (Go) ativo na porta %s [SEGURANÇA MÁXIMA]\n", port)
	fmt.Printf("📦 Base de Dados: SQLite (%s)\n", dbPath)
	if youtubeKey != "" {
		fmt.Printf("📺 YouTube Sync:  ATIVO (a cada %dh)\n", syncIntervalHours)
	} else {
		fmt.Printf("📺 YouTube Sync:  À espera de YOUTUBE_API_KEY no .env\n")
	}
	fmt.Printf("📄 .env carregado: PORT=%s, RATE_LIMIT=%d/min\n", port, rateLimit)
	fmt.Printf("🔒 Headers OWASP: Ativos (CSP, HSTS, X-Frame, X-Content-Type)\n")
	fmt.Printf("🔒 CORS Estrito:  %s\n", strings.Join(allowedOrigins, ", "))
	fmt.Printf("📍 Health Check:  http://localhost:%s/api/health\n", port)
	fmt.Printf("📍 Sync Status:   http://localhost:%s/api/sync/status\n", port)
	fmt.Println("==========================================================")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Erro ao iniciar o servidor Go: %v", err)
	}
}
