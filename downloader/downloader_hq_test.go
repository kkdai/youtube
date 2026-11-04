//go:build integration
// +build integration

package downloader

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownload_HighQuality(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// "How to: Enable captions in Google Meet" by Google Workspace (13s, up to 1080p)
	video, err := testDownloader.Client.GetVideoContext(ctx, "ezfPHVOf-NE")
	require.NoError(err)
	require.NoError(testDownloader.DownloadComposite(ctx, "", video, "hd1080", "mp4", ""))
}
