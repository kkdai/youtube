Download Youtube Video in Golang
==================

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/kkdai/youtube/master/LICENSE)  [![GoDoc](https://godoc.org/github.com/kkdai/youtube?status.svg)](https://godoc.org/github.com/kkdai/youtube)  [![Build Status](https://travis-ci.org/kkdai/youtube.svg?branch=master)](https://travis-ci.org/kkdai/youtube) [![](https://goreportcard.com/badge/github.com/kkdai/youtube)](https://goreportcard.com/badge/github.com/kkdai/youtube)


This package is a Youtube video download package, for more detail refer [https://github.com/rg3/youtube-dl](https://github.com/rg3/youtube-dl) for more download options.


## Overview
  * [Install](#install)
  * [Usage](#usage)
  * [Options](#options)
  * [Example: Download video from \[dotGo 2015 - Rob Pike - Simplicity is Complicated\]](#download-dotGo-2015-rob-pike-video)

## Install:
```shell
go get github.com/kkdai/youtube
```

OR

```shell
git clone https://github.com/kkdai/youtube.git
go run youtubedr/main.go
```

## Usage

### Use the binary directly
It's really simple to use, just get the video id from youtube url - ex: `https://www.youtube.com/watch?v=rFejpH_tAHM`, the video id is `rFejpH_tAHM`

```shell
$ youtubedr QAGDGja7kbs
$ youtubedr https://www.youtube.com/watch?v=rFejpH_tAHM
```

### Import this package in your golang program

```go
package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"
	"path/filepath"

	. "github.com/kkdai/youtube"
)

func main() {
	flag.Parse()
	log.Println(flag.Args())
	usr, _ := user.Current()
	currentDir := fmt.Sprintf("%v/Movies/youtubedr", usr.HomeDir)
	log.Println("download to dir=", currentDir)
	y := NewYoutube(true)
	arg := flag.Arg(0)
	if err := y.DecodeURL(arg); err != nil {
		fmt.Println("err:", err)
	}
	if err := y.StartDownload(filepath.Join(currentDir, "dl.mp4")); err != nil {
		fmt.Println("err:", err)
	}
}
```

## Options:

| option | type   | description                                                    | default value          |
| :----- | :----- | :------------------------------------------------------------- | :--------------------- |
| `-d`   | string | the output directory                                           | $HOME/Movies/youtubedr |
| `-o`   | string | the output file name ( ext will auto detect on default value ) | [video's title].ext    |
| `-d`   | string | the Socks 5 proxy (e.g. 10.10.10.10:7878)                      |                        |
| `-q`   | string | the output file quality (medium, hd720)                        |                        |

## Example:
 * ### download-dotGo-2015-rob-pike-video

    `go get github.com/kkdai/youtube/youtubedr`

    Download video from [dotGo 2015 - Rob Pike - Simplicity is Complicated](https://www.youtube.com/watch?v=rFejpH_tAHM)

    ```
    youtubedr https://www.youtube.com/watch?v=rFejpH_tAHM
    ```

 * ### Download video to specific folder and name

	`go get github.com/kkdai/youtube/youtubedr`

	Download video from [dotGo 2015 - Rob Pike - Simplicity is Complicated](https://www.youtube.com/watch?v=rFejpH_tAHM) to current directory and name the file to simplicity-is-complicated.mp4

	```
	youtubedr -d ./ -o simplicity-is-complicated.mp4 https://www.youtube.com/watch?v=rFejpH_tAHM
	```

 * ### Download video with specific quality

	`go get github.com/kkdai/youtube/youtubedr`

	Download video from [dotGo 2015 - Rob Pike - Simplicity is Complicated](https://www.youtube.com/watch?v=rFejpH_tAHM) with specific quality

	```
	youtubedr -q medium https://www.youtube.com/watch?v=rFejpH_tAHM
	```


## How it works

- Parse the video ID you input in URL
	- ex: `https://www.youtube.com/watch?v=rFejpH_tAHM`, the video id is `rFejpH_tAHM`
- Get video information via video id.
	- Use URL: `http://youtube.com/get_video_info?video_id=`
- Parse and decode video information.
	- Download URL in "url="
	- title in "title="
- Download video from URL
	- Need the string combination of "url"

## Inspired
- [https://github.com/ytdl-org/youtube-dl](https://github.com/ytdl-org/youtube-dl)
- [https://github.com/lepidosteus/youtube-dl](https://github.com/lepidosteus/youtube-dl)
- [拆解 Youtube 影片下載位置](http://hkgoldenmra.blogspot.tw/2013/05/youtube.html)
- [iawia002/annie](https://github.com/iawia002/annie)
- [How to get url from obfuscate video info: youtube video downloader with php](https://stackoverflow.com/questions/60607291/youtube-video-downloader-with-php)


## Project52
It is one of my [project 52](https://github.com/kkdai/project52).


## License
This package is licensed under MIT license. See LICENSE for details.
