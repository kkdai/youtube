package youtube

import "time"

var (
	_ DecipherOperationsCache = NewSimpleCache()
)

const defaultCacheExpiration = time.Minute * time.Duration(5)

type DecipherOperationsCache interface {
	Get(videoId string) []operation
	Set(video string, operations []operation)
}

type SimpleCache struct {
	videoID    string
	expiredAt  time.Time
	operations []operation
}

func NewSimpleCache() *SimpleCache {
	return &SimpleCache{}
}

// Get : get cache  when it has same video id and not expired
func (s SimpleCache) Get(videoId string) []operation {
	return s.GetCacheBefore(videoId, time.Now())
}

// GetCacheBefore : can pass time for testing
func (s SimpleCache) GetCacheBefore(videoId string, time time.Time) []operation {
	if videoId == s.videoID && s.expiredAt.After(time) {
		operations := make([]operation, len(s.operations))
		copy(operations, s.operations)
		return operations
	}
	return nil
}

// Set : set cache with default expiration
func (s *SimpleCache) Set(videoId string, operations []operation) {
	s.setWithExpiredTime(videoId, operations, time.Now().Add(defaultCacheExpiration))
}

func (s *SimpleCache) setWithExpiredTime(videoId string, operations []operation, time time.Time) {
	s.videoID = videoId
	s.operations = make([]operation, len(operations))
	copy(s.operations, operations)
	s.expiredAt = time
}
