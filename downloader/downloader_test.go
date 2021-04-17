package downloader

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDownloader = func() (dl Downloader) {
	dl.OutputDir = "download_test"
	dl.Debug = true
	return
}()

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
	assert.Len(video.Formats, 18)

	if assert.Greater(len(video.Formats), 0) {
		assert.NoError(testDownloader.Download(ctx, video, &video.Formats[0], ""))
	}
}

func TestYoutube_DownloadWithHighQualityFails(t *testing.T) {
	tests := []struct {
		name    string
		formats []youtube.Format
		message string
	}{
		{
			name:    "video format not found",
			formats: []youtube.Format{{ItagNo: 140}},
			message: "no video format found after filtering",
		},
		{
			name:    "audio format not found",
			formats: []youtube.Format{{ItagNo: 137, Quality: "hd1080", MimeType: "video/mp4", AudioChannels: 0}},
			message: "no audio format found after filtering",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			video := &youtube.Video{
				Formats: tt.formats,
			}

			err := testDownloader.DownloadComposite(context.Background(), "", video, "hd1080", "")
			assert.EqualError(t, err, tt.message)
		})
	}
}
