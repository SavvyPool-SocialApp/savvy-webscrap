package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/tealeg/xlsx"
)

type item struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Place   string `json:"place"`
	Type    string `json:"type"`
	Contact string `json:"contact"`
}

var (
	data       []item
	Xlfile     *xlsx.File
	initial    string
	count      int
	totalPages = 100
	sheetName  = "sheet1"
	fileName   = "data"
	visitLink  string
	tailLink   string
	params     string
	// visitLink  = "https://www.shiksha.com/engineering/colleges/b-tech-colleges-kerala"
	// visitLink = "https://www.shiksha.com/engineering/colleges/b-tech-colleges-delhi-other"
	// tailLink  = "?ct[]=74&ct[]=10653&ed[]=et_20&uaf[]=base_course&uaf[]=city&rf=filters"
)

func main() {

	fmt.Println("Enter the link head :")
	fmt.Scan(&visitLink)
	fmt.Println("Enter the link tail (Enter 0 if not) :")
	fmt.Scan(&tailLink)
	if tailLink == "0" {
		tailLink = ""
	}
	fmt.Println("File Name :")
	fmt.Scan(&fileName)
	fmt.Println("Sheet Name :")
	fmt.Scan(&sheetName)

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
	/////////
	c2 := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"),
		colly.AllowedDomains("www.shiksha.com"),
		colly.AllowURLRevisit(),
	)
	c2.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting :", r.URL.String())
	})
	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("on error :", e.Error())
	})
	c2.OnResponse(func(r *colly.Response) {
		fmt.Println("Response Code:", r.StatusCode)
	})
	var contact string
	c2.OnHTML(".cntct-gap", func(k *colly.HTMLElement) {
		contact = k.Text
		re := regexp.MustCompile(`\([^)]*\)`)
		contact = re.ReplaceAllString(contact, "")
		atIndex := strings.Index(contact, "@")
		if atIndex != -1 {
			contact = contact[atIndex+1:]
		}
	})
	c.OnHTML("._8165 ", func(h *colly.HTMLElement) {

		count++
		clgName := h.ChildText("h3[title]")
		if initial == clgName {
			fmt.Println("college repeated :", initial, clgName)
			// insert and save
			createAndSaveFile(sheetName, data)
			os.Exit(0)
		}
		place := h.ChildText("span._5588")
		types := h.ChildText("span")
		link := h.ChildAttr("a.ripple", "href")
		if strings.Contains(types, "Govt") {
			types = "Govt"
		} else if strings.Contains(types, "Pvt") {
			types = "Pvt"
		} else {
			types = ""
		}

		if count == 1 {
			initial = clgName
		}

		if link != "" {
			err := c2.Visit("https://www.shiksha.com" + link)
			if err != nil {
				fmt.Println("error in visiting subpage :", err.Error())
				// os.Exit(0)
				link = ""
			}
		}
		details := item{
			Name:    clgName,
			Place:   place,
			Type:    types,
			Contact: contact,
		}
		data = append(data, details)
		// make it empty
		contact = ""
	})

	for i := 0; i < totalPages; i++ {
		if i != 0 {
			params = strconv.Itoa(i + 1)
			params = "-" + params + tailLink
			fmt.Println("\t<><><><><><><><><><><><><><><><><><><><><><><><>\n\n\tVisiting Page :", i+1)
		} else {
			params = params + tailLink
		}
		err := c.Visit(visitLink + params)
		if err != nil {
			fmt.Println("error in visit :", err.Error())
			// os.Exit(0)
		}
	}
}

// func writeJSON(data []item) {
// 	file, err := json.MarshalIndent(data, "", " ")
// 	if err != nil {
// 		log.Println("Unable to create json file")
// 		return
// 	}

// 	os.WriteFile(fileName+".json", file, 0644)
// }

func createAndSaveFile(sheetname string, datas []item) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		panic(err)
	}

	row := sheet.AddRow()
	row.AddCell().SetValue("Name")
	row.AddCell().SetValue("District")
	row.AddCell().SetValue("Type")
	row.AddCell().SetValue("Contact")
	for _, i := range datas {
		row := sheet.AddRow()
		row.AddCell().SetValue(i.Name)
		row.AddCell().SetValue(i.Place)
		row.AddCell().SetValue(i.Type)
		row.AddCell().SetValue(i.Contact)
	}
	err = file.Save(fileName + ".xlsx")
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

// func writeTEXT(data []item) {
// 	file, err := os.Create(fileName + ".txt")
// 	if err != nil {
// 		fmt.Println("Error creating file:", err)
// 		return
// 	}
// 	defer file.Close()
// 	writer := bufio.NewWriter(file)

// 	for _, r := range data {
// 		writer.WriteString("Name :" + r.Name + "\n")
// 		writer.WriteString("District :" + r.Place + "\n")
// 		writer.WriteString("Type :" + r.Type + "\n \n")
// 	}

// 	err = writer.Flush()
// 	if err != nil {
// 		fmt.Println("Error flushing writer:", err)
// 		return
// 	}
// /}
