// Package post は記事の置き場所と命名規則を持つ。
//
// 命名規則 <name>-<slug> は config/_default/permalinks.toml の :contentbasename と
// .github/preview/changed-urls.sh が依存しているので、ここを唯一の定義とする。
package post

import "crypto/rand"

const (
	// Section は記事のセクション名。URLとHugoのcontent配下のパスの両方になる。
	Section = "post"
	// Dir は記事を置くディレクトリ。
	Dir = "content/" + Section
)

const (
	slugLength = 14
	charset    = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// DirName は <name>-<ランダムなslug> という記事ディレクトリ名を返す。
func DirName(name string) (string, error) {
	b := make([]byte, slugLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return name + "-" + string(b), nil
}
