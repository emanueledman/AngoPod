package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// FetchChannel busca informações e estatísticas do canal por ID (UC...) ou handle (@...)
func (c *Client) FetchChannel(channelIdentifier string) (*ChannelListResponse, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("YOUTUBE_API_KEY não configurada")
	}

	channelIdentifier = strings.TrimSpace(channelIdentifier)
	var endpoint string

	// Se começar com @ ou não for ID típico do YouTube (UC...), busca pelo handle
	if strings.HasPrefix(channelIdentifier, "@") || (!strings.HasPrefix(channelIdentifier, "UC") && len(channelIdentifier) < 20) {
		handle := channelIdentifier
		if !strings.HasPrefix(handle, "@") {
			handle = "@" + handle
		}
		endpoint = fmt.Sprintf(
			"https://www.googleapis.com/youtube/v3/channels?part=snippet,statistics,contentDetails&forHandle=%s&key=%s",
			url.QueryEscape(handle),
			url.QueryEscape(c.apiKey),
		)
	} else {
		endpoint = fmt.Sprintf(
			"https://www.googleapis.com/youtube/v3/channels?part=snippet,statistics,contentDetails&id=%s&key=%s",
			url.QueryEscape(channelIdentifier),
			url.QueryEscape(c.apiKey),
		)
	}

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição do canal YouTube: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro na API do YouTube (HTTP %d)", resp.StatusCode)
	}

	var data ChannelListResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("falha ao descodificar resposta do canal: %w", err)
	}

	if len(data.Items) == 0 {
		return nil, fmt.Errorf("nenhum canal encontrado no YouTube para '%s'", channelIdentifier)
	}

	return &data, nil
}

// FetchPlaylistVideoIDs busca os IDs dos vídeos recentes de uma playlist (geralmente uploads)
func (c *Client) FetchPlaylistVideoIDs(playlistID string, maxResults int) ([]string, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("YOUTUBE_API_KEY não configurada")
	}
	if maxResults <= 0 || maxResults > 50 {
		maxResults = 25
	}

	endpoint := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&playlistId=%s&maxResults=%d&key=%s",
		url.QueryEscape(playlistID),
		maxResults,
		url.QueryEscape(c.apiKey),
	)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição de playlist do YouTube: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro na API de playlist (HTTP %d)", resp.StatusCode)
	}

	var data PlaylistItemsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("falha ao descodificar playlist: %w", err)
	}

	var videoIDs []string
	for _, item := range data.Items {
		if id := item.Snippet.ResourceID.VideoID; id != "" {
			videoIDs = append(videoIDs, id)
		}
	}

	return videoIDs, nil
}

// FetchVideosDetails busca detalhes (duração, estatísticas) de múltiplos vídeos num único pedido
func (c *Client) FetchVideosDetails(videoIDs []string) (*VideoListResponse, error) {
	if c.apiKey == "" || len(videoIDs) == 0 {
		return &VideoListResponse{}, nil
	}

	idsParam := strings.Join(videoIDs, ",")
	endpoint := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails,statistics&id=%s&key=%s",
		url.QueryEscape(idsParam),
		url.QueryEscape(c.apiKey),
	)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição de vídeos do YouTube: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro na API de vídeos (HTTP %d)", resp.StatusCode)
	}

	var data VideoListResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("falha ao descodificar vídeos: %w", err)
	}

	return &data, nil
}

