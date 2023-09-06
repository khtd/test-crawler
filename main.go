package main

import (
	"test-crawler/builder"
	"test-crawler/config"
)

func main() {
	cfg := config.Load()
	crawlerApp := builder.NewBuilder().Build(cfg)
	crawlerApp.Run(cfg.StartLink)
}
