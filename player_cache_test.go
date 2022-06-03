package youtube

import (
	"testing"
	"time"
)

func TestPlayerCache(t *testing.T) {
	type args struct {
		setVideoID string
		getVideoID string
		expiredAt  string
		getCacheAt string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Get cache data with video id",
			args: args{
				setVideoID: "test",
				getVideoID: "test",
				expiredAt:  "2021-01-01 00:01:00",
				getCacheAt: "2021-01-01 00:00:00",
			},
			want: []byte("playerdata"),
		},
		{
			name: "Get nil when cache is expired",
			args: args{
				setVideoID: "test",
				getVideoID: "test",
				expiredAt:  "2021-01-01 00:00:00",
				getCacheAt: "2021-01-01 00:00:00",
			},
			want: nil,
		},
		{
			name: "Get nil when video id is not cached",
			args: args{
				setVideoID: "test",
				getVideoID: "not test",
				expiredAt:  "2021-01-01 00:00:01",
				getCacheAt: "2021-01-01 00:00:00",
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := playerCache{}
			timeFormat := "2006-01-02 15:04:05"
			expiredAt, _ := time.Parse(timeFormat, tt.args.expiredAt)
			s.setWithExpiredTime(tt.args.setVideoID, []byte("playerdata"), expiredAt)
			getCacheAt, _ := time.Parse(timeFormat, tt.args.getCacheAt)
			if got := s.GetCacheBefore(tt.args.getVideoID, getCacheAt); len(got) != len(tt.want) {
				t.Errorf("GetCacheBefore() = %v, want %v", got, tt.want)
			}
		})
	}
}
