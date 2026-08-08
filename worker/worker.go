package worker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Result struct {
	Line       string
	LineNumber int
	Path       string
}

func FindInFile(path string, searchTerm string, results chan<- Result) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file (%s): %v\n", path, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, searchTerm) {
			results <- Result{
				Line:       line,
				LineNumber: lineNumber,
				Path:       path,
			}
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file (%s): %v\n", path, err)
	}
}
