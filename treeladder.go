package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	// "path/filepath"
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

// Symbols - using Nerd Font symbols and keeping tree emoji
const (
	treeSymbol      = "🌳" // Tree for the main program (kept as emoji)
	folderSymbol    = ""  // Folder (nerd font)
	fileSymbol      = ""  // File (nerd font)
	questionSymbol  = ""  // Question (nerd font)
	successSymbol   = ""  // Success (nerd font)
	errorSymbol     = ""  // Error (nerd font)
	pathSymbol      = ""  // Navigation/path (nerd font)
	gitSymbol       = ""  // Git (nerd font)
	branchSymbol    = ""  // Branch (nerd font)
	constructSymbol = ""  // Construction/creation (nerd font)
	completeSymbol  = "🎉" // Celebration/completion (kept as emoji)
	warningSymbol   = ""  // Warning (nerd font)
	enterSymbol     = ""  // Input/enter (nerd font)
	skipSymbol      = "ﭥ" // Skip (nerd font)
	returnSymbol    = ""  // Return (nerd font)
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("%s Usage: treeladder create <repo_name>\n", warningSymbol)
		return
	}

	command := args[1]

	if command == "create" {
		if len(args) < 3 {
			fmt.Printf("%s Please provide a repository name\n", warningSymbol)
			return
		}

		fmt.Printf("%s %s TreeLadder%s\n", treeSymbol, colorGreen, colorReset)
		repoName := args[2]
		createRepo(repoName)
	} else {
		fmt.Printf("%s Unknown command. Use 'treeladder create <repo_name>'\n", warningSymbol)
	}
}

// Helper function to check if response is affirmative
func isAffirmative(response string) bool {
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "yes" || response == "y"
}

// Helper function to check if response is negative
func isNegative(response string) bool {
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "no" || response == "n"
}

func createRepo(repoName string) {
	// Create the root directory
	fmt.Printf("%s Creating repository: %s%s%s\n", constructSymbol, colorCyan, repoName, colorReset)
	err := os.Mkdir(repoName, 0o755)
	if err != nil {
		fmt.Printf("%s Error creating repository: %v\n", errorSymbol, err)
		return
	}

	fmt.Printf("%s Repository created successfully\n", successSymbol)

	// Change to the new directory
	err = os.Chdir(repoName)
	if err != nil {
		fmt.Printf("%s Error changing to repository directory: %v\n", errorSymbol, err)
		return
	}

	// Get the absolute path of the root directory
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s Error getting current directory: %v\n", errorSymbol, err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	// Display current path
	displayCurrentPath()

	// Ask if user wants to create files or folders
	fmt.Printf("%s Would you like to create any files and/or folders?\n(y/n): ", questionSymbol)
	response, _ := reader.ReadString('\n')

	if !isAffirmative(response) {
		fmt.Printf("%s No files or folders created. Exiting.\n", skipSymbol)
		return
	}

	// Create folders with recursive structure
	createFolders(reader, rootDir)

	// Create files at the root level
	createFiles(reader)

	// Display current path
	displayCurrentPath()

	// Ask if user wants to initialize git repository
	fmt.Printf("%s Would you like to initialize a git repository?\n(y/n): ", gitSymbol)
	response, _ = reader.ReadString('\n')

	if isAffirmative(response) {
		fmt.Printf("%s Initializing git repository...\n", gitSymbol)
		cmd := exec.Command("git", "init")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%s Error initializing git repository: %v\n", errorSymbol, err)
		} else {
			fmt.Printf("%s %s\n", successSymbol, strings.TrimSpace(string(output)))
		}
	}

	fmt.Printf("%s Project structure created successfully! %s\n", completeSymbol, treeSymbol)
}

func displayCurrentPath() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s Error getting current directory: %v\n", errorSymbol, err)
		return
	}

	// Get relative path for display
	// dir := filepath.Base(currentDir)

	fmt.Printf("%s %sCurrent path: %s%s\n", pathSymbol, colorCyan, currentDir, colorReset)
}

func createFolders(reader *bufio.Reader, rootDir string) {
	// Display current path
	displayCurrentPath()

	fmt.Printf("%s Would you like to create any folders?\n(y/n): ", folderSymbol)
	response, _ := reader.ReadString('\n')

	if !isAffirmative(response) {
		return
	}

	fmt.Printf("%s How many folders would you like to create? ", folderSymbol)
	countStr, _ := reader.ReadString('\n')
	countStr = strings.TrimSpace(countStr)
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		fmt.Printf("%s Invalid number of folders. Skipping.\n", warningSymbol)
		return
	}

	// Store current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s Error getting current directory: %v\n", errorSymbol, err)
		return
	}

	for i := 1; i <= count; i++ {
		fmt.Printf("%s Enter name for folder %d: ", enterSymbol, i)
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		if name == "" {
			fmt.Printf("%s Empty name provided for folder %d. Skipping.\n", skipSymbol, i)
			continue
		}

		err = os.Mkdir(name, 0o755)
		if err != nil {
			fmt.Printf("%s Error creating folder '%s': %v\n", errorSymbol, name, err)
			continue
		}
		fmt.Printf("%s %sCreated folder: %s%s\n", successSymbol, colorGreen, name, colorReset)

		// Display current path
		displayCurrentPath()

		// Ask if user wants to create content inside this folder
		fmt.Printf("%s Would you like to create content inside '%s%s%s'?\n(y/n): ", branchSymbol, colorGreen, name, colorReset)
		response, _ := reader.ReadString('\n')

		if isAffirmative(response) {
			// Change to the new folder
			fmt.Printf("%s Entering folder: %s%s%s\n", enterSymbol, colorGreen, name, colorReset)
			err = os.Chdir(name)
			if err != nil {
				fmt.Printf("%s Error changing to directory '%s': %v\n", errorSymbol, name, err)
				continue
			}

			// Recursively create folders in this folder
			createFolders(reader, rootDir)

			// Create files in this folder
			createFiles(reader)

			// Return to the parent directory
			fmt.Printf("%s Returning to parent folder\n", returnSymbol)
			err = os.Chdir(currentDir)
			if err != nil {
				fmt.Printf("%s Error returning to parent directory: %v\n", errorSymbol, err)
				// Try to force return to root directory
				os.Chdir(rootDir)
			}
		}
	}
}

func createFiles(reader *bufio.Reader) {
	// Display current path
	displayCurrentPath()

	fmt.Printf("%s Would you like to create any files?\n(y/n): ", fileSymbol)
	response, _ := reader.ReadString('\n')

	if !isAffirmative(response) {
		return
	}

	fmt.Printf("%s How many files would you like to create? ", fileSymbol)
	countStr, _ := reader.ReadString('\n')
	countStr = strings.TrimSpace(countStr)
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		fmt.Printf("%s Invalid number of files. Skipping.\n", warningSymbol)
		return
	}

	for i := 1; i <= count; i++ {
		fmt.Printf("%s Enter name for file %d: ", enterSymbol, i)
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		if name == "" {
			fmt.Printf("%s Empty name provided for file %d. Skipping.\n", skipSymbol, i)
			continue
		}

		// Create an empty file
		file, err := os.Create(name)
		if err != nil {
			fmt.Printf("%s Error creating file '%s': %v\n", errorSymbol, name, err)
		} else {
			file.Close()
			fmt.Printf("%s %sCreated file: %s%s\n", successSymbol, colorYellow, name, colorReset)
		}
	}
}
