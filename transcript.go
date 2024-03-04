package youtube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
	var str string
	for _, tr := range vt {
		str += tr.String() + "\n"
	}

	return str
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
// Not all videos have transcripts, only relatively new videos.
// If transcripts are disabled or not available, ErrTranscriptDisabled is returned.
func (c *Client) GetTranscriptCtx(ctx context.Context, video *Video, lang string) (VideoTranscript, error) {
	c.assureClient()

	if video == nil || video.ID == "" {
		return nil, fmt.Errorf("no video provided")
	}

	body, err := c.transcriptDataByInnertube(ctx, video.ID, lang)
	if err != nil {
		return nil, err
	}

	transcript, err := parseTranscript(body)
	if err != nil {
		return nil, err
	}

	return transcript, nil
}

func parseTranscript(body []byte) (VideoTranscript, error) {
	var resp transcriptResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if len(resp.Actions) > 0 {
		// Android client response
		if app := resp.Actions[0].AppSegment; app != nil {
			return getSegments(app)
		}

		// Web client response
		if web := resp.Actions[0].WebSegment; web != nil {
			return nil, fmt.Errorf("not implemented")
		}
	}

	return nil, ErrTranscriptDisabled
}

type segmenter interface {
	ParseSegments() []TranscriptSegment
}

func getSegments(f segmenter) (VideoTranscript, error) {
	if segments := f.ParseSegments(); len(segments) > 0 {
		return segments, nil
	}

	return nil, ErrTranscriptDisabled
}

// transcriptResp is the JSON structure as returned by the transcript API.
type transcriptResp struct {
	Actions []struct {
		AppSegment *appData `json:"elementsCommand"`
		WebSegment *webData `json:"updateEngagementPanelAction"`
	} `json:"actions"`
}

type appData struct {
	TEC struct {
		Args struct {
			ListArgs struct {
				Ow struct {
					InitialSeg []struct {
						TranscriptSegment struct {
							StartMs string `json:"startMs"`
							EndMs   string `json:"endMs"`
							Text    struct {
								String struct {
									// Content is the actual transctipt text
									Content string `json:"content"`
								} `json:"elementsAttributedString"`
							} `json:"snippet"`
							StartTimeText struct {
								String struct {
									// Content is the fomratted timestamp, e.g. '4:00'
									Content string `json:"content"`
								} `json:"elementsAttributedString"`
							} `json:"startTimeText"`
						} `json:"transcriptSegmentRenderer"`
					} `json:"initialSegments"`
				} `json:"overwrite"`
			} `json:"transformTranscriptSegmentListArguments"`
		} `json:"arguments"`
	} `json:"transformEntityCommand"`
}

func (s *appData) ParseSegments() []TranscriptSegment {
	rawSegments := s.TEC.Args.ListArgs.Ow.InitialSeg
	segments := make([]TranscriptSegment, 0, len(rawSegments))

	for _, segment := range rawSegments {
		startMs, _ := strconv.Atoi(segment.TranscriptSegment.StartMs)
		endMs, _ := strconv.Atoi(segment.TranscriptSegment.EndMs)

		segments = append(segments, TranscriptSegment{
			Text:       segment.TranscriptSegment.Text.String.Content,
			StartMs:    startMs,
			OffsetText: segment.TranscriptSegment.StartTimeText.String.Content,
			Duration:   endMs - startMs,
		})
	}

	return segments
}

type webData struct {
	Content struct {
		TR struct {
			Body struct {
				TBR struct {
					Cues []struct {
						Transcript struct {
							FormattedStartOffset struct {
								SimpleText string `json:"simpleText"`
							} `json:"formattedStartOffset"`
							Cues []struct {
								TranscriptCueRenderer struct {
									Cue struct {
										SimpleText string `json:"simpleText"`
									} `json:"cue"`
									StartOffsetMs string `json:"startOffsetMs"`
									DurationMs    string `json:"durationMs"`
								} `json:"transcriptCueRenderer"`
							} `json:"cues"`
						} `json:"transcriptCueGroupRenderer"`
					} `json:"cueGroups"`
				} `json:"transcriptSearchPanelRenderer"`
			} `json:"content"`
		} `json:"transcriptRenderer"`
	} `json:"content"`
}

func (s *webData) ParseSegments() []TranscriptSegment {
	// TODO: doesn't actually work now, check json.
	cues := s.Content.TR.Body.TBR.Cues
	segments := make([]TranscriptSegment, 0, len(cues))

	for _, s := range cues {
		formatted := s.Transcript.FormattedStartOffset.SimpleText
		segment := s.Transcript.Cues[0].TranscriptCueRenderer
		start, _ := strconv.Atoi(segment.StartOffsetMs)
		duration, _ := strconv.Atoi(segment.DurationMs)

		segments = append(segments, TranscriptSegment{
			Text:       segment.Cue.SimpleText,
			StartMs:    start,
			OffsetText: formatted,
			Duration:   duration,
		})
	}

	return segments
}
