package youtube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatList_FindByQuality(t *testing.T) {
	list := []Format{{
		ItagNo:       0,
		Quality:      "medium",
		QualityLabel: "360p",
	},
		{
			ItagNo:       1,
			Quality:      "large",
			QualityLabel: "480p",
		},
	}
	type args struct {
		quality string
	}
	tests := []struct {
		name string
		list FormatList
		args args
		want *Format
	}{
		{
			name: "find by quality, get correct one",
			list: list,
			args: args{
				quality: "medium",
			},
			want: &Format{
				ItagNo:       0,
				Quality:      "medium",
				QualityLabel: "360p",
			},
		},
		{
			name: "find by quality label, get correct one",
			list: list,
			args: args{
				quality: "480p",
			},
			want: &Format{
				ItagNo:       1,
				Quality:      "large",
				QualityLabel: "480p",
			},
		},
		{
			name: "find nothing with quality",
			list: list,
			args: args{
				quality: "small",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := tt.list.FindByQuality(tt.args.quality)
			assert.Equal(t, format, tt.want)
		})
	}
}

func TestFormatList_FindByItag(t *testing.T) {
	list := []Format{{
		ItagNo: 18,
	},
		{
			ItagNo: 135,
		},
	}
	type args struct {
		itagNo int
	}
	tests := []struct {
		name string
		list FormatList
		args args
		want *Format
	}{
		{
			name: "find itag 18",
			list: list,
			args: args{
				itagNo: 18,
			},
			want: &Format{
				ItagNo: 18,
			},
		},
		{
			name: "find itag 135",
			list: list,
			args: args{
				itagNo: 135,
			},
			want: &Format{
				ItagNo: 135,
			},
		},
		{
			name: "find nothing",
			list: list,
			args: args{
				itagNo: 9999,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := tt.list.FindByItag(tt.args.itagNo)
			assert.Equal(t, format, tt.want)
		})
	}
}

func TestFormatList_Type(t *testing.T) {
	list := []Format{{
		MimeType: "video/mp4; codecs=\"avc1.42001E, mp4a.40.2\"",
	},
	}
	type args struct {
		mimeType string
	}
	tests := []struct {
		name string
		list FormatList
		args args
		want FormatList
	}{
		{
			name: "find video",
			list: list,
			args: args{
				mimeType: "video/mp4; codecs=\"avc1.42001E, mp4a.40.2\"",
			},
			want: []Format{{
				MimeType: "video/mp4; codecs=\"avc1.42001E, mp4a.40.2\"",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := tt.list.Type("video")
			assert.Equal(t, format, tt.want)
		})
	}
}
