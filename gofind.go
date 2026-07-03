package main

import (
	"bufio"
	"fmt"
	"os"
)

// gofind -- the soon-to-be replacement for the Unix command "find"
// Plan:
// 1. Make the file emulate the core functionality of "find", e.g. to find the named file
// 2. Create test file
// 3. Compile binary and run
// 4. In the future, expand functionalities and maybe implement "grep"

func main() {
	fmt.Println("Hello World!")

	if err := handleInput(); err != nil {
		fmt.Printf("find failed, error: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func handleInput() error {
	fmt.Println("Provide file name with suffix")
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	fmt.Printf("You wrote: %s\n", scanner.Text())

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error: %v", err)
	}

	// Next, start to look for given file.
	file := scanner.Text()
	dir, err := os.Getwd()

	if err != nil {
		return fmt.Errorf("failed to get current work dir: %v", err)
	}

	entries, err := os.ReadDir(dir)

	if err != nil {
		return fmt.Errorf("failed to read dir entries: %v", err)
	}

	// Call this function recursively in case of nested folders.
	if err := helperFunction(entries, file, dir); err != nil {
		return fmt.Errorf("no match found, sorry!\n")
	}

	return nil
}

func helperFunction(entries []os.DirEntry, file, dir string) error {

	for i := range len(entries) {
		entry := entries[i]

		if entry.IsDir() {
			dir = dir + "/" + entry.Name()
			fmt.Println(dir)

			entriesNew, err := os.ReadDir(dir)

			if err != nil {
				return fmt.Errorf("failed to read dir entries: %v", err)
			}

			if err := helperFunction(entriesNew, file, dir); err != nil {
				return fmt.Errorf("no match")
			}

		}

		// Here we check the entry's name, at this point, entry must be a file.
		if file == entry.Name() {
			fmt.Printf("match found in directory: %s\n", dir)
			return nil
		}
	}

	return fmt.Errorf("no match")
}
