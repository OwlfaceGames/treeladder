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

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: treeladder create <repo_name>")
		return
	}

	command := args[1]

	if command == "create" {
		if len(args) < 3 {
			fmt.Println("Please provide a repository name")
			return
		}

		repoName := args[2]
		createRepo(repoName)
	} else {
		fmt.Println("Unknown command. Use 'treeladder create <repo_name>'")
	}
}

func createRepo(repoName string) {
	// Create the root directory
	err := os.Mkdir(repoName, 0o755)
	if err != nil {
		fmt.Printf("Error creating repository: %v\n", err)
		return
	}

	fmt.Printf("Created repository: %s\n", repoName)

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

	// Ask if user wants to create files or folders
	fmt.Print("Would you like to create any files and/or folders? (yes/no): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "yes" && response != "y" {
		fmt.Println("No files or folders created. Exiting.")
		return
	}

	// Create folders with recursive structure
	createFolders(reader, rootDir)

	// Create files at the root level
	createFiles(reader)

	// Ask if user wants to initialize git repository
	fmt.Print("Would you like to initialize a git repository? (yes/no): ")
	response, _ = reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response == "yes" || response == "y" {
		cmd := exec.Command("git", "init")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error initializing git repository: %v\n", err)
		} else {
			fmt.Println(strings.TrimSpace(string(output)))
		}
	}

	fmt.Println("Project structure created successfully!")
}

func createFolders(reader *bufio.Reader, rootDir string) {
	fmt.Print("Would you like to create any folders? (yes/no): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "yes" && response != "y" {
		return
	}

	fmt.Print("How many folders would you like to create? ")
	countStr, _ := reader.ReadString('\n')
	countStr = strings.TrimSpace(countStr)
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		fmt.Println("Invalid number of folders. Skipping.")
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
		fmt.Printf("Created folder: %s\n", name)

		// Ask if user wants to create content inside this folder
		fmt.Printf("Would you like to create content inside '%s'? (yes/no): ", name)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "yes" || response == "y" {
			// Change to the new folder
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
	fmt.Print("Would you like to create any files? (yes/no): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "yes" && response != "y" {
		return
	}

	fmt.Print("How many files would you like to create? ")
	countStr, _ := reader.ReadString('\n')
	countStr = strings.TrimSpace(countStr)
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		fmt.Println("Invalid number of files. Skipping.")
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
			fmt.Printf("Created file: %s\n", name)
		}
	}
}
