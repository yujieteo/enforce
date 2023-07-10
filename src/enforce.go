package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Create a dialog to select the project directory
	dialogFactory := &DirectoryDialogFactory{}
	dialog := dialogFactory.CreateDialog()
	projectPath, err := dialog.Browse()
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

	// Move files out of the selected directory into the main directory
	err = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			destPath := filepath.Join(projectPath, info.Name())
			moveOp := &MoveFileOperation{sourcePath: path, destPath: destPath}
			if err := moveOp.Execute(); err != nil {
				fmt.Println(err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	// Remove empty directories
	err = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			isEmpty, err := isDirectoryEmpty(path)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			if isEmpty {
				removeOp := &RemoveDirectoryOperation{dirPath: path}
				if err := removeOp.Execute(); err != nil {
					fmt.Println(err)
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	// Rename files in the main directory
	err = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			renameOp := &RenameFileOperation{filePath: path}
			if err := renameOp.Execute(); err != nil {
				fmt.Println(err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	// Create a directory structure
	projectDir := &RecursiveDirectory{Directory: &Directory{path: projectPath}}
	projectDir.AddOperation(&CreateDirectoryOperation{dirPath: filepath.Join(projectPath, "doc")})
	projectDir.AddOperation(&CreateDirectoryOperation{dirPath: filepath.Join(projectPath, "src")})
	projectDir.AddOperation(&CreateDirectoryOperation{dirPath: filepath.Join(projectPath, "job")})
	projectDir.AddOperation(&CreateDirectoryOperation{dirPath: filepath.Join(projectPath, "data")})
	projectDir.AddOperation(&CreateDirectoryOperation{dirPath: filepath.Join(projectPath, "ref")})
	projectDir.AddOperation(&CreateDirectoryOperation{dirPath: filepath.Join(projectPath, "media")})
	projectDir.AddOperation(&CreateDirectoryOperation{dirPath: filepath.Join(projectPath, "bin")})

	exampleDir := &RecursiveDirectory{Directory: &Directory{path: filepath.Join(projectPath, "doc", "report")}}
	projectDir.AddSubdirectory(exampleDir)

	for _, component := range []string{"doc", "src", "job", "data", "ref", "media", "bin"} {
		componentDir := &RecursiveDirectory{Directory: &Directory{path: filepath.Join(projectPath, component)}}
		projectDir.AddSubdirectory(componentDir)
	}

	// Move files to the project directory if the .git directory does not exist
	gitPath := filepath.Join(projectPath, ".git")
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		// Extract files to the project directory
		extractOp := &MoveFileOperation{
			sourcePath: projectPath,
			destPath:   projectPath,
		}
		projectDir.AddOperation(extractOp)

		// Rename files in the project directory
		renameOp := &RenameFileOperation{
			filePath: projectPath,
		}
		projectDir.AddOperation(renameOp)

		// Sort files in the project directory
		sorter := &FileSorter{
			FolderPath: projectPath,
		}
		projectDir.AddOperation(sorter)
	}

	// Execute all file operations
	err = projectDir.ExecuteOperations()
	if err != nil {
		fmt.Println("Error executing file operations:", err)
		return
	}

	// Initialize Git repository if it doesn't exist
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		cmd := exec.Command("git", "-C", projectPath, "init")
		err = cmd.Run()
		if err != nil {
			fmt.Println("Failed to initialize Git repository:", err)
			return
		}
		fmt.Println("Git repository initialized.")
	} else {
		fmt.Println("Git repository already exists. Files will not be sorted.")
	}

	// Create a .gitignore file
	textFileFactory := &TextFileFactory{
		ProjectPath: projectPath,
	}
	err = textFileFactory.CreateGitignore()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Program completed successfully.")
}
