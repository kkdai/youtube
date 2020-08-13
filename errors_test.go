package youtube

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		err      error
		expected string
	}{
		{ErrUnexpectedStatusCode(404), "unexpected status code: 404"},
		{ErrPlayabiltyStatus{"invalid", "for that reason"}, "cannot playback and download, status: invalid, reason: for that reason"},
		{ErrResponseStatus{}, "no response status found in the server's answer"},
		{ErrResponseStatus{Status: "foo"}, "response status: 'foo', no reason given"},
		{ErrResponseStatus{Status: "foo", Reason: "bar"}, "response status: 'foo', reason: 'bar'"},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.EqualError(t, tt.err, tt.expected)
		})
	}
}
