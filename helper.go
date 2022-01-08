package neugo

import (
	"io"
	"net/http"
	"regexp"
)

// Match all results using Regexp.
func matchMultiple(re *regexp.Regexp, content string) [][]string {
	return re.FindAllStringSubmatch(content, -1)
}

// Match one result using Regexp. Returns an empty string if no text matched.
func matchSingle(re *regexp.Regexp, content string) string {
	matched := re.FindAllStringSubmatch(content, -1)
	if len(matched) < 1 {
		return ""
	}
	return matched[0][1]
}

// Read and close the body of response.
func extractBody(resp *http.Response) string {
	res, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return string(res)
}
