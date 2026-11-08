package youtube

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Exercise both parsers and the JavaScript evaluator without a live player.
const testDecipherPlayer = `
var NT=[nf];
nf=function(a){return a.split("").reverse().join("")};
x.get("n"))&&(b=NT[0](b)||nf);
var ST={splice:function(a,b){a.splice(0,b)},reverse:function(a){a.reverse()}};
sf=function(a){a=a.split("");ST.splice(a,1);ST.reverse(a,0);return a.join("")};
`

type decipherTestTransport func(*http.Request) (*http.Response, error)

func (f decipherTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestGetStreamURLNParam(t *testing.T) {
	for _, cipher := range []bool{false, true} {
		name := "direct"
		if cipher {
			name = "signatureCipher"
		}
		t.Run(name, func(t *testing.T) {
			for _, rawQuery := range []string{
				"n=abc123&v=video-id&itag=140&tag=one&tag=two&token=a%2Bb%26c",
				"n=abc123&itag=140",
				"v=video-id&itag=140",
				"n=&v=video-id&itag=140",
			} {
				t.Run(rawQuery, func(t *testing.T) {
					want, err := url.ParseQuery(rawQuery)
					require.NoError(t, err)
					// Without a nonempty n, no n-function may be needed.
					player := testDecipherPlayer
					if want.Get("n") == "" {
						player = player[strings.Index(player, "var ST="):]
					} else {
						want.Set("n", "321cba")
					}
					client := Client{
						client: &WebClient,
						HTTPClient: &http.Client{Transport: decipherTestTransport(func(req *http.Request) (*http.Response, error) {
							var body string
							switch req.URL.Path {
							case "/embed/test-video":
								body = `<script src="/s/player/test/player_ias.vflset/en_US/base.js"></script>`
							case "/s/player/test/player_ias.vflset/en_US/base.js":
								body = player
							default:
								t.Fatalf("unexpected request: %s", req.URL)
							}
							return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
						})},
					}
					format := Format{URL: "https://media.example/videoplayback?" + rawQuery}
					if cipher {
						format.Cipher = url.Values{
							"url": {format.URL}, "s": {"abcdef"}, "sp": {"sig"},
						}.Encode()
						format.URL = ""
						want.Set("sig", "fedcb")
					}

					got, err := client.GetStreamURLContext(context.Background(), &Video{ID: "test-video"}, &format)
					require.NoError(t, err)
					parsed, err := url.Parse(got)
					require.NoError(t, err)
					require.Equal(t, "https", parsed.Scheme)
					require.Equal(t, "media.example", parsed.Host)
					require.Equal(t, "/videoplayback", parsed.Path)
					// Full equality also checks that v and repeated/escaped values survive,
					// and that no v parameter is added when it was absent.
					require.Equal(t, want, parsed.Query())
				})
			}
		})
	}
}

func TestDecryptNParamError(t *testing.T) {
	query := url.Values{"n": {"abc123"}, "v": {"video-id"}}
	got, err := new(Client).decryptNParam(nil, query)
	require.EqualError(t, err, "unable to decode nSig: unable to extract n-function name")
	require.Nil(t, got)
	require.Equal(t, url.Values{"n": {"abc123"}, "v": {"video-id"}}, query)
}
