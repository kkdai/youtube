package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"
)

type Video struct {
	ID       string
	Streams  []Stream
	Title    string
	Author   string
	Duration time.Duration
}

type Stream struct {
	ItagNo   int    `json:"itag"`
	URL      string `json:"url"`
	MimeType string `json:"mimeType"`
	Quality  string `json:"quality"`
	Cipher   string `json:"signatureCipher"`
}

type Format struct {
	Stream
	Bitrate          int    `json:"bitrate"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	LastModified     string `json:"lastModified"`
	ContentLength    string `json:"contentLength,omitempty"`
	QualityLabel     string `json:"qualityLabel"`
	ProjectionType   string `json:"projectionType"`
	AverageBitrate   int    `json:"averageBitrate,omitempty"`
	AudioQuality     string `json:"audioQuality"`
	ApproxDurationMs string `json:"approxDurationMs"`
	AudioSampleRate  string `json:"audioSampleRate"`
	AudioChannels    int    `json:"audioChannels"`
}

type AdaptiveFormat struct {
	Stream
	Bitrate   int `json:"bitrate"`
	Width     int `json:"width,omitempty"`
	Height    int `json:"height,omitempty"`
	InitRange struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"initRange"`
	IndexRange struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"indexRange"`
	LastModified     string `json:"lastModified"`
	ContentLength    string `json:"contentLength"`
	Fps              int    `json:"fps,omitempty"`
	QualityLabel     string `json:"qualityLabel,omitempty"`
	ProjectionType   string `json:"projectionType"`
	AverageBitrate   int    `json:"averageBitrate"`
	ApproxDurationMs string `json:"approxDurationMs"`
	ColorInfo        struct {
		Primaries               string `json:"primaries"`
		TransferCharacteristics string `json:"transferCharacteristics"`
		MatrixCoefficients      string `json:"matrixCoefficients"`
	} `json:"colorInfo,omitempty"`
	HighReplication bool   `json:"highReplication,omitempty"`
	AudioQuality    string `json:"audioQuality,omitempty"`
	AudioSampleRate string `json:"audioSampleRate,omitempty"`
	AudioChannels   int    `json:"audioChannels,omitempty"`
}

func (v *Video) FindStreamByQuality(quality string) *Stream {
	for i := range v.Streams {
		if v.Streams[i].Quality == quality {
			return &v.Streams[i]
		}
	}

	return nil
}

func (v *Video) FindStreamByItag(itagNo int) *Stream {
	for i := range v.Streams {
		if v.Streams[i].ItagNo == itagNo {
			return &v.Streams[i]
		}
	}
	return nil
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

	var prData PlayerResponseData
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

	// Get video download link
	streams, err := parseStreams(prData)
	if err != nil {
		return err
	}

	v.Streams = streams
	if len(v.Streams) == 0 {
		return errors.New("no Stream list found in the server's answer")
	}

	return nil
}

func parseStreams(prData PlayerResponseData) ([]Stream, error) {
	size := len(prData.StreamingData.Formats) + len(prData.StreamingData.AdaptiveFormats)
	streams := make([]Stream, 0, size)

	filterFormat := func(stream Stream) {
		if stream.MimeType == "" {
			// FIXME logging
			log.Printf("An error occurred while decoding one of the video's Stream's information: Stream %+v", stream)
			return
		}
		streams = append(streams, stream)
	}

	for _, format := range prData.StreamingData.Formats {
		filterFormat(format.Stream)
	}
	for _, format := range prData.StreamingData.AdaptiveFormats {
		filterFormat(format.Stream)
	}
	return streams, nil
}
