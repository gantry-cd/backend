package random

import (
	"crypto/rand"
	"errors"
)

var (
	ErrUnExpected = errors.New("unexpected error...")
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString はbcriptを使用して安全なランダムな文字列を生成します
func RandomString(length int) (string, error) {

	// 乱数を生成
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", ErrUnExpected
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}
