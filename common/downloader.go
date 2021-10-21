package common

import (
	"bytes"
	"fmt"
	"io"

	//"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)


func FetchUrl(_url string, c DiskCache) io.Reader {		// io.Reader or string
	// check if cache exists
	res := c.GetCache(_url)
	if len(res) == 0 {
		Throttle(5 * time.Second, _url)
		fmt.Println("downloading...", _url)
		resp, _ := http.Get(_url)
		b, _ := ioutil.ReadAll(resp.Body)
		s := string(b)
		c.SetCache(_url, s)
		defer resp.Body.Close()
		return bytes.NewReader(b)
		//return s	# string
	} else {
		fmt.Println("cached", _url)
		return bytes.NewReader([]byte(res["content"]))
		//return res["content"]
	}

}


// map of domains to last accessed time stamp for `Throttle`
var domains = make(map[string]time.Time)


// Throttle downloading by sleeping between requests to same domain
// `delay` amount of delay between downloads for each domain
func Throttle(delay time.Duration, _url string) {
	// fixme this need not be a closure but is an example
	func() {
		parts, _ := url.Parse(_url)
		host := parts.Host
		lastAccessed := domains[host]
		if delay > 0 && !lastAccessed.IsZero() {
			// delay if have accessed this domain recently
			sleepSecs := int64(delay) - (time.Now().Unix() - lastAccessed.Unix())
			if sleepSecs > 0 {
				time.Sleep(delay)
			}
		}
		domains[host] = time.Now()
	}()
}
