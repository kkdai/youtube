package youtube

import (
	"encoding/json"
	"regexp"
	"time"
)

const (
	playlistFetchURL string = "https://youtube.com/list_ajax?style=json&action_get_list=1" +
		"&list=%s&hl=en"
)

var (
	pidRegex      *regexp.Regexp = regexp.MustCompile("^[A-Za-z0-9_-]{24,34}$")
	pidInURLRegex *regexp.Regexp = regexp.MustCompile("[&?]list=([A-Za-z0-9_-]{24,34})(&.*)?$")
)

type PlaylistEntry struct {
	ID       string `json:"encrypted_id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Duration time.Duration
}

func (p *PlaylistEntry) UnmarshalJSON(b []byte) error {
	var wf struct {
		ID              string `json:"encrypted_id"`
		Title           string `json:"title"`
		Author          string `json:"author"`
		DurationSeconds int    `json:"duration_seconds"`
	}
	if err := json.Unmarshal(b, &wf); err != nil {
		return err
	}
	p.ID, p.Title, p.Author = wf.ID, wf.Title, wf.Author
	p.Duration = time.Second * time.Duration(wf.DurationSeconds)
	return nil
}

func (p PlaylistEntry) MarshalJSON() ([]byte, error) {
	var wf = struct {
		ID              string `json:"encrypted_id"`
		Title           string `json:"title"`
		Author          string `json:"author"`
		DurationSeconds int    `json:"duration_seconds"`
	}{
		p.ID,
		p.Title,
		p.Author,
		int(p.Duration.Seconds()),
	}
	return json.Marshal(wf)
}

type Playlist struct {
	ID     string
	Title  string           `json:"title"`
	Author string           `json:"author"`
	Videos []*PlaylistEntry `json:"video"`
}

func extractPlaylistID(url string) (string, error) {
	if pidRegex.Match([]byte(url)) {
		return url, nil
	}

	matches := pidInURLRegex.FindStringSubmatch(url)

	if matches != nil {
		return matches[1], nil
	}

	return "", ErrInvalidPlaylist
}
