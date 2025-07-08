# TreeLadder

TreeLadder is a command-line tool written in Go that helps developers quickly create file and directory structures for new projects. It provides an interactive interface to create nested folders and files, allowing you to scaffold project structures with minimal effort. The tool guides you through the creation process step by step, making it easy to build complex directory hierarchies without having to manually create each folder and file.

## Usage

* Create a new project: `treeladder create project_name`
* Follow the interactive prompts to create folders and files
* For each folder, you can choose to create nested content
* Navigate through the creation process by answering yes/no questions
* Check the version: `treeladder version` or `treeladder -v`

## Installation

* **From source**:
  * Clone the repository: `git clone https://github.com/owlfacegames/treeladder.git`
  * Build the binary: `cd treeladder && go build`
  * Move to a directory in your PATH: `mv treeladder /usr/local/bin/`

* **Using Go install**:
  * Run: `go install github.com/yourusername/treeladder@latest`

* **Pre-built binaries**:
  * Download the appropriate binary for your platform from the releases page
  * Extract the archive
  * Move the binary to a location in your PATH
