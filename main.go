package main

//go:generate statik

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	_ "github.com/mattn/siritori/statik"

	"github.com/rakyll/statik/fs"
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
)

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
	rs := []rune(text)
	r := rs[len(rs)-1]

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

func shiritori(text string) (string, error) {
	text = strings.Replace(text, "ー", "", -1)
	if rand.Int()%2 == 0 {
		text = hira2kana(text)
	} else {
		text = kana2hira(text)
	}
	return search(text)
}

func siritori(text string) (string, error) {
	rs := []rune(strings.TrimSpace(text))
	if len(rs) == 0 {
		return "", errors.New("なんやねん")
	}
	if rs[len(rs)-1] == 'ん' || rs[len(rs)-1] == 'ン' {
		return "", errors.New("出直して来い")
	}
	s, err := shiritori(text)
	if err != nil {
		return "", err
	}
	if s == "" {
		return "", errors.New("わかりません")
	}
	rs = []rune(s)
	if rs[len(rs)-1] == 'ん' || rs[len(rs)-1] == 'ン' {
		s += "\nあっ..."
	}
	return s, nil
}

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	result, err := siritori(strings.TrimSpace(string(b)))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(result)
}
