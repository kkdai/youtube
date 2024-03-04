package youtube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type filter struct {
	Quality       string
	ItagNo        int
	MimeType      string
	Language      string
	AudioChannels int
}

func (list FormatList) Filter(filter filter) FormatList {
	if filter.ItagNo > 0 {
		list = list.Itag(filter.ItagNo)
	}
	if filter.AudioChannels > 0 {
		list = list.AudioChannels(filter.AudioChannels)
	}
	if filter.Quality != "" {
		list = list.Quality(filter.Quality)
	}
	if filter.MimeType != "" {
		list = list.Type(filter.MimeType)
	}
	if filter.Language != "" {
		list = list.Language(filter.Language)
	}
	return list
}

func TestFormatList_Filter(t *testing.T) {
	t.Parallel()

	format1 := Format{
		ItagNo:       1,
		Quality:      "medium",
		QualityLabel: "360p",
	}

	format2 := Format{
		ItagNo:        2,
		Quality:       "large",
		QualityLabel:  "480p",
		MimeType:      `video/mp4; codecs="avc1.42001E, mp4a.40.2"`,
		AudioChannels: 1,
	}

	formatStereo := Format{
		ItagNo:        3,
		URL:           "stereo",
		AudioChannels: 2,
	}

	list := FormatList{
		format1,
		format2,
		formatStereo,
	}

	tests := []struct {
		name string
		args filter
		want []Format
	}{
		{
			name: "empty list with quality small",
			args: filter{
				Quality: "small",
			},
		},
		{
			name: "empty list with other itagNo",
			args: filter{
				ItagNo: 99,
			},
		},
		{
			name: "empty list with other mimeType",
			args: filter{
				MimeType: "other",
			},
		},
		{
			name: "empty list with other audioChannels",
			args: filter{
				AudioChannels: 7,
			},
		},
		{
			name: "audioChannels stereo",
			args: filter{
				AudioChannels: formatStereo.AudioChannels,
			},
			want: []Format{formatStereo},
		},
		{
			name: "find by medium quality",
			args: filter{
				Quality: "medium",
			},
			want: []Format{format1},
		},
		{
			name: "find by 480p",
			args: filter{
				Quality: "480p",
			},
			want: []Format{format2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formats := list.Filter(tt.args)

			if tt.want == nil {
				assert.Empty(t, formats)
			} else {
				assert.Equal(t, tt.want, []Format(formats))
			}
		})
	}
}

func TestFormatList_Sort(t *testing.T) {
	t.Parallel()

	list := FormatList{
		{Width: 512},
		{Width: 768, MimeType: "mp4"},
		{Width: 768, MimeType: "opus"},
	}

	list.Sort()

	assert.Equal(t, FormatList{
		{Width: 768, MimeType: "mp4"},
		{Width: 768, MimeType: "opus"},
		{Width: 512},
	}, list)
}
