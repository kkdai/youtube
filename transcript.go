package youtube

import (
	"context"
	"encoding/json"
	"encoding/xml"
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
		// if transcript by innertube failed, try to get from timedtext api
		return c.getCaptionTrackContext(ctx, video, lang)
	}

	transcript, err := parseTranscript(body)
	if err != nil {
		// if transcript isn't available, try to get from timedtext api
		if err == ErrTranscriptDisabled {
			return c.getCaptionTrackContext(ctx, video, lang)
		}
		return nil, err
	}

	return transcript, nil
}

func (c *Client) getCaptionTrackContext(ctx context.Context, video *Video, lang string) (VideoTranscript, error) {
	if len(video.CaptionTracks) == 0 {
		return nil, ErrTranscriptDisabled
	}

	captionsURL, err := video.getCaptionTrackURLByLanguage(lang)
	if err != nil {
		return nil, err
	}

	body, err := c.captionTrackDataByInnerTube(ctx, captionsURL)
	if err != nil {
		return nil, err
	}

	transcript, err := parseCaptionTrack(body)
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

func parseCaptionTrack(body []byte) (VideoTranscript, error) {
	var captions TimedText
	err := xml.Unmarshal(body, &captions)
	if err != nil {
		return nil, err
	}

	transcript := make(VideoTranscript, 0, len(captions.Body.Paragraphs))

	// calculate durations to match transcript api
	// if it's the last paragraph, use the duration from the paragraph
	// otherwise, calculate the duration from the next paragraph's start time
	for i, p := range captions.Body.Paragraphs {
		text := ""
		for _, s := range p.Segments {
			text += s.Text
		}
		if text == "" {
			lastTranscriptIdx := len(transcript) - 1
			// if the text is empty, skip this paragraph, but add duration to last segment
			if lastTranscriptIdx >= 0 && i < len(captions.Body.Paragraphs)-1 {
				transcript[lastTranscriptIdx].Duration += captions.Body.Paragraphs[i+1].Time - p.Time
			}
			continue
		}

		var duration int
		if i < len(captions.Body.Paragraphs)-1 {
			duration = captions.Body.Paragraphs[i+1].Time - p.Time
		} else {
			duration = p.Duration
		}

		offsetText := fmt.Sprintf("%d:%02d", p.Time/60000, p.Time/1000%60)

		transcript = append(transcript, TranscriptSegment{
			Text:       text,
			StartMs:    p.Time,
			OffsetText: offsetText,
			Duration:   duration,
		})
	}

	return transcript, nil
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

type TimedText struct {
	Name xml.Name `xml:"timedtext"`
	Body Body     `xml:"body"`
}

type Body struct {
	Paragraphs []Paragraph `xml:"p"`
}

type Paragraph struct {
	Time     int       `xml:"t,attr"`
	Duration int       `xml:"d,attr"`
	Segments []Segment `xml:"s"`
}

type Segment struct {
	Text string `xml:",chardata"`
}
