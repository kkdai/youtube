package youtube

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestParseViewCount(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{"1,234 views", 0}, // comma-formatted text isn't handled; expect 0 (unknown)
		{"12 views", 12},
		{"1.2K views", 1200},
		{"3.4M views", 3_400_000},
		{"2B views", 2_000_000_000},
		{"", 0},
		{"No views", 0},
	}

	for _, c := range cases {
		got := parseViewCount(c.in)
		if got != c.want {
			t.Errorf("parseViewCount(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestParseSearchResults(t *testing.T) {
	// Minimal synthetic response mimicking InnerTube's search shape,
	// with one videoRenderer and one non-video renderer that should be skipped.
	raw := `{
		"contents": {
			"twoColumnSearchResultsRenderer": {
				"primaryContents": {
					"sectionListRenderer": {
						"contents": [
							{
								"itemSectionRenderer": {
									"contents": [
										{
											"channelRenderer": {
												"channelId": "UC_should_be_skipped"
											}
										},
										{
											"videoRenderer": {
												"videoId": "abc123XYZ_9",
												"title": {
													"runs": [{"text": "Test Video Title"}]
												},
												"ownerText": {
													"runs": [{
														"text": "Test Channel",
														"navigationEndpoint": {
															"browseEndpoint": {"browseId": "UCabc123"}
														}
													}]
												},
												"lengthText": {"simpleText": "10:00"},
												"viewCountText": {"simpleText": "1.2M views"},
												"thumbnail": {
													"thumbnails": [
														{"url": "https://example.com/thumb.jpg", "width": 120, "height": 90}
													]
												},
												"detailedMetadataSnippets": [
													{"snippetText": {"runs": [{"text": "A short description."}]}}
												]
											}
										}
									]
								}
							}
						]
					}
				}
			}
		}
	}`

	results, err := parseSearchResults([]byte(raw))
	if err != nil {
		t.Fatalf("parseSearchResults returned error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result (channelRenderer should be skipped), got %d", len(results))
	}

	got := results[0]

	if got.ID != "abc123XYZ_9" {
		t.Errorf("ID = %q, want %q", got.ID, "abc123XYZ_9")
	}
	if got.Title != "Test Video Title" {
		t.Errorf("Title = %q, want %q", got.Title, "Test Video Title")
	}
	if got.Author != "Test Channel" {
		t.Errorf("Author = %q, want %q", got.Author, "Test Channel")
	}
	if got.ChannelID != "UCabc123" {
		t.Errorf("ChannelID = %q, want %q", got.ChannelID, "UCabc123")
	}
	if got.Duration != "10:00" {
		t.Errorf("Duration = %q, want %q", got.Duration, "10:00")
	}
	if got.ViewCountRaw != "1.2M views" {
		t.Errorf("ViewCountRaw = %q, want %q", got.ViewCountRaw, "1.2M views")
	}
	if got.Description != "A short description." {
		t.Errorf("Description = %q, want %q", got.Description, "A short description.")
	}
	if len(got.Thumbnails) != 1 || got.Thumbnails[0].URL != "https://example.com/thumb.jpg" {
		t.Errorf("Thumbnails = %+v, unexpected", got.Thumbnails)
	}
}

func TestParseSearchResults_EmptyOrMalformed(t *testing.T) {
	// Completely unrelated JSON shape — should not error, just return no results.
	raw := `{"unexpectedField": true}`

	results, err := parseSearchResults([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error on unrecognized shape, got: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestParseSearchResults_InvalidJSON(t *testing.T) {
	raw := `{not valid json`

	_, err := parseSearchResults([]byte(raw))
	if err == nil {
		t.Fatal("expected an error for invalid JSON, got nil")
	}
}

// TestVideoRendererToResult_MissingFields ensures a videoRenderer with
// only the required videoId field doesn't panic and yields sane zero values.
func TestVideoRendererToResult_MissingFields(t *testing.T) {
	raw := `{"videoRenderer": {"videoId": "onlyID1234"}}`

	var wrapper videoRendererWrapper
	if err := json.Unmarshal([]byte(raw), &wrapper); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if wrapper.VideoRenderer == nil {
		t.Fatal("expected non-nil VideoRenderer")
	}

	got := videoRendererToResult(wrapper.VideoRenderer)
	if got.ID != "onlyID1234" {
		t.Errorf("ID = %q, want %q", got.ID, "onlyID1234")
	}
	if got.Title != "" || got.Author != "" || got.Duration != "" {
		t.Errorf("expected empty optional fields, got %+v", got)
	}
}

// --- live network test, skipped by default ---
//
// Run with: RUN_LIVE_TESTS=1 go test -run TestSearch_Live -v ./...
// Mirrors the pattern used elsewhere in this package for tests that hit
// the real YouTube InnerTube API.

func TestSearch_Live(t *testing.T) {
	if os.Getenv("RUN_LIVE_TESTS") == "" {
		t.Skip("skipping live network test; set RUN_LIVE_TESTS=1 to run")
	}

	client := Client{}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	results, err := client.SearchContext(ctx, "dotGo 2015 Rob Pike")
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected at least one search result")
	}

	first := results[0]
	if first.ID == "" {
		t.Error("expected non-empty video ID on first result")
	}
	if first.Title == "" {
		t.Error("expected non-empty title on first result")
	}

	t.Logf("first result: %q by %q (%s)", first.Title, first.Author, first.Duration)
}