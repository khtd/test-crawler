package filter

import (
	"net/url"
	"strings"
	"test-crawler/common"
	"test-crawler/logging"
	"test-crawler/state"
)

type Filter struct {
	maxDepth           int
	state              *state.State
	logger             *logging.Logger
	linksToProcessChan chan<- string
	pendingCountChan   chan<- int
}

func NewFilter(
	maxDepth int,
	state *state.State,
	linksToProcessChan chan<- string,
	pendingCountChan chan<- int,
	logger *logging.Logger,
) *Filter {
	return &Filter{
		maxDepth:           maxDepth,
		linksToProcessChan: linksToProcessChan,
		pendingCountChan:   pendingCountChan,
		state:              state,
		logger:             logger,
	}
}

func (f *Filter) Filter(link common.NewLinks) {
	discovered := make([]string, 0)
	for _, newLink := range link.Discovered {
		newUrl, err := url.Parse(newLink)
		if err != nil {
			// TODO - process error on link procession
			f.pendingCountChan <- -1
			continue
		}

		newUrl.RawQuery = ""
		baseUrl, err := url.Parse(link.Base)
		if err != nil {
			// TODO - process error on link procession
			f.pendingCountChan <- -1
			continue
		}

		if !newUrl.IsAbs() {
			newUrl = baseUrl.ResolveReference(newUrl)
		}

		if newUrl.Host != baseUrl.Host {
			f.pendingCountChan <- -1
			continue
		}

		if absDiff(countUrlSegments(newUrl), countUrlSegments(baseUrl)) > f.maxDepth {
			f.pendingCountChan <- -1
			continue
		}

		docType := common.GetLinkFileType(newUrl.String())
		if docType == common.UNKNOWN {
			f.pendingCountChan <- -1
			continue
		}

		known := f.state.IsKnown(newUrl.String())
		if known {
			f.pendingCountChan <- -1
			continue
		}

		discovered = append(discovered, newUrl.String())
		f.linksToProcessChan <- newUrl.String()
	}

	f.state.MarkProcessed(link.Base, discovered, "")
}

func countUrlSegments(u *url.URL) int {
	path := u.Path
	if path == "/" {
		return 0
	}
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	return len(strings.Split(path, "/"))
}

func absDiff(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}
