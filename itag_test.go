// +build integration

package youtube

import (
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
	if err := y.DecodeURL(testVideoUrl); err != nil {
		t.Error("Cannot decode download url")
		return
	}
	itagInfo := y.GetItagInfo()
	cnt := 18
	gotCnt := len(itagInfo.Itags)
	if gotCnt != cnt {
		t.Errorf("get Itag info failed, want %v, but get %v", cnt, gotCnt)
	}
}
