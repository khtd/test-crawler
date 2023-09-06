package builder

import (
	"test-crawler/app"
	"test-crawler/config"
	"test-crawler/filter"
	"test-crawler/link_processor"
	"test-crawler/logging"
	"test-crawler/monitor"

	"test-crawler/state"
)

// Builder builds an app.
type Builder struct{}

// NewBuilder creates a new builder instance.
func NewBuilder() *Builder {
	return &Builder{}
}

// Build creates a new app (Crawler).
func (b *Builder) Build(cfg *config.Config) *app.CrawlerApp {
	currentState := state.NewState(logging.NewLogger("state"))
	currentState.Load()

	linksToProcessChan := make(chan string)
	pendingCountChan := make(chan int)

	monitor := monitor.NewMonitor(pendingCountChan, linksToProcessChan)
	filter := filter.NewFilter(cfg.Depth, currentState, linksToProcessChan, pendingCountChan, logging.NewLogger("filter"))
	processor := link_processor.NewLinkProcessor(currentState, filter, logging.NewLogger("processor"), pendingCountChan)

	return app.NewCrawlerApp(cfg.CrawlerThreads, linksToProcessChan, pendingCountChan, currentState, processor, monitor, logging.NewLogger("crawler"))
}
