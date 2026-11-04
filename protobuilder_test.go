package youtube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProtoBuilder(t *testing.T) {
	var pb ProtoBuilder

	require.NoError(t, pb.Varint(1, 128))
	require.NoError(t, pb.Varint(2, 1234567890))
	require.NoError(t, pb.Varint(3, 1234567890123456789))
	require.NoError(t, pb.String(4, "Hello"))
	require.NoError(t, pb.Bytes(5, []byte{1, 2, 3}))
	assert.Equal(t, "CIABENKF2MwEGJWCpu_HnoSRESIFSGVsbG8qAwECAw%3D%3D", pb.ToURLEncodedBase64())
}
