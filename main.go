package main

import (
	"time"

	"github.com/0xe8551ccb/github-crawler/github"
)

var initialURLs = []string{
	"https://github.com/niklasvh/html2canvas",
	"https://github.com/zhoubear/open-paperless",
	"https://github.com/uber-go/go-helix",
	"https://github.com/btraceio/btrace",
	"https://github.com/go-openapi/strfmt",
	"https://github.com/aio-libs/aiohttp",
	"https://github.com/loadimpact/k6",
	"https://github.com/prometheus/prometheus",
	"https://github.com/icsharpcode/WpfDesigner",
	"https://github.com/python/cpython",
	"https://github.com/pytest-dev/pytest",
	"https://github.com/jekyll/jekyll",
	"https://github.com/m4ll0k/Infoga",
	"https://github.com/equinusocio/material-theme",
}

func main() {
	var crawler = github.New("/Users/chris/Desktop/github-projects", 24, initialURLs...)
	// Set request interval
	crawler.RequestInterval = time.Millisecond * 100
	crawler.Run(true)
}
