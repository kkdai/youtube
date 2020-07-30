package downloader

type progress struct {
	contentLength     float64
	totalWrittenBytes float64
	downloadLevel     float64
}

func (dl *progress) Write(p []byte) (n int, err error) {
	n = len(p)
	dl.totalWrittenBytes = dl.totalWrittenBytes + float64(n)
	currentPercent := (dl.totalWrittenBytes / dl.contentLength) * 100
	if (dl.downloadLevel <= currentPercent) && (dl.downloadLevel < 100) {
		dl.downloadLevel++
	}
	return
}
