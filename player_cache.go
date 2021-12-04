package youtube

import "time"

const defaultCacheExpiration = time.Minute * time.Duration(5)

type playerCache struct {
	playerID  string
	expiredAt time.Time
	config    playerConfig
}

// Get : get cache  when it has same video id and not expired
func (s playerCache) Get(playerID string) playerConfig {
	return s.GetCacheBefore(playerID, time.Now())
}

// GetCacheBefore : can pass time for testing
func (s playerCache) GetCacheBefore(playerID string, time time.Time) playerConfig {
	if playerID == s.playerID && s.expiredAt.After(time) {
		return s.config
	}
	return nil
}

// Set : set cache with default expiration
func (s *playerCache) Set(playerID string, operations playerConfig) {
	s.setWithExpiredTime(playerID, operations, time.Now().Add(defaultCacheExpiration))
}

func (s *playerCache) setWithExpiredTime(playerID string, config playerConfig, time time.Time) {
	s.playerID = playerID
	s.config = config
	s.expiredAt = time
}
