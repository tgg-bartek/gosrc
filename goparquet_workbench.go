package main

// https://github.com/parquet-go/parquet-go/blob/main/example_test.go

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/parquet-go/parquet-go"
)

type Record struct {
	GameId                    int32   `parquet:"gid"`
	Season                    int32   `parquet:"season"`
	LocalDate                 string  `parquet:"local_date"`
	LocalTime                 string  `parquet:"local_time"`
	AmPm                      string  `parquet:"ampm"`
	IsDoubleheader            string  `parquet:"is_doubleheader"`
	GameNumber                int32   `parquet:"game_number"`
	Home                      string  `parquet:"home"`
	Away                      string  `parquet:"away"`
	Team                      string  `parquet:"team"`
	ResultType                string  `parquet:"result.type"`
	ResultEvent               string  `parquet:"result.event"`
	ResultEventType           string  `parquet:"result.eventType"`
	ResultDesc                string  `parquet:"result.description"`
	AtBatIndex                int     `parquet:"about.atBatIndex"`
	HalfInning                string  `parquet:"abount.halfInning"`
	Inning                    int     `parquet:"about.inning"`
	RunsCount                 float32 `parquet:"runs_count"`
	ErrorsCount               float32 `parquet:"errors_count"`
	IsHit                     float32 `parquet:"isHit"`
	HomeRuns                  int64   `parquet:"home.runs"`
	AwayRuns                  int64   `parquet:"away.runs"`
	HomeErrors                int64   `parquet:"home.errors"`
	AwayErrors                int64   `parquet:"away.errors"`
	HomeHits                  int64   `parquet:"home.hits"`
	AwayHits                  int64   `parquet:"away.hits"`
	MatchupBatterId           float64 `parquet:"matchup.batter.id"`
	MatchupBatterFullname     string  `parquet:"matchup.batter.fullName"`
	MatchupBatsideCode        string  `parquet:"matchup.batSide.code"`
	MatchupPitcherId          float64 `parquet:"matchup.pitcher.id"`
	MatchupPitcherFullname    string  `parquet:"matchup.pitcher.fullName"`
	MatchupPitchhandCode      string  `parquet:"matchup.pitchHand.code"`
	DetailsDescription        string  `parquet:"details.description"`
	DetailsIsscoringplay      float64 `parquet:"details.isScoringPlay"`
	DetailsIsout              bool    `parquet:"details.isOut"`
	CountBalls                int64   `parquet:"count.balls"`
	CountStrikes              int64   `parquet:"count.strikes"`
	CountOuts                 int64   `parquet:"count.outs"`
	Index                     int64   `parquet:"index"`
	Ispitch                   bool    `parquet:"isPitch"`
	Type                      string  `parquet:"type"`
	DetailsIsstrike           float64 `parquet:"details.isStrike"`
	DetailsIsball             float64 `parquet:"details.isBall"`
	DetailsTypeDescription    string  `parquet:"details.type.description"`
	PitchdataStartspeed       float64 `parquet:"pitchData.startSpeed"`
	PitchdataStrikezonetop    float64 `parquet:"pitchData.strikeZoneTop"`
	PitchdataStrikezonebottom float64 `parquet:"pitchData.strikeZoneBottom"`
	PitchdataCoordinatesX     float64 `parquet:"pitchData.coordinates.x"`
	PitchdataCoordinatesY     float64 `parquet:"pitchData.coordinates.y"`
	PitchdataZone             float64 `parquet:"pitchData.zone"`
	Pitchnumber               float64 `parquet:"pitchNumber"`
}

func main() {

	f, _ := os.Open("C:\\Users\\bartek\\go\\src\\gosrc\\data\\pbp-717432-c.parquet")
	// Now, we can read from the file.
	pf := parquet.NewReader(f)
	records := make([]Record, 0)
	for {
		var rec Record
		err := pf.Read(&rec)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		records = append(records, rec)
	}
	fmt.Println(records[0])

}
