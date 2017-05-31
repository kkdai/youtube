package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	. "github.com/kkdai/youtube"
)

func main() {
	flag.Parse()
	log.Println(flag.Args())
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	log.Println("download to dir=", currentDir)
	y := NewYoutube(true)
	arg := flag.Arg(0)
	y.DecodeURL(arg)
	y.StartDownload(currentDir)
}
