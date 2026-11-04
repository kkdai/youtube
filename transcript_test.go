package youtube

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranscript(t *testing.T) {
	video := &Video{ID: "9_MbW9FK1fA"}

	transcript, err := testClient.GetTranscript(video, "en")
	skipIfBlocked(t, err)
	require.NoError(t, err, "get transcript")
	require.Greater(t, len(transcript), 0, "no transcript segments found")

	for i, segment := range transcript {
		index := strconv.Itoa(i)

		require.NotEmpty(t, segment.Text, "text "+index)
		require.NotEmpty(t, segment.Duration, "duration "+index)
		require.NotEmpty(t, segment.OffsetText, "offset "+index)

		if i != 0 {
			require.NotEmpty(t, segment.StartMs, "startMs "+index)
		}
	}

	t.Log(transcript.String())
}

func TestTranscriptOtherLanguage(t *testing.T) {
	video := &Video{ID: "AXwDvYh2-uk"}

	transcript, err := testClient.GetTranscript(video, "id")
	skipIfBlocked(t, err)
	require.NoError(t, err, "get transcript")
	require.Greater(t, len(transcript), 0, "no transcript segments found")

	for i, segment := range transcript {
		index := strconv.Itoa(i)

		require.NotEmpty(t, segment.Text, "text "+index)
		require.NotEmpty(t, segment.Duration, "duration "+index)
		require.NotEmpty(t, segment.OffsetText, "offset "+index)

		if i != 0 {
			require.NotEmpty(t, segment.StartMs, "startMs "+index)
		}
	}

	t.Log(transcript.String())
}

func TestParseTranscript(t *testing.T) {
	// srv3 format: automatically generated captions consist of words (<s>)
	// and empty paragraphs for line breaks, manual ones contain the text directly
	body := []byte(`<?xml version="1.0" encoding="utf-8" ?><timedtext format="3">
<body>
<p t="0" d="4280" w="1"><s>So,</s><s t="120"> by</s><s t="280"> now</s></p>
<p t="1590" d="2690" w="1" a="1">
</p>
<p t="65000" d="1500">I&#39;m manual</p>
<p t="3723000" d="2000">an hour later</p>
</body>
</timedtext>`)

	transcript, err := parseTranscript(body)
	require.NoError(t, err)
	require.Equal(t, VideoTranscript{
		{Text: "So, by now", StartMs: 0, OffsetText: "0:00", Duration: 4280},
		{Text: "I'm manual", StartMs: 65000, OffsetText: "1:05", Duration: 1500},
		{Text: "an hour later", StartMs: 3723000, OffsetText: "1:02:03", Duration: 2000},
	}, transcript)
}

func TestParseTranscript_Empty(t *testing.T) {
	_, err := parseTranscript([]byte(`<timedtext format="3"><body></body></timedtext>`))
	require.ErrorIs(t, err, ErrTranscriptDisabled)
}

func TestFindCaptionTrack(t *testing.T) {
	tracks := []CaptionTrack{
		{LanguageCode: "de", BaseURL: "de"},
		{LanguageCode: "en", BaseURL: "en"},
		{LanguageCode: "en", Kind: "asr", BaseURL: "en-asr"},
	}

	require.Equal(t, "en-asr", findCaptionTrack(tracks, "en").BaseURL)
	require.Equal(t, "de", findCaptionTrack(tracks, "de").BaseURL)
	require.Nil(t, findCaptionTrack(tracks, "fr"))
}

func TestSrv3URL(t *testing.T) {
	require.Equal(t,
		"https://www.youtube.com/api/timedtext?fmt=srv3&lang=en&v=abc",
		srv3URL("https://www.youtube.com/api/timedtext?v=abc&lang=en&fmt=json3"),
	)
}
