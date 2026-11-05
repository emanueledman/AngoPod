package youtube

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"angopod-back/internal/models"
)

// ChannelToPodcast transforma a resposta do canal em models.Podcast
func ChannelToPodcast(resp *ChannelListResponse, cfg ChannelConfig) (*models.Podcast, string) {
	if resp == nil || len(resp.Items) == 0 {
		return nil, ""
	}
	item := resp.Items[0]

	subscribers, _ := strconv.ParseInt(item.Statistics.SubscriberCount, 10, 64)
	videoCount, _ := strconv.Atoi(item.Statistics.VideoCount)

	image := item.Snippet.Thumbnails.High.URL
	if image == "" {
		image = item.Snippet.Thumbnails.Medium.URL
	}
	if image == "" {
		image = item.Snippet.Thumbnails.Default.URL
	}

	publishedAt := item.Snippet.PublishedAt
	if t, err := time.Parse(time.RFC3339, publishedAt); err == nil {
		publishedAt = t.Format("2006-01-02")
	}

	pod := models.Podcast{
		ID:            item.ID,
		Slug:          cfg.Slug,
		Title:         item.Snippet.Title,
		Host:          cfg.Host,
		Category:      cfg.Category,
		Province:      cfg.Province,
		Language:      cfg.Language,
		Platform:      "YouTube",
		Format:        "Vídeo Podcast",
		Popularity:    calculatePopularity(subscribers),
		PublishedAt:   publishedAt,
		Status:        "Atualizado",
		Image:         image,
		Color:         cfg.Color,
		Description:   item.Snippet.Description,
		Listeners:     FormatCount(subscribers),
		EpisodesCount: videoCount,
		Rating:        4.9,
	}

	uploadsPlaylistID := item.ContentDetails.RelatedPlaylists.Uploads
	return &pod, uploadsPlaylistID
}

// VideoToEpisode transforma um item de vídeo em models.Episode
func VideoToEpisode(item struct {
	ID      string `json:"id"`
	Snippet struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		PublishedAt string `json:"publishedAt"`
		Thumbnails  struct {
			Medium struct{ URL string `json:"url"` } `json:"medium"`
			High   struct{ URL string `json:"url"` } `json:"high"`
		} `json:"thumbnails"`
	} `json:"snippet"`
	ContentDetails struct {
		Duration string `json:"duration"`
	} `json:"contentDetails"`
	Statistics struct {
		ViewCount string `json:"viewCount"`
		LikeCount string `json:"likeCount"`
	} `json:"statistics"`
}, pod *models.Podcast) models.Episode {
	views, _ := strconv.ParseInt(item.Statistics.ViewCount, 10, 64)

	date := item.Snippet.PublishedAt
	if t, err := time.Parse(time.RFC3339, date); err == nil {
		date = t.Format("02 Jan 2006")
	}

	image := item.Snippet.Thumbnails.High.URL
	if image == "" {
		image = item.Snippet.Thumbnails.Medium.URL
	}

	duration := parseISO8601Duration(item.ContentDetails.Duration)

	return models.Episode{
		ID:           item.ID,
		Title:        item.Snippet.Title,
		PodcastTitle: pod.Title,
		PodcastSlug:  pod.Slug,
		Show:         pod.Title,
		Time:         duration,
		Duration:     duration,
		Date:         date,
		Listens:      FormatCount(views),
		Trend:        "+12%",
		Image:        image,
		Tag:          pod.Category,
		AudioURL:     fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.ID),
	}
}

// IsFullEpisode verifica se o vídeo atinge a duração mínima para ser considerado podcast completo (e não corte/short)
func IsFullEpisode(isoDuration string, minMinutes int) bool {
	if minMinutes <= 0 {
		minMinutes = 30 // Padrão: 30 minutos
	}
	seconds := parseISO8601Seconds(isoDuration)
	return seconds >= (minMinutes * 60)
}

func parseISO8601Seconds(iso string) int {
	re := regexp.MustCompile(`PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`)
	matches := re.FindStringSubmatch(iso)
	if len(matches) < 4 {
		return 0
	}
	hours, _ := strconv.Atoi(matches[1])
	minutes, _ := strconv.Atoi(matches[2])
	seconds, _ := strconv.Atoi(matches[3])
	return (hours * 3600) + (minutes * 60) + seconds
}

// parseISO8601Duration converte "PT1H23M45S" -> "1h 23 min", "PT48M10S" -> "48 min"
func parseISO8601Duration(iso string) string {
	if iso == "" {
		return "45 min"
	}
	re := regexp.MustCompile(`PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`)
	matches := re.FindStringSubmatch(iso)
	if len(matches) < 4 {
		return "45 min"
	}

	hours := matches[1]
	minutes := matches[2]

	if hours != "" && minutes != "" {
		return fmt.Sprintf("%sh %s min", hours, minutes)
	}
	if hours != "" {
		return fmt.Sprintf("%sh", hours)
	}
	if minutes != "" {
		return fmt.Sprintf("%s min", minutes)
	}
	return "45 min"
}

// FormatCount converte 125000 -> "125K", 1400000 -> "1.4M"
func FormatCount(n int64) string {
	if n >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000.0)
	}
	if n >= 1_000 {
		return fmt.Sprintf("%dK", n/1_000)
	}
	return fmt.Sprintf("%d", n)
}

func calculatePopularity(subscribers int64) int {
	if subscribers > 500_000 {
		return 99
	}
	if subscribers > 100_000 {
		return 95
	}
	if subscribers > 50_000 {
		return 90
	}
	if subscribers > 10_000 {
		return 85
	}
	return 75
}
