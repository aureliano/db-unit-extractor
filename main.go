package main

import (
	"os"

	"github.com/aureliano/db-unit-extractor/cmd"
)

func main() {
	err := cmd.NewRootCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}
