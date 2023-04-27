package youtube

import (
	"encoding/base64"
	"math"

	sjson "github.com/bitly/go-simplejson"
)

type chunk struct {
	index int
	start int64
	end   int64
}

func getChunks(totalSize, chunkSize int64) []chunk {
	var chunks []chunk

	for i := 0; i < int(math.Ceil(float64(totalSize)/float64(chunkSize))); i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize - 1
		if end >= totalSize {
			end = totalSize - 1
		}

		chunks = append(chunks, chunk{i, start, end})
	}

	return chunks
}

func getFirstKeyJSON(j *sjson.Json) *sjson.Json {
	m, err := j.Map()
	if err != nil {
		return j
	}

	for key := range m {
		return j.Get(key)
	}

	return j
}

func isValidJSON(j *sjson.Json) bool {
	b, err := j.MarshalJSON()
	if err != nil {
		return false
	}

	if len(b) <= 4 {
		return false
	}

	return true
}

func sjsonGetText(j *sjson.Json, paths ...string) string {
	for _, path := range paths {
		if isValidJSON(j.Get(path)) {
			j = j.Get(path)
		}
	}

	if text, err := j.String(); err == nil {
		return text
	}

	if isValidJSON(j.Get("text")) {
		return j.Get("text").MustString()
	}

	if p := j.Get("runs"); isValidJSON(p) {
		var text string

		for i := 0; i < len(p.MustArray()); i++ {
			if textNode := p.GetIndex(i).Get("text"); isValidJSON(textNode) {
				text += textNode.MustString()
			}
		}

		return text
	}

	return ""
}

func getContinuation(j *sjson.Json) string {
	return j.GetPath("continuations").
		GetIndex(0).GetPath("nextContinuationData", "continuation").MustString()
}

func base64PadEnc(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func base64Enc(str string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(str))
}
