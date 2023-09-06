package link_processor

import (
	"fmt"
	"net/http"
	"strings"
	"test-crawler/common"
	"test-crawler/extractor"
	"test-crawler/filter"
	"test-crawler/logging"
	"test-crawler/state"
)

type LinkProcessor struct {
	httpClient       http.Client
	state            *state.State
	filter           *filter.Filter
	logger           *logging.Logger
	pendingCountChan chan<- int
}

func NewLinkProcessor(state *state.State, filter *filter.Filter, logger *logging.Logger, pendingCountChan chan<- int) *LinkProcessor {
	return &LinkProcessor{
		httpClient:       http.Client{},
		state:            state,
		filter:           filter,
		logger:           logger,
		pendingCountChan: pendingCountChan,
	}
}

func (p *LinkProcessor) Process(link string) string {
	if !p.state.TryToTake(link) {
		return "taken"
	}

	docType := common.GetLinkFileType(link)
	if docType == common.UNKNOWN {
		// TODO - process error on link procession
		return "unknown type"
	}

	response, err := http.Get(link)
	if err != nil {
		// TODO - process error on link procession
		return err.Error()
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		// TODO - process error on link procession
		return "response status code - " + fmt.Sprint(response.StatusCode)
	}

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text") && !strings.HasPrefix(contentType, "application/javascript") {
		// TODO - process error on link procession
		return "response status code - " + contentType
	}

	newLinks := extractor.ExtractLinks(response.Body, docType)
	go p.filter.Filter(common.NewLinks{Base: link, Discovered: newLinks})
	p.pendingCountChan <- len(newLinks)
	return ""
}
