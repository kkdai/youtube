package youtube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
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

	getVideoInfoURL string = "https://youtube.com/get_video_info?video_id=%s&eurl=https://youtube.googleapis.com/v/%s"
)

type Video struct {
	ID              string
	Title           string
	Description     string
	Author          string
	Duration        time.Duration
	Formats         FormatList
	DASHManifestURL string // URI of the DASH manifest file
	HLSManifestURL  string // URI of the HLS manifest file
}

func (v *Video) parseVideoInfo(body []byte) error {
	answer, err := url.ParseQuery(string(body))
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

	if err := v.isVideoFromInfoDownloadable(prData); err != nil {
		return err
	}

	return v.extractDataFromPlayerResponse(prData)
}

func (v *Video) isVideoFromInfoDownloadable(prData playerResponseData) error {
	return v.isVideoDownloadable(prData, false)
}

var playerResponsePattern = regexp.MustCompile(`var ytInitialPlayerResponse\s*=\s*(\{.+?\});`)

func (v *Video) parseVideoPage(body []byte) error {
	initialPlayerResponse := playerResponsePattern.FindSubmatch(body)
	if initialPlayerResponse == nil || len(initialPlayerResponse) < 2 {
		return errors.New("no ytInitialPlayerResponse found in the server's answer")
	}

	var prData playerResponseData
	if err := json.Unmarshal(initialPlayerResponse[1], &prData); err != nil {
		return fmt.Errorf("unable to parse player response JSON: %w", err)
	}

	if err := v.isVideoFromPageDownloadable(prData); err != nil {
		return err
	}

	return v.extractDataFromPlayerResponse(prData)
}

func (v *Video) isVideoFromPageDownloadable(prData playerResponseData) error {
	return v.isVideoDownloadable(prData, true)
}

func (v *Video) isVideoDownloadable(prData playerResponseData, isVideoPage bool) error {
	// Check if video is downloadable
	if prData.PlayabilityStatus.Status == "OK" {
		return nil
	}

	if !isVideoPage && !prData.PlayabilityStatus.PlayableInEmbed {
		return ErrNotPlayableInEmbed
	}

	return &ErrPlayabiltyStatus{
		Status: prData.PlayabilityStatus.Status,
		Reason: prData.PlayabilityStatus.Reason,
	}
}

func (v *Video) extractDataFromPlayerResponse(prData playerResponseData) error {
	v.Title = prData.VideoDetails.Title
	v.Description = prData.VideoDetails.ShortDescription
	v.Author = prData.VideoDetails.Author

	if seconds, _ := strconv.Atoi(prData.Microformat.PlayerMicroformatRenderer.LengthSeconds); seconds > 0 {
		v.Duration = time.Duration(seconds) * time.Second
	}

	// Assign Streams
	v.Formats = append(prData.StreamingData.Formats, prData.StreamingData.AdaptiveFormats...)

	if len(v.Formats) == 0 {
		return errors.New("no formats found in the server's answer")
	}

	v.HLSManifestURL = prData.StreamingData.HlsManifestURL
	v.DASHManifestURL = prData.StreamingData.DashManifestURL

	return nil
}

func (v *Video) FetchVideoInfo(ctx context.Context, c *Client) ([]byte, error) {
	// Circumvent age restriction to pretend access through googleapis.com
	url := fmt.Sprintf(getVideoInfoURL, v.ID, v.ID)
	return c.httpGetBodyBytes(ctx, url)
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
