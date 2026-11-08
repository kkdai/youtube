package youtube

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type playerParamsRoundTripper func(*http.Request) (*http.Response, error)

func (f playerParamsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func playerParamsHTTPClient(t *testing.T, handle func(*http.Request) string) *http.Client {
	t.Helper()
	return &http.Client{Transport: playerParamsRoundTripper(func(req *http.Request) (*http.Response, error) {
		var body string
		if req.URL.Path == "" || req.URL.Path == "/" {
			require.Equal(t, http.MethodGet, req.Method)
			require.Nil(t, req.Body)
			body = "\nytcfg.set({\"INNERTUBE_CONTEXT\":{\"client\":{\"visitorData\":\"test-visitor\"}}});"
		} else {
			body = handle(req)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})}
}

func TestClientPlayerParams(t *testing.T) {
	defaultInfo := DefaultClient
	androidInfo := AndroidClient
	require.Equal(t, androidInfo, defaultInfo)

	clients := []struct {
		name   string
		client Client
	}{
		{name: "default"},
		{name: "2AMB", client: Client{PlayerParams: "2AMB"}},
		{name: "opaque", client: Client{PlayerParams: "opaque +/=\n\"value\""}},
	}

	// Interleave requests from independent instances, then reuse each instance.
	for range 2 {
		for i := range clients {
			tc := &clients[i]
			t.Run(tc.name, func(t *testing.T) {
				requests := 0
				params := tc.client.PlayerParams
				tc.client.HTTPClient = playerParamsHTTPClient(t, func(req *http.Request) string {
					requests++
					require.Equal(t, http.MethodPost, req.Method)
					require.Equal(t, "/youtubei/v1/player", req.URL.Path)
					require.Equal(t, androidInfo.Key, req.URL.Query().Get("key"))
					require.Equal(t, "3", req.Header.Get("X-YouTube-Client-Name"))
					require.Equal(t, androidInfo.Version, req.Header.Get("X-YouTube-Client-Version"))
					require.Equal(t, androidInfo.UserAgent, req.UserAgent())
					require.Equal(t, "application/json", req.Header.Get("Content-Type"))
					require.Equal(t, "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8", req.Header.Get("Accept"))
					require.Equal(t, "https://youtube.com", req.Header.Get("Origin"))
					require.Equal(t, "navigate", req.Header.Get("Sec-Fetch-Mode"))
					require.Equal(t, "test-visitor", req.Header.Get("X-Goog-Visitor-Id"))

					var data map[string]any
					require.NoError(t, json.NewDecoder(req.Body).Decode(&data))
					require.Contains(t, data, "params") // Empty must be present, not omitted.
					require.Equal(t, params, data["params"])
					require.Equal(t, "9_MbW9FK1fA", data["videoId"])
					require.Equal(t, true, data["contentCheckOk"])
					require.Equal(t, true, data["racyCheckOk"])
					require.Equal(t, map[string]any{
						"contentPlaybackContext": map[string]any{"html5Preference": "HTML5_PREF_WANTS"},
					}, data["playbackContext"])
					info := data["context"].(map[string]any)["client"].(map[string]any)
					require.NotEmpty(t, info["visitorData"])
					delete(info, "visitorData")
					require.Equal(t, map[string]any{
						"hl": "en", "gl": "US", "timeZone": "UTC", "utcOffsetMinutes": float64(0),
						"clientName": "ANDROID", "clientVersion": androidInfo.Version, "userAgent": androidInfo.UserAgent,
					}, info)
					return `{}`
				})

				tc.client.assureClient()
				_, err := tc.client.videoDataByInnertube(context.Background(), "9_MbW9FK1fA")
				require.NoError(t, err)
				require.Equal(t, 1, requests)
				require.Same(t, &DefaultClient, tc.client.client)
			})
		}
	}
	require.Equal(t, defaultInfo, DefaultClient)
	require.Equal(t, androidInfo, AndroidClient)
}

func TestClientPlayerParamsDoesNotAffectPlaylist(t *testing.T) {
	client := Client{PlayerParams: "2AMB"}
	requests := 0
	client.HTTPClient = playerParamsHTTPClient(t, func(req *http.Request) string {
		requests++
		require.Equal(t, "/youtubei/v1/browse", req.URL.Path)
		var data innertubeRequest
		require.NoError(t, json.NewDecoder(req.Body).Decode(&data))
		require.Equal(t, "CgIQBg==", data.Params)
		switch requests {
		case 1:
			require.Equal(t, "VLPL59FEE129ADFF2B12", data.BrowseID)
			return `{
				"metadata": {"playlistHeaderRenderer": {"title": "Test playlist"}},
				"contents": {"singleColumnBrowseResultsRenderer": {"tabs": [{"tabRenderer": {"content": {
					"sectionListRenderer": {"contents": [{"playlistVideoListRenderer": {"contents": [
						{"playlistVideoRenderer": {"videoId": "first-video"}},
						{"continuationItemRenderer": {"continuationEndpoint": {"continuationCommand": {"token": "next-page"}}}}
					]}}]}
				}}}]}}
			}`
		case 2:
			require.Equal(t, "next-page", data.Continuation)
			return `{"onResponseReceivedActions":[{"appendContinuationItemsAction":{"continuationItems":[
				{"playlistVideoRenderer":{"videoId":"second-video"}}
			]}}]}`
		default:
			t.Fatalf("unexpected browse request %d", requests)
			return ""
		}
	})

	playlist, err := client.GetPlaylist("PL59FEE129ADFF2B12")
	require.NoError(t, err)
	require.Equal(t, 2, requests)
	require.Len(t, playlist.Videos, 2)
}

func TestClientPlayerParamsEmbeddedFallback(t *testing.T) {
	client := Client{PlayerParams: "2AMB"}
	var names []string
	client.HTTPClient = playerParamsHTTPClient(t, func(req *http.Request) string {
		require.Equal(t, "/youtubei/v1/player", req.URL.Path)
		var data innertubeRequest
		require.NoError(t, json.NewDecoder(req.Body).Decode(&data))
		require.Equal(t, "2AMB", data.Params)
		names = append(names, data.Context.Client.ClientName)
		if len(names) == 1 {
			return `{"playabilityStatus":{"status":"LOGIN_REQUIRED","reason":"Sign in to confirm your age"}}`
		}
		return `{"playabilityStatus":{"status":"OK"},"streamingData":{"adaptiveFormats":[
			{"itag":140,"mimeType":"audio/mp4","url":"https://example.com/audio"}
		]}}`
	})

	video, err := client.GetVideo("9_MbW9FK1fA")
	require.NoError(t, err)
	require.Len(t, video.Formats, 1)
	require.Equal(t, []string{"ANDROID", "WEB_EMBEDDED_PLAYER"}, names)
	require.Same(t, &DefaultClient, client.client)
	require.Equal(t, "2AMB", client.PlayerParams)
	require.Equal(t, AndroidClient, DefaultClient)
}
