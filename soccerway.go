package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gosrc/common"
	"os"
	"strings"
	"time"
)


var Header = []string{}
const homepage = "https://www.footywire.com/afl/footy/"


func check(e error) {
	if e != nil {
		panic(e)
	}
}


func parseSchedule(doc *goquery.Document, start time.Time, end time.Time) []string{
	urls := []string{}
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.Contains(href, "match_statistics") {
			urls = append(urls, href)
		}
	})
	dates := []string{}
	doc.Find("td[height=\"24\"]").Each(func(_ int, s *goquery.Selection) {
		text := strings.Trim(s.Text(), " ")[2:]
		dates = append(dates, text)
	})
	layout := "Mon 2 Jan 3:04pm 2006"	// Thu 18 Mar 7:25pm
	datesParsed := []time.Time{}
	for _, d := range dates[:len(urls)] {
		dateTime, _ := time.Parse(layout, d + " 2021")
		if dateTime.After(start) && dateTime.Before(end) {
			datesParsed = append(datesParsed, dateTime)
		}
		}
	return urls[:len(datesParsed)]

}

func parseScoringTables(url string, c common.DiskCache) map[int][][]string {
	//r, _ := regexp.Compile("[0-9]+")
	//matchId := r.FindString(url)
	//match := fetchUrl(homepage+url, matchId+".html")
	match := common.FetchUrl(homepage+url, c)
	matchDoc, _ := goquery.NewDocumentFromReader(match)

	if len(Header) == 0 {
		Header = getHeader(matchDoc)
	}

	var tables = make(map[int][][]string)
	matchDoc.Find("table[width=\"823\"]").Each(func(nTable int, tbl *goquery.Selection) {
		tableData := [][]string{}
		tbl.Find("td[height=\"18\"]").Each(func(nRow int, td *goquery.Selection) {
			rowData := []string{}
			row := td.Parent()
			row.Find("td").Each(func(nCell int, cell *goquery.Selection) {
			switch nCell {
				case 0:
				href, _ := cell.Find("a").Attr("href")
				value := cell.Text()
				rowData = append(rowData, href, value)
			default:
				value := cell.Text()
				rowData = append(rowData, value)
			}

			})
			tableData = append(tableData, rowData)
		})
		tables[nTable] = tableData
	})
	return tables

}


func getHeader(doc *goquery.Document) []string {
	// input is match doc
	tbl := doc.Find("table[width=\"823\"]").First()
	header := []string{}
	tbl.Find("td[class*=bnorm]").Each(func(i int, selection *goquery.Selection) {
		cell := selection.Text()
		if cell == "Player" {
			header = append(header, "Link", "Player")
		} else {
			header = append(header, cell)
		}
	})
	return header

}


func main() {
	url := "https://www.footywire.com/afl/footy/ft_match_list"
	c := common.DiskCache{Dir: "F:/godata", Expires: time.Hour * 1}

	data := common.FetchUrl(url, c)
	doc, err := goquery.NewDocumentFromReader(data)
	check(err)
	urls := parseSchedule(doc, time.Date(2021, 3, 18, 18, 0, 0, 0, time.UTC), time.Date(2021, 3, 20, 18, 0, 0, 0, time.UTC))

	//write data to csv
	f, err := os.Create("footywire-week1.csv")
	check(err)
	writer := csv.NewWriter(f)

	for nUrl, url := range urls {
		fmt.Println("-->", url)
		tables := parseScoringTables(url, c)
		if nUrl == 0 {
			writer.Write(Header)
		}
		for _, tbl := range tables {

			_ = writer.WriteAll(tbl)
		time.Sleep(5 * time.Second)
		}
		//break
	}
	defer f.Close()
}
