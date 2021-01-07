package youtube

import (
	"testing"
	"time"
)

func TestSimpleCache(t *testing.T) {
	type args struct {
		setVideoId string
		getVideoId string
		operations []operation
		expiredAt  string
		getCacheAt string
	}
	tests := []struct {
		name string
		args args
		want []operation
	}{
		{
			name: "Get cache data with video id",
			args: args{
				setVideoId: "test",
				getVideoId: "test",
				operations: []operation{func(bytes []byte) []byte { return nil }},
				expiredAt:  "2021-01-01 00:01:00",
				getCacheAt: "2021-01-01 00:00:00",
			},
			want: []operation{func(bytes []byte) []byte { return nil }},
		},
		{
			name: "Get nil when cache is expired",
			args: args{
				setVideoId: "test",
				getVideoId: "test",
				operations: []operation{func(bytes []byte) []byte { return nil }},
				expiredAt:  "2021-01-01 00:00:00",
				getCacheAt: "2021-01-01 00:00:00",
			},
			want: nil,
		},
		{
			name: "Get nil when video id is not cached",
			args: args{
				setVideoId: "test",
				getVideoId: "not test",
				operations: []operation{func(bytes []byte) []byte { return nil }},
				expiredAt:  "2021-01-01 00:00:01",
				getCacheAt: "2021-01-01 00:00:00",
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSimpleCache()
			timeFormat := "2006-01-02 15:04:05"
			expiredAt, _ := time.Parse(timeFormat, tt.args.expiredAt)
			s.setWithExpiredTime(tt.args.setVideoId, tt.args.operations, expiredAt)
			getCacheAt, _ := time.Parse(timeFormat, tt.args.getCacheAt)
			if got := s.GetCacheBefore(tt.args.getVideoId, getCacheAt); len(got) != len(tt.want) {
				t.Errorf("GetCacheBefore() = %v, want %v", got, tt.want)
			}
		})
	}
}
