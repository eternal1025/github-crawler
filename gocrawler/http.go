package gocrawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
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
}

// Find is a short-curt method to invoke corresponding method of Doc object
func (r *Response) Find(selector string) *goquery.Selection {
	return r.Doc.Find(selector)
}

func (r *Response) String() string {
	return fmt.Sprintf("<Response url=%s, status=%d>", r.Request.URL, r.StatusCode)
}
