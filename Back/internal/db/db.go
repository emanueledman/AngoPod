package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"angopod-back/internal/data"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB inicializa o banco SQLite e executa migrações automáticas
func InitDB(dbPath string) (*sql.DB, error) {
	if dbPath == "" {
		dbPath = "angopod.db"
	}

	// Garante que o diretório pai existe
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("falha ao criar pasta para base de dados: %w", err)
		}
	}

	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir conexão com SQLite: %w", err)
	}

	// Otimizações SQLite para concorrência e integridade
	if _, err := database.Exec("PRAGMA journal_mode = WAL; PRAGMA foreign_keys = ON;"); err != nil {
		log.Printf("⚠️ Aviso PRAGMA SQLite: %v", err)
	}

	DB = database

	if err := migrate(); err != nil {
		return nil, fmt.Errorf("falha nas migrações: %w", err)
	}

	// Se a tabela de podcasts estiver vazia, carrega o seed inicial
	seedInitialDataIfEmpty()

	return DB, nil
}

func migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS channel_configs (
		channel_id  TEXT PRIMARY KEY,
		slug        TEXT UNIQUE NOT NULL,
		title       TEXT NOT NULL,
		host        TEXT,
		category    TEXT,
		province    TEXT,
		language    TEXT DEFAULT 'Português',
		color       TEXT DEFAULT 'coral',
		active      BOOLEAN DEFAULT 1,
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS podcasts (
		id             TEXT PRIMARY KEY,
		slug           TEXT UNIQUE NOT NULL,
		rank           INTEGER DEFAULT 0,
		title          TEXT NOT NULL,
		host           TEXT,
		category       TEXT,
		province       TEXT,
		language       TEXT DEFAULT 'Português',
		platform       TEXT DEFAULT 'YouTube',
		format         TEXT DEFAULT 'Entrevistas',
		popularity     INTEGER DEFAULT 0,
		published_at   TEXT,
		status         TEXT DEFAULT 'Atualizado',
		image          TEXT,
		color          TEXT DEFAULT 'coral',
		description    TEXT,
		listeners      TEXT DEFAULT '0',
		episodes_count INTEGER DEFAULT 0,
		rating         REAL DEFAULT 5.0,
		channel_id     TEXT,
		synced_at      DATETIME,
		created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS episodes (
		id            TEXT PRIMARY KEY,
		podcast_id    TEXT NOT NULL,
		podcast_slug  TEXT NOT NULL,
		podcast_title TEXT NOT NULL,
		number        TEXT,
		title         TEXT NOT NULL,
		guest         TEXT,
		show          TEXT,
		time          TEXT,
		duration      TEXT,
		date          TEXT,
		listens       TEXT DEFAULT '0',
		trend         TEXT DEFAULT '+0%',
		image         TEXT,
		tag           TEXT,
		audio_url     TEXT,
		synced_at     DATETIME,
		created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(podcast_id) REFERENCES podcasts(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS provinces (
		name        TEXT PRIMARY KEY,
		count       INTEGER DEFAULT 0,
		x           INTEGER DEFAULT 0,
		y           INTEGER DEFAULT 0,
		description TEXT
	);

	CREATE TABLE IF NOT EXISTS sync_log (
		id              INTEGER PRIMARY KEY AUTOINCREMENT,
		started_at      DATETIME NOT NULL,
		finished_at     DATETIME,
		channels_synced INTEGER DEFAULT 0,
		episodes_synced INTEGER DEFAULT 0,
		status          TEXT,
		error_message   TEXT
	);
	`
	_, err := DB.Exec(schema)
	return err
}

func seedInitialDataIfEmpty() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM podcasts").Scan(&count)
	if err != nil || count > 0 {
		return
	}

	log.Println("🌱 Base de dados vazia. A semear dados iniciais de podcasts e províncias...")

	tx, err := DB.Begin()
	if err != nil {
		log.Printf("Erro ao iniciar transação de seed: %v", err)
		return
	}
	defer tx.Rollback()

	// Inserir Canais Configurados iniciais
	chanStmt, err := tx.Prepare(`
		INSERT INTO channel_configs (channel_id, slug, title, host, category, province, language, color, active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1)
		ON CONFLICT(channel_id) DO NOTHING
	`)
	if err == nil {
		defer chanStmt.Close()
		for _, c := range data.ConfiguredChannels {
			_, _ = chanStmt.Exec(c.ChannelID, c.Slug, c.Slug, c.Host, c.Category, c.Province, c.Language, c.Color)
		}
	}

	// As províncias são a geografia real de Angola (coordenadas para o mapa)
	provStmt, err := tx.Prepare(`
		INSERT INTO provinces (name, count, x, y, description)
		VALUES (?, 0, ?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET description=excluded.description
	`)
	if err == nil {
		defer provStmt.Close()
		for _, p := range data.Provinces {
			_, _ = provStmt.Exec(p.Name, p.X, p.Y, p.Description)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Erro ao efetivar configuração inicial: %v", err)
	} else {
		log.Println("✅ Base de dados inicializada limpa (sem dados fictícios).")
	}
}

// ClearAllData limpa todos os podcasts e episódios da base de dados
func ClearAllData() error {
	_, err := DB.Exec(`
		DELETE FROM episodes;
		DELETE FROM podcasts;
		DELETE FROM sync_log;
	`)
	return err
}

