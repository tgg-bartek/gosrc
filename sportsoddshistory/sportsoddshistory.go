// Data Scraper at https://www.sportsoddshistory.com/
// created 13-Jan-22

package main

import (
	"encoding/csv"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	//"fmt"
	"gosrc/common"
	"os"
)

const (
	homepage = "https://www.sportsoddshistory.com/"
	yearUrl  = "https://www.sportsoddshistory.com/nfl-game-season/?y={year}"
)

var header = []string{"Year", "Week", "Round", "Day", "Date", "Time (ET)", "", "Favorite", "Score",
	"Spread", "", "Underdog", "Over/Under", "Notes"}
var cache = common.DiskCache{Dir: "F:/_data/godata", Expires: -1}

func main() {

	// for i := 1971; i < 2010; i++ {
	// 	downloadSaveYear(i) // save year sheet at data/
	// 	time.Sleep(time.Minute * 3)
	// }
	// downloadSaveYear(1978)

}

// Download odds for single season (and save to a file)
func downloadSaveYear(year int) {
	url := common.FormatString(yearUrl, common.Template{"year": year})
	doc := getDoc(url)
	data := parseYearTables(doc)
	toCsv(data, common.FormatString("data/sportsoddshistory/{year}.csv", common.Template{"year": year}))
}

// Write data to CSV
func toCsv(data [][][]string, path string) {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(f)

	_ = writer.Write(header)
	for _, tbl := range data {
		_ = writer.WriteAll(tbl)
	}
}

// Get Document from URL
func getDoc(url string) *goquery.Document {
	reader := common.FetchUrl(url, cache)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		panic(err)
	}
	return doc

}

func parseYearTables(doc *goquery.Document) [][][]string { // map[int][][]string
	//var data = make(map[int][][]string)
	var data [][][]string
	doc.Find("table[class=\"soh1\"]").Each(func(n int, tbl *goquery.Selection) {
		//doc.Find("h3:contains(\"Regular Season - Week\"), h3:contains(\"Playoffs\")").Each(func(n int, h3 *goquery.Selection) {

		// skip first two (summary tables)
		if n > 1 {
			table := parseTable(tbl)
			//data[n] = table
			data = append(data, table)
		}
	})
	return data
}

func parseTable(tbl *goquery.Selection) [][]string {
	var table [][]string

	// Get Week name
	sib := tbl.Prev()
	if strings.Contains(sib.Text(), "BOLD") {
		sib = sib.Prev()
	}
	week := sib.Text()
	// Get Year
	r, _ := regexp.Compile("\\d{4}")
	year := r.FindString(week)

	tbl.Find("tr").Each(func(rNum int, tr *goquery.Selection) {
		if rNum == 0 {
			// header
		} else {
			var rowData []string
			rowData = append(rowData, year, cleanWeek(week))
			if cleanWeek(week) != "Playoffs" {
				rowData = append(rowData, "") // placeholder for non playoffs weeks
			}
			tr.Find("td").Each(func(_ int, td *goquery.Selection) {
				value := td.Text()
				rowData = append(rowData, value)
			})
			table = append(table, rowData)
		}

	})
	return table
}

func cleanWeek(week string) string {
	if strings.Contains(week, "Playoffs") {
		return "Playoffs"
	}
	r, _ := regexp.Compile("Week \\d+")
	match := r.FindString(week)
	return match
}
