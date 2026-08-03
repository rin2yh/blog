// Package slug は記事ディレクトリ名の末尾に付けるランダムなslugを生成する。
package slug

import "crypto/rand"

const (
	// Length は生成するslugの文字数。
	Length = 14
	// Charset はslugに使う文字の集合。
	Charset = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// Generate はLength文字のランダムなslugを返す。
func Generate() (string, error) {
	b := make([]byte, Length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = Charset[int(b[i])%len(Charset)]
	}
	return string(b), nil
}
