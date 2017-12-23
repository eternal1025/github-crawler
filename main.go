package main

import (
	"github.com/0xe8551ccb/github"
	"time"
)

var initialURLs = []string{
	"https://github.com/niklasvh/html2canvas",
	"https://github.com/zhoubear/open-paperless",
	"https://github.com/uber-go/go-helix",
	"https://github.com/btraceio/btrace",
	"https://github.com/tipsy/github-profile-summary",
	"https://github.com/tbroadley/github-spellcheck-cli",
	"https://github.com/go-openapi/strfmt",
	"https://github.com/aio-libs/aiohttp",
	"https://github.com/envoyproxy/envoy",
	"https://github.com/dgraph-io/dgraph",
	"https://github.com/OfficeDev/office-ui-fabric-react",
	"https://github.com/tensorflow/tensorflow",
	}

func main()  {
	var c = github.GitProjectCrawler{}
	c.Init("/Users/chris/Desktop/github-projects", 12, initialURLs...)
	// Set request interval
	c.RequestInterval = time.Second * 1
	c.Run(true)
}