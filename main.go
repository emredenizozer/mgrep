package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"mgrep/worker"

	"github.com/alexflint/go-arg"
)

var args struct {
	SearchTerm string `arg:"positional,required"`
	SearchDir  string `arg:"positional"`
}

func discoverDirs(jobs chan<- string, dir string) {
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Access error (%s): %v\n", path, err)
			return nil
		}
		if !d.IsDir() {
			jobs <- path
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Directory scan error: %v\n", err)
	}

	close(jobs)
}

func main() {
	arg.MustParse(&args)

	jobs := make(chan string, 100)
	results := make(chan worker.Result, 100)

	go discoverDirs(jobs, args.SearchDir)

	var workersWg sync.WaitGroup
	numWorkers := 10

	for i := 0; i < numWorkers; i++ {
		workersWg.Add(1)
		go func() {
			defer workersWg.Done()

			for path := range jobs {
				worker.FindInFile(path, args.SearchTerm, results)
			}
		}()
	}

	go func() {
		workersWg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Printf("%v[%v]: %v\n", r.Path, r.LineNumber, r.Line)
	}
}
