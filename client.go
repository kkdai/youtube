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

	// playerCache caches the JavaScript code of a player response
	playerCache playerCache
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
	body, err := c.videoDataByInnertube(ctx, id, Web)
	if err != nil {
		return nil, err
	}

	v := &Video{
		ID: id,
	}

	err = v.parseVideoInfo(body)
	// return early if all good
	if err == nil {
		return v, nil
	}

	// If the uploader has disabled embedding the video on other sites, parse video page
	if err == ErrNotPlayableInEmbed {
		// additional parameters are required to access clips with sensitiv content
		html, err := c.httpGetBodyBytes(ctx, "https://www.youtube.com/watch?v="+id+"&bpctr=9999999999&has_verified=1")
		if err != nil {
			return nil, err
		}

		return v, v.parseVideoPage(html)
	}

	// If the uploader marked the video as inappropriate for some ages, use embed player
	if err == ErrLoginRequired {
		bodyEmbed, errEmbed := c.videoDataByInnertube(ctx, id, EmbeddedClient)
		if errEmbed == nil {
			errEmbed = v.parseVideoInfo(bodyEmbed)
		}

		if errEmbed == nil {
			return v, nil
		}

		// private video clearly not age-restricted and thus should be explicit
		if errEmbed == ErrVideoPrivate {
			return v, errEmbed
		}

		// wrapping error so its clear whats happened
		return v, fmt.Errorf("can't bypass age restriction: %w", errEmbed)
	}

	// undefined error
	return v, err
}

type innertubeRequest struct {
	VideoID         string            `json:"videoId,omitempty"`
	BrowseID        string            `json:"browseId,omitempty"`
	Continuation    string            `json:"continuation,omitempty"`
	Context         inntertubeContext `json:"context"`
	PlaybackContext playbackContext   `json:"playbackContext,omitempty"`
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

type ClientType string

const (
	Web            ClientType = "WEB"
	EmbeddedClient ClientType = "WEB_EMBEDDED_PLAYER"
)

func (c *Client) videoDataByInnertube(ctx context.Context, id string, clientType ClientType) ([]byte, error) {
	config, err := c.getPlayerConfig(ctx, id)
	if err != nil {
		return nil, err
	}

	// fetch sts first
	sts, err := config.getSignatureTimestamp()
	if err != nil {
		return nil, err
	}

	data, keyToken := prepareInnertubeVideoData(id, sts, clientType)
	return c.httpPostBodyBytes(ctx, "https://www.youtube.com/youtubei/v1/player?key="+keyToken, data)
}

var innertubeClientInfo = map[ClientType]map[string]string{
	// might add ANDROID and other in future, but i don't see reason yet
	Web: {
		"version": "2.20210617.01.00",
		"key":     "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8",
	},
	EmbeddedClient: {
		"version": "1.19700101",
		// seems like same key works for both clients
		"key": "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8",
	},
}

func prepareInnertubeContext(clientType ClientType) (inntertubeContext, string) {
	cInfo, ok := innertubeClientInfo[clientType]
	if !ok {
		// if provided clientType not exist - use Web as fallback option
		clientType = Web
		cInfo = innertubeClientInfo[clientType]
	}

	return inntertubeContext{
		Client: innertubeClient{
			HL:            "en",
			GL:            "US",
			ClientName:    string(clientType),
			ClientVersion: cInfo["version"],
		},
	}, cInfo["key"]
}

func prepareInnertubeVideoData(videoID string, sts string, clientType ClientType) (innertubeRequest, string) {
	context, key := prepareInnertubeContext(clientType)

	return innertubeRequest{
		VideoID: videoID,
		Context: context,
		PlaybackContext: playbackContext{
			ContentPlaybackContext: contentPlaybackContext{
				SignatureTimestamp: sts,
			},
		},
	}, key
}

func prepareInnertubePlaylistData(ID string, continuation bool, clientType ClientType) (innertubeRequest, string) {
	context, key := prepareInnertubeContext(clientType)

	if continuation {
		return innertubeRequest{Context: context, Continuation: ID}, key
	}

	return innertubeRequest{Context: context, BrowseID: "VL" + ID}, key
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

	data, key := prepareInnertubePlaylistData(id, false, Web)
	body, err := c.httpPostBodyBytes(ctx, "https://www.youtube.com/youtubei/v1/browse?key="+key, data)
	if err != nil {
		return nil, err
	}

	p := &Playlist{ID: id}
	return p, p.parsePlaylistInfo(ctx, c, body)
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

// GetStreamContext returns the stream and the total size for a specific format with a context.
func (c *Client) GetStreamContext(ctx context.Context, video *Video, format *Format) (io.ReadCloser, int64, error) {
	url, err := c.GetStreamURL(video, format)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	r, w := io.Pipe()

	go c.download(req, w, format)

	return r, format.ContentLength, nil
}

func (c *Client) download(req *http.Request, w *io.PipeWriter, format *Format) {
	const chunkSize int64 = 10_000_000
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

	defer w.Close()
	//nolint:revive,errcheck
	if format.ContentLength == 0 {
		resp, err := c.httpDo(req)
		if err != nil {
			w.CloseWithError(err)
			return
		}

		defer resp.Body.Close()

		io.Copy(w, resp.Body)
		return
	}

	//nolint:revive,errcheck
	// load all the chunks
	for pos := int64(0); pos < format.ContentLength; {
		written, err := loadChunk(pos)
		if err != nil {
			w.CloseWithError(err)
			return
		}

		pos += written
	}
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

	uri, err := c.decipherURL(ctx, video.ID, cipher)
	if err != nil {
		return "", err
	}

	return uri, err
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

// httpPost does a HTTP POST request with a body, checks the response to be a 200 OK and returns it
func (c *Client) httpPost(ctx context.Context, url string, body interface{}) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
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

// httpPostBodyBytes reads the whole HTTP body and returns it
func (c *Client) httpPostBodyBytes(ctx context.Context, url string, body interface{}) ([]byte, error) {
	resp, err := c.httpPost(ctx, url, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
