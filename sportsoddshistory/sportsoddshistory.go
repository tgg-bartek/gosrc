// Data Scraper at https://www.sportsoddshistory.com/
// created 13-Jan-22

package main


import (
	"fmt"
	"gosrc/common"
	"github.com/PuerkitoBio/goquery"
)

const homepage = "https://www.sportsoddshistory.com/"


func main() {

	var url= "https://www.sportsoddshistory.com/nfl-game-season/?y=1980"

	cache := common.DiskCache{Dir: "F:/_data/godata", Expires: -1}
	reader := common.FetchUrl(url, cache)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		panic(err)
	}

	var data = make(map[int][][]string)
	doc.Find("table[class=\"soh1\"]").Each(func(n int, tbl *goquery.Selection) {
	//doc.Find("h3:contains(\"Regular Season - Week\"), h3:contains(\"Playoffs\")").Each(func(n int, h3 *goquery.Selection) {

		if n > 1 {
			table := parseTable(tbl)
			data[n] = table
		}

	})
	fmt.Println(len(data))

}


func parseTable(tbl *goquery.Selection) [][]string {
	var table [][]string
	tbl.Find("tr").Each(func(rNum int, tr *goquery.Selection) {
		if rNum == 0 {
			// header
		} else {
			var rowData []string
			tr.Find("td").Each(func(_ int, td *goquery.Selection) {
				value := td.Text()
				rowData = append(rowData, value)
			})
			table = append(table, rowData)
		}

	})
	return table
}