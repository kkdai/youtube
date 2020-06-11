// +build integration

package youtube

import (
	"reflect"
	"testing"
)

func TestYoutube_GetItagInfo(t *testing.T) {
	y := NewYoutube(false)
	if y == nil {
		t.Error("Cannot init object.")
		return
	}

	// url from issue #25
	testVideoUrl := "https://www.youtube.com/watch?v=rFejpH_tAHM"
	expected := ItagInfo{
		Title:  "dotGo 2015 - Rob Pike - Simplicity is Complicated",
		Author: "dotconferences",
		Itags: []Itag{
			{ItagNo: 18, Quality: "medium", Type: `video/mp4; codecs="avc1.42001E, mp4a.40.2"`},
			{ItagNo: 22, Quality: "hd720", Type: `video/mp4; codecs="avc1.64001F, mp4a.40.2"`},
		},
	}
	if err := y.DecodeURL(testVideoUrl); err != nil {
		t.Error("Cannot decode download url")
		return
	}
	itagInfo := y.GetItagInfo()
	if !reflect.DeepEqual(*itagInfo, expected) {
		t.Errorf("get Itag info failed")
	}
}
