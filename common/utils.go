package common

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	P "path"
	"strings"
)

type Template map[string]interface{}

// FormatString Python like string formatting with curly brackets {}
// src: https://stackoverflow.com/a/40811635
// Example
// var gameUrl = "https://statsapi.web.nhl.com/api/v1/game/{gameId}/feed/live"
// url := FormatString(gameUrl, Template{"gameId": "2012020660"})
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
	b, err := ioutil.ReadAll(stream)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// FileExists: checks if file exists at path exists

func FileExists(fp string) bool {
	_, err := os.Stat(fp)
	if os.IsNotExist(err) {
		return false // not found
	}
	return true // file found
}

// CreateDir creates directory (and sub dirs) if path does not exists
func CreateDir(p string) {

	dir, _ := P.Split(p)
	if dir != "" {
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
	}
}

// Inserting value to slice at index
// 0 <= index <= len(a)
// https://stackoverflow.com/a/61822301
func Insert(a []string, index int, value string) []string {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

// contains checks if a string is present in a slice
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
