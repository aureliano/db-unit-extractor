package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/aureliano/db-unit-extractor/cmd"
	_ "github.com/sijms/go-ora/v2"
)

func main() {
	cpuProfile := os.Getenv("CPU_PROFILE")
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	memProfile := os.Getenv("MEM_PROFILE")
	if memProfile != "" {
		f, err := os.Create(memProfile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.WriteHeapProfile(f)
		defer f.Close()
	}

	err := cmd.NewRootCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}
