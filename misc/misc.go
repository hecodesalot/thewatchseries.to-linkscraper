package misc

import (
	"flag"
	"fmt"
	"os"
)

//Website args[0] to be passed back to prog
var Website string

//Outfile args[1] to the passed back to prog
var Outfile string

//CheckArgs checks for correct args
func CheckArgs() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		Usage()
	}
	Website = args[0]
	Outfile = args[1]
	PrintBanner()
}

//PrintBanner Prints banner
func PrintBanner() {
	banner := `
---------------------------------------------------------------------
_______ _______  ______ _____ _______ _______  ______ _____  _____
|______ |______ |_____/   |   |______ |______ |_____/   |   |_____]
______| |______ |    \_ __|__ |______ ______| |    \_ __|__ |
---------------------------------------------------------------------
Name: 	SeriesRipper for thewatchseries.to
Author: andmuchmore 2016
Date:   22/06/16
---------------------------------------------------------------------
`
	fmt.Print(banner)
}

//Usage displays a usage message
func Usage() {
	banner := `
---------------------------------------------------------------------
_______ _______  ______ _____ _______ _______  ______ _____  _____
|______ |______ |_____/   |   |______ |______ |_____/   |   |_____]
______| |______ |    \_ __|__ |______ ______| |    \_ __|__ |
---------------------------------------------------------------------
Name: 	SeriesRipper for thewatchseries.to
Author: andmuchmore 2016
Date:   22/06/16
---------------------------------------------------------------------
Usage: crawl http://thewatchseries.to/link/ outfile.txt
---------------------------------------------------------------------
`
	fmt.Print(banner)
	os.Exit(0)
}
