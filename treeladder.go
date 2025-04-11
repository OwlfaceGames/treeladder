package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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

	reader := bufio.NewReader(os.Stdin)

	// Ask if user wants to create files or folders
	fmt.Print("Would you like to create any files and/or folders? (yes/no): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "yes" && response != "y" {
		fmt.Println("No files or folders created. Exiting.")
		return
	}

	// Create folders
	createItems("folders", reader)

	// Create files
	createItems("files", reader)

	fmt.Println("Project structure created successfully!")
}

func createItems(itemType string, reader *bufio.Reader) {
	fmt.Printf("Would you like to create any %s? (yes/no): ", itemType)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "yes" && response != "y" {
		return
	}

	fmt.Printf("How many %s would you like to create? ", itemType)
	countStr, _ := reader.ReadString('\n')
	countStr = strings.TrimSpace(countStr)
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		fmt.Printf("Invalid number of %s. Skipping.\n", itemType)
		return
	}

	for i := 1; i <= count; i++ {
		fmt.Printf("Enter name for %s %d: ", itemType[:len(itemType)-1], i)
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		if name == "" {
			fmt.Printf("Empty name provided for %s %d. Skipping.\n", itemType[:len(itemType)-1], i)
			continue
		}

		if itemType == "folders" {
			err = os.Mkdir(name, 0o755)
		} else {
			// Create an empty file
			file, err := os.Create(name)
			if err == nil {
				file.Close()
			}
		}

		if err != nil {
			fmt.Printf("Error creating %s '%s': %v\n", itemType[:len(itemType)-1], name, err)
		} else {
			fmt.Printf("Created %s: %s\n", itemType[:len(itemType)-1], name)
		}
	}
}
