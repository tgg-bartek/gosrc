// Data Scraper at https://www.sportsoddshistory.com/
// created 13-Jan-22

package main


import (
	"regexp"
	"strings"
	//"fmt"
	"gosrc/common"
	"github.com/PuerkitoBio/goquery"
	"encoding/csv"
	"os"
)

const (
	homepage = "https://www.sportsoddshistory.com/"
	yearUrl = "https://www.sportsoddshistory.com/nfl-game-season/?y={year}"
)

var header = []string{"Year", "Week", "Round", "Day", "Date", "Time (ET)", "", "Favorite", "Score", "Spread", "", "Underdog", "Over/Under", "Notes"}


func main() {

	cache := common.DiskCache{Dir: "F:/_data/godata", Expires: -1}

	url := common.FormatString(yearUrl, common.Template{"year": 1981})

	doc := getDoc(url, cache)
	data := parseYearTables(doc)
	//fmt.Print(len(data))

	toCsv(data, "data/1981.csv")

}


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


func getDoc(url string, cache common.DiskCache) *goquery.Document {
	reader := common.FetchUrl(url, cache)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil{
		panic(err)
	}
	return doc

}

func parseYearTables(doc *goquery.Document) [][][]string  {		// map[int][][]string
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
				rowData = append(rowData, "")	// placeholder for non playoffs weeks
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
