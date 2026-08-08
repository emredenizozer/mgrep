# mgrep - Concurrent Grep in Go

A lightweight, fast, and idiomatic Go implementation of a multi-threaded grep utility. `mgrep` recursively searches for a specific string across files in a given directory, utilizing Go's concurrency model (goroutines and channels) for maximum performance and efficiency.

## Features

- **Concurrent Execution:** Uses a worker pool pattern to process multiple files simultaneously.
- **Memory Efficient:** Streams results in real-time through channels rather than loading all matches into memory.
- **Idiomatic Go:** Safely manages goroutines using `sync.WaitGroup` and channel closures, preventing goroutine leaks and race conditions.
- **Fast Directory Traversal:** Utilizes the optimized `filepath.WalkDir` for fast and safe recursive directory reading.

## Installation

Ensure you have [Go](https://golang.org/doc/install) installed. Clone the repository and build the binary:

```bash
# Clone the repository
git clone https://github.com/emredenizozer/mgrep.git
cd mgrep

# Download dependencies
go mod tidy

# Build the binary (creates the 'mgrep' executable in the root)
go build -o mgrep
```

## Usage

The tool requires two positional arguments: the search term and the target directory.

**Using the compiled binary**
```bash
./mgrep <search-term> <target-directory>
```

### Example

Search for the word "Result" in "worker" directory:

```bash
./mgrep Result worker
```

**Output:**
```
worker/worker.go[10]: type Result struct {
worker/worker.go[16]: func FindInFile(path string, searchTerm string, results chan<- Result) {
worker/worker.go[30]:                   results <- Result{
```

## Architecture

1. **Producer:** A single goroutine uses `filepath.WalkDir` to traverse the directory structure and sends file paths to a buffered `jobs` channel.
2. **Worker Pool:** A pool of worker goroutines reads from the `jobs` channel. Each worker opens the file, scans it line by line, and sends any matching lines to a `results` channel.
3. **Consumer:** The main goroutine ranges over the `results` channel, printing matches to the standard output in real-time until all workers complete their tasks and the channel is closed.

## Dependencies

- [go-arg](https://github.com/alexflint/go-arg): Struct-based argument parsing in Go.