package db

import (
	"database/sql"
	"fmt"
	"strings"

	"angopod-back/internal/models"
)

// GetPodcasts busca podcasts aplicando filtros dinâmicos
func GetPodcasts(category, province, language, platform, search string) ([]models.Podcast, error) {
	query := `
		SELECT id, slug, rank, title, host, category, province, language,
		       platform, format, popularity, published_at, status, image,
		       color, description, listeners, episodes_count, rating
		FROM podcasts
		WHERE 1=1
	`
	var args []any

	if category != "" && !strings.EqualFold(category, "todos") && !strings.EqualFold(category, "todas") {
		query += " AND LOWER(category) = LOWER(?)"
		args = append(args, category)
	}

	if province != "" && !strings.EqualFold(province, "todas") && !strings.EqualFold(province, "todos") {
		query += " AND LOWER(province) = LOWER(?)"
		args = append(args, province)
	}

	if language != "" && !strings.EqualFold(language, "todos") && !strings.EqualFold(language, "todas") {
		query += " AND LOWER(language) = LOWER(?)"
		args = append(args, language)
	}

	if platform != "" && !strings.EqualFold(platform, "todas") && !strings.EqualFold(platform, "todos") {
		query += " AND LOWER(platform) = LOWER(?)"
		args = append(args, platform)
	}

	if search != "" {
		term := "%" + strings.ToLower(search) + "%"
		query += " AND (LOWER(title) LIKE ? OR LOWER(host) LIKE ? OR LOWER(description) LIKE ?)"
		args = append(args, term, term, term)
	}

	query += " ORDER BY rank ASC, popularity DESC"

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar podcasts: %w", err)
	}
	defer rows.Close()

	var result []models.Podcast
	for rows.Next() {
		var p models.Podcast
		var pubAt, desc, host, cat, prov, lang, plat, format, status, img, color, listeners sql.NullString
		err := rows.Scan(
			&p.ID, &p.Slug, &p.Rank, &p.Title, &host, &cat, &prov, &lang,
			&plat, &format, &p.Popularity, &pubAt, &status, &img,
			&color, &desc, &listeners, &p.EpisodesCount, &p.Rating,
		)
		if err != nil {
			return nil, err
		}
		p.Host = host.String
		p.Category = cat.String
		p.Province = prov.String
		p.Language = lang.String
		p.Platform = plat.String
		p.Format = format.String
		p.PublishedAt = pubAt.String
		p.Status = status.String
		p.Image = img.String
		p.Color = color.String
		p.Description = desc.String
		p.Listeners = listeners.String
		result = append(result, p)
	}

	return result, nil
}

// GetPodcastBySlug busca um podcast pelo slug
func GetPodcastBySlug(slug string) (*models.Podcast, error) {
	query := `
		SELECT id, slug, rank, title, host, category, province, language,
		       platform, format, popularity, published_at, status, image,
		       color, description, listeners, episodes_count, rating
		FROM podcasts
		WHERE slug = ?
		LIMIT 1
	`
	var p models.Podcast
	var pubAt, desc, host, cat, prov, lang, plat, format, status, img, color, listeners sql.NullString
	err := DB.QueryRow(query, slug).Scan(
		&p.ID, &p.Slug, &p.Rank, &p.Title, &host, &cat, &prov, &lang,
		&plat, &format, &p.Popularity, &pubAt, &status, &img,
		&color, &desc, &listeners, &p.EpisodesCount, &p.Rating,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	p.Host = host.String
	p.Category = cat.String
	p.Province = prov.String
	p.Language = lang.String
	p.Platform = plat.String
	p.Format = format.String
	p.PublishedAt = pubAt.String
	p.Status = status.String
	p.Image = img.String
	p.Color = color.String
	p.Description = desc.String
	p.Listeners = listeners.String

	return &p, nil
}

// UpsertPodcast insere ou atualiza os dados de um podcast sincronizado do YouTube
func UpsertPodcast(p models.Podcast, channelID string) error {
	query := `
		INSERT INTO podcasts (
			id, slug, rank, title, host, category, province, language,
			platform, format, popularity, published_at, status, image,
			color, description, listeners, episodes_count, rating, channel_id,
			synced_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT(id) DO UPDATE SET
			title=excluded.title,
			description=excluded.description,
			image=excluded.image,
			listeners=excluded.listeners,
			episodes_count=excluded.episodes_count,
			synced_at=CURRENT_TIMESTAMP,
			updated_at=CURRENT_TIMESTAMP
	`
	_, err := DB.Exec(query,
		p.ID, p.Slug, p.Rank, p.Title, p.Host, p.Category, p.Province, p.Language,
		p.Platform, p.Format, p.Popularity, p.PublishedAt, p.Status, p.Image,
		p.Color, p.Description, p.Listeners, p.EpisodesCount, p.Rating, channelID,
	)
	return err
}

