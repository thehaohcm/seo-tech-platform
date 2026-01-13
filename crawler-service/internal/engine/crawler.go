package engine

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/seo-tech-platform/crawler-service/internal/queue"
	"github.com/seo-tech-platform/crawler-service/pkg/config"
	"github.com/seo-tech-platform/crawler-service/pkg/logger"
)

type Crawler struct {
	config *config.Config
	logger *logger.Logger
	queue  *queue.RedisQueue
}

type CrawlJob struct {
	RunID     int    `json:"run_id"`
	ProjectID int    `json:"project_id"`
	StartURL  string `json:"start_url"`
	MaxPages  int    `json:"max_pages"`
}

type PageData struct {
	URL          string    `json:"url"`
	StatusCode   int       `json:"status_code"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	H1Tags       []string  `json:"h1_tags"`
	Links        []string  `json:"links"`
	LoadTime     int64     `json:"load_time_ms"`
	HTMLSnapshot string    `json:"html_snapshot"`
	Timestamp    time.Time `json:"timestamp"`
}

func NewCrawler(cfg *config.Config, logger *logger.Logger, queue *queue.RedisQueue) *Crawler {
	return &Crawler{
		config: cfg,
		logger: logger,
		queue:  queue,
	}
}

func (c *Crawler) createCollector(allowedDomain string) *colly.Collector {
	collector := colly.NewCollector(
		colly.AllowedDomains(allowedDomain),
		colly.MaxDepth(3),
		colly.Async(true),
	)

	// Set user agent
	collector.UserAgent = "SEO-Tech-Platform-Bot/1.0"

	// Set timeouts
	collector.SetRequestTimeout(30 * time.Second)

	// Limit parallelism
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	return collector
}

func (c *Crawler) ProcessJob(jobData string) error {
	var job CrawlJob
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		return fmt.Errorf("failed to parse job data: %w", err)
	}

	c.logger.Infof("Processing crawl job for RunID: %d, URL: %s", job.RunID, job.StartURL)

	// Extract domain from StartURL
	parsedURL, err := url.Parse(job.StartURL)
	if err != nil {
		return fmt.Errorf("failed to parse start URL: %w", err)
	}
	allowedDomain := parsedURL.Host

	c.logger.Infof("Crawling only URLs from domain: %s", allowedDomain)

	// Create a new collector for this job
	collector := c.createCollector(allowedDomain)

	// Setup callbacks
	c.setupCallbacks(collector, job.RunID, allowedDomain)

	// Start crawling
	if err := collector.Visit(job.StartURL); err != nil {
		return fmt.Errorf("failed to start crawling: %w", err)
	}

	// Wait for completion
	collector.Wait()

	c.logger.Infof("Completed crawl job for RunID: %d", job.RunID)
	return nil
}

func (c *Crawler) setupCallbacks(collector *colly.Collector, runID int, allowedDomain string) {
	// On HTML response
	collector.OnHTML("html", func(e *colly.HTMLElement) {
		pageData := c.extractPageData(e)
		pageData.Timestamp = time.Now()

		// Send to queue for analysis
		analysisJob := map[string]interface{}{
			"run_id":      runID,
			"url":         pageData.URL,
			"status_code": pageData.StatusCode,
			"title":       pageData.Title,
			"description": pageData.Description,
			"h1_tags":     pageData.H1Tags,
		}

		if err := c.queue.Push("analysis_queue", analysisJob); err != nil {
			c.logger.Errorf("Failed to push analysis job for %s: %v", pageData.URL, err)
		} else {
			c.logger.Infof("Extracted data from: %s", pageData.URL)
		}

		// Find and visit links (only from same domain)
		e.ForEach("a[href]", func(_ int, link *colly.HTMLElement) {
			href := link.Attr("href")
			if href != "" && c.isSameDomain(href, allowedDomain, e.Request.URL.String()) {
				link.Request.Visit(href)
			}
		})
	})

	// On request
	collector.OnRequest(func(r *colly.Request) {
		c.logger.Debugf("Visiting: %s", r.URL.String())
	})

	// On error
	collector.OnError(func(r *colly.Response, err error) {
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

// isSameDomain checks if a URL belongs to the same domain as the allowed domain
func (c *Crawler) isSameDomain(href string, allowedDomain string, baseURL string) bool {
	// Parse the href
	parsedHref, err := url.Parse(href)
	if err != nil {
		return false
	}

	// If it's a relative URL, it's same domain
	if parsedHref.Host == "" {
		return true
	}

	// If it's an absolute URL, check if host matches
	// Remove www. prefix for comparison
	hrefHost := strings.TrimPrefix(parsedHref.Host, "www.")
	allowedHost := strings.TrimPrefix(allowedDomain, "www.")

	isSame := hrefHost == allowedHost
	
	if !isSame {
		c.logger.Debugf("Skipping external URL: %s (domain: %s, allowed: %s)", href, parsedHref.Host, allowedDomain)
	}
	
	return isSame
}
