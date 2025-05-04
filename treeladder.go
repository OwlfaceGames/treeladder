package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

// Symbols - using actual Nerd Font icons and keeping tree emoji
const (
	treeSymbol      = "🌳🪜"     // Tree for the main program (kept as emoji)
	folderSymbol    = "\uf07b" // Folder nerd font icon
	fileSymbol      = "\uf15b" // File nerd font icon
	questionSymbol  = "\uf059" // Question nerd font icon
	successSymbol   = "\uf00c" // Success/check nerd font icon
	errorSymbol     = "\uf00d" // Error/x nerd font icon
	pathSymbol      = "\uf07c" // Open folder/path nerd font icon
	gitSymbol       = "\uf1d3" // Git nerd font icon
	branchSymbol    = "\ue725" // Branch nerd font icon
	constructSymbol = "\uf085" // Cog/construction nerd font icon
	completeSymbol  = "🎉"      // Celebration/completion (kept as emoji)
	warningSymbol   = "\uf071" // Warning nerd font icon
	enterSymbol     = "\uf054" // Right arrow/enter nerd font icon
	skipSymbol      = "\uf05e" // Skip/prohibited nerd font icon
	returnSymbol    = "\uf112" // Return arrow nerd font icon
	tmuxSymbol      = "\uf120" // Terminal symbol for tmux
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
		fmt.Printf("%s Usage: treeladder create <repo_name>\n", warningSymbol)
		return
	}

	command := args[1]

	if command == "version" || command == "--version" || command == "-v" {
		fmt.Printf("%s TreeLadder %s\n", treeSymbol, Version)
		fmt.Printf("Build date: %s\n", BuildDate)
		fmt.Printf("Git commit: %s\n", GitCommit)
		return
	}

	if command == "create" {
		if len(args) < 3 {
			fmt.Printf("%s Please provide a repository name\n", warningSymbol)
			return
		}

		fmt.Printf("%s %sTreeLadder%s\n", treeSymbol, colorGreen, colorReset)
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

// Helper function to check if a command exists
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
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

	// Check if tmux is installed
	if commandExists("tmux") {
		// Ask if user wants to create a tmux session
		fmt.Printf("%s Would you like to create a tmux session for this project?\n(y/n): ", tmuxSymbol)
		response, _ = reader.ReadString('\n')

		if isAffirmative(response) {
			createTmuxSession(repoName, rootDir)
		}
	}

	fmt.Printf("%s Project structure created successfully! %s\n", treeSymbol, completeSymbol)
}

func displayCurrentPath() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s Error getting current directory: %v\n", errorSymbol, err)
		return
	}

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

func createTmuxSession(projectName string, projectPath string) {
	fmt.Printf("%s Creating tmux session for project: %s%s%s\n", tmuxSymbol, colorGreen, projectName, colorReset)

	// Sanitize project name for tmux session name (remove spaces and special characters)
	sessionName := strings.ReplaceAll(projectName, " ", "_")
	sessionName = strings.ReplaceAll(sessionName, ".", "_")

	// Check if session already exists
	checkCmd := exec.Command("tmux", "has-session", "-t", sessionName)
	err := checkCmd.Run()

	if err == nil {
		fmt.Printf("%s Session '%s' already exists. Attaching...\n", warningSymbol, sessionName)
		attachCmd := exec.Command("tmux", "attach-session", "-t", sessionName)
		attachCmd.Stdin = os.Stdin
		attachCmd.Stdout = os.Stdout
		attachCmd.Stderr = os.Stderr
		attachCmd.Run()
		return
	}

	// Create a new session with the first window named "code"
	newSessionCmd := exec.Command("tmux", "new-session", "-s", sessionName, "-n", "code", "-d")
	err = newSessionCmd.Run()
	if err != nil {
		fmt.Printf("%s Error creating tmux session: %v\n", errorSymbol, err)
		return
	}

	// Change to the project directory in the first pane
	sendKeysCmd := exec.Command("tmux", "send-keys", "-t", sessionName+":code", "cd "+projectPath, "C-m")
	sendKeysCmd.Run()

	// Split the window vertically
	splitCmd := exec.Command("tmux", "split-window", "-h", "-t", sessionName+":code")
	splitCmd.Run()

	// Change to the project directory in the right pane
	sendKeysCmd = exec.Command("tmux", "send-keys", "-t", sessionName+":code.2", "cd "+projectPath, "C-m")
	sendKeysCmd.Run()

	// Ask if user wants to open an editor
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s Would you like to open an editor in the second pane? (vim/nvim/code/none)\n(vim/nvim/code/none): ", questionSymbol)
	editorResponse, _ := reader.ReadString('\n')
	editorResponse = strings.TrimSpace(strings.ToLower(editorResponse))

	// Open the chosen editor
	if editorResponse == "vim" || editorResponse == "nvim" || editorResponse == "code" {
		sendKeysCmd = exec.Command("tmux", "send-keys", "-t", sessionName+":code", editorResponse, "C-m")
		sendKeysCmd.Run()
	}
	// Create a second window named "term"
	newWindowCmd := exec.Command("tmux", "new-window", "-t", sessionName, "-n", "term")
	newWindowCmd.Run()

	// Change to the project directory in the term window
	sendKeysCmd = exec.Command("tmux", "send-keys", "-t", sessionName+":term", "cd "+projectPath, "C-m")
	sendKeysCmd.Run()

	// Create a third window named "other"
	newWindowCmd = exec.Command("tmux", "new-window", "-t", sessionName, "-n", "other")
	newWindowCmd.Run()

	// Stay at home directory for the "other" window
	homeDir, _ := os.UserHomeDir()
	sendKeysCmd = exec.Command("tmux", "send-keys", "-t", sessionName+":other", "cd "+homeDir, "C-m")
	sendKeysCmd.Run()

	// Select the first window
	selectWindowCmd := exec.Command("tmux", "select-window", "-t", sessionName+":code")
	selectWindowCmd.Run()

	fmt.Printf("%s Tmux session '%s' created. Attaching...\n", successSymbol, sessionName)

	// Attach to the session
	attachCmd := exec.Command("tmux", "attach-session", "-t", sessionName)
	attachCmd.Stdin = os.Stdin
	attachCmd.Stdout = os.Stdout
	attachCmd.Stderr = os.Stderr
	attachCmd.Run()
}
