package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/tealeg/xlsx"
)

var filename, sheetname string

type SchoolDetails struct {
	Code     string
	Name     string
	Place    string
	District string
	Contact  string
	Type     string
	Grade    string
}

var Data []SchoolDetails

func main() {

	fmt.Println("Enter the File Name :")
	fmt.Scan(&filename)
	fmt.Println("Enter the Sheet Name :")
	fmt.Scan(&sheetname)

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"),
		colly.AllowedDomains("sametham.kite.kerala.gov.in"),
		colly.AllowURLRevisit(),
		colly.MaxDepth(1),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting :", r.URL.String())

	})
	c.OnError(func(r *colly.Response, e error) {
		log.Println("on error :", e.Error())
	})
	c.OnResponse(func(r *colly.Response) {
		log.Println("Response Code :", r.StatusCode)
	})
	var count int
	c.OnHTML("tr", func(h *colly.HTMLElement) {
		count++
		if count > 2 {
			log.Println("Entered Row :", count-2)
			var code, name, place, district, ftype, grade, phone string
			h.ForEach("td", func(i int, h1 *colly.HTMLElement) {
				switch i {
				case 1:
					code = h1.Text
					if !strings.Contains(code, "HS") {
						code = ""
					}
				case 2:
					name = h1.Text
				case 3:
					place = h1.Text
				case 4:
					district = h1.Text
					re := regexp.MustCompile(`\([^)]*\)`)
					district = re.ReplaceAllString(district, "")
				case 7:
					ftype = h1.Text
				case 8:
					grade = h1.Text
				case 9:
					phone = h1.Text
				default:
					fmt.Print("")
				}
			})
			if code != "" {
				data := SchoolDetails{
					Code:     code,
					Name:     name,
					Place:    place,
					District: district,
					Type:     ftype,
					Grade:    grade,
					Contact:  phone,
				}
				Data = append(Data, data)
			}
		}
	})
	for param := 1; param < 15; param++ {
		err := c.Visit("https://sametham.kite.kerala.gov.in/search/districtWiseSchools/" + strconv.Itoa(param))
		if err != nil {
			log.Println("error in visiting:", err)
		}
	}
	fmt.Println("reached")
	writeXLSX(Data)

}

func writeXLSX(data []SchoolDetails) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet(sheetname)
	if err != nil {
		panic(err)
	}

	row := sheet.AddRow()
	row.AddCell().SetValue("Code")
	row.AddCell().SetValue("Name")
	row.AddCell().SetValue("Place")
	row.AddCell().SetValue("District")
	row.AddCell().SetValue("Financial Type")
	row.AddCell().SetValue("Contact")
	row.AddCell().SetValue("Grade")

	for _, r := range data {
		row := sheet.AddRow()
		row.AddCell().SetValue(r.Code)
		row.AddCell().SetValue(r.Name)
		row.AddCell().SetValue(r.Place)
		row.AddCell().SetValue(r.District)
		row.AddCell().SetValue(r.Type)
		row.AddCell().SetValue(r.Contact)
		row.AddCell().SetValue(r.Grade)
	}

	err = file.Save(filename + ".xlsx")
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
