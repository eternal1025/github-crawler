package github

import (
	"log"
	"path"
	"strings"
	"time"
	"github.com/PuerkitoBio/goquery"
	"github.com/0xe8551ccb/gocrawler"
	"github.com/0xe8551ccb/utils"
)

type GitProjectCrawler struct {
	gocrawler.GoCrawler
	storeLocation string
}

// Init 方法初始化爬虫，同时可以指定存储抓取内容的位置，最大的 worker 数量等
func (p *GitProjectCrawler) Init(storeLocation string, maxWorkers int, urls ...string) {
	p.storeLocation = storeLocation

	var requests []*gocrawler.Request
	for _, url := range urls {
		requests = append(requests, &gocrawler.Request{URL: url, Parser: p.ParseHome})
	}
	p.GoCrawler.Init("GitHub project crawler", maxWorkers, requests, p.ProcessItems)
}

// 解析项目主页，提取必要的信息
func (p *GitProjectCrawler) ParseHome(resp *gocrawler.Response) ([]*gocrawler.Request, []interface{}, error) {
	log.Printf("Parse project home from response: %s", resp)

	var project ProjectItem
	_, project.Name = path.Split(resp.Request.URL)
	project.URL = resp.Request.URL
	project.Watches = utils.StrToInt(resp.Find("a.social-count").Eq(0).Text())
	project.Stars = utils.StrToInt(resp.Find("a.social-count").Eq(1).Text())
	project.Forks = utils.StrToInt(resp.Find("a.social-count").Eq(2).Text())

	meta := make(map[string]interface{})
	meta["project"] = &project

	issueCountReq := gocrawler.Request{
		URL:    project.URL + "/issues",
		Parser: p.ParseIssueCount,
		Meta:   meta}
	openIssueReq := gocrawler.Request{
		URL:    project.URL + "/issues?q=is%3Aopen+is%3Aissue",
		Parser: p.ParseIssueItems,
		Meta:   meta}
	closedIssueReq := gocrawler.Request{
		URL:    project.URL + "/issues?q=is%3Aclosed+is%3Aissue",
		Parser: p.ParseIssueItems,
		Meta:   meta}

	requests := []*gocrawler.Request{&issueCountReq, &openIssueReq, &closedIssueReq}

	return requests, nil, nil
}

func (p *GitProjectCrawler) ParseIssueCount(resp *gocrawler.Response) ([]*gocrawler.Request, []interface{}, error) {
	log.Printf("Parse issue count from response: %s", resp)
	project, ok := resp.Request.Meta["project"]
	if !ok {
		return nil, nil, nil
	}

	switch proj := project.(type) {
	case *ProjectItem:
		openedCountText := resp.Find("div.table-list-header-toggle.states.float-left.pl-3 > a").First().Text()
		proj.OpenedIssuesCount = utils.StrToInt(strings.Replace(openedCountText, "Open", "", len(openedCountText)))

		closedCountText := resp.Find("div.table-list-header-toggle.states.float-left.pl-3 > a").Last().Text()
		proj.ClosedIssuesCount = utils.StrToInt(strings.Replace(closedCountText, "Closed", "", len(closedCountText)))
	}

	var items []interface{}
	items = append(items, project)

	return nil, items, nil
}

// 解析 issue 列表页面，当然还要处理分页问题
func (p *GitProjectCrawler) ParseIssueItems(resp *gocrawler.Response) ([]*gocrawler.Request, []interface{}, error) {
	log.Printf("Parse issue items from response: %s", resp)
	project, ok := resp.Request.Meta["project"]
	if !ok {
		return nil, nil, nil
	}

	isOpen := strings.Contains(strings.ToLower(resp.Request.URL), "open")
	projectName := project.(*ProjectItem).Name

	var items []interface{}

	resp.Find(`ul.js-navigation-container.js-active-navigation-container > li`).Each(func(i int, sel *goquery.Selection) {
		item := IssueItem{}
		item.IsOpen = isOpen
		item.ProjectName = projectName
		item.Title = strings.TrimSpace(sel.Find("div > div.float-left.col-9.p-2.lh-condensed > a").Text())

		if text, ok := sel.Attr("id"); ok {
			item.Number = "#" + strings.Replace(text, "issue_", "", len(text))
		}

		item.CommentCount = utils.StrToInt(sel.Find("a.muted-link > span").Text())
		if text, ok := sel.Find("relative-time").Attr("datetime"); ok {
			item.IssuedAt, _ = time.Parse(time.RFC3339, text)
		}

		item.Issuer = sel.Find("span.opened-by > a").Text()

		sel.Find("span.labels").Each(func(i int, sel *goquery.Selection) {
			item.Labels = append(item.Labels, sel.Find("a").Text())
		})

		items = append(items, &item)
	})

	var requests []*gocrawler.Request
	if nextPage, ok := resp.Find("a.next_page").Attr("href"); ok {
		requests = append(requests, &gocrawler.Request{
			URL:    "https://github.com" + nextPage,
			Parser: p.ParseIssueItems,
			Meta:   resp.Request.Meta,
		})
	}

	return requests, items, nil
}

// 解析 issue 详情页
func (p *GitProjectCrawler) ParseIssueDetails(resp *gocrawler.Response) ([]*gocrawler.Request, []interface{}, error) {
	log.Printf("Parse issue details from response: %s", resp)
	return nil, nil, nil
}

// 处理解析后的 items，可以选择存储，或者打印
func (p *GitProjectCrawler) ProcessItems(items []interface{}) {
	for _, item := range items {
		switch val := item.(type) {
		case *ProjectItem:
			log.Printf("Got project: %s\n", val)
			SaveProjectItem(p.storeLocation, val)
		case *IssueItem:
			log.Printf("Got issue: %s\n", val)
			SaveIssueItem(p.storeLocation, val)
		default:
			log.Printf("Got unknown item: %v\n", item)
		}
	}
}
