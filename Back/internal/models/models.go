package models

type Podcast struct {
	ID            string  `json:"id"`
	Slug          string  `json:"slug"`
	Rank          int     `json:"rank"`
	Title         string  `json:"title"`
	Host          string  `json:"host"`
	Category      string  `json:"category"`
	Province      string  `json:"province"`
	Language      string  `json:"language"`
	Platform      string  `json:"platform"`
	Format        string  `json:"format"`
	Popularity    int     `json:"popularity"`
	PublishedAt   string  `json:"publishedAt"`
	Status        string  `json:"status"`
	Image         string  `json:"image"`
	Color         string  `json:"color"`
	Description   string  `json:"description"`
	Listeners     string  `json:"listeners"`
	EpisodesCount int     `json:"episodesCount"`
	Rating        float64 `json:"rating"`
}

type Episode struct {
	ID           string `json:"id"`
	Number       string `json:"number,omitempty"`
	Title        string `json:"title"`
	PodcastTitle string `json:"podcastTitle"`
	PodcastSlug  string `json:"podcastSlug"`
	Guest        string `json:"guest,omitempty"`
	Show         string `json:"show,omitempty"`
	Time         string `json:"time,omitempty"`
	Duration     string `json:"duration,omitempty"`
	Date         string `json:"date,omitempty"`
	Listens      string `json:"listens,omitempty"`
	Trend        string `json:"trend,omitempty"`
	Image        string `json:"image,omitempty"`
	Tag          string `json:"tag,omitempty"`
	AudioURL     string `json:"audioUrl,omitempty"`
}

type Presenter struct {
	ID            string `json:"id"`
	Slug          string `json:"slug"`
	Name          string `json:"name"`
	Role          string `json:"role"`
	Quote         string `json:"quote"`
	Bio           string `json:"bio"`
	Story1        string `json:"story1"`
	Story2        string `json:"story2"`
	Photo         string `json:"photo"`
	BadgeNumber   string `json:"badgeNumber"`
	Followers     string `json:"followers"`
	Listeners     string `json:"listeners"`
	Conversations int    `json:"conversations"`
	PodcastSlug   string `json:"podcastSlug"`
	PodcastTitle  string `json:"podcastTitle"`
}

type Province struct {
	Name        string `json:"name"`
	Count       int    `json:"count"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Description string `json:"description,omitempty"`
}

type PulseHighlight struct {
	Label string `json:"label"`
	Title string `json:"title"`
	Href  string `json:"href"`
}

type MoverItem struct {
	Type       string `json:"type"`
	BadgeLabel string `json:"badgeLabel"`
	Value      string `json:"value"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Image      string `json:"image"`
	LinkHref   string `json:"linkHref"`
	LinkLabel  string `json:"linkLabel"`
}

type RankingResponse struct {
	TopShows     []Podcast   `json:"topShows"`
	TopEpisodes  []Episode   `json:"topEpisodes"`
	Movers       []MoverItem `json:"movers"`
	Weeks        []string    `json:"weeks"`
	UpdatedLabel string      `json:"updatedLabel"`
}

