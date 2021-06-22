package youtube

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Client offers methods to download video metadata and video streams.
type Client struct {
	// Debug enables debugging output through log package
	Debug bool

	// HTTPClient can be used to set a custom HTTP client.
	// If not set, http.DefaultClient will be used
	HTTPClient *http.Client

	// decipherOpsCache cache decipher operations
	decipherOpsCache DecipherOperationsCache
}

// GetVideo fetches video metadata
func (c *Client) GetVideo(url string) (*Video, error) {
	return c.GetVideoContext(context.Background(), url)
}

// GetVideoContext fetches video metadata with a context
func (c *Client) GetVideoContext(ctx context.Context, url string) (*Video, error) {
	id, err := ExtractVideoID(url)
	if err != nil {
		return nil, fmt.Errorf("extractVideoID failed: %w", err)
	}
	return c.videoFromID(ctx, id)
}

func (c *Client) videoFromID(ctx context.Context, id string) (*Video, error) {
	body, err := c.videoDataByInnertube(ctx, id)
	if err != nil {
		return nil, err
	}

	v := &Video{
		ID: id,
	}

	err = v.parseVideoInfo(body)

	// If the uploader has disabled embedding the video on other sites, parse video page
	if err == ErrNotPlayableInEmbed {
		html, err := c.httpGetBodyBytes(ctx, "https://www.youtube.com/watch?v="+id)
		if err != nil {
			return nil, err
		}

		return v, v.parseVideoPage(html)
	}

	return v, err
}

type innertubeRequest struct {
	VideoID         string            `json:"videoId"`
	Context         inntertubeContext `json:"context"`
	PlaybackContext playbackContext   `json:"playbackContext"`
}

type playbackContext struct {
	ContentPlaybackContext contentPlaybackContext `json:"contentPlaybackContext"`
}

type contentPlaybackContext struct {
	SignatureTimestamp string `json:"signatureTimestamp"`
}

type inntertubeContext struct {
	Client innertubeClient `json:"client"`
}

type innertubeClient struct {
	HL            string `json:"hl"`
	GL            string `json:"gl"`
	ClientName    string `json:"clientName"`
	ClientVersion string `json:"clientVersion"`
}

func (c *Client) videoDataByInnertube(ctx context.Context, id string) ([]byte, error) {
	// fetch sts first
	sts, err := c.getSignatureTimestamp(ctx, id)
	if err != nil {
		return nil, err
	}

	// seems like same token for all WEB clients
	//nolint:gosec
	const webToken = "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
	u := fmt.Sprintf("https://www.youtube.com/youtubei/v1/player?key=%s", webToken)

	data := innertubeRequest{
		VideoID: id,
		Context: inntertubeContext{
			Client: innertubeClient{
				HL:            "en",
				GL:            "US",
				ClientName:    "WEB",
				ClientVersion: "2.20210617.01.00",
			},
		},
		PlaybackContext: playbackContext{
			ContentPlaybackContext: contentPlaybackContext{
				SignatureTimestamp: sts,
			},
		},
	}

	reqData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}

	resp, err := c.httpDo(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return io.ReadAll(resp.Body)
}

// GetPlaylist fetches playlist metadata
func (c *Client) GetPlaylist(url string) (*Playlist, error) {
	return c.GetPlaylistContext(context.Background(), url)
}

// GetPlaylistContext fetches playlist metadata, with a context, along with a list of Videos, and some basic information
// for these videos. Playlist entries cannot be downloaded, as they lack all the required metadata, but
// can be used to enumerate all IDs, Authors, Titles, etc.
func (c *Client) GetPlaylistContext(ctx context.Context, url string) (*Playlist, error) {
	id, err := extractPlaylistID(url)
	if err != nil {
		return nil, fmt.Errorf("extractPlaylistID failed: %w", err)
	}
	requestURL := fmt.Sprintf(playlistFetchURL, id)
	resp, err := c.httpGet(ctx, requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := extractPlaylistJSON(resp.Body)
	if err != nil {
		return nil, err
	}
	p := &Playlist{ID: id}
	return p, json.Unmarshal(data, p)
}

func (c *Client) VideoFromPlaylistEntry(entry *PlaylistEntry) (*Video, error) {
	return c.videoFromID(context.Background(), entry.ID)
}

func (c *Client) VideoFromPlaylistEntryContext(ctx context.Context, entry *PlaylistEntry) (*Video, error) {
	return c.videoFromID(ctx, entry.ID)
}

// GetStream returns the stream and the total size for a specific format
func (c *Client) GetStream(video *Video, format *Format) (io.ReadCloser, int64, error) {
	return c.GetStreamContext(context.Background(), video, format)
}

// GetStream returns the stream and the total size for a specific format with a context.
func (c *Client) GetStreamContext(ctx context.Context, video *Video, format *Format) (io.ReadCloser, int64, error) {
	url, err := c.GetStreamURL(video, format)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	const chunkSize int64 = 10_000_000
	r, w := io.Pipe()

	// Loads a chunk a returns the written bytes.
	// Downloading in multiple chunks is much faster:
	// https://github.com/kkdai/youtube/pull/190
	loadChunk := func(pos int64) (int64, error) {
		req.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", pos, pos+chunkSize-1))

		resp, err := c.httpDo(req)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusPartialContent {
			return 0, ErrUnexpectedStatusCode(resp.StatusCode)
		}

		return io.Copy(w, resp.Body)
	}

	//nolint:revive,errcheck
	go func() {
		// load all the chunks
		for pos := int64(0); pos < format.ContentLength; {
			written, err := loadChunk(pos)
			if err != nil {
				w.CloseWithError(err)
				return
			}

			pos += written
		}

		w.Close()
	}()

	return r, format.ContentLength, nil
}

// GetStreamURL returns the url for a specific format
func (c *Client) GetStreamURL(video *Video, format *Format) (string, error) {
	return c.GetStreamURLContext(context.Background(), video, format)
}

// GetStreamURLContext returns the url for a specific format with a context
func (c *Client) GetStreamURLContext(ctx context.Context, video *Video, format *Format) (string, error) {
	if format.URL != "" {
		return format.URL, nil
	}

	cipher := format.Cipher
	if cipher == "" {
		return "", ErrCipherNotFound
	}

	return c.decipherURL(ctx, video.ID, cipher)
}

// httpDo sends an HTTP request and returns an HTTP response.
func (c *Client) httpDo(req *http.Request) (*http.Response, error) {
	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	if c.Debug {
		log.Println(req.Method, req.URL)
	}

	res, err := client.Do(req)

	if c.Debug && res != nil {
		log.Println(res.Status)
	}

	return res, err
}

// httpGet does a HTTP GET request, checks the response to be a 200 OK and returns it
func (c *Client) httpGet(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpDo(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, ErrUnexpectedStatusCode(resp.StatusCode)
	}
	return resp, nil
}

// httpGetBodyBytes reads the whole HTTP body and returns it
func (c *Client) httpGetBodyBytes(ctx context.Context, url string) ([]byte, error) {
	resp, err := c.httpGet(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
