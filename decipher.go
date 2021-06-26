package youtube

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

func (c *Client) decipherURL(ctx context.Context, cipher string) (string, error) {
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
	ops, err := c.getOpsFromCache(ctx)
	if err != nil {
		return "", err
	}

	// apply operations
	bs := []byte(queryParams.Get("s"))

	for _, op := range ops {
		switch op.Name {
		case OpReverse:
			bs = reverseFunc(bs)
		case OpSplice:
			bs = newSpliceFunc(op.Value)(bs)
		case OpSwap:
			bs = newSwapFunc(op.Value)(bs)
		}
	}

	decipheredURL := fmt.Sprintf("%s&%s=%s", queryParams.Get("url"), queryParams.Get("sp"), string(bs))
	return decipheredURL, nil
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
	actionsObjRegexp = regexp.MustCompile(fmt.Sprintf(
		"var (%s)=\\{((?:(?:%s%s|%s%s|%s%s),?\\n?)+)\\};", jsvarStr, jsvarStr, swapStr, jsvarStr, spliceStr, jsvarStr, reverseStr))

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

func parseOperations(player []byte) ([]Operation, error) {
	objResult := actionsObjRegexp.FindSubmatch(player)
	funcResult := actionsFuncRegexp.FindSubmatch(player)
	if len(objResult) < 3 || len(funcResult) < 2 {
		return nil, fmt.Errorf("error parsing signature tokens (#obj=%d, #func=%d)", len(objResult), len(funcResult))
	}

	obj := objResult[1]
	objBody := objResult[2]
	funcBody := funcResult[1]

	var reverseKey, spliceKey, swapKey string

	if result := reverseRegexp.FindSubmatch(objBody); len(result) > 1 {
		reverseKey = string(result[1])
	}
	if result := spliceRegexp.FindSubmatch(objBody); len(result) > 1 {
		spliceKey = string(result[1])
	}
	if result := swapRegexp.FindSubmatch(objBody); len(result) > 1 {
		swapKey = string(result[1])
	}

	regex, err := regexp.Compile(fmt.Sprintf("(?:a=)?%s\\.(%s|%s|%s)\\(a,(\\d+)\\)", obj, reverseKey, spliceKey, swapKey))
	if err != nil {
		return nil, err
	}

	var ops []Operation
	for _, s := range regex.FindAllSubmatch(funcBody, -1) {
		switch string(s[1]) {
		case reverseKey:
			ops = append(ops, Operation{Name: OpReverse})
		case swapKey:
			arg, _ := strconv.Atoi(string(s[2]))
			ops = append(ops, Operation{OpSwap, arg})
		case spliceKey:
			arg, _ := strconv.Atoi(string(s[2]))
			ops = append(ops, Operation{OpSplice, arg})
		}
	}
	return ops, nil
}

func (c *Client) getOpsFromCache(ctx context.Context) ([]Operation, error) {
	// if there is cache - we need to make sure its actual version
	if c.cache != nil {
		version, err := c.fetchPlayerVersion(ctx)
		if err != nil {
			// we couldn't fetch what's current player version, should we really exit here?..
			return nil, err
		}

		ops, ok := c.cache.getOps(version)
		if ok {
			return ops, nil
		}
	}

	// both empty cache and wrong version leads here
	_, err := c.cachePlayer(ctx)
	if err != nil {
		return nil, err
	}

	// we just re-cached player, no need to check version
	return c.cache.Operations, nil
}
