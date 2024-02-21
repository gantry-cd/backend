package url

import (
	"fmt"
	"strings"
)

// IncludeBasicAuth は指定された URL に Basic 認証情報を含めた URL を返す
// ex: https://example.com , user , pass -> https://user:pass@example.com
// ex: example.com , user , pass -> https://user:pass@example.com
func IncludeBasicAuth(url, user, pass string) string {
	urls := strings.Split(url, "//")
	if len(urls) < 2 {
		urls = append([]string{"https:"}, urls...)
	}
	return fmt.Sprintf("%s//%s:%s@%s", urls[0], user, pass, urls[1])
}
