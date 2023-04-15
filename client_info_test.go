package youtube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientInfos(t *testing.T) {
	clients := getClientInfos([]string{"mweb", "web", "android"})
	require.Len(t, clients, 3)

	assert.Equal(t, 0, int(clients[0].priority))
	assert.Equal(t, 30, int(clients[1].priority))
	assert.Equal(t, 40, int(clients[2].priority))
}
