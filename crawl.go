package main

import (
	"crawler/misc"
	"crawler/seriesripper"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	misc.CheckArgs()
	seriesripper.SubmittedURL = misc.Website
	seriesripper.OutputFileName = misc.Outfile
	seriesripper.CheckLinkType()
	elapsed := time.Since(start)
	fmt.Println("Extracted ", len(seriesripper.DownLinks), " in ", elapsed)
}
