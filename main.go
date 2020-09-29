package main

//go:generate statik

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ikawaha/kagome/tokenizer"
	_ "github.com/mattn/siritori/statik"

	"github.com/rakyll/statik/fs"
)

var (
	reIgnoreText = regexp.MustCompile(`[\[\]「」『』]`)
)

var upper = strings.NewReplacer(
	"ぁ", "あ",
	"ぃ", "い",
	"ぅ", "う",
	"ぇ", "え",
	"ぉ", "お",
	"ゃ", "や",
	"ゅ", "ゆ",
	"ょ", "よ",
	"ァ", "ア",
	"ィ", "イ",
	"ゥ", "ウ",
	"ェ", "エ",
	"ォ", "オ",
	"ャ", "ヤ",
	"ュ", "ユ",
	"ョ", "ヨ",
)

func isSpace(c []string) bool {
	return c[1] == "空白"
}

func kana2hira(s string) string {
	return strings.Map(func(r rune) rune {
		if 0x30A1 <= r && r <= 0x30F6 {
			return r - 0x0060
		}
		return r
	}, s)
}

func hira2kana(s string) string {
	return strings.Map(func(r rune) rune {
		if 0x3041 <= r && r <= 0x3096 {
			return r + 0x0060
		}
		return r
	}, s)
}

func search(text string) (string, error) {
	text = reIgnoreText.ReplaceAllString(text, " ")
	t := tokenizer.New()
	tokens := t.Tokenize(text)
	for _, tok := range tokens {
		c := tok.Features()
		if len(c) == 0 || isSpace(c) {
			continue
		}
		y := c[len(c)-1]
		if y == "*" {
			y = tok.Surface
		}
		text = y
	}

	if rand.Int()%2 == 0 {
		text = hira2kana(text)
	} else {
		text = kana2hira(text)
	}

	text = upper.Replace(strings.ReplaceAll(text, "ー", ""))
	rs := []rune(text)
	r := rs[len(rs)-1]

	if r == 'ん' || r == 'ン' {
		return "出直して来い", nil
	}

	statikFS, err := fs.New()
	if err != nil {
		return "", err
	}
	f, err := statikFS.Open("/dict.txt")
	if err != nil {
		return "", err
	}
	defer f.Close()
	buf := bufio.NewReader(f)

	words := []string{}
	for {
		b, _, err := buf.ReadLine()
		if err != nil {
			break
		}
		line := string(b)
		if ([]rune(line))[0] == r {
			words = append(words, line)
		}
	}
	if len(words) == 0 {
		return "", nil
	}
	return words[rand.Int()%len(words)], nil
}

func siritori(text string) (string, error) {
	text = strings.TrimSpace(text)
	rs := []rune(text)
	if len(rs) == 0 {
		return "", errors.New("なんやねん")
	}
	s, err := search(text)
	if err != nil {
		return "", err
	}
	if s == "" {
		return "", errors.New("わかりません")
	}
	rs = []rune(s)
	r := rs[len(rs)-1]
	if r == 'ん' || r == 'ン' {
		s += "\nあっ..."
	}
	return s, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	text := strings.Join(os.Args[1:], "")
	if len(os.Args) == 1 {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		text = string(b)
	}
	result, err := siritori(text)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(result)
}
