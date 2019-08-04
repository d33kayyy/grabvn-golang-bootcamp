package main

import (
	"bufio"
	"errors"
	"fmt"
	"grabvn-golang-bootcamp/week02/pkg"
	"os"
	"path/filepath"
	"sync"
)

// NOTE: func taken from https://blog.golang.org/pipelines
// walkFiles starts a goroutine to walk the directory tree at root and send the
// path of each regular file on the string channel. It sends the result of the
// walk on the error channel. If done is closed, walkFiles abandons its work.
func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)
	go func() {
		// Close the paths channel after Walk returns.
		defer close(paths)
		// No select needed for this send, since errc is buffered.
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			select {
			case paths <- path:
			case <-done:
				return errors.New("walk canceled")
			}
			return nil
		})
	}()
	return paths, errc
}

type result struct {
	counter map[string]int
	err     error
}

// Read and return a list of words in a file
func listWordInFile(path string) (words []string, err error) {
	file, err := os.Open(path)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords) // avoid buffer overflow when reading big file
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	return
}

// Count number of words in a file
func countWordInFile(path string) (map[string]int, error) {
	words, err := listWordInFile(path)
	m := make(map[string]int)
	for _, word := range words {
		m[word] += 1
	}
	return m, err
}

// Count number of words in each file then send result to a channel
func collectWordCounter(done <-chan struct{}, paths <-chan string, sink chan<- result) {
	for path := range paths {
		m, err := countWordInFile(path)
		select {
		case sink <- result{m, err}:
		case <-done:
			return
		}
	}
}

// Read all result from a channel and combine the values into a map
func combineResult(sink <-chan result) (map[string]int, error) {
	m := make(map[string]int)
	for r := range sink {
		if r.err != nil {
			return nil, r.err
		}
		m = pkg.Reduce(m, r.counter)
	}
	return m, nil
}

// List all readable files in folder and start goroutines to count words occurrence in each file.
// Result will be send to a channel from which data will be combined and return as a map
func countWordAll(root string) (map[string]int, error) {
	// countWordAll closes the done channel when it returns; it may do so before
	// receiving all the values from sink and errc.
	done := make(chan struct{})
	defer close(done)

	paths, errc := walkFiles(done, root)

	// Start a fixed number of goroutines to read and digest files.
	sink := make(chan result)
	var wg sync.WaitGroup

	// since we only read word by word, it is safe to have high number of workers (goroutines)
	const numWorker = 100
	wg.Add(numWorker)
	for i := 0; i < numWorker; i++ {
		go func() {
			collectWordCounter(done, paths, sink)
			defer wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(sink)
	}()

	// Check whether the Walk failed.
	if err := <-errc; err != nil {
		return nil, err
	}
	m, err := combineResult(sink)
	return m, err
}

func main() {
	m, err := countWordAll(os.Args[1])
	//m, err := countWordAll(".../../tmp")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(m)
}
