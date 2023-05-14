package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sqweek/dialog"
)

func main() {
	// Use dialog to select the project directory
	projectPath, err := dialog.Directory().Title("Select project directory").Browse()
	if err != nil {
		fmt.Println("Failed to select project directory:", err)
		return
	}

	// Validate the project path exists
	_, err = os.Stat(projectPath)
	if os.IsNotExist(err) {
		fmt.Println("Project path does not exist.")
		return
	}

	components := []string{"doc", "src", "job", "data", "ref", "eg"}

	for _, component := range components {
		componentPath := filepath.Join(projectPath, component)
		err := os.MkdirAll(componentPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// Create the files inside the directories
	docFilePath := filepath.Join(projectPath, "doc", "bib-file.bib")
	err = os.WriteFile(docFilePath, []byte("This is the bib file"), 0644)
	if err != nil {
		panic(err)
	}

	readmeFilePath := filepath.Join(projectPath, "doc", "README.md")
	err = os.WriteFile(readmeFilePath, []byte("This is the doc file"), 0644)
	if err != nil {
		panic(err)
	}

	// Create an example
	exampleName := "example1"
	exampleDir := filepath.Join(projectPath, "eg", exampleName)
	err = os.MkdirAll(exampleDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	for _, component := range components {
		componentPath := filepath.Join(exampleDir, component)
		err := os.MkdirAll(componentPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	exampleReadmeFilePath := filepath.Join(exampleDir, "doc", "README.md")
	err = os.WriteFile(exampleReadmeFilePath, []byte("This is the doc file for the example"), 0644)
	if err != nil {
		panic(err)
	}

	// Create the large data directory
	dataLargeDir := filepath.Join(projectPath, "data", "large")
	err = os.MkdirAll(dataLargeDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	removeEmptySubdirectories(projectPath)

	// Check if Git repository already exists
	gitPath := filepath.Join(projectPath, ".git")
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		// Initialize Git repository
		removeEmptySubdirectories(projectPath)

		sortFiles(projectPath)

		cmd := exec.Command("git", "-C", projectPath, "init")
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
		fmt.Println("Git repository initialized.")
	} else {
		removeEmptySubdirectories(projectPath)
		fmt.Println("Git repository already exists.")
	}

	for _, component := range components {
		componentPath := filepath.Join(projectPath, component)
		removeEmptySubdirectories(projectPath)
		err := os.MkdirAll(componentPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	errJob := os.MkdirAll(filepath.Join(projectPath, "job"), os.ModePerm)
	if errJob != nil {
		panic(err)
	}

	err = generateREADME(projectPath)
	if err != nil {
		fmt.Printf("Error generating README file: %v", err)
		return
	}

	fmt.Println("README file generated successfully.")
}

func sortFiles(folderPath string) {

	// Check if the current folder is the .git repository
	if filepath.Base(folderPath) == ".git" {
		fmt.Printf("Skipping .git repository: '%s'\n", folderPath)
		return
	}

	subdirectories, err := getSubdirectories(folderPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, subdir := range subdirectories {
		subdirPath := filepath.Join(folderPath, subdir)

		// Check if the current subdirectory is the .git directory
		if subdir == ".git" {
			fmt.Printf("Skipping .git repository: '%s'\n", subdirPath)
			continue
		}

		sortFiles(subdirPath)
	}

	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		extension := strings.ToLower(filepath.Ext(path))
		destFolder := ""

		switch extension {
		case ".pdf", ".djvu", ".epub", ".html", ".mkv", ".mp4":
			destFolder = filepath.Join("ref", strings.TrimSuffix(filepath.Base(path), extension))
		case ".docx", ".md", ".tex", ".txt", ".doc", ".pptx":
			destFolder = "doc"
		case ".ipynb":
			destFolder = filepath.Join("eg", strings.TrimSuffix(filepath.Base(path), extension))
		case ".py", ".go", ".inp", ".c", ".m", ".for", ".cpp", ".java":
			destFolder = filepath.Join("src", strings.TrimSuffix(filepath.Base(path), extension))
		default:
			destFolder = "data"
		}

		destFolderPath := filepath.Join(folderPath, destFolder)
		err = os.MkdirAll(destFolderPath, 0755)
		if err != nil {
			return err
		}

		destFilePath := filepath.Join(destFolderPath, filepath.Base(path))
		err = os.Rename(path, destFilePath)
		if err != nil {
			return err
		}

		fmt.Printf("Moved '%s' to '%s'\n", path, destFilePath)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func getSubdirectories(folderPath string) ([]string, error) {
	var subdirectories []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != folderPath {
			subdirectories = append(subdirectories, strings.TrimPrefix(path, folderPath+string(os.PathSeparator)))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return subdirectories, nil
}

func removeEmptySubdirectories(folderPath string) {
	subdirectories := make([]string, 0)

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != folderPath {
			subdirectories = append(subdirectories, path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// Remove empty subdirectories in reverse order
	for i := len(subdirectories) - 1; i >= 0; i-- {
		dir := subdirectories[i]
		isEmpty, err := isDirectoryEmpty(dir)
		if err != nil {
			log.Println(err)
			continue
		}

		if isEmpty {
			err := os.Remove(dir)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Printf("Removed empty directory: %s\n", dir)
			}
		}
	}
}

func isDirectoryEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == nil {
		// Directory is not empty
		return false, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		// Directory does not exist
		return false, nil
	}

	if errors.Is(err, io.EOF) {
		// Directory is empty
		return true, nil
	}

	return false, err
}

func generateREADME(projectPath string) error {
	readmePath := projectPath + "/README.md"
	readmeContent := []byte(`# Project Name

One Paragraph of project description goes here

## Table of Contents

- [Project Name](#project-name)
  - [Table of Contents](#table-of-contents)
  - [About](#about)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
  - [Usage](#usage)
  - [Contributing](#contributing)
  - [License](#license)
  - [Acknowledgements](#acknowledgements)

## About

Provide a brief introduction or overview of your project.

## Getting Started

Instructions on setting up and running the project.

### Prerequisites

List any software, libraries, or dependencies that need to be installed before running the project.

### Installation

Step-by-step instructions on how to install the project.

## Usage

Provide examples or instructions on how to use the project.

## Contributing

Explain how others can contribute to your project. Include guidelines for pull requests and code style.

## License

Mention the license under which the project is distributed (e.g., MIT License).

## Acknowledgements

Give credits to any external resources or individuals whose work has influenced your project.`)

	err := ioutil.WriteFile(readmePath, readmeContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write README file: %v", err)
	}

	return nil
}
