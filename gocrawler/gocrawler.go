// 非常粗糙的爬虫框架，极其简陋，以至于我都不知道能不能正常工作
// 以下功能：
// 1. 中间件
// 2. 登录组件
// 3. 反爬虫组件
// 均不支持
package gocrawler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"time"
	"sync"
)

type Request struct {
	triedCount int
	URL        string
	Parser     func(resp *Response) (requests []*Request, items []interface{}, err error)
	Meta       map[string]interface{}
}

func (r *Request) String() string {
	return fmt.Sprintf("<Request url=%s>", r.URL)
}

type Response struct {
	Doc        *goquery.Document
	Request    *Request
	StatusCode int
	Body       io.ReadCloser
}

// Find is a short-curt method to invoke corresponding method of Doc object
func (r *Response) Find(selector string) *goquery.Selection {
	return r.Doc.Find(selector)
}

func (r *Response) String() string {
	return fmt.Sprintf("<Response url=%s, status=%d>", r.Request.URL, r.StatusCode)
}

type GoCrawler struct {
	Name            string
	maxWorkers      int
	MaxTryCount     int
	idleCount       int
	worklist        chan []*Request
	tokens          chan struct{}
	seen            map[string]bool
	items           chan []interface{}
	itemHandler     func(items []interface{})
	initialRequests []*Request
	pendingWorkers  sync.Map
	anyWorkerStarted bool
}

func (c *GoCrawler) Init(name string, maxWorkers int, initialRequests []*Request, itemHandler func(items []interface{})) {
	c.Name = name
	c.maxWorkers = maxWorkers
	c.itemHandler = itemHandler
	c.initialRequests = initialRequests
	c.MaxTryCount = 12
}

func (c *GoCrawler) Run(stopWhenIdle bool) {
	if len(c.initialRequests) == 0 {
		log.Fatal("failed to start GoCrawler: empty initial request")
	}

	log.Printf("GoCrawler is running with %d initial requests", len(c.initialRequests))
	c.worklist = make(chan []*Request)
	c.seen = make(map[string]bool)
	c.tokens = make(chan struct{}, c.maxWorkers)
	c.items = make(chan []interface{})

	go func() {
		// feed initial requests for crawlers
		c.worklist <- c.initialRequests
	}()

	for {
		select {
		case requests := <-c.worklist:
			c.crawlRequests(requests)
		case items := <-c.items:
			if c.itemHandler != nil && len(items) > 0 {
				c.itemHandler(items)
			}
		default:
			if c.IsIdle() {
				if stopWhenIdle {
					log.Println("Shutdown GoCrawler gracefully~")
					return
				}
				log.Println("GoCrawler is idle, waiting to be feed with new requests~")
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// IsIdle 函数用于判断 Crawler 是否处于空闲状态
// 当 Crawler 为 Idle 时，必须满足如下几个条件：
// 1. 要处理的 `worklist` 为空
// 2. 要处理的 `items` 为空
// 3. 没有任何在等待中的 worker goroutine
func (c *GoCrawler) IsIdle() bool {
	// 如果没有任何 worker 启动过，则始终等待
	if !c.anyWorkerStarted {
		return false
	}

	hasPendingWorkers := false
	// check if there is any pending workers
	c.pendingWorkers.Range(func(key, value interface{}) bool {
		hasPendingWorkers = true
		return false
	})

	if hasPendingWorkers {
		return false
	}

	return true
}

func (c *GoCrawler) crawlRequests(requests []*Request) {
	for _, req := range requests {
		if c.seen[req.URL] {
			log.Printf("Ingore duplicate request: %s", req)
			continue
		}

		c.seen[req.URL] = true
		go func(r *Request) {
			// 标记一下，当任何一个 worker 启动过就可以标记了，不用考虑写冲突的问题
			c.anyWorkerStarted = true
			c.pendingWorkers.Store(r.URL, struct{}{})
			defer func() { c.pendingWorkers.Delete(r.URL) }()
			requests, items, err := c.crawl(r)
			if err != nil {
				log.Printf("Failed to crawl: %s", r)
				return
			}
			c.worklist <- requests
			c.items <- items
		}(req)
	}
	// sleep for a while, too many requests are not allowed
	time.Sleep(720 * time.Millisecond)
}

func (c *GoCrawler) crawl(req *Request) ([]*Request, []interface{}, error) {
	log.Printf("Crawling request %s", req)
	// acquire a token
	c.tokens <- struct{}{}
	// release token later
	defer func() { <-c.tokens }()

	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	r, err := client.Get(req.URL)
	if err != nil {
		log.Printf("Failed to fetch request %s:%s", req, err)
		c.seen[req.URL] = false
		return []*Request{req}, nil, nil
	}

	var resp Response
	resp.Request = req
	resp.StatusCode = r.StatusCode
	resp.Body = r.Body
	resp.Doc, err = goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		// We'll try this request later
		log.Printf("Invalid response %s", resp)
		return c.retry(req)
	}

	if req.Parser != nil {
		return req.Parser(&resp)
	}
	return nil, nil, fmt.Errorf("missing request parser: %s", req)
}

func (c *GoCrawler) retry(req *Request) ([]*Request, []interface{}, error) {
	if req.triedCount > c.MaxTryCount {
		log.Printf("Max tries exceed, ignore request %s", req)
		return nil, nil, nil
	}
	req.triedCount += 1
	c.seen[req.URL] = false
	delay := time.Duration(req.triedCount*req.triedCount) * time.Second
	log.Printf("Retry request %s after %d seconds: +%d", req, delay, req.triedCount)
	time.Sleep(delay)
	return []*Request{req}, nil, nil
}

func (c *GoCrawler) String() string {
	return fmt.Sprintf("GoCrawler(name=%s, workers=%d, worklist=%d)",
		c.Name, c.maxWorkers, len(c.worklist))
}
