package main

import (
	"os"

	"github.com/aureliano/db-unit-extractor/cmd"
	_ "github.com/sijms/go-ora/v2"
)

func main() {
	err := cmd.NewRootCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}
