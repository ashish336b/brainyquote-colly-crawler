package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gocolly/colly"
)

type quotes struct {
	Quote     string
	Author    string
	SourceURL string
}

func main() {
	items := []quotes{}
	c := colly.NewCollector(
		colly.AllowedDomains("brainyquote.com"),
	)

	// Find list of all author and go inside each author page
	c.OnHTML("#authorColumns > div  a", func(h *colly.HTMLElement) {
		link := h.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", h.Text, link)
		c.Visit(h.Request.AbsoluteURL(link))
	})

	// scrape quote for each author
	c.OnHTML(".grid-item.qb", func(h *colly.HTMLElement) {
		quote := h.DOM.Find(".b-qt > div").Text()
		author := h.DOM.Find(".bq-aut").Text()
		// get url of page
		url := h.Request.URL.String()
		items = append(items, quotes{quote, author, url})
	})

	// add paginated url to queue if any
	c.OnHTML(".pagination a", func(h *colly.HTMLElement) {
		link := h.Attr("href")
		c.Visit(h.Request.AbsoluteURL(link))
	})

	// log visiting urls
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://www.brainyquote.com")

	time := time.Now().Unix()
	fmt.Print(time)
	file, _ := json.MarshalIndent(items, "", " ")

	_ = ioutil.WriteFile(fmt.Sprintf("%d.json", time), file, 0644)
}
