// A command line tool for clean zsh cdr history file
// This tool removes non exisitng directories from the history file

package main

import (
	"flag"
	"os"

	ccdrh "github.com/sssho/ccdrh/src"
)

type flagVar struct {
	cacheFile string
}

func parseFlags() *flagVar {
	var flags flagVar
	flag.StringVar(&flags.cacheFile, "f", "", "Path to cdr history file")
	flag.Parse()

	return &flags
}

func main() {
	flags := parseFlags()
	os.Exit(ccdrh.Run(flags.cacheFile))
}
