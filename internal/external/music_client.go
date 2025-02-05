// internal/external/music_client.go
package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"song-library-test-task/internal/service" // to use service.SongInfo
)

type musicInfoClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewMusicInfoClient(baseURL string, timeout time.Duration) service.ExternalClient {
	return &musicInfoClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *musicInfoClient) FetchSongInfo(ctx context.Context, groupName, songTitle string) (*service.SongInfo, error) {
	// Example GET request: baseURL/info?group=groupName&song=songTitle
	endpoint := fmt.Sprintf("%s/info", c.baseURL)

	// Build query
	u, _ := url.Parse(endpoint)
	q := u.Query()
	q.Set("group", groupName)
	q.Set("song", songTitle)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var data struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &service.SongInfo{
		ReleaseDate: data.ReleaseDate,
		Text:        data.Text,
		Link:        data.Link,
	}, nil
}
