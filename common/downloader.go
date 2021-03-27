package common

import (
	"io"
	"io/ioutil"
	"net/http"
	"bytes"
	"fmt"
)


func FetchUrl(url string, c DiskCache) io.Reader {
	// check if common exists
	res := Get(url, c)
	if len(res) == 0 {
		fmt.Println("downloading...", url)
		resp, _ := http.Get(url)
		b, _ := ioutil.ReadAll(resp.Body)
		s := string(b)
		Set(url, s, c)
		defer resp.Body.Close()
		return bytes.NewReader(b)

	} else {
		fmt.Println("cached", url)
		return bytes.NewReader([]byte(res["content"]))
	}

}


//class Throttle(object):
//    """ Throttle downloading by sleeping between requests to same domain """
//
//    def __init__(self, delay):
//        # amount of delay between downloads for each domain
//        self.delay = delay
//        # timestamp of when a domain was last accessed
//        self.domains = {}
//
//    def wait(self, url):
//        """ Delay if have accessed this domain recently """
//        domain = urlsplit(url).netloc
//        last_accessed = self.domains.get(domain)
//        if self.delay > 0 and last_accessed is not None:
//            sleep_secs = self.delay - (datetime.utcnow() - last_accessed).seconds
//            if sleep_secs > 0:
//                time.sleep(sleep_secs)
//        self.domains[domain] = datetime.utcnow()
