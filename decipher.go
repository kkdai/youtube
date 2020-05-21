package youtube

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func (y *Youtube) parseDecipherOpsAndArgs() (operations []string, args []int, err error) {
	// try to get whole page
	client, err := y.getHTTPClient()
	if err != nil {
		return nil, nil, fmt.Errorf("get http client error=%s", err)
	}

	if y.VideoID == "" {
		return nil, nil, fmt.Errorf("video id is empty , err=%s", err)
	}
	embedUrl := fmt.Sprintf("https://youtube.com/embed/%s?hl=en", y.VideoID)

	embeddedPageResp, err := client.Get(embedUrl)
	if err != nil {
		return nil, nil, err
	}
	defer embeddedPageResp.Body.Close()

	if embeddedPageResp.StatusCode != 200 {
		return nil, nil, err
	}

	embeddedPageBodyBytes, err := ioutil.ReadAll(embeddedPageResp.Body)
	if err != nil {
		return nil, nil, err
	}
	embeddedPage := string(embeddedPageBodyBytes)

	playerConfigPattern := regexp.MustCompile(`yt\.setConfig\({'PLAYER_CONFIG':(.*)}\);`)
	playerConfig := playerConfigPattern.FindString(embeddedPage)

	basejsPattern := regexp.MustCompile(`"js":"\\/s\\/player(.*)base\.js`)
	// eg: "js":\"\/s\/player\/f676c671\/player_ias.vflset\/en_US\/base.js
	escapedBasejsUrl := basejsPattern.FindString(playerConfig)
	// eg: ["js", "\/s\/player\/f676c671\/player_ias.vflset\/en_US\/base.js]
	arr := strings.Split(escapedBasejsUrl, ":\"")
	basejsUrl := "https://youtube.com" + strings.ReplaceAll(arr[len(arr)-1], "\\", "")
	basejsUrlResp, err := client.Get(basejsUrl)
	if err != nil {
		return nil, nil, err
	}

	defer basejsUrlResp.Body.Close()
	if basejsUrlResp.StatusCode != 200 {
		return nil, nil, err
	}

	basejsBodyBytes, err := ioutil.ReadAll(basejsUrlResp.Body)
	if err != nil {
		return nil, nil, err
	}
	basejs := string(basejsBodyBytes)

	// regex to get name of decipher function
	decipherFuncNamePattern := regexp.MustCompile(`(\w+)=function\(\w+\){(\w+)=(\w+)\.split\(\x22{2}\);.*?return\s+(\w+)\.join\(\x22{2}\)}`)

	// Ft=function(a){a=a.split("");Et.vw(a,2);Et.Zm(a,4);Et.Zm(a,46);Et.vw(a,2);Et.Zm(a,34);Et.Zm(a,59);Et.cn(a,42);return a.join("")} => get Ft
	arr = decipherFuncNamePattern.FindStringSubmatch(basejs)
	funcName := arr[1]
	decipherFuncBodyPattern := regexp.MustCompile(fmt.Sprintf(`[^h\.]%s=function\(\w+\)\{(.*?)\}`, funcName))

	// eg: get a=a.split("");Et.vw(a,2);Et.Zm(a,4);Et.Zm(a,46);Et.vw(a,2);Et.Zm(a,34);Et.Zm(a,59);Et.cn(a,42);return a.join("")
	arr = decipherFuncBodyPattern.FindStringSubmatch(basejs)
	decipherFuncBody := arr[1]

	// FuncName in Body => get Et
	funcNameInBodyRegex := regexp.MustCompile(`(\w+).\w+\(\w+,\d+\);`)
	arr = funcNameInBodyRegex.FindStringSubmatch(decipherFuncBody)
	funcNameInBody := arr[1]
	decipherDefBodyRegex := regexp.MustCompile(fmt.Sprintf(`var\s+%s=\{(\w+:function\(\w+(,\w+)?\)\{(.*?)\}),?\};`, funcNameInBody))
	re := regexp.MustCompile(`\r?\n`)
	basejs = re.ReplaceAllString(basejs, "")
	arr1 := decipherDefBodyRegex.FindStringSubmatch(basejs)

	// eg:  vw:function(a,b){a.splice(0,b)},cn:function(a){a.reverse()},Zm:function(a,b){var c=a[0];a[0]=a[b%a.length];a[b%a.length]=c}
	decipherDefBody := arr1[1]
	// eq:  [ a=a.split("") , Et.vw(a,2) , Et.Zm(a,4) , Et.Zm(a,46) , Et.vw(a,2) , Et.Zm(a,34), Et.Zm(a,59) , Et.cn(a,42) , return a.join("") ]
	decipherFuncs := strings.Split(decipherFuncBody, ";")

	var funcSeq []string
	var funcArgs []int

	for _, v := range decipherFuncs {
		// calledFuncNameRegex \w+(?:.|\[)("?\w+(?:")?)\]?\(
		// eg: Et.vw(a,2) => vw
		calledFuncNameRegex, err := regexp.Compile(`\w+(?:.|\[)("?\w+(?:")?)\]?\(`)
		if err != nil {
			return nil, nil, err
		}
		arr := calledFuncNameRegex.FindStringSubmatch(v)
		if len(arr) < 1 || arr[1] == "" {
			continue
		}
		calledFuncName := arr[1]

		// Et.vw(a,2) => [a, 2]
		funcArgRegex := regexp.MustCompile(`\(\w+,(\d+)\)`)

		// splice
		spliceFuncPattern := fmt.Sprintf(`%s:\bfunction\b\([a],b\).(\breturn\b)?.?\w+\.splice`, calledFuncName)
		if regexp.MustCompile(spliceFuncPattern).MatchString(decipherDefBody) {
			arr := funcArgRegex.FindStringSubmatch(v)
			arg, err := strconv.Atoi(arr[1])
			if err != nil {
				return nil, nil, err
			}
			funcSeq = append(funcSeq, "splice")
			funcArgs = append(funcArgs, arg)
			continue
		}

		// Swap
		swapFuncPattern := fmt.Sprintf(`%s:\bfunction\b\(\w+\,\w\).\bvar\b.\bc=a\b`, calledFuncName)
		if regexp.MustCompile(swapFuncPattern).MatchString(decipherDefBody) {
			arr := funcArgRegex.FindStringSubmatch(v)
			arg, err := strconv.Atoi(arr[1])
			if err != nil {
				return nil, nil, err
			}
			funcSeq = append(funcSeq, "swap")
			funcArgs = append(funcArgs, arg)
			continue
		}

		// Reverse
		reverseFuncPattern := fmt.Sprintf(`%s:\bfunction\b\(\w+\)`, calledFuncName)
		if regexp.MustCompile(reverseFuncPattern).MatchString(decipherDefBody) {
			arr := funcArgRegex.FindStringSubmatch(v)
			arg, err := strconv.Atoi(arr[1])
			if err != nil {
				return nil, nil, err
			}
			funcSeq = append(funcSeq, "reverse")
			funcArgs = append(funcArgs, arg)
			continue
		}
	}
	return funcSeq, funcArgs, nil
}

func (y *Youtube) decipher(cipher string) (string, error) {
	queryParams, err := url.ParseQuery(cipher)
	if err != nil {
		return "", err
	}
	cipherMap := make(map[string]string)
	for key, value := range queryParams {
		cipherMap[key] = strings.Join(value, "")
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

	s := cipherMap["s"]
	bs := []byte(s)
	splice := func(b int) {
		bs = bs[b:]
	}
	swap := func(b int) {
		pos := b % len(bs)
		bs[0], bs[pos] = bs[pos], bs[0]
	}
	reverse := func(options ...interface{}) {
		l, r := 0, len(bs)-1
		for l < r {
			bs[l], bs[r] = bs[r], bs[l]
			l++
			r--
		}
	}
	operations, args, err := y.parseDecipherOpsAndArgs()
	if err != nil {
		return "", err
	}
	for i, op := range operations {
		switch op {
		case "splice":
			splice(args[i])
		case "swap":
			swap(args[i])
		case "reverse":
			reverse(args[i])
		}
	}
	cipherMap["s"] = string(bs)

	decipheredUrl := fmt.Sprintf("%s&%s=%s", cipherMap["url"], cipherMap["sp"], cipherMap["s"])
	return decipheredUrl, nil
}
