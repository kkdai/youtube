package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// SearchResult represents a single public video returned by a search query.
type SearchResult struct {
	ID           string
	Title        string
	Author       string
	ChannelID    string
	Duration     string // e.g. "12:34" as rendered by YouTube; empty for live streams
	ViewCountRaw string // e.g. "1.2M views", raw text as returned by YouTube
	Thumbnails   Thumbnails
	Description  string
}

// searchInnertubeRequest mirrors innertubeRequest but adds the query field
// used by the /search endpoint.
type searchInnertubeRequest struct {
	Query   string            `json:"query"`
	Context inntertubeContext `json:"context"`
	Params  string            `json:"params,omitempty"`
}

// searchParamsVideosOnly restricts InnerTube search results to the "Video"
// content type (Filter -> Type -> Video), matching what the web client sends.
const searchParamsVideosOnly = "EgIQAQ%3D%3D"

// Search performs a video-only search and returns matching public videos.
func (c *Client) Search(query string) ([]SearchResult, error) {
	return c.SearchContext(context.Background(), query)
}

// SearchContext performs a video-only search with a context.
func (c *Client) SearchContext(ctx context.Context, query string) ([]SearchResult, error) {
	c.assureClient()

	data := searchInnertubeRequest{
		Query:   query,
		Context: prepareInnertubeContext(*c.client),
		Params:  searchParamsVideosOnly,
	}

	body, err := c.httpPostBodyBytes(ctx, "https://www.youtube.com/youtubei/v1/search?key="+c.client.Key, data)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}

	results, err := parseSearchResults(body)
	if err != nil {
		return nil, fmt.Errorf("parsing search results failed: %w", err)
	}

	return results, nil
}

// --- response parsing ---
//
// InnerTube search responses nest results under a chain of renderer
// objects. We decode into a generic structure and walk it, keeping only
// videoRenderer entries, since itemSectionRenderer.contents mixes several
// renderer types (channelRenderer, shelfRenderer, adSlotRenderer, etc.).

type searchResponse struct {
	Contents struct {
		TwoColumnSearchResultsRenderer struct {
			PrimaryContents struct {
				SectionListRenderer struct {
					Contents []searchSection `json:"contents"`
				} `json:"sectionListRenderer"`
			} `json:"primaryContents"`
		} `json:"twoColumnSearchResultsRenderer"`
	} `json:"contents"`
}

type searchSection struct {
	ItemSectionRenderer struct {
		Contents []json.RawMessage `json:"contents"`
	} `json:"itemSectionRenderer"`
}

type videoRendererWrapper struct {
	VideoRenderer *videoRenderer `json:"videoRenderer"`
}

type videoRenderer struct {
	VideoID string `json:"videoId"`
	Title   struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
	} `json:"title"`
	OwnerText struct {
		Runs []struct {
			Text               string `json:"text"`
			NavigationEndpoint struct {
				BrowseEndpoint struct {
					BrowseID string `json:"browseId"`
				} `json:"browseEndpoint"`
			} `json:"navigationEndpoint"`
		} `json:"runs"`
	} `json:"ownerText"`
	LengthText struct {
		SimpleText string `json:"simpleText"`
	} `json:"lengthText"`
	ViewCountText struct {
		SimpleText string `json:"simpleText"`
	} `json:"viewCountText"`
	Thumbnail struct {
		Thumbnails []struct {
			URL    string `json:"url"`
			Width  uint   `json:"width"`
			Height uint   `json:"height"`
		} `json:"thumbnails"`
	} `json:"thumbnail"`
	DetailedMetadataSnippets []struct {
		SnippetText struct {
			Runs []struct {
				Text string `json:"text"`
			} `json:"runs"`
		} `json:"snippetText"`
	} `json:"detailedMetadataSnippets"`
}

func parseSearchResults(body []byte) ([]SearchResult, error) {
	var resp searchResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	var results []SearchResult

	sections := resp.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents
	for _, section := range sections {
		for _, raw := range section.ItemSectionRenderer.Contents {
			var wrapper videoRendererWrapper
			if err := json.Unmarshal(raw, &wrapper); err != nil {
				// Skip renderer types we don't know how to parse rather
				// than fail the whole search.
				continue
			}
			vr := wrapper.VideoRenderer
			if vr == nil || vr.VideoID == "" {
				// Not a video result (channelRenderer, shelfRenderer,
				// adSlotRenderer, etc.) — skip it.
				continue
			}

			results = append(results, videoRendererToResult(vr))
		}
	}

	return results, nil
}

func videoRendererToResult(vr *videoRenderer) SearchResult {
	var title string
	if len(vr.Title.Runs) > 0 {
		var b strings.Builder
		for _, run := range vr.Title.Runs {
			b.WriteString(run.Text)
		}
		title = b.String()
	}

	var author, channelID string
	if len(vr.OwnerText.Runs) > 0 {
		author = vr.OwnerText.Runs[0].Text
		channelID = vr.OwnerText.Runs[0].NavigationEndpoint.BrowseEndpoint.BrowseID
	}

	var description string
	if len(vr.DetailedMetadataSnippets) > 0 {
		var b strings.Builder
		for _, run := range vr.DetailedMetadataSnippets[0].SnippetText.Runs {
			b.WriteString(run.Text)
		}
		description = b.String()
	}

	var thumbs Thumbnails
	for _, t := range vr.Thumbnail.Thumbnails {
		thumbs = append(thumbs, Thumbnail{URL: t.URL, Width: t.Width, Height: t.Height})
	}

	return SearchResult{
		ID:           vr.VideoID,
		Title:        title,
		Author:       author,
		ChannelID:    channelID,
		Duration:     vr.LengthText.SimpleText,
		ViewCountRaw: vr.ViewCountText.SimpleText,
		Thumbnails:   thumbs,
		Description:  description,
	}
}

// parseViewCount is a best-effort helper to turn "1.2M views" style text
// into an approximate integer. Returns 0 if parsing fails — callers should
// treat that as "unknown", not "zero views".
func parseViewCount(raw string) int64 {
	s := strings.ToLower(raw)
	s = strings.TrimSuffix(s, " views")
	s = strings.TrimSuffix(s, " view")
	s = strings.TrimSpace(s)

	multiplier := int64(1)
	switch {
	case strings.HasSuffix(s, "k"):
		multiplier = 1_000
		s = strings.TrimSuffix(s, "k")
	case strings.HasSuffix(s, "m"):
		multiplier = 1_000_000
		s = strings.TrimSuffix(s, "m")
	case strings.HasSuffix(s, "b"):
		multiplier = 1_000_000_000
		s = strings.TrimSuffix(s, "b")
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(f * float64(multiplier))
}
