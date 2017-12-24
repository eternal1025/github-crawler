// 非常粗糙的爬虫框架，极其简陋，以至于我都不知道能不能正常工作
// 以下功能：
// 1. 中间件
// 2. 登录组件
// 3. 反爬虫组件
// 均不支持
// 后续扩展时，应该抽出一个 engine 来，单独负责爬虫请求调度吧
package gocrawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"sync"
	"time"
)

type HandleFunc func()

type GoCrawler struct {
	// 每个 Crawler 尽可能搞个个性化的名字吧，虽然我还没用到它
	Name                string
	MaxTryCount         int
	RequestInterval     time.Duration
	idleCount           int
	maxWorkers          int
	worklist            chan []*Request
	tokens              chan struct{}
	seen                map[string]bool
	items               chan []interface{}
	itemHandler         func(items []interface{})
	initialRequests     []*Request
	pendingWorkersCount int
	pendingMutex        sync.Mutex // 只有写才会产生竞争
	anyWorkerStarted    bool
}

func (c *GoCrawler) Init(name string, maxWorkers int, initialRequests []*Request, itemHandler func(items []interface{})) {
	c.Name = name
	c.MaxTryCount = 12
	if int(c.RequestInterval) == 0 {
		c.RequestInterval = 1000 * time.Millisecond
	}

	c.maxWorkers = maxWorkers
	c.itemHandler = itemHandler
	c.initialRequests = initialRequests
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
	// 如果没有任何 worker 启动过，则始终等待（针对爬虫第一次启动时，特殊处理）
	if !c.anyWorkerStarted {
		return false
	}

	return c.pendingWorkersCount == 0
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
			c.incrWorker()
			defer c.decrWorker()
			requests, items, err := c.crawl(r)
			if err != nil {
				log.Printf("Failed to crawl %s: %s", r, err)
				return
			}
			c.worklist <- requests
			c.items <- items
		}(req)
	}
	// sleep for a while, it's not polite to open too many requests at the same time~
	time.Sleep(c.RequestInterval)
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
		return c.retry(req)
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
		log.Printf("Invalid response %s", &resp)
		return c.retry(req)
	}

	if req.Parser != nil {
		return req.Parser(&resp)
	}
	return nil, nil, fmt.Errorf("missing request parser: %s", req)
}

func (c *GoCrawler) retry(req *Request) ([]*Request, []interface{}, error) {
	if req.triedCount > c.MaxTryCount {
		return nil, nil, fmt.Errorf("max tries exceed, ignore request %s", req)
	}
	req.triedCount += 1
	c.seen[req.URL] = false
	delay := time.Duration(req.triedCount*req.triedCount) * time.Second
	log.Printf("Retry request %s after %d seconds: +%d", req, int(delay.Seconds()), req.triedCount)
	time.Sleep(delay)
	return []*Request{req}, nil, nil
}

func (c *GoCrawler) incrWorker() {
	c.pendingMutex.Lock()
	defer c.pendingMutex.Unlock()
	c.pendingWorkersCount += 1
}

func (c *GoCrawler) decrWorker() {
	c.pendingMutex.Lock()
	defer c.pendingMutex.Unlock()
	if c.pendingWorkersCount > 0 {
		c.pendingWorkersCount -= 1
	}
}
func (c *GoCrawler) String() string {
	return fmt.Sprintf("GoCrawler(name=%s, workers=%d, worklist=%d)",
		c.Name, c.maxWorkers, len(c.worklist))
}
