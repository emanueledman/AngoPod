package db

import (
	"database/sql"
	"fmt"

	"angopod-back/internal/models"
)

// GetEpisodes busca episódios, opcionalmente filtrados por podcastSlug
func GetEpisodes(podcastSlug string) ([]models.Episode, error) {
	query := `
		SELECT id, podcast_id, podcast_slug, podcast_title, number, title,
		       guest, show, time, duration, date, listens, trend, image, tag, audio_url
		FROM episodes
	`
	var args []any
	if podcastSlug != "" {
		query += " WHERE podcast_slug = ?"
		args = append(args, podcastSlug)
	}
	query += " ORDER BY date DESC, created_at DESC"

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar episódios: %w", err)
	}
	defer rows.Close()

	var result []models.Episode
	for rows.Next() {
		var ep models.Episode
		var podID sql.NullString
		var num, guest, show, t, dur, dt, listens, tr, img, tag, aurl sql.NullString
		err := rows.Scan(
			&ep.ID, &podID, &ep.PodcastSlug, &ep.PodcastTitle, &num, &ep.Title,
			&guest, &show, &t, &dur, &dt, &listens, &tr, &img, &tag, &aurl,
		)
		if err != nil {
			return nil, err
		}
		ep.Number = num.String
		ep.Guest = guest.String
		ep.Show = show.String
		ep.Time = t.String
		ep.Duration = dur.String
		ep.Date = dt.String
		ep.Listens = listens.String
		ep.Trend = tr.String
		ep.Image = img.String
		ep.Tag = tag.String
		ep.AudioURL = aurl.String

		result = append(result, ep)
	}

	return result, nil
}

// GetEpisodeByID busca um episódio pelo ID
func GetEpisodeByID(id string) (*models.Episode, error) {
	query := `
		SELECT id, podcast_id, podcast_slug, podcast_title, number, title,
		       guest, show, time, duration, date, listens, trend, image, tag, audio_url
		FROM episodes
		WHERE id = ?
		LIMIT 1
	`
	var ep models.Episode
	var podID sql.NullString
	var num, guest, show, t, dur, dt, listens, tr, img, tag, aurl sql.NullString
	err := DB.QueryRow(query, id).Scan(
		&ep.ID, &podID, &ep.PodcastSlug, &ep.PodcastTitle, &num, &ep.Title,
		&guest, &show, &t, &dur, &dt, &listens, &tr, &img, &tag, &aurl,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	ep.Number = num.String
	ep.Guest = guest.String
	ep.Show = show.String
	ep.Time = t.String
	ep.Duration = dur.String
	ep.Date = dt.String
	ep.Listens = listens.String
	ep.Trend = tr.String
	ep.Image = img.String
	ep.Tag = tag.String
	ep.AudioURL = aurl.String

	return &ep, nil
}

// UpsertEpisode insere ou atualiza um episódio vindo do YouTube
func UpsertEpisode(ep models.Episode, podcastID string) error {
	query := `
		INSERT INTO episodes (
			id, podcast_id, podcast_slug, podcast_title, number, title,
			guest, show, time, duration, date, listens, trend, image, tag, audio_url, synced_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(id) DO UPDATE SET
			title=excluded.title,
			duration=excluded.duration,
			time=excluded.time,
			image=excluded.image,
			audio_url=excluded.audio_url,
			synced_at=CURRENT_TIMESTAMP
	`
	_, err := DB.Exec(query,
		ep.ID, podcastID, ep.PodcastSlug, ep.PodcastTitle, ep.Number, ep.Title,
		ep.Guest, ep.Show, ep.Time, ep.Duration, ep.Date, ep.Listens, ep.Trend,
		ep.Image, ep.Tag, ep.AudioURL,
	)
	return err
}

