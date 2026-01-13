package engine

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/seo-tech-platform/crawler-service/pkg/config"
	"github.com/seo-tech-platform/crawler-service/pkg/logger"
)

type Crawler struct {
	collector *colly.Collector
	config    *config.Config
	logger    *logger.Logger
}

type CrawlJob struct {
	RunID     int    `json:"run_id"`
	ProjectID int    `json:"project_id"`
	StartURL  string `json:"start_url"`
	MaxPages  int    `json:"max_pages"`
}

type PageData struct {
	URL          string            `json:"url"`
	StatusCode   int               `json:"status_code"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	H1Tags       []string          `json:"h1_tags"`
	Links        []string          `json:"links"`
	LoadTime     int64             `json:"load_time_ms"`
	HTMLSnapshot string            `json:"html_snapshot"`
	Timestamp    time.Time         `json:"timestamp"`
}

func NewCrawler(cfg *config.Config, logger *logger.Logger) *Crawler {
	c := colly.NewCollector(
		colly.AllowedDomains(),
		colly.MaxDepth(3),
		colly.Async(true),
	)

	// Set user agent
	c.UserAgent = "SEO-Tech-Platform-Bot/1.0"

	// Set timeouts
	c.SetRequestTimeout(30 * time.Second)

	// Limit parallelism
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	return &Crawler{
		collector: c,
		config:    cfg,
		logger:    logger,
	}
}

func (c *Crawler) ProcessJob(jobData string) error {
	var job CrawlJob
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		return fmt.Errorf("failed to parse job data: %w", err)
	}

	c.logger.Infof("Processing crawl job for RunID: %d, URL: %s", job.RunID, job.StartURL)

	// Setup callbacks
	c.setupCallbacks(job.RunID)

	// Start crawling
	if err := c.collector.Visit(job.StartURL); err != nil {
		return fmt.Errorf("failed to start crawling: %w", err)
	}

	// Wait for completion
	c.collector.Wait()

	c.logger.Infof("Completed crawl job for RunID: %d", job.RunID)
	return nil
}

func (c *Crawler) setupCallbacks(runID int) {
	// On HTML response
	c.collector.OnHTML("html", func(e *colly.HTMLElement) {
		pageData := c.extractPageData(e)
		pageData.Timestamp = time.Now()

		// TODO: Send to queue for analysis
		data, _ := json.Marshal(pageData)
		c.logger.Infof("Extracted data from: %s", pageData.URL)
		c.logger.Debugf("Page data: %s", string(data))

		// Find and visit links
		e.ForEach("a[href]", func(_ int, link *colly.HTMLElement) {
			href := link.Attr("href")
			if href != "" {
				link.Request.Visit(href)
			}
		})
	})

	// On request
	c.collector.OnRequest(func(r *colly.Request) {
		c.logger.Debugf("Visiting: %s", r.URL.String())
	})

	// On error
	c.collector.OnError(func(r *colly.Response, err error) {
		c.logger.Errorf("Error visiting %s: %v", r.Request.URL, err)
	})
}

func (c *Crawler) extractPageData(e *colly.HTMLElement) *PageData {
	// Extract title
	title := e.ChildText("title")

	// Extract meta description
	description := e.ChildAttr("meta[name='description']", "content")

	// Extract H1 tags
	var h1Tags []string
	e.ForEach("h1", func(_ int, h1 *colly.HTMLElement) {
		h1Tags = append(h1Tags, h1.Text)
	})

	// Extract all links
	var links []string
	e.ForEach("a[href]", func(_ int, link *colly.HTMLElement) {
		href := link.Attr("href")
		if href != "" {
			links = append(links, href)
		}
	})

	return &PageData{
		URL:         e.Request.URL.String(),
		StatusCode:  e.Response.StatusCode,
		Title:       title,
		Description: description,
		H1Tags:      h1Tags,
		Links:       links,
	}
}
