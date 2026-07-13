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
// Dynamic plan:
// Add support for regex
// Add support for flags (e.g. --help)
// --help only works if gofind called in gofind repo, must be worked out on
// Create test file
// Compile binary and run
// In the future, expand functionalities and maybe implement "grep"

// Slice of matches and color codes for better output.
var (
	matches []string
	Red     = "\033[31m"
	Reset   = "\033[0m"
)

func main() {
	// If binary is only called w/o further inputs, print help as default
	if len(os.Args) == 1 {
		handleHelp()
		os.Exit(0)
	} else if os.Args[1] == "--help" {
		handleHelp()
		os.Exit(0)
	}

	// If flag is included, check first it's in position [1]
	// If too many inputs are given, exit with code 0
	// NOTE: Currently does NOT support file name with literal '-' included
	var findFile string
	if len(os.Args) == 3 && !strings.Contains(os.Args[2], "-") {
		findFile = os.Args[2]
	} else if len(os.Args) > 3 {
		fmt.Println("too many arguments passed, only one flag is supported")
		os.Exit(0)
	} else if len(os.Args) == 2 && !strings.Contains(os.Args[1], "-") {
		findFile = os.Args[1]
	}

	t := time.Now()
	if err := handleInput(findFile); err != nil {
		fmt.Printf("find failed, error: %v", err)
		e := time.Now()
		fmt.Println("time elapsed: " + strconv.FormatFloat(e.Sub(t).Seconds(), 'f', -1, 64) + " seconds")
		os.Exit(1)
	}

	e := time.Now()
	// If flag for time was used, print time elapsed, else skip
	if os.Args[1] == "-t" || os.Args[1] == "--time-elapsed" {
		fmt.Println("time elapsed: " + strconv.FormatFloat(e.Sub(t).Seconds(), 'f', -1, 64) + " seconds")
	}

	os.Exit(0)
}

func handleHelp() {
	data, err := os.ReadFile("stdout/help.txt")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	fmt.Println(string(data))
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
