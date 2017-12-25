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
	"https://github.com/gin-gonic/gin",
	"https://github.com/jinzhu/gorm",
	"https://github.com/alphazero/Go-Redis",
	"https://github.com/garyburd/redigo",
	"https://github.com/hashicorp/nomad",
	"https://github.com/golang/dep",
	"https://github.com/facebookincubator/create-react-app",
	"https://github.com/loadimpact/k6",
	"https://github.com/tldr-pages/tldr",
	"https://github.com/primefaces/primeng",
	"https://github.com/google/go-cmp",
	"https://github.com/prometheus/prometheus",
	"https://github.com/icsharpcode/WpfDesigner",
	"https://github.com/python/cpython",
	"https://github.com/pytest-dev/pytest",
	"https://github.com/jekyll/jekyll",
	"https://github.com/m4ll0k/Infoga",
	"https://github.com/equinusocio/material-theme",
	"https://github.com/aa112901/remusic",
}

func main() {
	var c = github.GitProjectCrawler{}
	c.Init("/Users/chris/Desktop/github-projects", 24, initialURLs...)
	// Set request interval
	c.RequestInterval = time.Millisecond * 100
	c.Run(true)
}
