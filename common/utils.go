package common

import (
	"fmt"
	"io"
	"strings"
)



type Template map[string]interface{}

// FormatString Python like string formatting with curly brackets {}
// src: https://stackoverflow.com/a/40811635
// Example
// var gameUrl = "https://statsapi.web.nhl.com/api/v1/game/{gameId}/feed/live"
// url := formatString(gameUrl, Template{"gameId": "2012020660"})
func FormatString(s string, t Template) string {
	args, i := make([]string, len(t)*2), 0
	for k, v := range t {
		args[i] = "{" + k + "}"
		args[i+1] = fmt.Sprint(v)
		i += 2
	}
	return strings.NewReplacer(args...).Replace(s)
}


func ReaderToString(stream io.Reader) string {
	b, err := io.ReadAll(stream)
	if err != nil {
		panic(err)
	}
	return string(b)
}
