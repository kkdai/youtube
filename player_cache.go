package youtube

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
)

// CacheStorage should be a path to file, available to read and write
// By setting empty string as value you can disable backup-ing
var CacheStorage = filepath.Join(os.TempDir(), "youtubedr.json")

type playerCache struct {
	Version            string
	SignatureTimestamp string
	Operations         []Operation
}

type Operation struct {
	Type  int
	Value int
}

const (
	OpReverse = iota
	OpSplice
	OpSwap
)

func (c *playerCache) setSts(sts string) *playerCache {
	c.SignatureTimestamp = sts
	return c
}

func (c *playerCache) getSts(version string) (string, bool) {
	if c.Version != version {
		return "", false
	}
	return c.SignatureTimestamp, true
}

func (c *playerCache) addOps(op ...Operation) *playerCache {
	c.Operations = append(c.Operations, op...)
	return c
}

func (c *playerCache) getOps(version string) ([]Operation, bool) {
	if c.Version != version {
		return nil, false
	}
	return c.Operations, true
}

// loadLocalCache reads CacheStorage and parses it into Client.cache.
// If CacheStorage is empty - it does nothing
func (c *Client) loadLocalCache() bool {
	if CacheStorage == "" {
		return false
	}

	localCacheData, err := os.ReadFile(CacheStorage)
	if err != nil || len(localCacheData) <= 2 {
		return false
	}

	cache := new(playerCache)
	if err := json.Unmarshal(localCacheData, cache); err != nil {
		return false
	}

	c.cache = cache
	return true
}

// storeCacheLocally encodes Client.cache to json and writes to CacheStorage.
// If CacheStorage is empty - it does nothing
func storeCacheLocally(cache *playerCache) error {
	if CacheStorage == "" {
		return nil
	}

	cacheData, err := json.Marshal(*cache)
	if err != nil {
		return err
	}

	f, err := os.Create(CacheStorage)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	if _, err = f.Write(cacheData); err != nil {
		return err
	}
	return nil
}

// getSignatureTimestamp tries to get sts value from cache if its available.
// Otherwise, loads completely fresh player config and returns sts from it
func (c *Client) getSignatureTimestamp(ctx context.Context) (string, error) {
	// if there is no cache - try to load local one first
	if c.cache == nil {
		// result doesn't really matter at this point
		// because we check cache right after
		_ = c.loadLocalCache()
	}

	// if there is cache - we need to make sure its actual version
	if c.cache != nil {
		version, err := c.fetchPlayerVersion(ctx)
		if err != nil {
			// we couldn't fetch what's current player version, should we really exit here?..
			return "", err
		}

		sts, ok := c.cache.getSts(version)
		// if sts is empty means cache is invalid so we need to update it anyways
		if ok && sts != "" {
			return sts, nil
		}
	}

	// if we're here means either there is no cache and no backup
	// or current cache's version is outdated

	// so we need to update it
	_, err := c.cachePlayer(ctx)
	if err != nil {
		return "", err
	}

	// it should be fresh, so no need to check version
	return c.cache.SignatureTimestamp, nil
}

// getOperations tries to get deciphering operations from cache if its available.
// Otherwise, loads completely fresh player config and returns operations from it
func (c *Client) getOperations(ctx context.Context) ([]Operation, error) {
	// if there is no cache - try to load local one first
	if c.cache == nil {
		// result doesn't really matter at this point
		// because we check cache right after
		_ = c.loadLocalCache()
	}

	// if there is cache - we need to make sure its actual version
	if c.cache != nil {
		version, err := c.fetchPlayerVersion(ctx)
		if err != nil {
			// we couldn't fetch what's current player version, should we really exit here?..
			return nil, err
		}

		ops, ok := c.cache.getOps(version)
		// if sts is empty means cache is invalid so we need to update it anyways
		if ok && ops != nil {
			return ops, nil
		}
	}

	// if we're here means either there is no cache and no backup
	// or current cache's version is outdated

	// so we need to update it
	_, err := c.cachePlayer(ctx)
	if err != nil {
		return nil, err
	}

	// it should be fresh, so no need to check version
	return c.cache.Operations, nil
}
