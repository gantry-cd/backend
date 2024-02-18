package branch

import (
	"net/url"
	"strings"
)

const (
	// Transpile1123 は、Branch名をURLエンコードしたのち%を.p.に置換する
	MagicPercent = ".p."
)

// Transpile は、Branch名をURLエンコードしたのち%を.p.に置換する
func Transpile1123(s string) string {
	return strings.ReplaceAll(url.QueryEscape(s), "%", MagicPercent)
}

// TranspileBranchName は、Transpile1123を元に戻す
func TranspileBranchName(s string) (string, error) {
	return url.QueryUnescape(strings.ReplaceAll(s, MagicPercent, "%"))
}
