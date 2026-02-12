package youtube

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Cookies map[string]*http.Cookie

const cookieDomain = ".youtube.com"

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (c Cookies) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if !c.matchesDomain(u) {
		return
	}

	for _, cookie := range cookies {
		c.setCookie(cookie)
	}
}

func (c Cookies) setCookie(cookie *http.Cookie) {
	if cookie.Domain != cookieDomain || cookie.Path != "/" {
		return
	}

	if cookie.Expires.IsZero() {
		slog.Info("delete cookie", "name", cookie.Name)
		delete(c, cookie.Name)
	} else {
		slog.Info("set cookie", "name", cookie.Name)
		c[cookie.Name] = cookie
	}
}

func (c Cookies) matchesDomain(u *url.URL) bool {
	return strings.HasSuffix(u.Host, cookieDomain)
}

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (c Cookies) Cookies(u *url.URL) (result []*http.Cookie) {
	if !c.matchesDomain(u) {
		return nil
	}

	log.Println("asking for", u)
	for _, cookie := range c {
		result = append(result, cookie)
	}

	return
}

func readCookies(path string) (http.CookieJar, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jar := Cookies{}
	scanner := bufio.NewScanner(file)

	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, ".") {
			continue
		}

		fields := strings.Split(line, "\t")
		ts, _ := strconv.ParseInt(fields[4], 10, 64)

		if length := len(fields); length < 7 {
			return nil, fmt.Errorf("not enough fields in cookie file expected >= 7, is = %d", length)
		}

		jar.setCookie(&http.Cookie{
			Domain:  fields[0],
			Path:    fields[2],
			Expires: time.Unix(int64(ts), 0),
			Name:    fields[5],
			Value:   fields[6],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return jar, nil
}
