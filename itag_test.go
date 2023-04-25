package youtube

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYoutube_GetItagInfo(t *testing.T) {
	require := require.New(t)
	client := Client{}

	// url from issue #25
	url := "https://www.youtube.com/watch?v=rFejpH_tAHM"
	video, err := client.GetVideo(url)
	require.NoError(err)
	require.GreaterOrEqual(len(video.Formats), 24)
}
