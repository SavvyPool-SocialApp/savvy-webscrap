package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/tealeg/xlsx"
)

type item struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Place string `json:"place"`
	Type  string `json:"type"`
}

var (
	count      int
	datas      []item
	totalPages = 4
	// visitLink  = "https://www.shiksha.com/engineering/colleges/b-tech-colleges-kerala"
	visitLink = "https://www.shiksha.com/engineering/colleges/b-tech-colleges-delhi-other"
	tailLink  = "?ct[]=74&ct[]=10653&ed[]=et_20&uaf[]=base_course&uaf[]=city&rf=filters"
	params    = ""
)

func main() {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"),
		colly.AllowedDomains("www.shiksha.com"),
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
	c.OnHTML("._8165 ", func(h *colly.HTMLElement) {

		count++
		clgName := h.ChildText("h3[title]")
		place := h.ChildText("span._5588")
		types := h.ChildText("span")

		if strings.Contains(types, "Govt") {
			types = "Govt"
		} else if strings.Contains(types, "Pvt") {
			types = "Pvt"
		} else {
			types = ""
		}

		item := item{
			Count: count,
			Name:  clgName,
			Place: place,
			Type:  types,
		}
		datas = append(datas, item)

	})

	for i := 0; i < totalPages; i++ {
		if i != 0 {
			params = strconv.Itoa(i + 1)
			params = "-" + params + tailLink
		}else{
			params = params + tailLink
		}
		err := c.Visit(visitLink + params)
		if err != nil {
			fmt.Println("error in visit :", err.Error())
		}
	}

	writeJSON(datas)
	writeXLSX(datas)
	writeTEXT(datas)
}

func writeJSON(data []item) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	os.WriteFile("data.json", file, 0644)
}
func writeXLSX(data []item) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		panic(err)
	}

	row := sheet.AddRow()
	row.AddCell().SetValue("Name")
	row.AddCell().SetValue("District")
	row.AddCell().SetValue("Type")

	for _, r := range data {
		row := sheet.AddRow()
		row.AddCell().SetValue(r.Name)
		row.AddCell().SetValue(r.Place)
		row.AddCell().SetValue(r.Type)
	}

	err = file.Save("data.xlsx")
	if err != nil {
		panic(err)
	}
}

func writeTEXT(data []item) {
	file, err := os.Create("data.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	for _, r := range data {
		writer.WriteString("Name :" + r.Name + "\n")
		writer.WriteString("District :" + r.Place + "\n")
		writer.WriteString("Type :" + r.Type + "\n \n")
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}
}
