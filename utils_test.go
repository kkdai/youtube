package youtube

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetChunks1(t *testing.T) {
	require := require.New(t)
	chunks := getChunks(13, 5)

	require.Len(chunks, 3)
	require.EqualValues(0, chunks[0].start)
	require.EqualValues(4, chunks[0].end)
	require.EqualValues(5, chunks[1].start)
	require.EqualValues(9, chunks[1].end)
	require.EqualValues(10, chunks[2].start)
	require.EqualValues(12, chunks[2].end)
}

func TestGetChunks_length(t *testing.T) {
	require := require.New(t)
	require.Len(getChunks(10, 9), 2)
	require.Len(getChunks(10, 10), 1)
	require.Len(getChunks(10, 11), 1)
}
