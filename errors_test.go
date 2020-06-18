package youtube

import (
	"fmt"
	"strings"
	"testing"
)

func TestErrDecodingStreamInfo_Error(t *testing.T) {
	streamPos := 1
	err := ErrDecodingStreamInfo{
		streamPos: streamPos,
	}
	substr := fmt.Sprintf(`stream %d`, streamPos)
	if got := err.Error(); !strings.Contains(got, substr) {
		t.Errorf("Error() = %v should contain %v", got, substr)
	}
}
