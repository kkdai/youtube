package youtube

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func parseM3U8(r io.Reader) ([]Format, error) {
	var result []Format

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#EXT-X-STREAM-INF:") {
			continue
		}

		// TODO parse line
		scanner.Scan()
		url := scanner.Text()

		itag := extractURLcomponent(url, "itag")
		itagNo, _ := strconv.Atoi(itag)

		result = append(result, Format{
			URL:    url,
			ItagNo: itagNo,
		})
	}

	return result, scanner.Err()
}

func extractURLcomponent(url, arg string) string {
	i := strings.Index(url, "/"+arg+"/")
	if i < 0 {
		return ""
	}
	i += len(arg) + 2

	j := strings.Index(url[i:], "/")
	if j < 0 {
		return ""
	}

	return url[i : i+j]
}
