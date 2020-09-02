// A command line tool for clean zsh cdr history file
// This tool removes non exisitng directories from the history file

package main

import (
	"os"

	ccdrh "github.com/sssho/ccdrh/src"
)

func main() {
	os.Exit(ccdrh.Run())
}
