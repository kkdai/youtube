package youtube

import (
	"strings"

	sjson "github.com/bitly/go-simplejson"
)

type chunk struct {
	start int64
	end   int64
	data  chan []byte
}

func getChunks(totalSize, chunkSize int64) []chunk {
	var chunks []chunk

	if chunkSize <= 0 {
		chunkSize = Size2Mb
	}

	// googlevideo rejects a single request that covers the entire file.
	if totalSize > 1 && chunkSize >= totalSize {
		chunkSize = totalSize / 2
	}

	for start := int64(0); start < totalSize; start += chunkSize {
		end := min(chunkSize+start-1, totalSize-1)

		chunks = append(chunks, chunk{start, end, make(chan []byte, 1)})
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
		var text strings.Builder

		for i := 0; i < len(p.MustArray()); i++ {
			if textNode := p.GetIndex(i).Get("text"); isValidJSON(textNode) {
				text.WriteString(textNode.MustString())
			}
		}

		return text.String()
	}

	return ""
}

func getContinuation(j *sjson.Json) string {
	return j.GetPath("continuations").
		GetIndex(0).GetPath("nextContinuationData", "continuation").MustString()
}
