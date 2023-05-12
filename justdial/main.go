package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/tealeg/xlsx"
)

type data struct {
	Name    string
	Address string
	Phone   string
}

func main() {
	var datas []data
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"),
		colly.AllowedDomains("justdial.com", "www.justdial.com"),
		colly.AllowURLRevisit(),
		colly.MaxDepth(1),
	)

	url := "https://www.justdial.com/Kasaragod/Institutes/nct-10268288"
	place := "Kasargod"

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting :", r.URL.String())

	})
	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("on error :", e.Error())
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Status Code:", r.StatusCode)
	})
	c.OnHTML(".resultbox_info", func(h *colly.HTMLElement) {
		fmt.Println("reached")
		name := h.ChildAttr("h2.resultbox_title", "title")
		address := h.ChildText(".resultbox_address")
		phone := h.ChildAttr("a.colorFFF", "href")
		if address != "" {
			address = fmt.Sprintf("%s, %s", address, place)
		} else {
			address = place
		}
		phone = strings.TrimPrefix(phone, "tel:")
		fmt.Println(
			"Name :", name,
			"\nAddress :", address,
			"\nPhone :", phone,
		)
		datas = append(datas, data{
			Name:    name,
			Address: address,
			Phone:   phone,
		})

	})
	err := c.Visit(url)
	if err != nil {
		fmt.Println("found error while visiting", err)
	}

	// push into excel
	writeXLSX(datas)
}

func writeXLSX(data []data) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("sheetname")
	if err != nil {
		panic(err)
	}

	row := sheet.AddRow()
	row.AddCell().SetValue("Name")
	row.AddCell().SetValue("Place")
	row.AddCell().SetValue("Phone")

	for _, r := range data {
		row := sheet.AddRow()
		row.AddCell().SetValue(r.Name)
		row.AddCell().SetValue(r.Address)
		row.AddCell().SetValue(r.Phone)
	}

	err = file.Save("kasargod" + ".xlsx")
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
