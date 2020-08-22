package youtube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type LoadStatus uint8

const (
	// Data has not been loaded, and therefore nothing but the ID is available.
	NotLoaded LoadStatus = iota
	// ID, Title, Author, Duration are available, but this video cannot be downloaded as the extra
	// metadata required has not been loaded.
	WeaklyLoaded
	// Data is fully loaded, this  video can be downloaded.
	FullyLoaded

	getVideoInfoURL string = "https://youtube.com/get_video_info?video_id={video}&eurl={eurl}"
	eurlURL         string = "https://youtube.googleapis.com/v/{video}"
)

type Video struct {
	ID       string
	Title    string
	Author   string
	Duration time.Duration
	Formats  FormatList
	Loaded   LoadStatus
}

func (v *Video) parseVideoInfo(info string) error {
	answer, err := url.ParseQuery(info)
	if err != nil {
		return err
	}

	status := answer.Get("status")
	if status != "ok" {
		return &ErrResponseStatus{
			Status: status,
			Reason: answer.Get("reason"),
		}
	}

	// read the streams map
	playerResponse := answer.Get("player_response")
	if playerResponse == "" {
		return errors.New("no player_response found in the server's answer")
	}

	var prData playerResponseData
	if err := json.Unmarshal([]byte(playerResponse), &prData); err != nil {
		return fmt.Errorf("unable to parse player response JSON: %w", err)
	}

	v.Title = prData.VideoDetails.Title
	v.Author = prData.VideoDetails.Author

	if seconds, _ := strconv.Atoi(prData.Microformat.PlayerMicroformatRenderer.LengthSeconds); seconds > 0 {
		v.Duration = time.Duration(seconds) * time.Second
	}

	// Check if video is downloadable
	if prData.PlayabilityStatus.Status != "OK" {
		return &ErrPlayabiltyStatus{
			Status: prData.PlayabilityStatus.Status,
			Reason: prData.PlayabilityStatus.Reason,
		}
	}

	// Assign Streams
	v.Formats = append(prData.StreamingData.Formats, prData.StreamingData.AdaptiveFormats...)

	if len(v.Formats) == 0 {
		return errors.New("no formats found in the server's answer")
	}

	return nil
}

func (v *Video) FetchVideoInfo(ctx context.Context, c *Client) ([]byte, error) {
	// Circumvent age restriction to pretend access through googleapis.com
	url := strings.Replace(getVideoInfoURL, "{eurl}", eurlURL, 1)
	url = strings.Replace(url, "{video}", v.ID, -1)

	resp, err := c.httpGet(ctx, url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (v *Video) LoadInfo(ctx context.Context, c *Client) error {
	body, err := v.FetchVideoInfo(ctx, c)

	if err != nil {
		return err
	}

	err = v.parseVideoInfo(string(body))

	if err != nil {
		v.Loaded = NotLoaded
	} else {
		v.Loaded = FullyLoaded
	}

	return err
}
