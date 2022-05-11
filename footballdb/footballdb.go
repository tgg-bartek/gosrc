// Data Scraper at footballdb.com
// created 3-May-22

package main

import (
	"fmt"
	"gosrc/common"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

const (
	yearUrl  = "https://www.footballdb.com/games/index.html?lg=NFL&yr={year}" // schedule
	cacheDir = "F:/_data/godata/footballdb"
)

type fixture struct {
	Week      string `json:"week"`
	Date      string `json:"date"`
	Home      string `json:"home"`
	Away      string `json:"away"`
	HomeScore string `json:"homeScore"`
	AwayScore string `json:"awayScore"`
	Info      string `json:"info"`
	Url       string `json:"url"`
	// TableNum  int
}

// func getWeeks(c *colly.Collector) *[]string {
// 	wks := &[]string{}
// 	c.OnHTML("div.ltbluediv", func(div *colly.HTMLElement) {
// 		weekName := div.DOM.Find("span").Text()
// 		*wks = append(*wks, weekName)
// 	})
// 	return wks
// }

func getWeeks(c *colly.Collector, wks *[]string) {
	c.OnHTML("div.ltbluediv", func(div *colly.HTMLElement) {
		weekName := div.DOM.Find("span").Text()
		*wks = append(*wks, weekName)
	})
}

func parseSchedRow(tr *colly.HTMLElement) fixture {
	f := fixture{}
	tr.ForEach("td", func(n int, td *colly.HTMLElement) {
		if n == 0 {
			f.Date = td.DOM.Find("span.hidden-xs").Text()
		}
		if n == 1 {
			f.Away = td.DOM.Find("span.hidden-xs").Text()
		}
		if n == 2 {
			f.AwayScore = td.Text
		}
		if n == 3 {
			f.Home = td.DOM.Find("span.hidden-xs").Text()
		}
		if n == 4 {
			f.HomeScore = td.Text
		}
		if n == 5 {
			f.Info = td.Text
		}
		if n == 6 {
			f.Url = td.ChildAttr("a", "href")
		}
	})
	return f
}

func parseScheduleYear(c *colly.Collector, yearTable *map[int][]fixture) {
	nTable := 0
	c.OnHTML("table tbody", func(table *colly.HTMLElement) {
		var weekTable []fixture
		table.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
			f := parseSchedRow(tr)
			// f.TableNum = nTable
			// *yearTable = append(*yearTable, f)
			weekTable = append(weekTable, f)
		})
		// fmt.Println(weekTable[0], "***")
		(*yearTable)[nTable] = weekTable // note the () to dereference https://stackoverflow.com/a/36463641/4279824
		nTable++
	})
}

func getScheduleYear(c *colly.Collector, year int) []fixture {
	var wks []string
	getWeeks(c, &wks)
	yearTable := make(map[int][]fixture)
	parseScheduleYear(c, &yearTable)

	url := common.FormatString(yearUrl, common.Template{"year": year})
	c.Visit(url)

	// Add week Name to every fixture
	var yearData []fixture
	for weekNum, weekData := range yearTable {
		weekName := wks[weekNum]
		for _, fix := range weekData {
			fix.Week = weekName
			yearData = append(yearData, fix)
		}
	}
	return yearData
}

func parseTeamStats(c *colly.Collector, gameUrl string) [][]string {
	table := [][]string{}
	added := false
	c.OnHTML("div#divBox_team table", func(t *colly.HTMLElement) {
		header := []string{}
		t.ForEach("thead td span.hidden-xs", func(i int, td *colly.HTMLElement) {
			if 1 > 0 {
				if td.Text != "" {
					header = append(header, td.Text)
				}
			}
		})
		if !added {
			table = append(table, header)
			added = true
		}

		t.ForEach("tbody tr", func(itr int, tr *colly.HTMLElement) {
			var field []string
			tr.ForEach("td", func(itd int, td *colly.HTMLElement) {
				field = append(field, td.Text)
			})
			table = append(table, field)
		})

	})
	c.Visit(gameUrl)
	return table
}

func parsePlayerStats(c *colly.Collector, gameUrl string) (map[string][]string, map[string][][]string) {
	var headers = make(map[string][]string)
	var data = make(map[string][][]string)

	c.OnHTML("div#divBox_stats", func(div *colly.HTMLElement) {

		// seen := []string{}

		boxes := div.DOM.Find("div.divider")
		boxes.Each(func(i int, s *goquery.Selection) {
			stat_name := s.Find("h2").Text()

			sec := s.NextUntil("div.divider")
			sec.Find("table").Each(func(tNum int, tbl *goquery.Selection) {
				// Header of each data type
				head := tbl.Find("thead th")
				header := []string{}
				head.Each(func(iTh int, th *goquery.Selection) {
					if iTh == 0 {
						value := th.Find("span.hidden-xs").Text()
						header = append(header, value)
					} else {
						value := th.Text()
						header = append(header, value)
					}
				})
				teamName := header[0]
				header[0] = "Team"
				common.Insert(header, 1, "Player")
				common.Insert(header, 2, "PlayerLink")

				_, ok := headers[stat_name]
				if !ok {
					headers[stat_name] = header
				}

				// Header not added to data so add and flag the adding
				// if !common.Contains(seen, stat_name) {
				// 	data[stat_name] = header
				// 	seen = append(seen, stat_name)
				// }

				// Get table data
				body := tbl.Find("tbody tr")
				tableData := [][]string{}
				body.Each(func(iTr int, tr *goquery.Selection) {
					row := []string{}
					row = append(row, teamName) // append team name for each player row
					tr.Find("td").Each(func(iTd int, td *goquery.Selection) {
						if iTd == 0 {
							span := td.Find("span.hidden-xs")
							player_name := span.Text()
							player_link, _ := span.Find("a").Attr("href")
							row = append(row, player_name, player_link)
						} else {
							value := td.Text()
							row = append(row, value)

						}
					})
					tableData = append(tableData, row)
				})
				tables, ok := data[stat_name]
				if ok {
					tables = append(tables, tableData...)
					data[stat_name] = tables
				} else {
					data[stat_name] = tableData
				}

			})
		})

	})
	c.Visit(gameUrl)
	return headers, data
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.CacheDir(cacheDir),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)")
		r.Headers.Set("referer", "https://www.footballdb.com/games/")
		r.Headers.Set("authority", "www.footballdb.com")
		r.Headers.Set("scheme", "https")
	})

	// // Year Schedule
	// sched := getScheduleYear(c, 2009)
	// fmt.Println(len(sched))
	// for _, f := range sched {
	// 	fmt.Println(f)
	// 	break
	// }

	gameUrl := "https://www.footballdb.com/games/boxscore/tennessee-titans-vs-pittsburgh-steelers-2009091001"
	// Team Stats
	// table := parseTeamStats(c, gameUrl)
	// fmt.Println(table)

	// Player Stats
	headers, data := parsePlayerStats(c, gameUrl)
	fmt.Println(headers["Passing"])
	for _, row := range data["Passing"] {
		fmt.Println(row)
	}

}
