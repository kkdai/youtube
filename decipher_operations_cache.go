package youtube

import "time"

var (
	_ DecipherOperationsCache = NewSimpleCache()
)

const defaultCacheExpiration = time.Minute * time.Duration(5)

type DecipherOperationsCache interface {
	Get(videoID string) []DecipherOperation
	Set(video string, operations []DecipherOperation)
}

type SimpleCache struct {
	videoID    string
	expiredAt  time.Time
	operations []DecipherOperation
}

func NewSimpleCache() *SimpleCache {
	return &SimpleCache{}
}

// Get : get cache  when it has same video id and not expired
func (s SimpleCache) Get(videoID string) []DecipherOperation {
	return s.GetCacheBefore(videoID, time.Now())
}

// GetCacheBefore : can pass time for testing
func (s SimpleCache) GetCacheBefore(videoID string, time time.Time) []DecipherOperation {
	if videoID == s.videoID && s.expiredAt.After(time) {
		operations := make([]DecipherOperation, len(s.operations))
		copy(operations, s.operations)
		return operations
	}
	return nil
}

// Set : set cache with default expiration
func (s *SimpleCache) Set(videoID string, operations []DecipherOperation) {
	s.setWithExpiredTime(videoID, operations, time.Now().Add(defaultCacheExpiration))
}

func (s *SimpleCache) setWithExpiredTime(videoID string, operations []DecipherOperation, time time.Time) {
	s.videoID = videoID
	s.operations = make([]DecipherOperation, len(operations))
	copy(s.operations, operations)
	s.expiredAt = time
}
