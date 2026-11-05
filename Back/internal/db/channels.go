package db

import (
	"angopod-back/internal/youtube"
)

// GetConfiguredChannels busca todos os canais ativos na base de dados
func GetConfiguredChannels() ([]youtube.ChannelConfig, error) {
	rows, err := DB.Query(`
		SELECT channel_id, slug, COALESCE(host, ''), COALESCE(category, 'Geral'),
		       COALESCE(province, 'Luanda'), COALESCE(language, 'Português'), COALESCE(color, 'coral')
		FROM channel_configs
		WHERE active = 1
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []youtube.ChannelConfig
	for rows.Next() {
		var c youtube.ChannelConfig
		if err := rows.Scan(&c.ChannelID, &c.Slug, &c.Host, &c.Category, &c.Province, &c.Language, &c.Color); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

// AddChannel insere ou ativa um novo canal do YouTube para sincronização
func AddChannel(cfg youtube.ChannelConfig) error {
	_, err := DB.Exec(`
		INSERT INTO channel_configs (channel_id, slug, title, host, category, province, language, color, active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1)
		ON CONFLICT(channel_id) DO UPDATE SET
			slug=excluded.slug,
			host=excluded.host,
			category=excluded.category,
			province=excluded.province,
			language=excluded.language,
			color=excluded.color,
			active=1
	`, cfg.ChannelID, cfg.Slug, cfg.Slug, cfg.Host, cfg.Category, cfg.Province, cfg.Language, cfg.Color)
	return err
}

// DeleteChannel remove um canal e seus dados (podcasts e episódios associados)
func DeleteChannel(channelIDOrSlug string) error {
	// Remove episódios associados ao podcast
	_, _ = DB.Exec(`
		DELETE FROM episodes
		WHERE podcast_id = ? OR podcast_slug = ?
		   OR podcast_id IN (SELECT id FROM podcasts WHERE channel_id = ? OR slug = ?)
	`, channelIDOrSlug, channelIDOrSlug, channelIDOrSlug, channelIDOrSlug)

	// Remove o podcast
	_, _ = DB.Exec(`DELETE FROM podcasts WHERE channel_id = ? OR slug = ? OR id = ?`, channelIDOrSlug, channelIDOrSlug, channelIDOrSlug)

	// Remove da tabela de canais configurados
	_, err := DB.Exec(`DELETE FROM channel_configs WHERE channel_id = ? OR slug = ?`, channelIDOrSlug, channelIDOrSlug)

	return err
}

