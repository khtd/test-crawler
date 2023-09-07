package config

import "flag"

type Config struct {
	StartLink      string
	CrawlerThreads int
	Depth          int
}

func Load() *Config {
	link := flag.String("link", "https://www.scrapethissite.com/pages/", "start point for crawling")
	crawlerThreads := flag.Int("threads", 15, "number of crawler threads to spawn")
	depth := flag.Int("depth", 1, "restrict max depth to crawl. -1 mean there are no restrictions")
	flag.Parse()

	return &Config{
		StartLink:      *link,
		CrawlerThreads: *crawlerThreads,
		Depth:          *depth,
	}
}
