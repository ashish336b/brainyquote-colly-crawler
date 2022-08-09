package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	c := colly.NewCollector()

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

	c.Visit("https://www.brainyquote.com/")

	writeToJSON(items)
	writeToCSV(items)
}

func writeToJSON(quoteData []quotes) {
	time := time.Now().Unix()
	file, _ := json.MarshalIndent(quoteData, "", " ")

	_ = ioutil.WriteFile(fmt.Sprintf("%d.json", time), file, 0644)
	_ = ioutil.WriteFile("data.json", file, 0644)
}

func writeToCSV(quoteData []quotes) {
	jsonDataFromFile, err := ioutil.ReadFile("./data.json")

	if err != nil {
		fmt.Println(err)
	}

	// Unmarshal JSON data
	err = json.Unmarshal([]byte(jsonDataFromFile), &quoteData)

	if err != nil {
		fmt.Println(err)
	}

	csvFile, err := os.Create("./data.csv")

	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)

	writer.Write([]string{"Quote", "Author", "Source URL"})
	for _, usance := range quoteData {
		var row []string
		row = append(row, usance.Quote)
		row = append(row, usance.Author)
		row = append(row, usance.SourceURL)
		writer.Write(row)
	}

	// remember to flush!
	writer.Flush()
}
