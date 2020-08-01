package youtube

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type YoutubeTest struct {
	Note         string `json:"note"`
	URL          string `json:"url"`
	OnlyMatching bool   `json:"only_matching"`
	InfoDict     struct {
		ID           string      `json:"id"`
		Ext          string      `json:"ext"`
		Title        string      `json:"title"`
		Uploader     string      `json:"uploader"`
		UploaderID   string      `json:"uploader_id"`
		UploaderURL  string      `json:"uploader_url"`
		ChannelID    string      `json:"channel_id"`
		ChannelURL   string      `json:"channel_url"`
		UploadDate   string      `json:"upload_date"`
		Description  string      `json:"description"`
		Categories   []string    `json:"categories"`
		Tags         []string    `json:"tags"`
		Duration     *int        `json:"duration"`
		ViewCount    interface{} `json:"view_count"`
		LikeCount    interface{} `json:"like_count"`
		DislikeCount interface{} `json:"dislike_count"`
		StartTime    int         `json:"start_time"`
		EndTime      int         `json:"end_time"`
	} `json:"info_dict"`
}

func TestYoutubeDL(t *testing.T) {

	jsonFile, err := os.Open("testdata/tests.json")
	require.NoError(t, err)
	defer jsonFile.Close()

	var tests []YoutubeTest
	require.NoError(t, json.NewDecoder(jsonFile).Decode(&tests))

	for _, tc := range tests {
		if tc.OnlyMatching {
			continue
		}

		t.Run(tc.URL, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			if tc.Note != "" {
				log.Println(tc.Note)
			}

			video, err := testClient.GetVideo(tc.URL)
			require.NoError(err)
			require.NotNil(video)
			assert.Equal(tc.InfoDict.ID, video.ID)
			assert.Equal(tc.InfoDict.Uploader, video.Uploader)

			if title := tc.InfoDict.Title; strings.HasPrefix(title, "md5:") {
				md5Sum := md5.Sum([]byte(video.Title))
				hexSum := hex.EncodeToString(md5Sum[:])
				assert.Equal(title[4:], hexSum, "title: %v", video.Title)
			} else {
				assert.Equal(tc.InfoDict.Title, video.Title)
			}

			if tc.InfoDict.Duration != nil {
				assert.EqualValues(*tc.InfoDict.Duration, video.Duration.Seconds())
			}
		})
	}
}

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

			var format *Format
			if tc.itagNo > 0 {
				format = video.Formats.FindByItag(tc.itagNo)
				require.NotNil(format)
			} else {
				format = &video.Formats[0]
			}

			url, err := testClient.getStreamURL(ctx, video, format)
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
