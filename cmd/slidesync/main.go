// slidesync は rin2yh/slides のデッキを走査し、blogにまだ記事が無いものを
// externalUrl付きの記事として生成する。
//
//	go run ./cmd/slidesync <slides-repo-path>
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"blog/internal/slug"
)

const (
	defaultBaseURL    = "https://slides.rin2yh.com"
	defaultContentDir = "content/post"
	deckGlob          = "slides/*.md"
	dateLayout        = "2006-01-02"
	// 発表日は日付までしか持たないので、JSTの0時に固定する。
	timeSuffix = "T00:00:00+09:00"
)

// deck は slides リポジトリのデッキ1枚分。
type deck struct {
	name  string // ファイル名のstem。URLの末尾になる
	title string
	date  string // YYYY-MM-DD
}

// externalURLRe は記事のfrontmatterからexternalUrlの値を取り出す。
var externalURLRe = regexp.MustCompile(`(?m)^\s*externalUrl\s*=\s*['"]([^'"]+)['"]`)

func main() {
	baseURL := flag.String("base", defaultBaseURL, "スライドの公開URLのベース")
	contentDir := flag.String("content", defaultContentDir, "記事を探して生成するディレクトリ")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: go run ./cmd/slidesync [flags] <slides-repo-path>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(flag.Arg(0), *baseURL, *contentDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(slidesRepo, baseURL, contentDir string) error {
	decks, err := readDecks(slidesRepo)
	if err != nil {
		return err
	}

	published, err := collectExternalURLs(contentDir)
	if err != nil {
		return err
	}

	missing := missingDecks(decks, published, baseURL)
	if len(missing) == 0 {
		fmt.Println("記事が無いデッキはありません")
		return nil
	}

	for _, d := range missing {
		path, err := writePost(contentDir, d, deckURL(baseURL, d.name))
		if err != nil {
			return err
		}
		fmt.Printf("Created: %s (%s)\n", path, d.title)
	}
	return nil
}

// readDecks は slides リポジトリのデッキを読み、発表日の新しい順に返す。
func readDecks(slidesRepo string) ([]deck, error) {
	paths, err := filepath.Glob(filepath.Join(slidesRepo, filepath.FromSlash(deckGlob)))
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("デッキが見つかりません: %s", filepath.Join(slidesRepo, deckGlob))
	}

	decks := make([]deck, 0, len(paths))
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		d, err := parseDeck(filepath.Base(p), string(b))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		decks = append(decks, d)
	}

	sort.Slice(decks, func(i, j int) bool {
		if decks[i].date != decks[j].date {
			return decks[i].date > decks[j].date
		}
		return decks[i].name < decks[j].name
	})
	return decks, nil
}

// parseDeck はデッキのfrontmatterからtitleとdateを取り出す。
// どちらかが欠けていたり日付の形式が違えばエラーにする。黙って一覧から落とさないため。
func parseDeck(filename, content string) (deck, error) {
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	fm, err := frontmatter(content)
	if err != nil {
		return deck{}, err
	}

	title := fm["title"]
	if title == "" {
		return deck{}, fmt.Errorf("frontmatterにtitleがありません")
	}

	date := fm["date"]
	if date == "" {
		return deck{}, fmt.Errorf("frontmatterにdateがありません")
	}
	if _, err := time.Parse(dateLayout, date); err != nil {
		return deck{}, fmt.Errorf("dateが%s形式ではありません: %q", dateLayout, date)
	}

	return deck{name: name, title: title, date: date}, nil
}

// frontmatter は先頭の --- で挟まれたブロックをkey/valueとして読む。
// Marpのfrontmatterはフラットなので、ネストは扱わない。
func frontmatter(content string) (map[string]string, error) {
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")

	start := -1
	for i, l := range lines {
		if strings.TrimSpace(l) == "" {
			continue
		}
		if strings.TrimSpace(l) == "---" {
			start = i
		}
		break
	}
	if start < 0 {
		return nil, fmt.Errorf("frontmatterがありません")
	}

	fm := map[string]string{}
	for _, l := range lines[start+1:] {
		if strings.TrimSpace(l) == "---" {
			return fm, nil
		}
		key, value, ok := strings.Cut(l, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		fm[key] = unquote(strings.TrimSpace(value))
	}
	return nil, fmt.Errorf("frontmatterが閉じていません")
}

// unquote は前後を囲む同じ引用符を1組だけ外す。
func unquote(s string) string {
	if len(s) < 2 {
		return s
	}
	q := s[0]
	if (q == '\'' || q == '"') && s[len(s)-1] == q {
		return s[1 : len(s)-1]
	}
	return s
}

// collectExternalURLs は記事ディレクトリ配下のexternalUrlを集める。
// 記事のパスやファイル名ではなくURLで突き合わせるので、記事を移動しても壊れない。
func collectExternalURLs(contentDir string) (map[string]bool, error) {
	urls := map[string]bool{}

	err := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, m := range externalURLRe.FindAllStringSubmatch(string(b), -1) {
			urls[normalizeURL(m[1])] = true
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return urls, nil
}

// missingDecks は記事がまだ無いデッキだけを返す。
func missingDecks(decks []deck, published map[string]bool, baseURL string) []deck {
	var missing []deck
	for _, d := range decks {
		if published[normalizeURL(deckURL(baseURL, d.name))] {
			continue
		}
		missing = append(missing, d)
	}
	return missing
}

func deckURL(baseURL, name string) string {
	return strings.TrimSuffix(baseURL, "/") + "/" + name + "/"
}

// normalizeURL は末尾のスラッシュの有無を吸収する。
func normalizeURL(u string) string {
	return strings.TrimSuffix(u, "/")
}

// writePost は既存のexternal記事と同じfrontmatterを持つ記事を作る。本文は空。
func writePost(contentDir string, d deck, url string) (string, error) {
	s, err := slug.Generate()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(contentDir, fmt.Sprintf("%s-%s", d.name, s))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	path := filepath.Join(dir, "index.md")
	if err := os.WriteFile(path, []byte(renderPost(d, url)), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// renderPost は記事のfrontmatterを組み立てる。
// キーの順序は既存のexternal記事に合わせる。
func renderPost(d deck, url string) string {
	var b strings.Builder
	b.WriteString("+++\n")
	fmt.Fprintf(&b, "date = %s\n", tomlString(d.date+timeSuffix))
	b.WriteString("draft = false\n")
	fmt.Fprintf(&b, "title = %s\n", tomlString(d.title))
	b.WriteString("categories = ['tech', 'external']\n")
	b.WriteString("tags = ['slide']\n")
	fmt.Fprintf(&b, "externalUrl = %s\n", tomlString(url))
	b.WriteString("+++\n")
	return b.String()
}

// tomlString は既存記事に合わせてシングルクォートのリテラル文字列にする。
// リテラル文字列はシングルクォートを含められないので、その場合だけ基本文字列に落とす。
func tomlString(s string) string {
	if !strings.Contains(s, "'") {
		return "'" + s + "'"
	}
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`)
	return `"` + r.Replace(s) + `"`
}
