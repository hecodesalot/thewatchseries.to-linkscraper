package webstuff

import (
	"crypto/tls"
	"net/http"

	"github.com/jackdanger/collectlinks"
)

//ErrorURLs collects urls that fail
var ErrorURLs []string

//GetWebpage grabs the source code of a webpage
func GetWebpage(url string) []string {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := http.Client{Transport: transport}
	resp, err := client.Get(url)
	if err != nil {
		ErrorURLs = append(ErrorURLs, url)
		return []string{"error"}
	}

	defer resp.Body.Close()

	links := collectlinks.All(resp.Body)

	return links
}
