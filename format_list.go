package youtube

import (
	"sort"
	"strconv"
	"strings"
)

type FormatList []Format

func (list FormatList) FindByQuality(quality string) *Format {
	for i := range list {
		if list[i].Quality == quality || list[i].QualityLabel == quality {
			return &list[i]
		}
	}
	return nil
}

func (list FormatList) FindByItag(itagNo int) *Format {
	for i := range list {
		if list[i].ItagNo == itagNo {
			return &list[i]
		}
	}
	return nil
}

// FindByType returns mime type of video which only audio or video
func (list FormatList) FindByType(t string) FormatList {
	var fl FormatList
	for i := range list {
		if strings.Contains(list[i].MimeType, t) {
			fl = append(fl, list[i])
		}
	}
	return fl
}

// FilterQuality returns a new FormatList filtered by quality, quality label or itag,
// but not audio quality
func (list FormatList) FilterQuality(quality string) FormatList {
	var fl FormatList
	for _, f := range list {
		itag, _ := strconv.Atoi(quality)
		if itag == f.ItagNo || strings.Contains(f.Quality, quality) || strings.Contains(f.QualityLabel, quality) {
			fl = append(fl, f)
		}
	}
	return fl
}

// FilterByAudioChannels returns a new FormatList filtered by the matching AudioChannels
func (list FormatList) FilterByAudioChannels(n int) FormatList {
	var fl FormatList
	for _, f := range list {
		if f.AudioChannels == n {
			fl = append(fl, f)
		}
	}
	return fl
}

func (v *Video) FilterQuality(quality string) {
	v.Formats = v.Formats.FilterQuality(quality)
	//v.AudioFormats = v.AudioFormats.FilterQuality(quality)
	//v.VideoFormats = v.VideoFormats.FilterQuality(quality)
	v.Formats.SortFormats()
}

// SortFormats sort all Formats fields
func (list FormatList) SortFormats() {
	sort.SliceStable(list, func(i, j int) bool {
		return sortFormat(i, j, list)
	})
}

// sortFormat sorts video by resolution, FPS, codec (av01, vp9, avc1), bitrate
// sorts audio by codec (mp4, opus), channels, bitrate, sample rate
func sortFormat(i int, j int, formats FormatList) bool {
	// Sort by Width
	if formats[i].Width == formats[j].Width {
		// Sort by FPS
		if formats[i].FPS == formats[j].FPS {
			if formats[i].FPS == 0 && formats[i].AudioChannels > 0 && formats[j].AudioChannels > 0 {
				// Audio
				// Sort by codec
				codec := map[int]int{}
				for _, index := range []int{i, j} {
					if strings.Contains(formats[index].MimeType, "mp4") {
						codec[index] = 1
					} else if strings.Contains(formats[index].MimeType, "opus") {
						codec[index] = 2
					}
				}
				if codec[i] == codec[j] {
					// Sort by Audio Channel
					if formats[i].AudioChannels == formats[j].AudioChannels {
						// Sort by Audio Bitrate
						if formats[i].Bitrate == formats[j].Bitrate {
							// Sort by Audio Sample Rate
							return formats[i].AudioSampleRate > formats[j].AudioSampleRate
						}
						return formats[i].Bitrate > formats[j].Bitrate
					}
					return formats[i].AudioChannels > formats[j].AudioChannels
				}
				return codec[i] < codec[j]
			}
			// Video
			// Sort by codec
			codec := map[int]int{}
			for _, index := range []int{i, j} {
				if strings.Contains(formats[index].MimeType, "av01") {
					codec[index] = 1
				} else if strings.Contains(formats[index].MimeType, "vp9") {
					codec[index] = 2
				} else if strings.Contains(formats[index].MimeType, "avc1") {
					codec[index] = 3
				}
			}
			if codec[i] == codec[j] {
				// Sort by Audio Bitrate
				return formats[i].Bitrate > formats[j].Bitrate
			}
			return codec[i] < codec[j]
		}
		return formats[i].FPS > formats[j].FPS
	}
	return formats[i].Width > formats[j].Width
}
