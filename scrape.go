package main

import (
	"github.com/gocolly/colly"
)

//TODO: make a data model, DATA: NextCheck, []items;

type LuxionItem struct {
	Name string
	Link string
}

func Scrape() []LuxionItem {
	c := colly.NewCollector()

	items := make([]LuxionItem, 0, 6)

	c.OnHTML("div.row:nth-child(2)", func(e *colly.HTMLElement) {
		e.ForEach("div.col-6", func(_ int, row *colly.HTMLElement) {
			link, _ := row.DOM.Find("a").Attr("href")
			item := LuxionItem{
				Name: row.DOM.Find("div.card-header").Find("a").Text(),
				Link: link,
			}
			items = append(items, item)
		})
	})
	c.Visit("https://trovesaurus.com/luxion")
	return items
}

func Join(items []LuxionItem) string {
	str := ""
	for i := 0; i < len(items); i++ {
		if i == len(items)-1 {
			str += items[i].Name
			continue
		}
		str += items[i].Name + ", "
	}
	return str
}
