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
	"strings"
	"time"

	"blog/internal/post"
)

const (
	// スライドの公開URL。config/_default/menu.toml のSlidesリンクと対になる。
	defaultBaseURL = "https://slides.rin2yh.com"
	deckGlob       = "slides/*.md"
	dateLayout     = "2006-01-02"
	// 発表日は日付までしか持たないので、JSTの0時に固定する。
	timeSuffix = "T00:00:00+09:00"
)

type deck struct {
	name  string // ファイル名のstem。URLの末尾になる
	title string
	date  string // YYYY-MM-DD
}

// content/post の記事はすべてTOML frontmatterなので、その形だけを見る。
var externalURLRe = regexp.MustCompile(`(?m)^\s*externalUrl\s*=\s*['"]([^'"]+)['"]`)

func main() {
	baseURL := flag.String("base", defaultBaseURL, "スライドの公開URLのベース")
	contentDir := flag.String("content", post.Dir, "記事を探して生成するディレクトリ")
	summaryPath := flag.String("summary", "", "生成した記事の一覧をmarkdownの箇条書きで書き出す先")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: go run ./cmd/slidesync [flags] <slides-repo-path>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(flag.Arg(0), *baseURL, *contentDir, *summaryPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(slidesRepo, baseURL, contentDir, summaryPath string) error {
	decks, err := readDecks(slidesRepo)
	if err != nil {
		return err
	}

	published, err := collectExternalURLs(contentDir)
	if err != nil {
		return err
	}

	var summary strings.Builder
	for _, d := range decks {
		url := deckURL(baseURL, d.name)
		if published[strings.TrimSuffix(url, "/")] {
			continue
		}
		path, err := writePost(contentDir, d, url)
		if err != nil {
			return err
		}
		fmt.Printf("Created: %s (%s)\n", path, d.title)
		fmt.Fprintf(&summary, "- [%s](%s) (%s)\n", d.title, url, d.date)
	}

	if summary.Len() == 0 {
		fmt.Println("記事が無いデッキはありません")
	}
	// 生成が無ければ空ファイルを書く。workflowはこのファイルの中身だけを見る。
	if summaryPath != "" {
		return os.WriteFile(summaryPath, []byte(summary.String()), 0o644)
	}
	return nil
}

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
	return decks, nil
}

// parseDeck はtitleとdateを取り出す。欠けていたらエラーにする。黙って一覧から落とさないため。
func parseDeck(filename, content string) (deck, error) {
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

	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	return deck{name: name, title: title, date: date}, nil
}

// frontmatter は先頭の --- で挟まれたブロックを読む。Marpのfrontmatterはフラット。
func frontmatter(content string) (map[string]string, error) {
	body, ok := strings.CutPrefix(strings.ReplaceAll(content, "\r\n", "\n"), "---\n")
	if !ok {
		return nil, fmt.Errorf("frontmatterがありません")
	}

	fm := map[string]string{}
	for _, l := range strings.Split(body, "\n") {
		if strings.TrimSpace(l) == "---" {
			return fm, nil
		}
		key, value, ok := strings.Cut(l, ":")
		if !ok || strings.TrimSpace(key) == "" {
			continue
		}
		fm[strings.TrimSpace(key)] = unquote(strings.TrimSpace(value))
	}
	return nil, fmt.Errorf("frontmatterが閉じていません")
}

// unquote は囲みの引用符を1組だけ外す。strings.Trimは両端から何個でも外してしまい、
// strconv.Unquoteはシングルクォートをruneリテラル扱いするので、どちらも使えない。
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

// collectExternalURLs は掲載済みのURLを集める。パスやファイル名ではなくURLで
// 突き合わせるので、記事を移動しても壊れない。
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
			urls[strings.TrimSuffix(m[1], "/")] = true
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func deckURL(baseURL, name string) string {
	return strings.TrimSuffix(baseURL, "/") + "/" + name + "/"
}

func writePost(contentDir string, d deck, url string) (string, error) {
	name, err := post.DirName(d.name)
	if err != nil {
		return "", err
	}

	dir := filepath.Join(contentDir, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	path := filepath.Join(dir, "index.md")
	if err := os.WriteFile(path, []byte(renderPost(d, url)), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// renderPost の出力は既存のexternal記事に合わせてある。archetypes/external.md とは
// draft と categories と tags が違うので、片方を直すときはもう片方も見ること。
func renderPost(d deck, url string) string {
	return fmt.Sprintf(`+++
date = %s
draft = false
title = %s
categories = ['tech', 'external']
tags = ['slide']
externalUrl = %s
+++
`, tomlString(d.date+timeSuffix), tomlString(d.title), tomlString(url))
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
