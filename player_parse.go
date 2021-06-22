package youtube

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
)

var basejsPattern = regexp.MustCompile(`(/s/player/\w+/player_ias.vflset/\w+/base.js)`)

// we may use \d{5} instead of \d+ since currently its 5 digits, but i can't be sure it will be 5 digits always
var signatureRegexp = regexp.MustCompile(`(?m)(?:^|,)(?:signatureTimestamp:)(\d+)`)

func (c *Client) fetchPlayerConfig(ctx context.Context, videoID string) ([]byte, error) {
	embedURL := fmt.Sprintf("https://youtube.com/embed/%s?hl=en", videoID)
	embedBody, err := c.httpGetBodyBytes(ctx, embedURL)
	if err != nil {
		return nil, err
	}

	// example: /s/player/f676c671/player_ias.vflset/en_US/base.js
	escapedBasejsURL := string(basejsPattern.Find(embedBody))
	if escapedBasejsURL == "" {
		log.Println("playerConfig:", string(embedBody))
		return nil, errors.New("unable to find basejs URL in playerConfig")
	}

	return c.httpGetBodyBytes(ctx, "https://youtube.com"+escapedBasejsURL)
}

func (c *Client) getSignatureTimestamp(ctx context.Context, videoID string) (string, error) {
	basejsBody, err := c.fetchPlayerConfig(ctx, videoID)
	if err != nil {
		return "", err
	}

	result := signatureRegexp.FindSubmatch(basejsBody)
	if result == nil {
		return "", ErrSignatureTimestampNotFound
	}

	return string(result[1]), nil
}
