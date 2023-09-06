package common

import (
	"net/url"
	"path"
)

func GetLinkFileType(link string) DocType {
	u, err := url.Parse(link)
	if err != nil {
		// TODO - process error on link procession
		return UNKNOWN
	}

	p := u.Path
	if p == "" {
		p = "/index.html"
	}
	ext := path.Ext(path.Base(p))
	if ext == "" {
		ext = ".html"
	}

	switch ext {
	case ".html":
		return HTML
	case ".css":
		return CSS
	case ".js":
		return JS
	case ".txt":
		return PLAIN
	default:
		return UNKNOWN
	}
}
