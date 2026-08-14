// gofind — the soon-to-be replacement for the Unix command "find"
package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
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
	skip      []string
	time      bool
	help      bool
	file      string
	literal   bool
}

// parseStruct reads input flags as populates ParameterFlags.
func parseStruct() (*ParameterFlags, error) {
	pF := ParameterFlags{}

	for i, arg := range os.Args {
		// Skip first argument (i.e. gofind)
		if i == 0 {
			continue
		}

		// When last arg is the file to be found populate field
		// Handle also literal and non-literal file names
		// Handle case when source dir is given and index on source, break and populate pF.file
		if !pF.literal && i == len(os.Args)-1 && !strings.Contains(arg, "-") {
			pF.file = arg
			break
		} else if pF.literal && i == len(os.Args)-1 {
			pF.file = arg
			break
		}

		// Continue if previous flag was -s or -S.
		// This way the next argument will not be handled by the switch-case default.
		if os.Args[i-1] == "-s" || os.Args[i-1] == "--source" || os.Args[i-1] == "-S" || os.Args[i-1] == "--skip" {
			continue
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
		case "-S":
			pF.skip = parseSkipInput(i)
		case "--skip":
			pF.skip = parseSkipInput(i)
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
			fmt.Printf("unknown flag: %v\n", arg)
			return nil, errors.New("invalid flag input")
		}
	}
	return &pF, nil
}

// Helper function for parsing the input when -S/--skip flag is used.
func parseSkipInput(i int) []string {
	var skips []string

	for file := range strings.SplitSeq(os.Args[i+1], " ") {
		if file == "" {
			continue
		}
		skips = append(skips, file)
	}
	return skips
}

// Main function of gofind.
func main() {
	// Start by parsing inputs
	params, err := parseStruct()
	if params == nil {
		fmt.Printf("struct parsing failed: %v\n", err)
		os.Exit(0)
	}
	// If only prompted with gofind or with help flag, exit with code 0.
	if params.help || len(os.Args) == 1 {
		handleHelp()
		os.Exit(0)
	}
	findFile := params.file
	// In case of regular expressions including asterisks (*) in input
	if strings.Contains(findFile, "*") && len(findFile) == 1 {
		fmt.Println("will not search for every file")
		os.Exit(0)
	} else if strings.Contains(findFile, "*") && !params.literal {
		t := time.Now()

		// Use handleRegex for regex searches
		if err := handleRegex(findFile, params); err != nil {
			fmt.Printf("find failed, error: %v\n", err)
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
		fmt.Printf("find failed, error: %v\n", err)
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

// handleHelp prints out contents of help.txt to terminal
func handleHelp() {
	data, err := helpFile.ReadFile("stdout/help.txt")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	fmt.Println(string(data))
}

// handleRegex searches for given file with regex.
func handleRegex(input string, params *ParameterFlags) error {
	// In case of source dir given, change working dir to source dir.
	if params.source != "" {
		if err := os.Chdir(params.source); err != nil {
			fmt.Printf("error: %v\n", err)
			return errors.New("failed to get source dir")
		}
	}

	var pattern string
	// Iterate over the given input to construct final pattern.
	for seq := range strings.SplitSeq(input, "*") {
		if strings.EqualFold(seq, "") {
			pattern = pattern + ".*"
		} else {
			pattern = pattern + seq
		}
	}
	// Final constructed regex pattern (e.g. "*.mp4").
	re, err := regexp.Compile(pattern)
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
	// "no match" means search has completed, regardless if found or not.
	case "no match":
		if len(matches) == 0 {
			fmt.Println("found no matches")
			return nil
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

// Helper function that recursively searches through folder tree until match(es) found.
func handleRegexHelper(entries []os.DirEntry, dir string, re *regexp.Regexp, params *ParameterFlags) error {
	for i := range len(entries) {
		entry := entries[i]

		if entry.IsDir() {
			if slices.Index(params.skip, entry.Name()) != -1 {
				if params.verbose {
					fmt.Printf("skipping directory: %s\n", dir+"/"+entry.Name())
				}
				continue
			}

			// Use when searching for directories
			if params.directory && re.MatchString(entry.Name()) {
				dir = dir + "/" + entry.Name()
				matches = append(matches, dir)
				dir = strings.TrimSuffix(dir, entry.Name())
				dir = strings.TrimSuffix(dir, "/")
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

// handleInput searches for given input (FILE).
func handleInput(input string, params *ParameterFlags) error {
	// In case of source dir given, change working dir to source dir.
	if params.source != "" {
		if err := os.Chdir(params.source); err != nil {
			fmt.Printf("error: %v\n", err)
			return errors.New("failed to get source dir")
		}
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

	// Call this function recursively in case of nested folders.
	switch err := helperFunction(entries, dir, params); err.Error() {
	// "no match" means earch has completed, regardless if found or not.
	case "no match":
		if len(matches) == 0 {
			fmt.Println("found no matches")
			return nil
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

// Helper function that recursively searches through folder tree until match(es) found.
func helperFunction(entries []os.DirEntry, dir string, params *ParameterFlags) error {
	for i := range len(entries) {
		entry := entries[i]

		if entry.IsDir() {
			if slices.Index(params.skip, entry.Name()) != -1 {
				if params.verbose {
					fmt.Printf("skipping directory: %s\n", dir+"/"+entry.Name())
				}
				continue
			}

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
