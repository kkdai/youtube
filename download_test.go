// +build integration

package youtube

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownload_Regular(t *testing.T) {
	curDir, _ := os.Getwd()
	outputDir := filepath.Join(curDir, downloadToDir)
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
			itagNo:     0,
			quality:    "",
		},
		{
			// Video from issue #25
			name:       "quality:medium",
			url:        "https://www.youtube.com/watch?v=54e6lBE3BoQ",
			outputFile: "medium_test.mp4",
			itagNo:     0,
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
		fmt.Println("download to " + outputDir + "\\" + tc.outputFile)
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			y := NewYoutube(true, false)
			require.NoError(y.DecodeURL(tc.url))
			require.NoError(y.StartDownload(outputDir, tc.outputFile, tc.quality, tc.itagNo))
		})
	}
}

func TestDownload_HighQuality(t *testing.T) {
	require := require.New(t)
	y := NewYoutube(false, false)

	// url from issue #21
	testVideoId := "n3kPvBCYT3E"
	require.NoError(y.DecodeURL(testVideoId))

	curDir, _ := os.Getwd()
	outputDir := filepath.Join(curDir, downloadToDir)
	fmt.Println("download to " + outputDir + "\\" + "Silhouette Eurobeat Remix")
	require.NoError(y.StartDownloadWithHighQuality(outputDir, "", "hd1080"))
}

func TestDownload_WhenPlayabilityStatusIsNotOK(t *testing.T) {
	y := NewYoutube(false, false)

	testcases := []struct {
		issue   string
		videoId string
	}{
		{
			issue:   "issue#65",
			videoId: "9ja-K2FslBU",
		},
		{
			issue:   "issue#59",
			videoId: "nINQjT7Zr9w",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.issue, func(t *testing.T) {
			assert.Error(t, y.DecodeURL(tc.videoId))
		})
	}
}
