package youtube

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func (y *Client) decipherURL(ctx context.Context, videoId string, cipher string) (string, error) {
	queryParams, err := url.ParseQuery(cipher)
	if err != nil {
		return "", err
	}

	/* eg:
	    extract decipher from  https://youtube.com/s/player/4fbb4d5b/player_ias.vflset/en_US/base.js

	    var Mt={
		splice:function(a,b){a.splice(0,b)},
		reverse:function(a){a.reverse()},
		EQ:function(a,b){var c=a[0];a[0]=a[b%a.length];a[b%a.length]=c}};

		a=a.split("");
		Mt.splice(a,3);
		Mt.EQ(a,39);
		Mt.splice(a,2);
		Mt.EQ(a,1);
		Mt.splice(a,1);
		Mt.EQ(a,35);
		Mt.EQ(a,51);
		Mt.splice(a,2);
		Mt.reverse(a,52);
		return a.join("")
	*/

	operations, err := y.parseDecipherOps(ctx, videoId)
	if err != nil {
		return "", err
	}

	// apply operations
	bs := []byte(queryParams.Get("s"))
	for _, op := range operations {
		bs = op(bs)
	}

	decipheredUrl := fmt.Sprintf("%s&%s=%s", queryParams.Get("url"), queryParams.Get("sp"), string(bs))
	return decipheredUrl, nil
}

const (
	jsvarStr   = "[a-zA-Z_\\$][a-zA-Z_0-9]*"
	reverseStr = ":function\\(a\\)\\{" +
		"(?:return )?a\\.reverse\\(\\)" +
		"\\}"
	spliceStr = ":function\\(a,b\\)\\{" +
		"a\\.splice\\(0,b\\)" +
		"\\}"
	swapStr = ":function\\(a,b\\)\\{" +
		"var c=a\\[0\\];a\\[0\\]=a\\[b(?:%a\\.length)?\\];a\\[b(?:%a\\.length)?\\]=c(?:;return a)?" +
		"\\}"
)

var (
	playerConfigPattern = regexp.MustCompile(`yt\.setConfig\({'PLAYER_CONFIG':(.*)}\);`)
	basejsPattern       = regexp.MustCompile(`"js":"\\/s\\/player(.*)base\.js`)

	actionsObjRegexp = regexp.MustCompile(fmt.Sprintf(
		"var (%s)=\\{((?:(?:%s%s|%s%s|%s%s),?\\n?)+)\\};", jsvarStr, jsvarStr, reverseStr, jsvarStr, spliceStr, jsvarStr, swapStr))

	actionsFuncRegexp = regexp.MustCompile(fmt.Sprintf(
		"function(?: %s)?\\(a\\)\\{"+
			"a=a\\.split\\(\"\"\\);\\s*"+
			"((?:(?:a=)?%s\\.%s\\(a,\\d+\\);)+)"+
			"return a\\.join\\(\"\"\\)"+
			"\\}", jsvarStr, jsvarStr, jsvarStr))

	reverseRegexp = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, reverseStr))
	spliceRegexp  = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, spliceStr))
	swapRegexp    = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, swapStr))
)

func (y *Client) parseDecipherOps(ctx context.Context, videoId string) (operations []operation, err error) {
	if videoId == "" {
		return nil, errors.New("video id is empty")
	}

	embedUrl := fmt.Sprintf("https://youtube.com/embed/%s?hl=en", videoId)
	embeddedPageResp, err := y.httpGet(ctx, embedUrl)
	if err != nil {
		return nil, err
	}
	defer embeddedPageResp.Body.Close()

	embeddedPageBodyBytes, err := ioutil.ReadAll(embeddedPageResp.Body)
	if err != nil {
		return nil, err
	}

	playerConfig := playerConfigPattern.Find(embeddedPageBodyBytes)

	// eg: "js":\"\/s\/player\/f676c671\/player_ias.vflset\/en_US\/base.js
	escapedBasejsUrl := string(basejsPattern.Find(playerConfig))
	// eg: ["js", "\/s\/player\/f676c671\/player_ias.vflset\/en_US\/base.js]
	arr := strings.Split(escapedBasejsUrl, ":\"")
	basejsUrl := "https://youtube.com" + strings.ReplaceAll(arr[len(arr)-1], "\\", "")
	basejsUrlResp, err := y.httpGet(ctx, basejsUrl)
	if err != nil {
		return nil, err
	}
	defer basejsUrlResp.Body.Close()

	basejsBodyBytes, err := ioutil.ReadAll(basejsUrlResp.Body)
	if err != nil {
		return nil, err
	}

	bodyString := string(basejsBodyBytes)
	objResult := actionsObjRegexp.FindStringSubmatch(bodyString)
	funcResult := actionsFuncRegexp.FindStringSubmatch(bodyString)
	if len(objResult) < 3 || len(funcResult) < 2 {
		return nil, errors.New("error parsing signature tokens")
	}

	obj := objResult[1]
	objBody := objResult[2]
	funcBody := funcResult[1]

	var reverseKey, spliceKey, swapKey string

	if result := reverseRegexp.FindStringSubmatch(objBody); len(result) > 1 {
		reverseKey = result[1]
	}
	if result := spliceRegexp.FindStringSubmatch(objBody); len(result) > 1 {
		spliceKey = result[1]
	}
	if result := swapRegexp.FindStringSubmatch(objBody); len(result) > 1 {
		swapKey = result[1]
	}

	regex, err := regexp.Compile(fmt.Sprintf("(?:a=)?%s\\.(%s|%s|%s)\\(a,(\\d+)\\)", obj, reverseKey, spliceKey, swapKey))
	if err != nil {
		return nil, err
	}

	var ops []operation
	for _, s := range regex.FindAllStringSubmatch(funcBody, -1) {
		switch s[1] {
		case reverseKey:
			ops = append(ops, reverseFunc)
		case swapKey:
			arg, _ := strconv.Atoi(s[2])
			ops = append(ops, newSwapFunc(arg))
		case spliceKey:
			arg, _ := strconv.Atoi(s[2])
			ops = append(ops, newSpliceFunc(arg))
		}
	}
	return ops, nil
}
