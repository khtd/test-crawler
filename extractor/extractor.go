package extractor

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"test-crawler/common"

	"github.com/tdewolff/parse/css"
	"golang.org/x/net/html"
)

func ExtractLinks(document io.ReadCloser, docType common.DocType) []string {
	switch docType {
	case common.CSS:
		return extractFromCss(document)
	case common.HTML:
		return extractFromHTML(document)
	case common.JS:
		return extractFromPlain(document)
	case common.PLAIN:
		return extractFromPlain(document)
	default:
		return extractFromPlain(document)
	}
}

func extractFromPlain(document io.ReadCloser) []string {
	links := make([]string, 0)
	linkReg, err := regexp.Compile(`(?P<protocol>(https?:\/\/(www\.)?)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,4}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*))`)
	if err != nil {
		fmt.Println(err)
		// TODO - process error on link procession
	}

	s := bufio.NewScanner(document)
	for s.Scan() {
		found := linkReg.FindAllString(s.Text(), -1)
		if found != nil {
			links = append(links, found...)
		}
	}
	return links
}

func extractFromCss(document io.ReadCloser) []string {
	links := make([]string, 0)
	l := css.NewLexer(document)
	for {
		tt, text := l.Next()
		switch tt {
		case css.ErrorToken:
			return links
		case css.AtKeywordToken:
			if string(text) == "@import" {
				for {
					tt, text := l.Next()
					if tt == css.StringToken {
						str := string(text)
						links = append(links, str[1:len(str)-1])
						break
					} else if tt == css.URLToken {
						str := string(text)
						links = append(links, str[5:len(str)-2])
						break
					}
				}
			}
		case css.URLToken:
			str := string(text)
			links = append(links, str[5:len(str)-2])
		}
	}
}

func extractFromHTML(document io.ReadCloser) []string {
	links := make([]string, 0)
	z := html.NewTokenizer(document)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return links
		}
		if tt == html.StartTagToken {
			token := z.Token()
			switch token.DataAtom.String() {
			case "a":
				link, found := extractAttrFromToken(token, "href")
				if found {
					links = append(links, link)
				}
			case "script":
				link, found := extractAttrFromToken(token, "src")
				if found {
					links = append(links, link)
				}
			case "link":
				link, found := extractAttrFromToken(token, "href")
				if found {
					links = append(links, link)
				}
			}
		}
	}
}

func extractAttrFromToken(token html.Token, key string) (string, bool) {
	for _, attr := range token.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}
