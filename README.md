Download Youtube Video in Golang
==================

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/kkdai/youtube/master/LICENSE)  [![GoDoc](https://godoc.org/github.com/kkdai/youtube?status.svg)](https://godoc.org/github.com/kkdai/youtube)  [![Build Status](https://travis-ci.org/kkdai/youtube.svg?branch=master)](https://travis-ci.org/kkdai/youtube) [![](https://goreportcard.com/badge/github.com/kkdai/youtube)](https://goreportcard.com/badge/github.com/kkdai/youtube)



This package is a Youtube video download package, for more detail refer [https://github.com/rg3/youtube-dl](https://github.com/rg3/youtube-dl) for more download option.


How it works
---------------

- Parse the video ID you input in URL
	- ex: `https://www.youtube.com/watch?v=rFejpH_tAHM`, the video id is `rFejpH_tAHM`
- Get video information via video id.
	- Use URL: `http://youtube.com/get_video_info?video_id=`
- Parse and decode video information.
	- Download URL in "url="
	- title in "title="
	- Need sinature in "sig="
- Download video from URL
	- Need the string combination of "url+sig"
 
Install
---------------
`go get github.com/kkdai/youtube`


Usage
---------------

```go

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	. "github.com/kkdai/youtube"
)

func main() {
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	log.Println("download to dir=", currentDir)
	y := NewYoutube()
	y.DecodeURL("https://www.youtube.com/watch?v=rFejpH_tAHM")
	y.StartDownload(currentDir)
}
```

Use the binary directly
---------------
`go get github.com/kkdai/youtube/youtubedr`

Download video from [dotGo 2015 - Rob Pike - Simplicity is Complicated](https://www.youtube.com/watch?v=rFejpH_tAHM)

```
youtubedr https://www.youtube.com/watch?v=rFejpH_tAHM
```


Inspired
---------------

- [https://github.com/lepidosteus/youtube-dl](https://github.com/lepidosteus/youtube-dl)
- [拆解 Youtube 影片下載位置](http://hkgoldenmra.blogspot.tw/2013/05/youtube.html)

Project52
---------------

It is one of my [project 52](https://github.com/kkdai/project52).


License
---------------

This package is licensed under MIT license. See LICENSE for details.

