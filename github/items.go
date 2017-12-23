// 目标是抓取 GitHub 指定项目
package github

import (
	"encoding/json"
	"fmt"
	"time"
)

type IssueItem struct {
	ProjectName  string    `json:"project_name"`
	Number       string    `json:"number"`
	Title        string    `json:"title"`
	Issuer       string    `json:"issuer"`
	IsOpen       bool      `json:"is_open"`
	IssuedAt     time.Time `json:"issued_at"`
	CommentCount int       `json:"comment_count"`
	Labels       []string  `json:"labels"`
}

func (i *IssueItem) String() string {
	buf, err := json.Marshal(i)
	if err != nil {
		return fmt.Sprintf("%s: %d comments", i.Title, i.CommentCount)
	}
	return string(buf)
}

type ProjectItem struct {
	Name              string `json:"name"`
	URL               string `json:"url"`
	Stars             int    `json:"stars"`
	Forks             int    `json:"forks"`
	Watches           int    `json:"watches"`
	OpenedIssuesCount int    `json:"opened_issues_count"`
	ClosedIssuesCount int    `json:"closed_issues_count"`
}

func (p *ProjectItem) String() string {
	buf, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("%s: %d watches, %d stars, %d forks", p.Name, p.Watches, p.Stars, p.Watches)
	}
	return string(buf)
}
