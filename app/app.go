package app

import (
	"sync"
	"test-crawler/link_processor"
	"test-crawler/logging"
	"test-crawler/monitor"
	"test-crawler/state"
)

type CrawlerApp struct {
	linksToProcessChan chan string
	pendingCountChan   chan int
	state              *state.State
	monitor            *monitor.Monitor
	processor          *link_processor.LinkProcessor
	logger             *logging.Logger
	threads            int
}

func NewCrawlerApp(
	threads int,
	linksToProcessChan chan string,
	pendingCountChan chan int,
	state *state.State,
	processor *link_processor.LinkProcessor,
	monitor *monitor.Monitor,
	logger *logging.Logger,
) *CrawlerApp {
	return &CrawlerApp{
		linksToProcessChan: linksToProcessChan,
		pendingCountChan:   pendingCountChan,
		state:              state,
		processor:          processor,
		monitor:            monitor,
		logger:             logger,
		threads:            threads,
	}
}

func (c *CrawlerApp) Run(startLink string) {
	c.logger.Info("start crawling")

	go c.monitor.Start()

	unprocessedLinks := c.state.GetUnprocessed(startLink)
	if len(unprocessedLinks) > 0 {
		go func(unprocessedLinks []string) {
			for _, link := range unprocessedLinks {
				c.linksToProcessChan <- link
				c.pendingCountChan <- 1
			}
		}(unprocessedLinks)

		var wg sync.WaitGroup

		for i := 0; i < c.threads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for link := range c.linksToProcessChan {
					err := c.processor.Process(link)
					if err != "" {
						c.state.MarkProcessed(link, make([]string, 0), err)
					}
					c.pendingCountChan <- -1
				}
			}()
		}

		wg.Wait()
	}

	c.logger.Info("finished crawling")
	c.state.Save()
	c.state.Print()
}
