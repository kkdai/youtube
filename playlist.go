package youtube

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

const (
	playlistFetchURL string = "https://youtube.com/list_ajax?style=json&action_get_list=1" +
		"&list={playlist}&hl=en"
)

var (
	idRegex *regexp.Regexp = regexp.MustCompile("[&?]list=([^&]{34})")
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

// Parse information provided from youtube on this playlist; only basic information.
func (p *Playlist) parsePlaylistResponse(info string) error {
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
			Loaded:   WeaklyLoaded, // Videos are only partially loaded with data from playlist.
		})
	}

	p.Videos = videos
	return nil
}

// Fetch, from youtube, the information for this playlist (Author, Title, Description, etc...) along
// with a list of videos and their basic information, such as ID, Title, Author. These videos cannot
// be downloaded until more information is loaded!
func (p *Playlist) FetchPlaylistInfo(ctx context.Context, c *Client) ([]byte, error) {
	requestURL := strings.Replace(playlistFetchURL, "{playlist}", p.ID, 1)
	resp, err := c.httpGet(ctx, requestURL)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// Fetch information on this playlist from youtube, and then - providing there was no error - parse
// this information.
func (p *Playlist) LoadInfo(ctx context.Context, c *Client) error {
	body, err := p.FetchPlaylistInfo(ctx, c)

	if err != nil {
		return err
	}

	return p.parsePlaylistResponse(string(body))
}

func extractPlaylistID(url string) (string, error) {
	var id string

	if len(url) == 34 {
		id = url
	} else {
		matches := idRegex.FindStringSubmatch(url)

		if len(matches) == 2 {
			id = matches[1]
		} else {
			return "", ErrPlaylistIDMinLength
		}
	}

	return id, nil
}
