package youtube

import (
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
func (list FormatList) FindByType(t string) []Format {
	var f []Format
	for i := range list {
		if strings.Contains(list[i].MimeType, t) {
			f = append(f, list[i])
		}
	}
	return f
}
