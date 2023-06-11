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
	var file *os.File
	var err error
	if cpuProfile != "" {
		file, err = os.Create(cpuProfile)
		if err != nil {
			log.Fatal(err)
		}

		if err = pprof.StartCPUProfile(file); err != nil {
			log.Fatal(err)
		}
	}

	memProfile := os.Getenv("MEM_PROFILE")
	if memProfile != "" {
		file, err = os.Create(memProfile)
		if err != nil {
			log.Fatal(err)
		}

		if err = pprof.WriteHeapProfile(file); err != nil {
			log.Fatal(err)
		}
	}

	err = cmd.NewRootCommand().Execute()
	if cpuProfile != "" {
		pprof.StopCPUProfile()
	}

	if file != nil {
		_ = file.Close()
	}

	if err != nil {
		os.Exit(1)
	}
}
