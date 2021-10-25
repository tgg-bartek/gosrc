// Data scraper at footywire.com
// single year supported only
// saves data to CSV. supports command line arguments

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gosrc/common"
	"os"
	P "path"
	"strings"
	"time"
)


const homepage = "https://www.footywire.com/afl/footy/"
const cacheDir = "F:/_data/godata"

// Header slice to be populated during the program call
// scoreboard header
var Header []string


func check(e error) {
	if e != nil {
		panic(e)
	}
}


// parses year schedule from doc. limit output with date range (required)
func parseSchedule(doc *goquery.Document, start time.Time, end time.Time) []string {
	var urls []string
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.Contains(href, "match_statistics") {
			urls = append(urls, href)
		}
	})
	var dates []string
	doc.Find("td[height=\"24\"]").Each(func(_ int, s *goquery.Selection) {
		text := strings.Trim(s.Text(), " ")[2:]
		dates = append(dates, text)
	})
	layout := "Mon 2 Jan 3:04pm 2006"	// Thu 18 Mar 7:25pm
	var datesParsed []time.Time
	for _, d := range dates[:len(urls)] {
		dateTime, _ := time.Parse(layout, d + " 2021")
		if dateTime.After(start) && dateTime.Before(end) {
			datesParsed = append(datesParsed, dateTime)
		}
	}
	return urls[:len(datesParsed)]

}

// parses match scoreboard (player stats) both teams combined
func parseScoringTables(url string, c common.DiskCache) map[int][][]string {
	//r, _ := regexp.Compile("[0-9]+")
	//matchId := r.FindString(url)
	match := common.FetchUrl(homepage+url, c)
	matchDoc, _ := goquery.NewDocumentFromReader(match)

	if len(Header) == 0 {
		Header = getHeader(matchDoc)
	}

	var tables = make(map[int][][]string)
	matchDoc.Find("table[width=\"823\"]").Each(func(nTable int, tbl *goquery.Selection) {
		var tableData [][]string

		tbl.Find("td[height=\"18\"]").Each(func(nRow int, td *goquery.Selection) {
			var rowData []string
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


// gets scoreboard header
func getHeader(doc *goquery.Document) []string {
	// input is match doc
	tbl := doc.Find("table[width=\"823\"]").First()
	var header []string
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


// input date time in UTC (use date time same as footywire.com but a different format)
func parseTime(ts string) time.Time {
	layout := "2006-01-02 15:04"
	d, _ := time.Parse(layout, ts)
	return d
}


// create directory if not exists
func createDir(p string) {

	dir, _ := P.Split(p)
	if dir != "" {
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return
			}
		}
	}
}

// footywire.com data scraper
// Download "player data" to CSV. Supports disk cache. Uses datetime range to limit output. Example:
// -s "2021-04-01 18:00" -e "2021-04-02 20:00" -p "data/footywire/week1.csv"
func main() {

	// parse command line args
	var start string
	var end string
	var path string

	flag.StringVar(&start, "s", "start", "older end of time range YYY-mm-dd H:M (use quotes)")
	flag.StringVar(&end, "e", "end", "newer end of time range YYY-mm-dd H:M (use quotes)")
	flag.StringVar(&path, "p", "path", "file path to save data")
	flag.Parse()

	// scrape data
	url := "https://www.footywire.com/afl/footy/ft_match_list"
	c := common.DiskCache{Dir: cacheDir, Expires: -1}
	// schedule cache need to update while season is in progress
	cSchedule := common.DiskCache{Dir: cacheDir, Expires: time.Hour * 24}

	data := common.FetchUrl(url, cSchedule)
	doc, err := goquery.NewDocumentFromReader(data)
	check(err)

	urls := parseSchedule(doc, parseTime(start), parseTime(end))

	if path == "path" {
		// flag was not provided, use default path
		path = "data/.footywire/data.csv"
	}
	createDir(path)
	f, err := os.Create(path)
	defer f.Close()
	check(err)
	writer := csv.NewWriter(f)

	for nUrl, url := range urls {
		fmt.Println("-->", url)
		tables := parseScoringTables(url, c)
		if nUrl == 0 {
			err := writer.Write(Header)
			if err != nil {
				return
			}
		}
		for _, tbl := range tables {
			_ = writer.WriteAll(tbl)}
	}
}
