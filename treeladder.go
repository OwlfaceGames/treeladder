package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorRed    = "\033[31m"
	colorPurple = "\033[35m"
)

// Version information - will be set during build
var (
	Version   = "dev"
	BuildDate = "unknown"
	GitCommit = "unknown"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: treeladder create <repo_name>\n")
		return
	}

	command := args[1]

	if command == "version" || command == "--version" || command == "-v" {
		fmt.Printf("TreeLadder %s\n", Version)
		fmt.Printf("Build date: %s\n", BuildDate)
		fmt.Printf("Git commit: %s\n", GitCommit)
		return
	}

	if command == "create" {
		if len(args) < 3 {
			fmt.Printf("Please provide a repository name\n")
			return
		}

		fmt.Printf("%sTreeLadder%s\n", colorGreen, colorReset)
		repoName := args[2]
		createRepo(repoName)
	} else {
		fmt.Printf("Unknown command. Use 'treeladder create <repo_name>'\n")
	}
}

// Helper function to check if response is affirmative
func isAffirmative(response string) bool {
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "yes" || response == "y"
}

func createRepo(repoName string) {
	// Create the root directory
	fmt.Printf("Creating repository: %s%s%s\n", colorCyan, repoName, colorReset)
	err := os.Mkdir(repoName, 0o755)
	if err != nil {
		fmt.Printf("Error creating repository: %v\n", err)
		return
	}

	fmt.Printf("Repository created successfully\n")

	// Change to the new directory
	err = os.Chdir(repoName)
	if err != nil {
		fmt.Printf("Error changing to repository directory: %v\n", err)
		return
	}

	// Get the absolute path of the root directory
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	// Display current path
	displayCurrentPath()

	// Ask if user wants to create files or folders
	fmt.Printf("Would you like to create any files and/or folders?\n(y/n): ")
	response, _ := reader.ReadString('\n')

	if !isAffirmative(response) {
		fmt.Printf("No files or folders created. Exiting.\n")
		return
	}

	// Create folders with recursive structure
	createFolders(reader, rootDir)

	// Create files at the root level
	createFiles(reader)

	fmt.Printf("Project structure created successfully!\n")
}

func displayCurrentPath() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	fmt.Printf("%sCurrent path: %s%s\n", colorCyan, currentDir, colorReset)
}

func createFolders(reader *bufio.Reader, rootDir string) {
	// Display current path
	displayCurrentPath()

	fmt.Printf("Would you like to create any folders?\n(y/n): ")
	response, _ := reader.ReadString('\n')

	if !isAffirmative(response) {
		return
	}

	fmt.Printf("How many folders would you like to create? ")
	countStr, _ := reader.ReadString('\n')
	countStr = strings.TrimSpace(countStr)
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		fmt.Printf("Invalid number of folders. Skipping.\n")
		return
	}

	// Store current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	for i := 1; i <= count; i++ {
		fmt.Printf("Enter name for folder %d: ", i)
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		if name == "" {
			fmt.Printf("Empty name provided for folder %d. Skipping.\n", i)
			continue
		}

		err = os.Mkdir(name, 0o755)
		if err != nil {
			fmt.Printf("Error creating folder '%s': %v\n", name, err)
			continue
		}
		fmt.Printf("%sCreated folder: %s%s\n", colorGreen, name, colorReset)

		// Display current path
		displayCurrentPath()

		// Ask if user wants to create content inside this folder
		fmt.Printf("Would you like to create content inside '%s%s%s'?\n(y/n): ", colorGreen, name, colorReset)
		response, _ := reader.ReadString('\n')

		if isAffirmative(response) {
			// Change to the new folder
			fmt.Printf("Entering folder: %s%s%s\n", colorGreen, name, colorReset)
			err = os.Chdir(name)
			if err != nil {
				fmt.Printf("Error changing to directory '%s': %v\n", name, err)
				continue
			}

			// Recursively create folders in this folder
			createFolders(reader, rootDir)

			// Create files in this folder
			createFiles(reader)

			// Return to the parent directory
			fmt.Printf("Returning to parent folder\n")
			err = os.Chdir(currentDir)
			if err != nil {
				fmt.Printf("Error returning to parent directory: %v\n", err)
				// Try to force return to root directory
				os.Chdir(rootDir)
			}
		}
	}
}

func createFiles(reader *bufio.Reader) {
	// Display current path
	displayCurrentPath()

	fmt.Printf("Would you like to create any files?\n(y/n): ")
	response, _ := reader.ReadString('\n')

	if !isAffirmative(response) {
		return
	}

	fmt.Printf("How many files would you like to create? ")
	countStr, _ := reader.ReadString('\n')
	countStr = strings.TrimSpace(countStr)
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		fmt.Printf("Invalid number of files. Skipping.\n")
		return
	}

	for i := 1; i <= count; i++ {
		fmt.Printf("Enter name for file %d: ", i)
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		if name == "" {
			fmt.Printf("Empty name provided for file %d. Skipping.\n", i)
			continue
		}

		// Create an empty file
		file, err := os.Create(name)
		if err != nil {
			fmt.Printf("Error creating file '%s': %v\n", name, err)
		} else {
			file.Close()
			fmt.Printf("%sCreated file: %s%s\n", colorYellow, name, colorReset)
		}
	}
}
