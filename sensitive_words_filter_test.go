package sensitive_words_filter

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

var tests = []struct {
	level             int
	skipDistance      int
	text              string
	replaceText       string
	hasSensitiveWords bool
	out               string
}{
	{FilterLevelLow, 1, "你好丑", "帅", true, "你好帅"},
	{FilterLevelMiddle, 1, "你好丑", "帅", true, "你好帅"},
	{FilterLevelHight, 1, "你好丑", "帅", true, "你好帅"},
	{FilterLevelLow, 1, "你好傻逼", "*", true, "你好**"},
	{FilterLevelMiddle, 3, "你好傻  逼", "*", true, "你好****"},
	{FilterLevelHight, 1, "你好傻 逼", "*", true, "你好* 逼"},
	{FilterLevelHight, 1, "你好傻逼啊", "*", true, "你好**啊"},
	{FilterLevelHight, 1, "你好", "*", false, "你好"},
}

func TestFilter(t *testing.T) {
	filter := GetInstance()
	filter.Build([]string{"傻逼", "丑", "傻子", "傻"})
	for i, tt := range tests {
		filter.SetLevel(tt.level)
		filter.SetSkipDistance(tt.skipDistance)
		filter.SetReplaceText(tt.replaceText)
		text, hasSensitiveWords := filter.Filter(tt.text)
		if text != tt.out || hasSensitiveWords != tt.hasSensitiveWords {
			t.Errorf("%d . %q => %q, wanted: %q", i, tt.text, text, tt.out)
		}
	}
}

func LoadDict(path string) []string {
	var lines []string
	f, err := os.Open(path)
	if err != nil {
		return lines
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func TestFilterThousandWords(t *testing.T) {
	filter := GetInstance()
	lines := LoadDict("./res/words.txt")
	filter.Build(lines)
	filter.SetLevel(FilterLevelLow)
	filter.SetSkipDistance(3)
	filter.SetReplaceText("*")
	testStr := "澳门 赌场在线充值，日本性感-女优?在线发牌，无限畅玩，提现秒到账！"
	text, _ := filter.Filter(testStr)
	fmt.Println(text)
}
