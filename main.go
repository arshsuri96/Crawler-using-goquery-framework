package main

import (
	"fmt"
	"log"
	"net/http" //make requests
	"strings"
	"math/rand"
	"time"
	"github.com/PuerkitoBio/goquery" //response from the page, converts response into a readable document for golang to read

)

type SeoData struct {
	URL             string
	metaDescription string
	Title           string
	H1              string
	StatuCode       int
}

type parser interface {
}

type DefaultParser struct {

}


//this looks like our browser to the server
var userAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
}

func randomUserAgent() string{
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

// ta
func extractSiteMapURL(startURL string) []string {
	Worklist := make(chan []string)
	toCrawl := []string{}
	var n int
	n++
	go func(worklist <- []string{startURL})()

	for ; n>0 ; n--

	list := <- Worklist //worklist is a channel where we are publishing all the links to be scrapped 
	for _, link := range list{
	go func(link string){ //we are using go routines to make requests simultaneously, not one by one [ concurrency ]
		response, err = makeRequest(link)
		if err != nil {
			log.Printf("Error recieving URL: %s", link)
		}
		url, _ = extractURLs(response) 
		if err!= nil {
			log.Printf("Error extracting document from response, URL: %s", link)
		}
		sitemapFiles, pages := isSiteMap(urls)
		if sitemapFiles != nil {
			worklist <- sitemapFiles
		}  
		for _, page = range pages {
			toCrawl = append(toCrawl, page)
		}
	}(link)
	}
	return toCrawl //slice of urls to crawl
}
//checks if the page is a siteMap, it needs to have .XML format
func isSiteMap(url []string)([]string, []string){
	 siteMapFiles := []string{}
	 pages := []string{}
	 for _, page := range url{
		 foundSiteMap == true {
			 fmt.Println("Found siteMap", page)
			 siteMapFiles = append(siteMapFiles, page)
		 }else{
			 pages = append(pages, page)
		 }
	 }

	 return siteMapFiles, pages
	}

func makeRequest(url string)(*http.Response, error) {
	http.Client{
		Timeout : 10*time.Second, 
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User.Agent", randomUserAgent())
	if err!=nil {	 
		return nil, err
	}

	//if all is okay we make a request
	res, err := client.Do(req)
	if err!= nil {
		return nil, err
	}
	return res, nil
	
}


// in scrape page, we crawl the page and then we will get SEOdata from it
func scrapeURL(url string, parser Parser)(SeoData, error){
	res, err := crawlPage(url)
	if err != nil {
		return SeoData{}, err 
	}
	data, err := parser.getSEOData(res)
	if err!= nil {
		//empty slice of strings SeoData will be returned
		return SeoData{}, err
	}
	return data, nil
}

func crawlPage() {

}


//to scrape a page, we need to crawl a page and get SEOdata from it
func scrapPage(url string, parser Parser)(SeoData, error) {
	res, err := crawlPage(url)
	if err != nil {
		return SeoData{}, err
	}
	data, err := parser.getSEOData(res)
	if err!=nil{
		return SeoData{}, error
	}
	return data, nil
}

// its going to take url as input and return SEOData
func scrapeSiteMap(url string) []SeoData {
	results := extractSitemapURL(url)
	res := scrapeURL(results)

	//the result of scrapeURL will be of the format SEOdata
	return res

}

func main() {
	p := DefaultParser{}
	results := scrapeSiteMap(" ")
	for _, res := range results {

	}

}
