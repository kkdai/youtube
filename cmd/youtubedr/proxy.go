package main

import (
	"net/http"
	"net/url"

	"golang.org/x/net/http/httpproxy"
)

// Go's proxy config only uses HTTPS_PROXY for https:// requests.
// The CLI help advertises HTTP_PROXY, so treat it as the HTTPS fallback too.
func proxyFromEnvironmentWithHTTPSFallback() func(*http.Request) (*url.URL, error) {
	cfg := httpproxy.FromEnvironment()
	if cfg.HTTPSProxy == "" && cfg.HTTPProxy != "" {
		cfg.HTTPSProxy = cfg.HTTPProxy
	}

	proxyFunc := cfg.ProxyFunc()
	return func(r *http.Request) (*url.URL, error) {
		return proxyFunc(r.URL)
	}
}
