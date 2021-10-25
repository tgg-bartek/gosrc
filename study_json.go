package main

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/tidwall/gjson"
	//"gosrc/common"
	"os"
	"time"
)

//var cache = common.DiskCache{
//	Dir: "F:/godata",
//	Expires: -1,
//}
//
//var gameUrl = "https://statsapi.web.nhl.com/api/v1/game/{gameId}/feed/live"


func main() {
	//--- JSON Input
	//url := common.FormatString(gameUrl, common.Template{"gameId": "2012020660"})
	//stream := common.FetchUrl(url, cache)
	//s := common.ReaderToString(stream)
	//data := []byte(s)

	bytes, _ := os.ReadFile("/data/nhl-small.json")
	//data := []byte(`{"link": "/api/v1/game/2012020660/feed/live"}`)
	var x gameFile
	if err := json.Unmarshal(bytes, &x); err != nil {
		panic(err)
	}
	spew.Dump(x.GameData.Players)
	//for k := range x.GameData.Players {
	//	fmt.Println(k)
	//}
}


type gameFile struct {
	Copyright string `json:"copyright"`
	GamePk int `json:"gamePk"`
	Link string `json:"link"`
	MetaData metaData `json:"metaData"`
	GameData gameData `json:"gameData"`
}

type metaData struct {
	Wait int 	`json:"wait"`
	Timestamp string `json:"timeStamp"`
}

type gameData struct {
	Players map[string]player	`json:"players"`

}

type player struct {
	Id int `json:"id"`
	FullName string `json:"fullName"`
	Link string `json:"link"`
	Weight int `json:"weight"`
	Rookie bool `json:"rookie"`
}



///////////////////////////////////////////////////////////////////////////////


func gjsonExample(s string) {
	value := gjson.Get(s, "gameData.game")
	fmt.Println(value)
}


func chanExample() {
	c := make(chan string)
	go getTextPol(c)
	go getTextEng(c)

	x, y := <-c, <-c
	fmt.Println(x, y)
}

func getTextPol(c chan string) {
	text, _ := os.ReadFile("data/polish.txt.txt")
	c <- string(text)
}


func getTextEng(c chan string) {
	text, _ := os.ReadFile("data/english.txt")
	time.Sleep(time.Millisecond * 100)
	c <- string(text)
}




