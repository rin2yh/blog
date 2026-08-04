// slidesync は rin2yh/slides のデッキを走査し、blogにまだ記事が無いものを
// externalUrl付きの記事として生成する。
//
// 作った記事はmarkdownの箇条書きで -out に書き出す。何も作らなければ空になる。
// 進捗は標準出力、異常は標準エラーに出す。
//
//	go run ./cmd/slidesync <slides-repo-path>
package main

import (
	"flag"
	"fmt"
	"io"
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
	outPath := flag.String("out", "", "作った記事の一覧をmarkdownの箇条書きで書き出す先")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: go run ./cmd/slidesync [flags] <slides-repo-path>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	// 一覧を使わないなら捨てる。標準出力は進捗を出すために空けておく
	list := io.Discard
	if *outPath != "" {
		f, err := os.Create(*outPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		list = f
	}

	if err := run(flag.Arg(0), *baseURL, *contentDir, list); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// run は記事を作り、作った分をmarkdownの箇条書きで list に書く。
func run(slidesRepo, baseURL, contentDir string, list io.Writer) error {
	paths, err := filepath.Glob(filepath.Join(slidesRepo, deckGlob))
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		return fmt.Errorf("デッキが見つかりません: %s", filepath.Join(slidesRepo, deckGlob))
	}

	published, err := collectExternalURLs(contentDir)
	if err != nil {
		return err
	}

	created := 0
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		d, err := parseDeck(filepath.Base(p), string(b))
		if err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}

		url := strings.TrimSuffix(baseURL, "/") + "/" + d.name + "/"
		if published[strings.TrimSuffix(url, "/")] {
			continue
		}

		name, err := post.DirName(d.name)
		if err != nil {
			return err
		}
		dir := filepath.Join(contentDir, name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte(renderPost(d, url)), 0o644); err != nil {
			return err
		}

		created++
		fmt.Printf("Created: %s\n", dir)
		fmt.Fprintf(list, "- [%s](%s) (%s)\n", d.title, url, d.date)
	}

	// 件数はPRをレビューするときの手掛かりになる。URLの組み立て方が
	// slides側と食い違うと、ここが急に全デッキ分に跳ねる
	fmt.Printf("デッキ%d件のうち%d件の記事を作りました\n", len(paths), created)
	return nil
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
