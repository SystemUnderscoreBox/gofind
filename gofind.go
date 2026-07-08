package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// gofind — the soon-to-be replacement for the Unix command "find"
// Plan:
// 2. Create test file
// 3. Compile binary and run
// 4. In the future, expand functionalities and maybe implement "grep"

// Variable slice of matches and color codes for better output.
var (
	matches []string
	Red     = "\033[31m"
	Reset   = "\033[0m"
)

func main() {
	t := time.Now()
	if err := handleInput(os.Args[1]); err != nil {
		fmt.Printf("find failed, error: %v", err)
		e := time.Now()
		fmt.Println("time elapsed: " + strconv.FormatFloat(e.Sub(t).Seconds(), 'f', -1, 64) + " seconds")
		os.Exit(1)
	}

	e := time.Now()
	fmt.Println("time elapsed: " + strconv.FormatFloat(e.Sub(t).Seconds(), 'f', -1, 64) + " seconds")

	os.Exit(0)
}

func handleInput(input string) error {
	fmt.Printf("You wrote: %s\n", input)

	// Next, start to look for given file.
	dir, err := os.Getwd()

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return errors.New("failed to get current work dir")
	}

	entries, err := os.ReadDir(dir)

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return errors.New("failed to read dir entries")
	}

	// Call this function recursively in case of nested folders.
	switch err := helperFunction(entries, input, dir); err.Error() {
	// "no match" means complete search has completed, regardless if found or not.
	case "no match":
		if len(matches) == 0 {
			return errors.New("found no matches\n")
		}

		for match := range len(matches) {
			coloredMatch := Red + matches[match] + "/" + input + Reset
			fmt.Printf("found: %s\n", coloredMatch)
		}
	// If for some reason a failure happens, return err.
	case "fail":
		return err
	}

	return nil
}

func helperFunction(entries []os.DirEntry, input, dir string) error {
	for i := range len(entries) {
		entry := entries[i]

		if entry.IsDir() {
			dir = dir + "/" + entry.Name()
			fmt.Println(dir)

			entriesNew, err := os.ReadDir(dir)

			// Skip to next entry in case ReadDir fails.
			// Prints out error.
			if err != nil {
				fmt.Printf("error: %v\n", err)
				dir = strings.TrimSuffix(dir, entry.Name())
				dir = strings.TrimSuffix(dir, "/")
				continue
			}

			if len(entriesNew) == 0 {
				dir = strings.TrimSuffix(dir, entry.Name())
				dir = strings.TrimSuffix(dir, "/")
				continue
			}

			switch err := helperFunction(entriesNew, input, dir); err.Error() {
			// If folder has no matching files, go back and go to next entry.
			case "no match":
				dir = strings.TrimSuffix(dir, entry.Name())
				dir = strings.TrimSuffix(dir, "/")
				continue
			// If for some reason a failure happens, return err.
			case "fail":
				return err
			}
		}

		// Here we check the entry's name, at this point, entry must be a file.
		if input == entry.Name() {
			matches = append(matches, dir)
			continue
		}
	}

	return errors.New("no match")
}
