package youtube

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
)

var playerVersionRegexp = regexp.MustCompile(`(?:/s/player/)(\w+)/`)

// fetchPlayerVersion looks at youtube for actual version of player
func (c *Client) fetchPlayerVersion(ctx context.Context) (string, error) {
	// embed url without videoID works as well and weights slightly less
	embedBody, err := c.httpGetBodyBytes(ctx, "https://www.youtube.com/embed")
	if err != nil {
		return "", err
	}

	playerVersionData := playerVersionRegexp.FindSubmatch(embedBody)
	if playerVersionData == nil {
		if c.Debug {
			log.Println("playerConfig: ", string(embedBody))
		}
		return "", errors.New("unable to find basejs URL in playerConfig")
	}

	return string(playerVersionData[1]), nil
}

// getPlayerConfig generates player config url for given version and fetches it's content
func (c *Client) getPlayerConfig(ctx context.Context, version string) ([]byte, error) {
	// example: /s/player/f676c671/player_ias.vflset/en_US/base.js
	playerURL := fmt.Sprintf("https://www.youtube.com/s/player/%s/player_ias.vflset/en_US/base.js", version)
	return c.httpGetBodyBytes(ctx, playerURL)
}

// cachePlayer fetches new player config and caches it.
// If CacheStorage is set - it also tries to backup cache locally for faster access
func (c *Client) cachePlayer(ctx context.Context) (*playerCache, error) {
	version, err := c.fetchPlayerVersion(ctx)
	if err != nil || version == "" {
		if c.Debug {
			log.Printf("fetchPlayerVersion failed: version: %s, err: %s\n", version, err)
		}
		return nil, err
	}

	//c.cache = &PlayerCache{version: version}
	newCache := &playerCache{Version: version}

	player, err := c.getPlayerConfig(ctx, version)
	if err != nil {
		return nil, err
	}

	sts, err := parseSts(player)
	if err != nil {
		return nil, err
	}

	ops, err := parseOperations(player)
	if err != nil {
		return nil, err
	}

	c.cache = newCache.setSts(sts).addOps(ops...)

	// each time we load new cache - try to store it locally
	if err := storeCacheLocally(c.cache); err != nil {
		if c.Debug {
			log.Printf("failed to store cache locally: %s", err)
		}
	}

	return c.cache, nil
}

// we may use \d{5} instead of \d+ since currently its 5 digits, but i can't be sure it will be 5 digits always
var signatureRegexp = regexp.MustCompile(`(?m)(?:^|,)(?:signatureTimestamp:)(\d+)`)

func parseSts(player []byte) (string, error) {
	result := signatureRegexp.FindSubmatch(player)
	if result == nil {
		return "", ErrSignatureTimestampNotFound
	}

	return string(result[1]), nil
}
