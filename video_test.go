package youtube

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownload_Regular(t *testing.T) {
	ctx := context.Background()

	testcases := []struct {
		name       string
		url        string
		outputFile string
		itagNo     int
		quality    string
	}{
		{
			// Video from issue #25
			name:       "default",
			url:        "https://www.youtube.com/watch?v=54e6lBE3BoQ",
			outputFile: "default_test.mp4",
			quality:    "",
		},
		{
			// Video from issue #25
			name:       "quality:medium",
			url:        "https://www.youtube.com/watch?v=54e6lBE3BoQ",
			outputFile: "medium_test.mp4",
			quality:    "medium",
		},
		{
			name: "without-filename",
			url:  "https://www.youtube.com/watch?v=n3kPvBCYT3E",
		},
		{
			name:       "Format",
			url:        "https://www.youtube.com/watch?v=54e6lBE3BoQ",
			outputFile: "muxedstream_test.mp4",
			itagNo:     18,
		},
		{
			name:       "AdaptiveFormat_video",
			url:        "https://www.youtube.com/watch?v=54e6lBE3BoQ",
			outputFile: "adaptiveStream_video_test.m4v",
			itagNo:     134,
		},
		{
			name:       "AdaptiveFormat_audio",
			url:        "https://www.youtube.com/watch?v=54e6lBE3BoQ",
			outputFile: "adaptiveStream_audio_test.m4a",
			itagNo:     140,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			video, err := testClient.GetVideoContext(ctx, tc.url)
			require.NoError(err)

			var stream *Stream
			if tc.itagNo > 0 {
				stream = video.FindStreamByItag(tc.itagNo)
				require.NotNil(stream)
			} else {
				stream = &video.Streams[0]
			}

			url, err := testClient.getStreamURL(ctx, video, stream)
			require.NoError(err)
			require.NotEmpty(url)
		})
	}
}

func TestDownload_WhenPlayabilityStatusIsNotOK(t *testing.T) {
	testcases := []struct {
		issue   string
		videoID string
		err     string
	}{
		{
			issue:   "issue#65",
			videoID: "9ja-K2FslBU",
			err:     `status: ERROR`,
		},
		{
			issue:   "issue#59",
			videoID: "nINQjT7Zr9w",
			err:     `status: LOGIN_REQUIRED`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.issue, func(t *testing.T) {
			_, err := testClient.GetVideo(tc.videoID)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.err)
		})
	}
}
