package youtube

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseM3U8(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	file, err := os.Open("testdata/index.m3u8")

	require.NoError(err)
	assert.NotNil(file)
	defer file.Close()

	result, err := parseM3U8(file)
	require.NoError(err)
	require.Len(result, 4)

	assert.Equal(229, result[0].ItagNo)
}
