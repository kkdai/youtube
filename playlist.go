package youtube

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"
)

const (
	playlist_fetch_url string = "https://youtube.com/list_ajax?style=json&action_get_list=1" +
		"&list={playlist}&hl=en"
)

var (
	id_regex *regexp.Regexp = regexp.MustCompile("[&?]list=([^&]{34})")
)

type playlistResponse struct {
	Title          string                   `json:"title"`
	Author         string                   `json:"author"`
	ResponseVideos []*playlistResponseVideo `json:"video"`
}

type playlistResponseVideo struct {
	ID       string `json:"encrypted_id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Duration int    `json:"duration_seconds"`
}

type Playlist struct {
	ID     string
	Title  string
	Author string
	Videos []*Video
}

func (p *Playlist) parsePlaylistResponse(info []byte) error {
	resp := new(playlistResponse)

	if err := json.Unmarshal([]byte(info), resp); err != nil {
		return err
	}

	p.Title = resp.Title
	p.Author = resp.Author
	var videos []*Video

	for _, v := range resp.ResponseVideos {
		d := time.Second * time.Duration(v.Duration)
		videos = append(videos, &Video{
			Title:    v.Title,
			Author:   v.Author,
			ID:       v.ID,
			Duration: d,
		})
	}

	p.Videos = videos
	return nil
}

func extractPlaylistID(url string) (string, error) {
	var id string

	if len(url) == 34 {
		id = url
	} else {
		matches := id_regex.FindStringSubmatch(id)

		if len(matches) == 2 {
			id = matches[1]
		} else {
			return "", ErrPlaylistIDMinLength
		}
	}

	if strings.Count(id, " ")+strings.Count(id, "\t") == len(id) {
		return "", ErrPlaylistIDEmpty
	}

	return id, nil
}
