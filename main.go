package main

import "github.com/0xe8551ccb/github"

var initURLs  = []string{
	"https://github.com/niklasvh/html2canvas",
	"https://github.com/zhoubear/open-paperless",
	"https://github.com/uber-go/go-helix",
	"https://github.com/btraceio/btrace",
	"https://github.com/tipsy/github-profile-summary",
	"https://github.com/tbroadley/github-spellcheck-cli",
	"https://github.com/go-openapi/strfmt",
	}

func main()  {
	var c = github.GitProjectCrawler{}
	c.Init("projects", 4, initURLs...)
	c.Run()
}