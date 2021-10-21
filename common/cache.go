// Package common: disk cache implementation for the internet
package common

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)


type DiskCache struct {
	Dir string
	Expires time.Duration  // set to -1 for expire never
}


func (cache DiskCache) urlToPath(_url string) string {
	components, _ := url.Parse(_url)
	fileName := components.Host + components.Path + components.RawQuery + components.Fragment
	fileName = strings.ReplaceAll(fileName, "//", "/")

	if strings.HasSuffix(fileName, "/") {
		fileName += "index.html"
	}

	// replace invalid characters
	r, _ := regexp.Compile("[^-/0-9a-zA-Z.,;_ ]")
	fileName = r.ReplaceAllString(fileName, "_")
	fileNameFormatted := ""
	for _, part := range strings.Split(fileName, "/") {
		if len(part) > 255 {
			fileNameFormatted += part[:255] + "/"
		} else {
			fileNameFormatted += part + "/"
		}
	}
	return strings.ReplaceAll(filepath.Join(cache.Dir, fileNameFormatted), "\\", "/")

}


func (cache DiskCache) hasExpired(timestamp time.Time) bool {
	// return whether this timestamp has expires.
	if cache.Expires == -1 {
		return false
	}
	return time.Now().After(timestamp.Add(cache.Expires))

}

func (cache DiskCache) SetCache(_url string, x string) {
	filePath := cache.urlToPath(_url)
	dir, _ := path.Split(filePath)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	data := map[string]string{
		"timestamp": time.Now().Format(time.RFC3339),
		"content": x,
	}
	jdata, _:= json.Marshal(data)
	ioutil.WriteFile(filePath, jdata, 0644)

}


func (cache DiskCache) GetCache(_url string) map[string]string {
	filePath := cache.urlToPath(_url)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return	nil	// not found
	}
	jdata, _ := ioutil.ReadFile(filePath)
	var dat map[string]string
	json.Unmarshal(jdata, &dat)
	ts, _ := time.Parse(time.RFC3339, dat["timestamp"])
	//if hasExpired(ts, cache) {
	if cache.hasExpired(ts) {
		return nil
	}
	return dat
}
