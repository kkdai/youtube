package youtube

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func ExampleClient_GetStream() {
	video, err := testClient.GetVideo("https://www.youtube.com/watch?v=9_MbW9FK1fA")
	if err != nil {
		panic(err)
	}

	// Typically youtube only provides separate streams for video and audio.
	// If you want audio and video combined, take a look a the downloader package.
	formats := video.Formats.Quality("medium")
	reader, _, err := testClient.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}

	// do something with the reader

	reader.Close()
}

func TestExtractVideoId(t *testing.T) {
	testcases := []struct {
		input      string
		expectedID string
	}{
		{
			input:      "https://www.youtube.com/watch?v=9_MbW9FK1fA",
			expectedID: "9_MbW9FK1fA",
		},
		{
			input:      "https://www.youtube.com/watch?v=-D2IEZbn5Xs&list=PLQUru4nFApg_ocT-XYXFg50_l4V2BRtwi&index=15",
			expectedID: "-D2IEZbn5Xs",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			id, err := ExtractVideoID(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expectedID, id)
		})
	}

}

func TestSimpleTest(t *testing.T) {
	video, err := testClient.GetVideo("https://www.youtube.com/watch?v=9_MbW9FK1fA")
	require.NoError(t, err, "get body")

	_, err = testClient.GetTranscript(video, "en")
	require.NoError(t, err, "get transcript")

	// Typically youtube only provides separate streams for video and audio.
	// If you want audio and video combined, take a look a the downloader package.
	formats := video.Formats.Quality("hd1080")
	require.NotEmpty(t, formats)

	start := time.Now()
	reader, _, err := testClient.GetStream(video, &formats[0])
	require.NoError(t, err, "get stream")

	t.Log("Duration Milliseconds: ", time.Since(start).Milliseconds())

	// do something with the reader
	b, err := io.ReadAll(reader)
	require.NoError(t, err, "read body")

	t.Log("Downloaded ", len(b))
}

func TestDownload_Regular(t *testing.T) {
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
		{
			// Video from issue #138
			name:       "NotPlayableInEmbed",
			url:        "https://www.youtube.com/watch?v=gr-IqFcNExY",
			outputFile: "not_playable_in_embed.mp4",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			video, err := testClient.GetVideo(tc.url)
			require.NoError(err)

			formats := video.Formats
			if tc.itagNo > 0 {
				formats = formats.Itag(tc.itagNo)
				require.NotEmpty(formats)
			}

			url, err := testClient.GetStreamURL(video, &video.Formats[0])
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
			videoID: "yZIXLfi8CZQ",
			err:     ErrVideoPrivate.Error(),
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

// See https://github.com/kkdai/youtube/pull/238
func TestDownload_SensitiveContent(t *testing.T) {
	_, err := testClient.GetVideo("MS91knuzoOA")
	require.EqualError(t, err, "can't bypass age restriction: embedding of this video has been disabled")
}

func TestParse_PublishDate(t *testing.T) {
	testcases := []struct {
		videoID     string
		publishDate time.Time
	}{
		{"eRqCe_VHs6M", time.Date(2020, time.July, 11, 13, 29, 21, 0, time.UTC)},
	}

	for _, tc := range testcases {
		t.Run(tc.videoID, func(t *testing.T) {
			got, err := testClient.GetVideo(tc.videoID)
			if err != nil {
				t.Fatalf("video parsing failed")
			}
			require.Equal(t, tc.publishDate, got.PublishDate)
		})
	}
}
