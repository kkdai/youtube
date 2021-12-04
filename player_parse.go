package youtube

import (
	"context"
	"errors"
	"fmt"
	"regexp"
)

type playerConfig []byte

var basejsPattern = regexp.MustCompile(`(/s/player/\w+/player_ias.vflset/\w+/base.js)`)

// we may use \d{5} instead of \d+ since currently its 5 digits, but i can't be sure it will be 5 digits always
var signatureRegexp = regexp.MustCompile(`(?m)(?:^|,)(?:signatureTimestamp:)(\d+)`)

func (c *Client) getPlayerConfig(ctx context.Context, videoID string) (playerConfig, error) {

	embedURL := fmt.Sprintf("https://youtube.com/embed/%s?hl=en", videoID)
	embedBody, err := c.httpGetBodyBytes(ctx, embedURL)
	if err != nil {
		return nil, err
	}

	// example: /s/player/f676c671/player_ias.vflset/en_US/base.js
	escapedBasejsURL := string(basejsPattern.Find(embedBody))
	if escapedBasejsURL == "" {
		return nil, errors.New("unable to find basejs URL in playerConfig")
	}

	config := c.playerCache.Get(escapedBasejsURL)
	if config != nil {
		return config, nil
	}

	config, err = c.httpGetBodyBytes(ctx, "https://youtube.com"+escapedBasejsURL)
	if err != nil {
		return nil, err
	}

	c.playerCache.Set(escapedBasejsURL, config)
	return config, nil
}
