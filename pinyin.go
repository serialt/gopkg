package sugar

import (
	"fmt"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

// 汉字转拼音，当多音字时候只去第一个
func Pinyin(chinese string) string {
	ss := fmt.Sprint(pinyin.LazyConvert(chinese, nil))
	ss = strings.ReplaceAll(ss, "[", "")
	ss = strings.ReplaceAll(ss, "]", "")
	ss = strings.ReplaceAll(ss, " ", "")
	return ss
}
