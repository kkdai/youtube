package youtube

import (
	"fmt"
	"net/url"
	"path"
	"slices"
	"strings"
)

type Thumbnails []Thumbnail

// Possible thumbnail names in order of preference.
var thumbnailNames = [...]string{
	"maxresdefault",
	"hq720",
	"sddefault",
	"hqdefault",
	"0",
	"mqdefault",
	"default",
	"sd1",
	"sd2",
	"sd3",
	"hq1",
	"hq2",
	"hq3",
	"mq1",
	"mq2",
	"mq3",
	"1",
	"2",
	"3",
}

// Resolutions of potential thumbnails.
// See: https://stackoverflow.com/a/20542029
var thumbnailResolutions = map[string][2]uint{
	"maxresdefault": {1920, 1080},
	"hq720":         {1280, 720},
	"sddefault":     {640, 480},
	"sd3":           {640, 480},
	"sd2":           {640, 480},
	"sd1":           {640, 480},
	"hqdefault":     {480, 360},
	"hq3":           {480, 360},
	"hq2":           {480, 360},
	"hq1":           {480, 360},
	"0":             {480, 360},
	"mqdefault":     {320, 180},
	"mq3":           {320, 180},
	"mq2":           {320, 180},
	"mq1":           {320, 180},
	"default":       {120, 90},
	"1":             {120, 90},
	"2":             {120, 90},
	"3":             {120, 90},
}

var thumbnailExtensions = [...]string{
	"_live.webp",
	"_live.jpg",
	".webp",
	".jpg",
}

// PossibleThumbnails returns a list of known possible thumbnail URLs.
func PossibleThumbnails(videoID string) Thumbnails {
	thumbnails := make(Thumbnails, 0, len(thumbnailNames)*len(thumbnailExtensions))
	for _, name := range thumbnailNames {
		for _, ext := range thumbnailExtensions {
			thumbnailSize := thumbnailResolutions[name]
			thumbnail := Thumbnail{
				Width:  thumbnailSize[0],
				Height: thumbnailSize[1],
			}

			if strings.HasSuffix(ext, ".webp") {
				thumbnail.URL = fmt.Sprintf("https://i.ytimg.com/vi_webp/%s/%s%s", videoID, name, ext)
			} else {
				thumbnail.URL = fmt.Sprintf("https://i.ytimg.com/vi/%s/%s%s", videoID, name, ext)
			}

			thumbnails = append(thumbnails, thumbnail)
		}
	}
	return thumbnails
}

// Extended returns an extended list of possible thumbnails including some not
// returned in the video response. These additional thumbnails may or may
// not be available, and resolution information may not be accurate.
func (t Thumbnails) Extended(videoID string) Thumbnails {
	possible := PossibleThumbnails(videoID)
	extended := make([]Thumbnail, len(t)+len(possible))
	copy(extended, t)
	copy(extended[len(t):], possible)
	return extended
}

// Sort sorts the thumbnail list, abiding by the same sorting as used by yt-dlp.
func (t Thumbnails) Sort() {
	slices.SortStableFunc(t, cmpThumbnails)
}

// FilterExt removes thumbnails that do not have the given extension.
func (t Thumbnails) FilterExt(ext ...string) Thumbnails {
	return slices.DeleteFunc(t, func(thumbnail Thumbnail) bool {
		u, err := url.Parse(thumbnail.URL)
		if err != nil {
			return true
		}
		return !slices.Contains(ext, path.Ext(u.Path))
	})
}

// FilterLive removes thumbnails that do not match the provided live status.
func (t Thumbnails) FilterLive(live bool) Thumbnails {
	return slices.DeleteFunc(t, func(thumbnail Thumbnail) bool {
		name := path.Base(thumbnail.URL)
		parts := strings.SplitN(name, "_", 2)
		l := len(parts) > 1 && strings.HasPrefix(parts[1], "live")
		return l != live
	})
}

// MinWidth filters out thumbnails with greater than desired width.
func (t Thumbnails) MinWidth(w uint) Thumbnails {
	return slices.DeleteFunc(t, func(thumbnail Thumbnail) bool {
		return thumbnail.Width < w
	})
}

// MaxWidth filters out thumbnails with less than desired width.
func (t Thumbnails) MaxWidth(w uint) Thumbnails {
	return slices.DeleteFunc(t, func(thumbnail Thumbnail) bool {
		return thumbnail.Width > w
	})
}

// MinHeight filters out thumbnails with greater than desired height.
func (t Thumbnails) MinHeight(w uint) Thumbnails {
	return slices.DeleteFunc(t, func(thumbnail Thumbnail) bool {
		return thumbnail.Height < w
	})
}

// MaxHeight filters out thumbnails with less than desired height.
func (t Thumbnails) MaxHeight(w uint) Thumbnails {
	return slices.DeleteFunc(t, func(thumbnail Thumbnail) bool {
		return thumbnail.Height > w
	})
}

func cmpThumbnails(t1, t2 Thumbnail) int {
	score1 := scoreThumbnail(t1)
	score2 := scoreThumbnail(t2)
	return score2 - score1
}

func scoreThumbnail(t Thumbnail) int {
	nameIndex := 0
	for i, name := range thumbnailNames {
		if strings.Contains(t.URL, name) {
			nameIndex = i
			break
		}
	}

	score := -2 * nameIndex

	if strings.Contains(t.URL, ".webp") {
		score += 1
	}

	return score
}
