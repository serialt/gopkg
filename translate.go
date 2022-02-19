package gopkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func TranslateCh2En(text string) (string, error) {
	url := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=zh-cn&tl=en&dt=t&q=%s", url.QueryEscape(text))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//返回的json反序列化比较麻烦, 直接字符串拆解
	ss := string(bs)
	ss = strings.ReplaceAll(ss, "[", "")
	ss = strings.ReplaceAll(ss, "]", "")
	ss = strings.ReplaceAll(ss, "null,", "")
	ss = strings.Trim(ss, `"`)
	ps := strings.Split(ss, `","`)
	return ps[0], nil
}

func TranslateEn2Ch(text string) (string, error) {
	url := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=en&tl=zh-cn&dt=t&q=%s", url.QueryEscape(text))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//返回的json反序列化比较麻烦, 直接字符串拆解
	ss := string(bs)
	ss = strings.ReplaceAll(ss, "[", "")
	ss = strings.ReplaceAll(ss, "]", "")
	ss = strings.ReplaceAll(ss, "null,", "")
	ss = strings.Trim(ss, `"`)
	ps := strings.Split(ss, `","`)
	return ps[0], nil
}
