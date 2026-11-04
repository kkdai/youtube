package youtube

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	ErrTranscriptDisabled = errors.New("transcript is disabled on this video")
)

// TranscriptSegment is a single transcipt segment spanning a few milliseconds.
type TranscriptSegment struct {
	// Text is the transcipt text.
	Text string `json:"text"`

	// StartMs is the start timestamp in ms.
	StartMs int `json:"offset"`

	// OffsetText e.g. '4:00'.
	OffsetText string `json:"offsetText"`

	// Duration the transcript segment spans in ms.
	Duration int `json:"duration"`
}

func (tr TranscriptSegment) String() string {
	return tr.OffsetText + " - " + strings.TrimSpace(tr.Text)
}

type VideoTranscript []TranscriptSegment

func (vt VideoTranscript) String() string {
	var str strings.Builder
	for _, tr := range vt {
		str.WriteString(tr.String() + "\n")
	}

	return str.String()
}

// GetTranscript fetches the video transcript if available.
//
// Not all videos have transcripts, only relatively new videos.
// If transcripts are disabled or not available, ErrTranscriptDisabled is returned.
func (c *Client) GetTranscript(video *Video, lang string) (VideoTranscript, error) {
	return c.GetTranscriptCtx(context.Background(), video, lang)
}

// GetTranscriptCtx fetches the video transcript if available.
//
// The transcript is read from the video's caption track in the given language,
// preferring the automatically generated one. The caption tracks of the video
// are used if present, otherwise they are fetched.
//
// Not all videos have transcripts, only relatively new videos.
// If transcripts are disabled or not available, ErrTranscriptDisabled is returned.
func (c *Client) GetTranscriptCtx(ctx context.Context, video *Video, lang string) (VideoTranscript, error) {
	c.assureClient()

	if video == nil || video.ID == "" {
		return nil, fmt.Errorf("no video provided")
	}

	tracks := video.CaptionTracks
	if len(tracks) == 0 {
		var err error
		if tracks, err = c.captionTracks(ctx, video.ID); err != nil {
			return nil, err
		}
	}

	track := findCaptionTrack(tracks, lang)
	if track == nil {
		return nil, ErrTranscriptDisabled
	}

	body, err := c.httpGetBodyBytes(ctx, srv3URL(track.BaseURL))
	if err != nil {
		return nil, err
	}

	return parseTranscript(body)
}

// captionTracks fetches the caption tracks of a video from the player API.
func (c *Client) captionTracks(ctx context.Context, id string) ([]CaptionTrack, error) {
	body, err := c.videoDataByInnertube(ctx, id)
	if err != nil {
		return nil, err
	}

	var prData playerResponseData
	if err := json.Unmarshal(body, &prData); err != nil {
		return nil, err
	}

	// e.g. ErrBotCheck, so the caller doesn't assume missing captions
	if err := (&Video{}).isVideoDownloadable(prData, false); err != nil {
		return nil, err
	}

	return prData.Captions.PlayerCaptionsTracklistRenderer.CaptionTracks, nil
}

// findCaptionTrack returns the caption track for the language, preferring the
// automatically generated one ("asr").
func findCaptionTrack(tracks []CaptionTrack, lang string) *CaptionTrack {
	var found *CaptionTrack

	for i := range tracks {
		if tracks[i].LanguageCode != lang {
			continue
		}
		if tracks[i].Kind == "asr" {
			return &tracks[i]
		}
		if found == nil {
			found = &tracks[i]
		}
	}

	return found
}

// srv3URL requests the caption track in the srv3 format, which parseTranscript expects.
func srv3URL(baseURL string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	q := u.Query()
	q.Set("fmt", "srv3")
	u.RawQuery = q.Encode()

	return u.String()
}

// timedText is the srv3 XML format of a caption track.
type timedText struct {
	Paragraphs []struct {
		StartMs    int    `xml:"t,attr"`
		DurationMs int    `xml:"d,attr"`
		Text       string `xml:",chardata"`
		// Words of automatically generated captions
		Words []struct {
			Text string `xml:",chardata"`
		} `xml:"s"`
	} `xml:"body>p"`
}

func parseTranscript(body []byte) (VideoTranscript, error) {
	var tt timedText
	if err := xml.Unmarshal(body, &tt); err != nil {
		return nil, err
	}

	segments := make(VideoTranscript, 0, len(tt.Paragraphs))

	for _, p := range tt.Paragraphs {
		var text strings.Builder
		text.WriteString(p.Text)
		for _, w := range p.Words {
			text.WriteString(w.Text)
		}

		// skip line breaks of automatically generated captions
		content := strings.TrimSpace(text.String())
		if content == "" {
			continue
		}

		segments = append(segments, TranscriptSegment{
			Text:       content,
			StartMs:    p.StartMs,
			OffsetText: formatOffset(p.StartMs),
			Duration:   p.DurationMs,
		})
	}

	if len(segments) == 0 {
		return nil, ErrTranscriptDisabled
	}

	return segments, nil
}

// formatOffset formats a timestamp like YouTube does, e.g. "4:05" or "1:02:03".
func formatOffset(ms int) string {
	seconds := ms / 1000
	hours, minutes, seconds := seconds/3600, seconds%3600/60, seconds%60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}

	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
