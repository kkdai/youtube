package youtube

import (
	"context"
	"fmt"
	"io/ioutil"
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
}

// GetVideo fetches video metadata
func (c *Client) GetVideo(url string) (*Video, error) {
	return c.GetVideoContext(context.Background(), url)
}

// GetVideoContext fetches video metadata with a context
func (c *Client) GetVideoContext(ctx context.Context, url string) (*Video, error) {
	id, err := extractVideoID(url)
	if err != nil {
		return nil, fmt.Errorf("extractVideoID failed: %w", err)
	}

	v := &Video{
		ID: id,
	}

	return v, v.LoadInfo(ctx, c)
}

// Fetch playlist metadata
func (c *Client) GetPlaylist(url string) (*Playlist, error) {
	return c.GetPlaylistContext(context.Background(), url)
}

// Fetch playlist metadata, with a context, along with a list of Videos, and some basic information
// for these videos. These videos cannot be downloaded, as they lack all the required metadata, but
// can be used to enumerate all IDs, Authors, Titles, etc.
func (c *Client) GetPlaylistContext(ctx context.Context, url string) (*Playlist, error) {
	id, err := extractPlaylistID(url)
	if err != nil {
		return nil, fmt.Errorf("extractPlaylistID failed: %w", err)
	}

	p := &Playlist{ID: id}
	return p, p.LoadInfo(ctx, c)
}

// GetStream returns the HTTP response for a specific format
func (c *Client) GetStream(video *Video, format *Format) (*http.Response, error) {
	return c.GetStreamContext(context.Background(), video, format)
}

// GetStreamContext returns the HTTP response for a specific format with a context
func (c *Client) GetStreamContext(ctx context.Context, video *Video, format *Format) (*http.Response, error) {
	url, err := c.GetStreamURLContext(ctx, video, format)
	if err != nil {
		return nil, err
	}

	return c.httpGet(ctx, url)
}

// GetStreamURL returns the url for a specific format
func (c *Client) GetStreamURL(video *Video, format *Format) (string, error) {
	return c.GetStreamURLContext(context.Background(), video, format)
}

// GetStreamURL returns the url for a specific format with a context
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

// httpGet does a HTTP GET request, checks the response to be a 200 OK and returns it
func (c *Client) httpGet(ctx context.Context, url string) (resp *http.Response, err error) {
	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	if c.Debug {
		log.Println("GET", url)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, ErrUnexpectedStatusCode(resp.StatusCode)
	}

	return
}

// httpGetBodyBytes reads the whole HTTP body and returns it
func (c *Client) httpGetBodyBytes(ctx context.Context, url string) ([]byte, error) {
	resp, err := c.httpGet(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
