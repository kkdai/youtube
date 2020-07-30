package downloader

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/kkdai/youtube"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testClient = youtube.Client{Debug: true}

	testDownloader = func() (dl Downloader) {
		dl.OutputDir = "download_test"
		dl.Client = testClient
		return
	}()
)

func TestMain(m *testing.M) {

	exitCode := m.Run()
	// the following code doesn't work under debugger, please delete download files manually
	if err := os.RemoveAll(testDownloader.OutputDir); err != nil {
		log.Fatal(err.Error())
	}
	os.Exit(exitCode)
}

func TestDownload_FirstStream(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	ctx := context.Background()

	// youtube-dl test video
	video, err := testDownloader.Client.GetVideoContext(ctx, "BaW_jenozKc")
	require.NoError(err)
	require.NotNil(video)

	assert.Equal(`youtube-dl test video "'/\ä↭𝕐`, video.Title)
	assert.Equal(`Philipp Hagemeister`, video.Author)
	assert.Equal(10*time.Second, video.Duration)
	assert.Len(video.Streams, 18)

	if assert.Greater(len(video.Streams), 0) {
		assert.NoError(testDownloader.Download(ctx, "", video, &video.Streams[0]))
	}
}

func TestYoutube_DownloadWithHighQualityFails(t *testing.T) {
	tests := []struct {
		name    string
		streams []youtube.Stream
		message string
	}{
		{
			name:    "video Stream not found",
			streams: []youtube.Stream{{ItagNo: 140}},
			message: "no Stream video/mp4 for hd1080 found",
		},
		{
			name:    "audio Stream not found",
			streams: []youtube.Stream{{ItagNo: 137}},
			message: "no Stream audio/mp4 for hd1080 found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			video := &youtube.Video{
				Streams: tt.streams,
			}

			err := testDownloader.DownloadWithHighQuality(context.Background(), "", video, "hd1080")
			assert.EqualError(t, err, tt.message)
		})
	}
}
