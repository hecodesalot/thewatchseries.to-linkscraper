package seriesripper

import (
	"crawler/webstuff"
	"encoding/base64"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"gopkg.in/cheggaaa/pb.v1"
)

// SeriesLink is the check string for a series link
var SeriesLink = "serie"

// SeasonLink is the link for a season link
var SeasonLink = "season"

// EpisodeLink is the check string for an episode link
var EpisodeLink = "episode"

// DownloadLink is the download link
var DownloadLink = "/cale.html?r="

//DomainURL current working domain
var DomainURL = "http://thewatchseries.to"

// SubmittedURL is the link to start with
var SubmittedURL string

// OutputFileName is the name of the output file
var OutputFileName string

// LinkType is the type of link (series/episode/download)
var LinkType int

//SeriesArray array of links for sorting
var SeriesArray = []string{}

//DownLinks just of all unsorted downloads
var DownLinks []string

var eplinks []string

//CheckLinkType check what sort of link has been submitted
func CheckLinkType() {
	linktype := 0
	fmt.Println()
	CheckURL := strings.Replace(SubmittedURL, DomainURL, "", -1)
	if strings.Contains(CheckURL, SeriesLink) == true {
		linktype = 1

	} else if strings.Contains(CheckURL, SeasonLink) == true {
		linktype = 2

	} else if strings.Contains(CheckURL, EpisodeLink) == true {
		linktype = 3
	}

	switch linktype {
	case 1:
		fmt.Println("Detected Series Link!")
		currentLinks := GetSeasonsLinks(SubmittedURL)
		fmt.Println("Found ", len(currentLinks), " seasons!")
		doSeasonsLinks(currentLinks)
		processEpsList()
		if len(webstuff.ErrorURLs) > 0 {
			fmt.Println("Retrying urls with errors")
			for _, faillink := range webstuff.ErrorURLs {
				GetDownloadLink(faillink, CleanStrings(faillink))
			}
		}
		SortLinkArray()
	case 2:
		fmt.Println("Detected Season Link!")
		doSeasonsLinks([]string{SubmittedURL})
		processEpsList()
		if len(webstuff.ErrorURLs) > 0 {
			fmt.Println("Retrying urls with errors")
			for _, faillink := range webstuff.ErrorURLs {
				GetDownloadLink(faillink, CleanStrings(faillink))
			}
		}
		SortLinkArray()
	case 3:
		fmt.Println("Detected Episode Link!")
		GetDownloadLink(SubmittedURL, CleanStrings(SubmittedURL))
		SortLinkArray()
	}
}

//GetSeasonsLinks grab all the links for a seasons
func GetSeasonsLinks(url string) []string {
	links := webstuff.GetWebpage(SubmittedURL)
	var seriesLinks []string
	for _, link := range links {
		if CheckLink(link, SeasonLink) == true {
			seriesLinks = append(seriesLinks, link)
		}
	}
	return seriesLinks
}

func doSeasonsLinks(urls []string) {
	var wg sync.WaitGroup
	for _, eps := range urls {
		wg.Add(1)
		go GetEpisodes(eps, &wg)
	}
	wg.Wait()
}

func processEpsList() {
	jobs := make(chan int)
	results := make(chan int)
	length := len(eplinks)
	fmt.Println("Proccessing Episodes...")
	bar := pb.StartNew(length)
	for w := 1; w <= length; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= length; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= length; a++ {
		<-results
		bar.Increment()
	}
	bar.FinishPrint("Complete")
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		buildLink := "http://thewatchseries.to" + eplinks[id-1]
		GetDownloadLink(buildLink, eplinks[id-1])
		time.Sleep(time.Second)
		results <- j * 2
	}
}

//GetEpisodes grabs a list of all episodes
func GetEpisodes(url string, wg *sync.WaitGroup) {
	links := webstuff.GetWebpage(url)
	var episodeLinks []string
	for _, link := range links {
		if CheckLink(link, EpisodeLink) == true {
			episodeLinks = append(episodeLinks, link)
		}
	}
	eplinks = append(eplinks, episodeLinks...)
	wg.Done()
}

//GetDownloadLink grabs a list of all download links
func GetDownloadLink(url string, name string) {
	links := webstuff.GetWebpage(url)
	var downloadLinks []string
	for _, link := range links {
		if CheckLink(link, DownloadLink) == true {
			downloadLinks = append(downloadLinks, DecodeDownloadLink(CleanStrings(link)))
		}
	}
	BuildLinkArray(name, downloadLinks)
	DownLinks = append(DownLinks, downloadLinks...)
}

//DecodeDownloadLink fixes string and base64 decodes link
func DecodeDownloadLink(base64String string) string {
	cleanString := strings.Replace(base64String, "/cale.html?r=", "", -1)
	decoded, _ := base64.StdEncoding.DecodeString(cleanString)
	return string(decoded)
}

//CleanStrings Cleans strings for output
func CleanStrings(inString string) string {
	newString := inString
	newString = strings.Replace(newString, "_", " ", -1)
	newString = strings.Replace(newString, "-", " ", -1)
	newString = strings.Replace(newString, "/episode/", "", -1)
	newString = strings.Replace(newString, ".html", "", -1)
	newString = strings.Replace(newString, "/cale?r=", "", -1)
	newString = strings.Replace(newString, DomainURL, "", -1)
	return newString
}

//WriteLinksFile writes link to an output file
func WriteLinksFile(currentLink string) {
	if _, err := os.Stat(OutputFileName); os.IsNotExist(err) {
		os.Create(OutputFileName)
	}
	f, err := os.OpenFile(OutputFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	f.WriteString(currentLink + "\n")
	f.Close()
}

//CheckLink Check if link contains a sub string
func CheckLink(url string, checkStr string) bool {
	if strings.Contains(url, checkStr) == true {
		return true
	}
	return false
}

//BuildLinkArray adds links to array for sorting
func BuildLinkArray(name string, value []string) {
	name = CleanStrings(name)
	data := append([]string{name}, value...)
	concatString := strings.Join(data, "\n")
	SeriesArray = append(SeriesArray, concatString)
}

//SortLinkArray sorts SeriesArray
func SortLinkArray() {
	sort.Strings(SeriesArray)
	for _, sLinks := range SeriesArray {
		WriteLinksFile(sLinks)
	}
}
