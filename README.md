GitHub Project Crawler
==========================

A simple crawler written in Golang 1.9.

# Simple entry

```
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
```


# Features

1. Auto retry any failed requests;
2. Auto shutdown when there'is no pending request;
3. Simple auto-throttling strategies.

# How to run it?

1. Clone this project to your local disk:

    ```
    git clone https://githu.com/0xe8551ccb/github-crawler
    ```

2. Install glide:

    ```
    # On Mac
    brew install glide

    # Or
    curl https://glide.sh/get | sh
    ```

3. Install requirements:

    ```
    glide install
    ```

4. Start crawler:

    ```
    go run main.go
    ```

# Snapshot

![](http://blog.chriscabin.com/wp-content/uploads/2017/12/Snip20171223_17.png)