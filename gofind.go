package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// gofind -- the soon-to-be replacement for the Unix command "find"
// Plan:
// 1. Make the file emulate the core functionality of "find", e.g. to find the named file
// 1.1 Currently finds the first instance of file name, must report all instances
// 1.2 At the end of each run, signal SIGSEGV gets outputed. Figure that out later.
// 2. Create test file
// 3. Compile binary and run
// 4. In the future, expand functionalities and maybe implement "grep"

func main() {
	if err := handleInput(os.Args[1]); err != nil {
		fmt.Printf("find failed, error: %v", err)
		os.Exit(1)
	}

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
	if err := helperFunction(entries, input, dir); err != nil {
		return errors.New("no match found, sorry!\n")
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

			if err != nil {
				fmt.Printf("error: %v\n", err)
				return errors.New("failed to read dir entries")
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
			fmt.Printf("match found in directory: %s\n", dir)
			return nil
		}
	}

	return errors.New("no match")
}
