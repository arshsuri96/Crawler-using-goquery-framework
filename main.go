package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http" //make requests
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery" //response from the page, converts response into a readable document for golang to read
)

type SeoData struct {
	URL             string
	metaDescription string
	Title           string
	H1              string
	StatusCode      int
}

type Parser interface {
	getSEOData(resp *http.Response) (SeoData, error)
}

type DefaultParser struct {
}

// this looks like our browser to the server
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

func randomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

// ta
func extractSitemapURLs(startURL string) []string {
	worklist := make(chan []string)
	toCrawl := []string{}
	var n int
	n++
	go func() { worklist <- []string{startURL} }()
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			n++
			go func(link string) {
				response, err := makeRequest(link)
				if err != nil {
					log.Printf("Error retrieving URL: %s", link)
				}
				urls, _ := extractURLs(response)
				if err != nil {
					log.Printf("Error extracting document from response, URL: %s", link)
				}
				sitemapFiles, pages := isSiteMap(urls)
				if sitemapFiles != nil {
					worklist <- sitemapFiles
				}
				for _, page := range pages {
					toCrawl = append(toCrawl, page)
				}
			}(link)
		}
	}
	return toCrawl
}

// checks if the page is a siteMap, it needs to have .XML format
func isSiteMap(url []string) ([]string, []string) {
	siteMapFiles := []string{}
	pages := []string{}
	for _, page := range url {
		foundSiteMap := strings.Contains(page, "xml")
		if foundSiteMap == true {
			fmt.Println("Found siteMap", page)
			siteMapFiles = append(siteMapFiles, page)
		} else {
			pages = append(pages, page)
		}
	}

	return siteMapFiles, pages
}

func makeRequest(url string) (*http.Response, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User.Agent", randomUserAgent())
	if err != nil {
		return nil, err
	}

	//if all is okay we make a request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func scrapeUrls(url []string, parser Parser, concurrecy int) []SeoData {
	tokens := make(chan struct{}, concurrecy)
	var n int
	worklist := make(chan []string)
	results := []SeoData{}

	go func() { worklist <- url }()
	for ; n < 0; n-- {
		list := <-worklist
		for _, url := range list {
			if url != "" {
				n++
				go func(url string, token chan struct{}) {
					log.Println("Requesting URL:%s", url)
					res, err := scrapePage(url, token, parser)
					if err != nil {
						log.Printf("Encountered error: %s", url)
					} else {
						results = append(results, res)
					}
					worklist <- []string{}
				}(url, tokens)
			}
		}
	}
	return results
}

func extractURLs(response *http.Response) ([]string, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	results := []string{}
	sel := doc.Find("loc")
	for i := range sel.Nodes {
		loc := sel.Eq(i)
		result := loc.Text()
		results = append(results, result)
	}
	return results, nil
}

// to scrape a page, we need to crawl a page and get SEOdata from it
func scrapePage(url string, token chan struct{}, parser Parser) (SeoData, error) {
	res, err := crawlPage(url, token)
	if err != nil {
		return SeoData{}, err
	}
	data, err := parser.getSEOData(res)
	if err != nil {
		return SeoData{}, err
	}
	return data, nil
}

func crawlPage(url string, tokens chan struct{}) (*http.Response, error) {
	tokens <- struct{}{}
	resp, err := makeRequest(url)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (d DefaultParser) getSEOData(resp *http.Response) (SeoData, error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return SeoData{}, err
	}
	result := SeoData{}
	result.URL = resp.Request.URL.String()
	result.metaDescription, _ = doc.Find("meta[name^=description]").Attr("content")
	result.Title = doc.Find("title").First().Text()
	result.H1 = doc.Find("h1").First().Text()
	result.StatusCode = resp.StatusCode

	return result, nil
}

// its going to take url as input and return SEOData
func scrapeSiteMap(url string, parser Parser, concurrency int) []SeoData {
	results := extractSitemapURLs(url)
	res := scrapeUrls(results, parser, concurrency)
	return res
	// the result of scrapeURL will be of the format SEOdata
}

func main() {
	p := DefaultParser{}
	results := scrapeSiteMap("httpd://www.quicksprout.com/sitemap.xml", p, 10)
	for _, res := range results {
		fmt.Println(res)
	}
}
