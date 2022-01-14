// Data scraper at footywire.com (Player Data)
// single year supported only. supports disk cache.
// saves data to CSV. supports command line arguments.
// Example (command line):
// go run footywire.go -s "2021-04-01 18:00" -e "2021-04-02 20:00" -p "data/footywire/week1.csv"
// For datetime args use UTC tz in format Y-m-d H:M

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gosrc/common"
	"os"
	"regexp"
	"strings"
	"time"
)


// Disk Cache constants
const homepage = "https://www.footywire.com/afl/footy/"
const cacheDir = "F:/_data/godata"

// Header player stats header to be populated during the program call
var Header []string


func check(e error) {
	if e != nil {
		panic(e)
	}
}


func main() {


	// Parse command line args
	var (
		start string
		end string
		path string
	)

	flag.StringVar(&start, "s", "start", "older end of time range YYY-mm-dd H:M (use quotes)")
	flag.StringVar(&end, "e", "end", "newer end of time range YYY-mm-dd H:M (use quotes)")
	flag.StringVar(&path, "p", "path", "file path to save data")
	flag.Parse()

	// Scrape data
	url := "https://www.footywire.com/afl/footy/ft_match_list"
	c := common.DiskCache{Dir: cacheDir, Expires: -1}
	// schedule cache need to update while season has not ended
	cSchedule := common.DiskCache{Dir: cacheDir, Expires: time.Hour * 24}

	data := common.FetchUrl(url, cSchedule)
	doc, err := goquery.NewDocumentFromReader(data)
	check(err)

	urls := parseSchedule(doc, parseTime(start), parseTime(end))
	if path == "path" {
		// flag was not provided, use default path
		path = "data/.footywire/data.csv"
	}

	// Write to CSV
	common.CreateDir(path)
	f, err := os.Create(path)
	defer f.Close()
	check(err)
	writer := csv.NewWriter(f)

	for nUrl, url := range urls {
		fmt.Println("-->", url)
		tables := parseScoringTables(url, c)
		if nUrl == 0 {
			err := writer.Write(Header)
			check(err)
		}
		for _, tbl := range tables {
			err = writer.WriteAll(tbl)
			check(err)
		}
	}
}


func getYear(doc *goquery.Document) string {
	yearTag := doc.Find("span[class=\"hltitle\"]")
	r, _ := regexp.Compile("[0-9]+$")
	val := r.FindString(yearTag.Text())
	return val
}


// parses year schedule from doc. limit output with date range (required)
func parseSchedule(doc *goquery.Document, start time.Time, end time.Time) []string {
	y := getYear(doc)

	layout := "Mon 2 Jan 3:04pm 2006"	// Thu 18 Mar 7:25pm
	var urls []string
	doc.Find("td[id=\"contentpagecell\"]").Find("a[href*=\"match_statistics\"]").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		tr := s.Closest("tr")		// parent tr
		timeString := tr.Find("td[height=\"24\"]").Text()
		timeString = strings.Trim(timeString, "\xc2\xa0")
		dateTime, err := time.Parse(layout, timeString + " " + y)
		check(err)
		// if url in date range, take
		if dateTime.After(start) && dateTime.Before(end) {
			urls = append(urls, href)
		}
	})
	return urls

}

// parses match scoreboard (player stats) both teams combined
func parseScoringTables(url string, c common.DiskCache) map[int][][]string {
	//r, _ := regexp.Compile("[0-9]+")
	//matchId := r.FindString(url)
	match := common.FetchUrl(homepage+url, c)
	matchDoc, _ := goquery.NewDocumentFromReader(match)

	// Populate Header for CSV output
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


// parses command line arg date time
func parseTime(ts string) time.Time {
	layout := "2006-01-02 15:04"
	d, _ := time.Parse(layout, ts)
	return d
}
