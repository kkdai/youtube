package youtube

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranscript(t *testing.T) {
	video := &Video{ID: "9_MbW9FK1fA"}

	transcript, err := testClient.GetTranscript(video)
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
