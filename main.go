package main

import (
	"github.com/0xe8551ccb/github"
)

var initialURLs = []string{
	//"https://github.com/niklasvh/html2canvas",
	//"https://github.com/zhoubear/open-paperless",
	//"https://github.com/uber-go/go-helix",
	//"https://github.com/btraceio/btrace",
	//"https://github.com/tipsy/github-profile-summary",
	"https://github.com/tbroadley/github-spellcheck-cli",
	//"https://github.com/go-openapi/strfmt",
	}

func main()  {
	var c = github.GitProjectCrawler{}
	c.Init("/Users/chris/Desktop/github-projects", 8, initialURLs...)
	c.Run(true)
}