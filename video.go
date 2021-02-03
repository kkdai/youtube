package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"time"
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
	return v.parseVideoInfoOrPage(body, false)
}

func (v *Video) parseVideoPage(body []byte) error {
	return v.parseVideoInfoOrPage(body, true)
}

func (v *Video) parseVideoInfoOrPage(body []byte, isVideoPage bool) error {
	var prData playerResponseData

	if !isVideoPage {
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

		if err := json.Unmarshal([]byte(playerResponse), &prData); err != nil {
			return fmt.Errorf("unable to parse player response JSON: %w", err)
		}
	} else {
		playerResponsePattern := regexp.MustCompile(`var ytInitialPlayerResponse\s*=\s*(\{.+?\});`)

		initialPlayerResponse := playerResponsePattern.FindSubmatch(body)
		if initialPlayerResponse == nil || len(initialPlayerResponse) < 2 {
			return errors.New("no ytInitialPlayerResponse found in the server's answer")
		}

		if err := json.Unmarshal(initialPlayerResponse[1], &prData); err != nil {
			return fmt.Errorf("unable to parse player response JSON: %w", err)
		}
	}

	v.Title = prData.VideoDetails.Title
	v.Description = prData.VideoDetails.ShortDescription
	v.Author = prData.VideoDetails.Author

	if seconds, _ := strconv.Atoi(prData.Microformat.PlayerMicroformatRenderer.LengthSeconds); seconds > 0 {
		v.Duration = time.Duration(seconds) * time.Second
	}

	// Check if video is downloadable
	if prData.PlayabilityStatus.Status != "OK" {
		if !isVideoPage && !prData.PlayabilityStatus.PlayableInEmbed {
			return ErrNotPlayableInEmbed
		}

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

	v.HLSManifestURL = prData.StreamingData.HlsManifestURL
	v.DASHManifestURL = prData.StreamingData.DashManifestURL

	return nil
}
