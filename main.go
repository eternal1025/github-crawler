package main

import "github.com/0xe8551ccb/github"

var initURLs  = []string{
	"https://github.com/niklasvh/html2canvas",
	}

func main()  {
	var c = github.GitProjectCrawler{}
	c.Init("projects", 1, initURLs...)
	c.Run()
}