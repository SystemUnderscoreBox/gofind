// gofind — the soon-to-be replacement for the Unix command "find"
package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// Slice of matches
	matches []string
	Red     = "\033[31m"
	Reset   = "\033[0m"

	//go:embed stdout/help.txt
	helpFile embed.FS
)

// Struct containing all available flags
type ParameterFlags struct {
	verbose   bool
	directory bool
	source    string
	time      bool
	help      bool
	file      string
	literal   bool
}

func parseStruct() *ParameterFlags {
	pF := ParameterFlags{}

	for i, arg := range os.Args {
		// Skip first argument (i.e. gofind)
		if i == 0 {
			continue
		}

		// When last arg is the file to be found populate field
		// Handle also literal and non-literal file names
		if !pF.literal && i == len(os.Args)-1 && !strings.Contains(arg, "-") {
			pF.file = arg
			break
		} else if pF.literal && i == len(os.Args)-1 {
			pF.file = arg
			break
		}

		switch arg {
		case "-v":
			pF.verbose = true
		case "--verbose":
			pF.verbose = true
		case "-d":
			pF.directory = true
		case "--directory":
			pF.directory = true
		case "-s":
			pF.source = os.Args[i+1]
		case "--source":
			pF.source = os.Args[i+1]
		case "-t":
			pF.time = true
		case "--time-elapsed":
			pF.time = true
		case "-h":
			pF.help = true
		case "--help":
			pF.help = true
		case "-l":
			pF.literal = true
		case "--literal":
			pF.literal = true
		default:
			fmt.Println("unknown input: use --help for list of flags")
			return nil
		}
	}

	return &pF
}

func main() {
	// Start by parsing inputs
	params := parseStruct()
	if params == nil {
		os.Exit(0)
	}

	if params.help {
		handleHelp()
		os.Exit(0)
	}

	findFile := params.file
	// In case of regular expressions including one asterisk (*) in input
	// Works (currently) for simple queries, e.g. *.pdf
	if strings.Contains(findFile, "*") && len(findFile) == 1 {
		fmt.Println("will not search for every file")
		os.Exit(0)
	} else if strings.Contains(findFile, "*") {
		t := time.Now()

		// Use handleRegex for regex searches
		if err := handleRegex(findFile, params); err != nil {
			fmt.Printf("find failed, error: %v", err)
			e := time.Now()
			if params.time {
				fmt.Println("time elapsed: " + strconv.FormatFloat(e.Sub(t).Seconds(), 'f', -1, 64) + " seconds")
			}
			os.Exit(1)
		}

		e := time.Now()
		if params.time {
			fmt.Println("time elapsed: " + strconv.FormatFloat(e.Sub(t).Seconds(), 'f', -1, 64) + " seconds")
		}

		os.Exit(0)
	}

	t := time.Now()
	if err := handleInput(findFile, params); err != nil {
		fmt.Printf("find failed, error: %v", err)
		e := time.Now()
		if params.time {
			fmt.Println("time elapsed: " + strconv.FormatFloat(e.Sub(t).Seconds(), 'f', -1, 64) + " seconds")
		}
		os.Exit(1)
	}

	e := time.Now()
	if params.time {
		fmt.Println("time elapsed: " + strconv.FormatFloat(e.Sub(t).Seconds(), 'f', -1, 64) + " seconds")
	}

	os.Exit(0)
}

func handleHelp() {
	data, err := helpFile.ReadFile("stdout/help.txt")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	fmt.Println(string(data))
}

func handleRegex(input string, params *ParameterFlags) error {
	seqs := strings.Split(input, "*")

	// TODO: Currently only supports one instance (e.g. "*.mp4")
	re, err := regexp.Compile(".*(" + seqs[1] + ")")

	if err != nil {
		return err
	}

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

	// Call this function recursively.
	switch err := handleRegexHelper(entries, dir, re, params); err.Error() {
	// "no match" means complete search has completed, regardless if found or not.
	case "no match":
		if len(matches) == 0 {
			return errors.New("found no matches\n")
		}

		for match := range len(matches) {
			coloredMatch := Red + matches[match] + Reset
			fmt.Printf("found: %s\n", coloredMatch)
		}
	// If for some reason a failure happens, return err.
	case "fail":
		return err
	}

	return nil
}

func handleRegexHelper(entries []os.DirEntry, dir string, re *regexp.Regexp, params *ParameterFlags) error {
	for i := range len(entries) {
		entry := entries[i]

		if entry.IsDir() {
			// Use when searching for directories
			if params.directory && params.file == entry.Name() {
				matches = append(matches, dir)
			}

			dir = dir + "/" + entry.Name()

			if params.verbose {
				fmt.Println(dir)
			}

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

			switch err := handleRegexHelper(entriesNew, dir, re, params); err.Error() {
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
		str := re.FindString(entry.Name())
		if !params.directory && str == entry.Name() {
			dir = dir + "/" + re.FindString(entry.Name())
			matches = append(matches, dir)
			dir = strings.TrimSuffix(dir, re.FindString(entry.Name()))
			dir = strings.TrimSuffix(dir, "/")
			continue
		}
	}

	return errors.New("no match")
}

func handleInput(input string, params *ParameterFlags) error {
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
	switch err := helperFunction(entries, dir, params); err.Error() {
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

func helperFunction(entries []os.DirEntry, dir string, params *ParameterFlags) error {
	for i := range len(entries) {
		entry := entries[i]

		if entry.IsDir() {
			// Use when searching for directories
			if params.directory && params.file == entry.Name() {
				matches = append(matches, dir)
			}

			dir = dir + "/" + entry.Name()

			if params.verbose {
				fmt.Println(dir)
			}

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

			switch err := helperFunction(entriesNew, dir, params); err.Error() {
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
		if !params.directory && params.file == entry.Name() {
			matches = append(matches, dir)
			continue
		}
	}

	return errors.New("no match")
}
