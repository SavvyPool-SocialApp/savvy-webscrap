package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/tealeg/xlsx"
)

type data struct {
	College string `json:"college"`
	Place   string `json:"place"`
	Link    string `json:"link"`
}

var s int
var arr []data

func main() {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"),
		colly.AllowedDomains("www.4icu.org"),
		colly.AllowURLRevisit(),
	)
	c2 := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"),
		colly.AllowedDomains("www.4icu.org"),
		colly.AllowURLRevisit(),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting :", r.URL.String())
	})
	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("on error :", e.Error())
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response Code :", r.StatusCode)
	})
	c.OnHTML(".panel", func(h *colly.HTMLElement) {
		h.ForEach("tr", func(i int, n *colly.HTMLElement) {
			if i > 1 {
				var clg, place, link string
				val := strings.Split(n.Text, "\n")
				clg = val[2]
				place = val[3]
				link = n.ChildAttr("a", "href")
				arr = append(arr, data{
					College: clg,
					Place:   place,
					Link:    link,
				})
			}
		})

	})
	c2.OnHTML("a[itemprop=url],[target=blank]", func(h *colly.HTMLElement) {
		link := h.Attr("href")
		arr[s].Link = link
	})
	err := c.Visit("https://www.4icu.org/us/")
	if err != nil {
		fmt.Println(err)
	}

	for i, val := range arr {
		s = i
		c2.Visit("https://www.4icu.org/" + val.Link)

		fmt.Println(s, " ", arr[s].Link)
	}

	writeXLSX(arr)

}

func writeXLSX(data []data) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("sheet1")
	if err != nil {
		panic(err)
	}

	row := sheet.AddRow()
	row.AddCell().SetValue("Name")
	row.AddCell().SetValue("Place")
	row.AddCell().SetValue("Link")

	for _, r := range data {
		row := sheet.AddRow()
		row.AddCell().SetValue(r.College)
		row.AddCell().SetValue(r.Place)
		row.AddCell().SetValue(r.Link)
	}

	err = file.Save("us" + ".xlsx")
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
