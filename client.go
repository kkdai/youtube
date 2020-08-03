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

	// Circumvent age restriction to pretend access through googleapis.com
	eurl := "https://youtube.googleapis.com/v/" + id
	resp, err := c.httpGet(ctx, "https://youtube.com/get_video_info?video_id="+id+"&eurl="+eurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	v := &Video{
		ID: id,
	}

	return v, v.parseVideoInfo(string(body))
}

// GetStream returns the HTTP response for a specific stream
func (c *Client) GetStream(video *Video, stream *Stream) (*http.Response, error) {
	return c.GetStreamContext(context.Background(), video, stream)
}

// GetStreamContext returns the HTTP response for a specific stream with a context
func (c *Client) GetStreamContext(ctx context.Context, video *Video, stream *Stream) (*http.Response, error) {
	url, err := c.getStreamUrl(ctx, video, stream)
	if err != nil {
		return nil, err
	}

	return c.httpGet(ctx, url)
}

func (c *Client) getStreamUrl(ctx context.Context, video *Video, stream *Stream) (string, error) {
	if stream.URL != "" {
		return stream.URL, nil
	}

	cipher := stream.Cipher
	if cipher == "" {
		return "", ErrCipherNotFound
	}

	return c.decipherURL(ctx, video.ID, cipher)
}

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
