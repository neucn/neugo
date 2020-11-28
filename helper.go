package neugo

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

var (
	errorNoMatched = errors.New("no matched")
)

// 使用正则在文本中匹配出所有结果，如果没有匹配到则返回error
func matchMultiple(re *regexp.Regexp, content string) ([][]string, error) {
	matched := re.FindAllStringSubmatch(content, -1)
	if len(matched) < 1 {
		return nil, errorNoMatched
	}
	return matched, nil
}

// 使用正则在文本中匹配出一个结果并取出，如果没有匹配到则返回error
func matchSingle(re *regexp.Regexp, content string) (string, error) {
	matched := re.FindAllStringSubmatch(content, -1)
	if len(matched) < 1 {
		return "", errorNoMatched
	}
	return matched[0][1], nil
}

// 读取响应体，并关闭resp.Body
func extractBody(resp *http.Response) string {
	res, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return string(res)
}
