// 非常粗糙的爬虫框架，极其简陋，以至于我都不知道能不能正常工作
// 以下功能：
// 1. 中间件
// 2. 登录组件
// 3. 反爬虫组件
// 4. 失败重试
// 5. 友好提示
// 均不支持
package gocrawler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

type Request struct {
	URL    string
	Parser func(resp *Response) (requests []*Request, items []interface{}, err error)
	Meta   map[string]interface{}
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
	worklist        chan []*Request
	tokens          chan struct{}
	seen            map[string]bool
	items           chan []interface{}
	itemHandler     func(items []interface{})
	initialRequests []*Request
}

func (c *GoCrawler) Init(name string, maxWorkers int, initialRequests []*Request, itemHandler func(items []interface{})) {
	c.Name = name
	c.maxWorkers = maxWorkers
	c.itemHandler = itemHandler
	c.initialRequests = initialRequests
}

func (c *GoCrawler) Run() {
	log.Printf("Go crawler is running with %d initial requests", len(c.initialRequests))
	c.worklist = make(chan []*Request)
	c.seen = make(map[string]bool)
	c.tokens = make(chan struct{}, c.maxWorkers)
	c.items = make(chan []interface{})

	go func() {
		// feed initial requests for crawlers
		c.worklist <- c.initialRequests
	}()

	// process items in another goroutine
	go func() {
		for items := range c.items {
			if c.itemHandler != nil && len(items) > 0 {
				c.itemHandler(items)
			}
		}
	}()

	for list := range c.worklist {
		for _, req := range list {
			if c.seen[req.URL] {
				log.Printf("Ingore duplicate request: %s", req)
				continue
			}

			c.seen[req.URL] = true
			go func(r *Request) {
				requests, items, err := c.crawl(r)
				if err != nil {
					log.Printf("Failed to crawl: %s", r)
					return
				}
				c.worklist <- requests
				c.items <- items
			}(req)
		}
	}

	log.Println("Go crawler stopped")
}

func (c *GoCrawler) crawl(req *Request) (requests []*Request, items []interface{}, err error) {
	log.Printf("Crawling request %s", req)
	// acquire a token
	c.tokens <- struct{}{}
	// release token later
	defer func() { <-c.tokens }()
	r, err := http.Get(req.URL)
	if err != nil {
		log.Printf("Failed to fetch request %s:%s", req, err)
	}

	var resp Response
	resp.Request = req
	resp.StatusCode = r.StatusCode
	resp.Body = r.Body
	resp.Doc, err = goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return nil, nil, err
	}

	if req.Parser != nil {
		return req.Parser(&resp)
	}
	return nil, nil, fmt.Errorf("missing request parser: %s", req)
}

func (c *GoCrawler) String() string {
	return fmt.Sprintf("GoCrawler(name=%s, workers=%d, worklist=%d)",
		c.Name, c.maxWorkers, len(c.worklist))
}
