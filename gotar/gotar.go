package main

import (
	"flag"
	. "fmt"
	"goarchive"
	"os"
)

var (
	verbose = flag.Bool("v", false, "verbose")
)

func usage() {
	Println("Usage:", "gotar [flags] [file]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.Arg(0) == "" {
		usage()
	}
	file := flag.Arg(0)
	zip, _ := goarchive.NewZip(file)
	zip.Verbose = *verbose
	if err := zip.Decompress("./"); err != nil {
		Printf("%v\n", err)
	}
}
