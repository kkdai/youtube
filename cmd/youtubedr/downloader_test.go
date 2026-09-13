package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProxyConfigFromEnvironment_AllProxyFallback(t *testing.T) {
	t.Setenv("HTTP_PROXY", "")
	t.Setenv("HTTPS_PROXY", "")
	t.Setenv("http_proxy", "")
	t.Setenv("https_proxy", "")
	t.Setenv("ALL_PROXY", "")
	t.Setenv("all_proxy", "http://127.0.0.1:17891")

	cfg := proxyConfigFromEnvironment()
	require.Equal(t, "http://127.0.0.1:17891", cfg.HTTPProxy)
	require.Equal(t, "http://127.0.0.1:17891", cfg.HTTPSProxy)
}

func TestProxyConfigFromEnvironment_HTTPSProxyWins(t *testing.T) {
	t.Setenv("HTTP_PROXY", "")
	t.Setenv("http_proxy", "")
	t.Setenv("https_proxy", "")
	t.Setenv("HTTPS_PROXY", "http://127.0.0.1:1")
	t.Setenv("all_proxy", "http://127.0.0.1:2")

	cfg := proxyConfigFromEnvironment()
	require.Equal(t, "http://127.0.0.1:1", cfg.HTTPSProxy)
}
