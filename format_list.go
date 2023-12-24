package youtube

import (
	"sort"
	"strconv"
	"strings"
)

type FormatList []Format

// Type returns a new FormatList filtered by itag
func (list FormatList) Select(f func(Format) bool) (result FormatList) {
	for i := range list {
		if f(list[i]) {
			result = append(result, list[i])
		}
	}
	return result
}

// Type returns a new FormatList filtered by itag
func (list FormatList) Itag(itagNo int) FormatList {
	return list.Select(func(f Format) bool {
		return f.ItagNo == itagNo
	})
}

// Type returns a new FormatList filtered by mime type
func (list FormatList) Type(value string) FormatList {
	return list.Select(func(f Format) bool {
		return strings.Contains(f.MimeType, value)
	})
}

// Type returns a new FormatList filtered by display name
func (list FormatList) Language(displayName string) FormatList {
	return list.Select(func(f Format) bool {
		return f.LanguageDisplayName() == displayName
	})
}

// Quality returns a new FormatList filtered by quality, quality label or itag,
// but not audio quality
func (list FormatList) Quality(quality string) FormatList {
	itag, _ := strconv.Atoi(quality)

	return list.Select(func(f Format) bool {
		return itag == f.ItagNo || strings.Contains(f.Quality, quality) || strings.Contains(f.QualityLabel, quality)
	})
}

// AudioChannels returns a new FormatList filtered by the matching AudioChannels
func (list FormatList) AudioChannels(n int) FormatList {
	return list.Select(func(f Format) bool {
		return f.AudioChannels == n
	})
}

// AudioChannels returns a new FormatList filtered by the matching AudioChannels
func (list FormatList) WithAudioChannels() FormatList {
	return list.Select(func(f Format) bool {
		return f.AudioChannels > 0
	})
}

// FilterQuality reduces the format list to formats matching the quality
func (v *Video) FilterQuality(quality string) {
	v.Formats = v.Formats.Quality(quality)
	v.Formats.Sort()
}

// Sort sorts all formats fields
func (list FormatList) Sort() {
	sort.SliceStable(list, func(i, j int) bool {
		return sortFormat(i, j, list)
	})
}

// sortFormat sorts video by resolution, FPS, codec (av01, vp9, avc1), bitrate
// sorts audio by default, codec (mp4, opus), channels, bitrate, sample rate
func sortFormat(i int, j int, formats FormatList) bool {

	// Sort by Width
	if formats[i].Width == formats[j].Width {
		// Format 137 downloads slowly, give it less priority
		// see https://github.com/kkdai/youtube/pull/171
		switch 137 {
		case formats[i].ItagNo:
			return false
		case formats[j].ItagNo:
			return true
		}

		// Sort by FPS
		if formats[i].FPS == formats[j].FPS {
			if formats[i].FPS == 0 && formats[i].AudioChannels > 0 && formats[j].AudioChannels > 0 {
				// Audio
				// Sort by default
				if (formats[i].AudioTrack == nil && formats[j].AudioTrack == nil) || (formats[i].AudioTrack != nil && formats[j].AudioTrack != nil && formats[i].AudioTrack.AudioIsDefault == formats[j].AudioTrack.AudioIsDefault) {
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
				} else if formats[i].AudioTrack != nil && formats[i].AudioTrack.AudioIsDefault {
					return true
				}
				return false
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
