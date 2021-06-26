package youtube

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
)

var playerVersionRegexp = regexp.MustCompile(`(?:/s/player/)(\w+)/`)

func (c *Client) fetchPlayerVersion(ctx context.Context) (string, error) {
	// url without videoID works as well and weight slightly less
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

func (c *Client) getPlayerConfig(ctx context.Context, version string) ([]byte, error) {
	// example: /s/player/f676c671/player_ias.vflset/en_US/base.js
	playerURL := fmt.Sprintf("https://www.youtube.com/s/player/%s/player_ias.vflset/en_US/base.js", version)
	return c.httpGetBodyBytes(ctx, playerURL)
}

func (c *Client) cachePlayer(ctx context.Context) (*PlayerCache, error) {
	version, err := c.fetchPlayerVersion(ctx)
	if err != nil || version == "" {
		return nil, err
	}

	c.cache = &PlayerCache{Version: version}

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

	return c.cache.setSts(sts).addOps(ops...), nil
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

func (c *Client) getSignatureTimestamp(ctx context.Context) (string, error) {
	// if there is cache - we need to make sure its actual version
	if c.cache != nil {
		version, err := c.fetchPlayerVersion(ctx)
		if err != nil {
			// we couldn't fetch what's current player version, should we really exit here?..
			return "", err
		}

		sts, ok := c.cache.getSts(version)
		if ok {
			return sts, nil
		}
	}

	// both empty cache and wrong version leads here
	_, err := c.cachePlayer(ctx)
	if err != nil {
		return "", err
	}

	// we just re-cached player, no need to check version
	return c.cache.SignatureTimestamp, nil
}
