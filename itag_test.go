// +build integration

package youtube

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYoutube_GetItagInfo(t *testing.T) {
	require := require.New(t)
	y := NewYoutube(false, false)

	// url from issue #25
	testVideoUrl := "https://www.youtube.com/watch?v=rFejpH_tAHM"
	require.NoError(y.DecodeURL(testVideoUrl))

	itagInfo := y.GetStreamInfo()

	require.Len(itagInfo.Streams, 18)
}
