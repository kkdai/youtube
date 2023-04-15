package youtube

import (
	"strings"

	"golang.org/x/exp/slices"
)

const embedUrl = "https://www.youtube.com/" // Can be any valid URL
const defaultLanguage = "en"

var defaultClients = []string{"android", "web"}
var webClient = clientInfos[0]
var embeddedClient = clientInfos[0]

// client info for the innertube API
type clientInfo struct {
	name             string
	priority         byte
	innertubeContext innertubeContext
	innertubeApiKey  string
	requireJsPlayer  bool
}

func init() {
	prepareClientInfos()
}

func prepareClientInfos() {
	for i := range clientInfos {
		client := clientInfos[i]
		ctx := &client.innertubeContext

		if strings.HasSuffix(client.name, "_embedded") {
			ctx.ThirdParty.EmbedURL = embedUrl
		}

		if ctx.Client.HL == "" {
			ctx.Client.HL = defaultLanguage
		}
	}
}

func getClientInfos(allowedClients []string) (result []clientInfo) {
	for i := range clientInfos {
		if slices.Contains(allowedClients, clientInfos[i].name) {
			result = append(result, clientInfos[i])
		}
	}

	slices.SortFunc(result, func(a, b clientInfo) bool {
		return a.priority < b.priority
	})

	return
}

type innertubeRequest struct {
	VideoID         string           `json:"videoId,omitempty"`
	BrowseID        string           `json:"browseId,omitempty"`
	Continuation    string           `json:"continuation,omitempty"`
	Context         innertubeContext `json:"context"`
	PlaybackContext playbackContext  `json:"playbackContext,omitempty"`
}

type playbackContext struct {
	ContentPlaybackContext contentPlaybackContext `json:"contentPlaybackContext"`
}

type contentPlaybackContext struct {
	SignatureTimestamp string `json:"signatureTimestamp"`
}

type innertubeContext struct {
	Client     innertubeClient `json:"client"`
	ThirdParty innertubeThirdparty
}

type innertubeThirdparty struct {
	EmbedURL string `json:"embedUrl"`
}

type innertubeClient struct {
	HL                string `json:"hl"`
	GL                string `json:"gl"`
	Clientname        string `json:"clientName"`
	Clientversion     string `json:"clientVersion"`
	Devicemodel       string `json:"deviceModel"`
	Useragent         string `json:"userAgent"`
	Thirdparty        string `json:"thirdParty"`
	Androidsdkversion byte   `json:"androidsdkversion"`
}
