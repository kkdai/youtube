package youtube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtoBuilder(t *testing.T) {
	var pb ProtoBuilder

	pb.Varint(1, 128)
	pb.Varint(2, 1234567890)
	pb.Varint(3, 1234567890123456789)
	pb.String(4, "Hello")
	pb.Bytes(5, []byte{1, 2, 3})
	assert.Equal(t, "CIABENKF2MwEGJWCpu_HnoSRESIFSGVsbG8qAwECAw%3D%3D", pb.ToUrlEncodedBase64())
}
