package main

// Dependency:
// go install github.com/parquet-go/parquet-go
// Example: https://github.com/parquet-go/parquet-go/blob/main/example_test.go

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/parquet-go/parquet-go"
)

type Record struct {
	GidMlb                        int64   `parquet:"gid_mlb"`
	GidRef                        string  `parquet:"gid_ref"`
	Season                        int64   `parquet:"season"`
	LocalDate                     string  `parquet:"local_date"`
	LocalTime                     string  `parquet:"local_time"`
	Ampm                          string  `parquet:"ampm"`
	IsDoubleheader                string  `parquet:"is_doubleheader"`
	GameNumber                    int64   `parquet:"game_number"`
	Home                          string  `parquet:"home"`
	Away                          string  `parquet:"away"`
	Team                          string  `parquet:"team"`
	ResultType                    string  `parquet:"result.type"`
	ResultEvent                   string  `parquet:"result.event"`
	ResultEventtype               string  `parquet:"result.eventType"`
	ResultDescription             string  `parquet:"result.description"`
	AboutAtbatindex               float64 `parquet:"about.atBatIndex"`
	AboutHalfinning               string  `parquet:"about.halfInning"`
	AboutInning                   float64 `parquet:"about.inning"`
	Index                         float64 `parquet:"index"`
	MatchupBatterId               float64 `parquet:"matchup.batter.id"`
	MatchupBatterFullname         string  `parquet:"matchup.batter.fullName"`
	MatchupBatsideCode            string  `parquet:"matchup.batSide.code"`
	MatchupPitcherId              float64 `parquet:"matchup.pitcher.id"`
	MatchupPitcherFullname        string  `parquet:"matchup.pitcher.fullName"`
	MatchupPitchhandCode          string  `parquet:"matchup.pitchHand.code"`
	DetailsDescription            string  `parquet:"details.description"`
	DetailsIsstrike               bool    `parquet:"details.isStrike"`
	DetailsIsball                 bool    `parquet:"details.isBall"`
	DetailsIsscoringplay          bool    `parquet:"details.isScoringPlay"`
	DetailsTypeDescription        string  `parquet:"details.type.description"`
	CountBalls                    float64 `parquet:"count.balls"`
	CountStrikes                  float64 `parquet:"count.strikes"`
	CountOuts                     float64 `parquet:"count.outs"`
	RunsCount                     float64 `parquet:"runs_count"`
	ErrorsCount                   float64 `parquet:"errors_count"`
	Ishit                         float64 `parquet:"isHit"`
	HomeRuns                      float64 `parquet:"home.runs"`
	AwayRuns                      float64 `parquet:"away.runs"`
	HomeErrors                    float64 `parquet:"home.errors"`
	AwayErrors                    float64 `parquet:"away.errors"`
	HomeHits                      float64 `parquet:"home.hits"`
	AwayHits                      float64 `parquet:"away.hits"`
	PitchdataStartspeed           float64 `parquet:"pitchData.startSpeed"`
	PitchdataStrikezonetop        float64 `parquet:"pitchData.strikeZoneTop"`
	PitchdataStrikezonebottom     float64 `parquet:"pitchData.strikeZoneBottom"`
	PitchdataCoordinatesX         float64 `parquet:"pitchData.coordinates.x"`
	PitchdataCoordinatesY         float64 `parquet:"pitchData.coordinates.y"`
	PitchdataZone                 float64 `parquet:"pitchData.zone"`
	Pitchnumber                   float64 `parquet:"pitchNumber"`
	Ispitch                       bool    `parquet:"isPitch"`
	Type                          string  `parquet:"type"`
	HitdataTrajectory             string  `parquet:"hitData.trajectory"`
	HitdataHardness               string  `parquet:"hitData.hardness"`
	HitdataCoordinatesCoordx      float64 `parquet:"hitData.coordinates.coordX"`
	HitdataCoordinatesCoordy      float64 `parquet:"hitData.coordinates.coordY"`
	Runners0MovementStart         string  `parquet:"runners.0.movement.start"`
	Runners0MovementEnd           string  `parquet:"runners.0.movement.end"`
	Runners0MovementOutnumber     float64 `parquet:"runners.0.movement.outNumber"`
	Runners0DetailsEventtype      string  `parquet:"runners.0.details.eventType"`
	Runners0DetailsRunnerId       float64 `parquet:"runners.0.details.runner.id"`
	Runners0DetailsRunnerFullname string  `parquet:"runners.0.details.runner.fullName"`
	Runners0DetailsPlayindex      float64 `parquet:"runners.0.details.playIndex"`
	Runners1MovementStart         string  `parquet:"runners.1.movement.start"`
	Runners1MovementEnd           string  `parquet:"runners.1.movement.end"`
	Runners1MovementOutnumber     float64 `parquet:"runners.1.movement.outNumber"`
	Runners1DetailsEventtype      string  `parquet:"runners.1.details.eventType"`
	Runners1DetailsRunnerId       float64 `parquet:"runners.1.details.runner.id"`
	Runners1DetailsRunnerFullname string  `parquet:"runners.1.details.runner.fullName"`
	Runners1DetailsPlayindex      float64 `parquet:"runners.1.details.playIndex"`
	Runners2MovementStart         string  `parquet:"runners.2.movement.start"`
	Runners2MovementEnd           string  `parquet:"runners.2.movement.end"`
	Runners2MovementOutnumber     float64 `parquet:"runners.2.movement.outNumber"`
	Runners2DetailsEventtype      string  `parquet:"runners.2.details.eventType"`
	Runners2DetailsRunnerId       float64 `parquet:"runners.2.details.runner.id"`
	Runners2DetailsRunnerFullname string  `parquet:"runners.2.details.runner.fullName"`
	Runners2DetailsPlayindex      float64 `parquet:"runners.2.details.playIndex"`
	Runners3MovementStart         string  `parquet:"runners.3.movement.start"`
	Runners3MovementEnd           string  `parquet:"runners.3.movement.end"`
	Runners3MovementOutnumber     float64 `parquet:"runners.3.movement.outNumber"`
	Runners3DetailsEventtype      string  `parquet:"runners.3.details.eventType"`
	Runners3DetailsRunnerId       float64 `parquet:"runners.3.details.runner.id"`
	Runners3DetailsRunnerFullname string  `parquet:"runners.3.details.runner.fullName"`
	Runners3DetailsPlayindex      float64 `parquet:"runners.3.details.playIndex"`
	Runners4MovementStart         string  `parquet:"runners.4.movement.start"`
	Runners4MovementEnd           string  `parquet:"runners.4.movement.end"`
	Runners4MovementOutnumber     float64 `parquet:"runners.4.movement.outNumber"`
	Runners4DetailsEventtype      string  `parquet:"runners.4.details.eventType"`
	Runners4DetailsRunnerId       float64 `parquet:"runners.4.details.runner.id"`
	Runners4DetailsRunnerFullname string  `parquet:"runners.4.details.runner.fullName"`
	Runners4DetailsPlayindex      float64 `parquet:"runners.4.details.playIndex"`
	Runners5MovementStart         string  `parquet:"runners.5.movement.start"`
	Runners5MovementEnd           string  `parquet:"runners.5.movement.end"`
	Runners5MovementOutnumber     float64 `parquet:"runners.5.movement.outNumber"`
	Runners5DetailsEventtype      string  `parquet:"runners.5.details.eventType"`
	Runners5DetailsRunnerId       float64 `parquet:"runners.5.details.runner.id"`
	Runners5DetailsRunnerFullname string  `parquet:"runners.5.details.runner.fullName"`
	Runners5DetailsPlayindex      float64 `parquet:"runners.5.details.playIndex"`
	Runners6MovementStart         string  `parquet:"runners.6.movement.start"`
	Runners6MovementEnd           string  `parquet:"runners.6.movement.end"`
	Runners6MovementOutnumber     float64 `parquet:"runners.6.movement.outNumber"`
	Runners6DetailsEventtype      string  `parquet:"runners.6.details.eventType"`
	Runners6DetailsRunnerId       float64 `parquet:"runners.6.details.runner.id"`
	Runners6DetailsRunnerFullname string  `parquet:"runners.6.details.runner.fullName"`
	Runners6DetailsPlayindex      float64 `parquet:"runners.6.details.playIndex"`
	Runners7MovementStart         string  `parquet:"runners.7.movement.start"`
	Runners7MovementEnd           string  `parquet:"runners.7.movement.end"`
	// Runners7MovementOutnumber     float64 `parquet:"runners.7.movement.outNumber"`
	Runners7DetailsEventtype      string  `parquet:"runners.7.details.eventType"`
	Runners7DetailsRunnerId       float64 `parquet:"runners.7.details.runner.id"`
	Runners7DetailsRunnerFullname string  `parquet:"runners.7.details.runner.fullName"`
	Runners7DetailsPlayindex      float64 `parquet:"runners.7.details.playIndex"`
	// Runners8MovementStart         string  `parquet:"runners.8.movement.start"`
	Runners8MovementEnd string `parquet:"runners.8.movement.end"`
	// Runners8MovementOutnumber     float64 `parquet:"runners.8.movement.outNumber"`
	Runners8DetailsEventtype      string  `parquet:"runners.8.details.eventType"`
	Runners8DetailsRunnerId       float64 `parquet:"runners.8.details.runner.id"`
	Runners8DetailsRunnerFullname string  `parquet:"runners.8.details.runner.fullName"`
	Runners8DetailsPlayindex      float64 `parquet:"runners.8.details.playIndex"`
}

func main() {

	f, _ := os.Open("C:\\Users\\bartek\\go\\src\\gosrc\\data\\pbp-2015.parquet")
	// Now, we can read from the file.
	pf := parquet.NewReader(f)
	// records := make([]Record, 0)
	count := 0
	for {
		var rec Record
		err := pf.Read(&rec)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// records = append(records, rec)

		if count == 3 {
			break
		}
		fmt.Println(rec)
		count += 1

	}

}
