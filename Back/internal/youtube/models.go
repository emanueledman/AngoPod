package youtube

// ChannelConfig define os metadados editoriais que associamos a um canal do YouTube
type ChannelConfig struct {
	ChannelID string `json:"channelId"`
	Slug      string `json:"slug"`
	Host      string `json:"host"`
	Category  string `json:"category"`
	Province  string `json:"province"`
	Language  string `json:"language"`
	Color     string `json:"color"`
}

// Respostas da YouTube Data API v3
type ChannelListResponse struct {
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Thumbnails  struct {
				Default struct{ URL string `json:"url"` } `json:"default"`
				Medium  struct{ URL string `json:"url"` } `json:"medium"`
				High    struct{ URL string `json:"url"` } `json:"high"`
			} `json:"thumbnails"`
			PublishedAt string `json:"publishedAt"`
		} `json:"snippet"`
		Statistics struct {
			SubscriberCount string `json:"subscriberCount"`
			VideoCount      string `json:"videoCount"`
			ViewCount       string `json:"viewCount"`
		} `json:"statistics"`
		ContentDetails struct {
			RelatedPlaylists struct {
				Uploads string `json:"uploads"`
			} `json:"relatedPlaylists"`
		} `json:"contentDetails"`
	} `json:"items"`
}

type PlaylistItemsResponse struct {
	Items []struct {
		Snippet struct {
			ResourceID struct {
				VideoID string `json:"videoId"`
			} `json:"resourceId"`
		} `json:"snippet"`
	} `json:"items"`
}

type VideoListResponse struct {
	Items []struct {
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
			Duration string `json:"duration"` // Formato ISO 8601 ex: PT45M12S
		} `json:"contentDetails"`
		Statistics struct {
			ViewCount string `json:"viewCount"`
			LikeCount string `json:"likeCount"`
		} `json:"statistics"`
	} `json:"items"`
}

