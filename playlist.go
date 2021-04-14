package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	sjson "github.com/bitly/go-simplejson"
	"golang.org/x/net/html"
)

const (
	playlistFetchURL string = "https://www.youtube.com/playlist?list=%s&hl=en"
	// The following are used in tests but also for fetching test data
	testPlaylistResponseDataFile = "./testdata/playlist_test_data.html"
	testPlaylistID               = "PL59FEE129ADFF2B12"
)

var (
	playlistIDRegex    = regexp.MustCompile("^[A-Za-z0-9_-]{24,34}$")
	playlistInURLRegex = regexp.MustCompile("[&?]list=([A-Za-z0-9_-]{24,34})(&.*)?$")
)

type Playlist struct {
	ID     string
	Title  string
	Author string
	Videos []*PlaylistEntry
}

type PlaylistEntry struct {
	ID       string
	Title    string
	Author   string
	Duration time.Duration
}

func extractPlaylistID(url string) (string, error) {
	if playlistIDRegex.Match([]byte(url)) {
		return url, nil
	}

	matches := playlistInURLRegex.FindStringSubmatch(url)

	if matches != nil {
		return matches[1], nil
	}

	return "", ErrInvalidPlaylist
}

func extractPlaylistJSON(r io.Reader) ([]byte, error) {
	const prefix = "var ytInitialData ="

	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	var data []byte
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" && n.FirstChild != nil {
			script := n.FirstChild.Data
			if strings.HasPrefix(script, prefix) {
				script = strings.TrimPrefix(script, prefix)
				data = []byte(strings.Trim(script, ";"))
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return data, nil
}

// structs for playlist extraction

// Title: metadata.playlistMetadataRenderer.title | sidebar.playlistSidebarRenderer.items[0].playlistSidebarPrimaryInfoRenderer.title.runs[0].text
// Author: sidebar.playlistSidebarRenderer.items[1].playlistSidebarSecondaryInfoRenderer.videoOwner.videoOwnerRenderer.title.runs[0].text

// Videos: contents.twoColumnBrowseResultsRenderer.tabs[0].tabRenderer.content.sectionListRenderer.contents[0].itemSectionRenderer.contents[0].playlistVideoListRenderer.contents
// ID: .videoId
// Title: title.runs[0].text
// Author: .shortBylineText.runs[0].text
// Duration: .lengthSeconds

// TODO?: Author thumbnails: sidebar.playlistSidebarRenderer.items[0].playlistSidebarPrimaryInfoRenderer.thumbnailRenderer.playlistVideoThumbnailRenderer.thumbnail.thumbnails
// TODO? Video thumbnails: .thumbnail.thumbnails

func (p *Playlist) UnmarshalJSON(b []byte) (err error) {
	var j *sjson.Json
	j, err = sjson.NewJson(b)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("JSON parsing error: %v", r)
		}
	}()
	p.Title = j.GetPath("metadata", "playlistMetadataRenderer", "title").MustString()
	p.Author = j.GetPath("sidebar", "playlistSidebarRenderer", "items").GetIndex(1).
		GetPath("playlistSidebarSecondaryInfoRenderer", "videoOwner", "videoOwnerRenderer", "title", "runs").
		GetIndex(0).Get("text").MustString()
	vJSON, err := j.GetPath("contents", "twoColumnBrowseResultsRenderer", "tabs").GetIndex(0).
		GetPath("tabRenderer", "content", "sectionListRenderer", "contents").GetIndex(0).
		GetPath("itemSectionRenderer", "contents").GetIndex(0).
		GetPath("playlistVideoListRenderer", "contents").MarshalJSON()

	var vids []*videosJSONExtractor
	if err := json.Unmarshal(vJSON, &vids); err != nil {
		return err
	}
	p.Videos = make([]*PlaylistEntry, 0, len(vids))
	for _, v := range vids {
		if v.Renderer == nil { // items such as continuationItemRenderer can mess things up in that array
			continue
		}
		p.Videos = append(p.Videos, v.PlaylistEntry())
	}
	return nil
}

type videosJSONExtractor struct {
	Renderer *struct {
		ID       string   `json:"videoId"`
		Title    withRuns `json:"title"`
		Author   withRuns `json:"shortBylineText"`
		Duration string   `json:"lengthSeconds"`
	} `json:"playlistVideoRenderer"`
}

func (vje videosJSONExtractor) PlaylistEntry() *PlaylistEntry {
	ds, err := strconv.Atoi(vje.Renderer.Duration)
	if err != nil {
		panic("invalid video duration: " + vje.Renderer.Duration)
	}
	return &PlaylistEntry{
		ID:       vje.Renderer.ID,
		Title:    vje.Renderer.Title.String(),
		Author:   vje.Renderer.Author.String(),
		Duration: time.Second * time.Duration(ds),
	}
}

type withRuns struct {
	Runs []struct {
		Text string `json:"text"`
	} `json:"runs"`
}

func (wr withRuns) String() string {
	if len(wr.Runs) > 0 {
		return wr.Runs[0].Text
	}
	return ""
}
